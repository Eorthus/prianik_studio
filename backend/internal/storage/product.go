package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"

	"pryanik_studio/internal/models"

	"github.com/jmoiron/sqlx"
)

// GetProducts возвращает список товаров с пагинацией и фильтрацией
func (r *PostgresRepository) GetProducts(ctx context.Context, filter models.ProductFilter) (models.ProductList, error) {
	var result models.ProductList
	result.Page = filter.Page
	result.PageSize = filter.PageSize

	// Базовый запрос для подсчета общего количества товаров
	countQuery := `
    SELECT COUNT(*) 
    FROM products p
    JOIN product_translations pt ON p.id = pt.product_id
    WHERE pt.language = $1
    `

	// Базовый запрос для выборки товаров
	query := `
    SELECT p.id, p.category_id, p.subcategory_id, p.created_at, p.updated_at,
           pt.name, pt.description, pt.price, pt.currency
    FROM products p
    JOIN product_translations pt ON p.id = pt.product_id
    WHERE pt.language = $1
    `

	// Добавляем условия фильтрации
	args := []interface{}{filter.Language}
	argCount := 1

	if filter.CategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND p.category_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND p.category_id = $%d", argCount)
		args = append(args, *filter.CategoryID)
	}

	if filter.SubcategoryID != nil {
		argCount++
		query += fmt.Sprintf(" AND p.subcategory_id = $%d", argCount)
		countQuery += fmt.Sprintf(" AND p.subcategory_id = $%d", argCount)
		args = append(args, *filter.SubcategoryID)
	}

	if filter.Search != nil && *filter.Search != "" {
		argCount++
		searchTerm := "%" + *filter.Search + "%"
		query += fmt.Sprintf(" AND (pt.name ILIKE $%d OR pt.description ILIKE $%d)", argCount, argCount)
		countQuery += fmt.Sprintf(" AND (pt.name ILIKE $%d OR pt.description ILIKE $%d)", argCount, argCount)
		args = append(args, searchTerm)
	}

	// Добавляем сортировку по цене, если указана
	// Изменяем для сортировки по цене из product_translations, а не из products
	if filter.SortByPrice != nil {
		sortDirection := "ASC"
		if *filter.SortByPrice == "desc" {
			sortDirection = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY pt.price %s", sortDirection) // Сортируем по pt.price вместо p.price
	} else {
		// По умолчанию сортируем по ID
		query += " ORDER BY p.id DESC"
	}

	// Получаем общее количество товаров
	var totalItems int
	err := r.db.GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при получении общего количества товаров")
		return result, fmt.Errorf("ошибка при получении общего количества товаров: %w", err)
	}

	// Рассчитываем общее количество страниц
	result.TotalItems = totalItems
	result.TotalPages = int(math.Ceil(float64(totalItems) / float64(filter.PageSize)))

	// Добавляем пагинацию
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	// Запрос товаров
	var products []struct {
		ID            int64        `db:"id"`
		CategoryID    int64        `db:"category_id"`
		SubcategoryID int64        `db:"subcategory_id"`
		CreatedAt     sql.NullTime `db:"created_at"`
		UpdatedAt     sql.NullTime `db:"updated_at"`
		Name          string       `db:"name"`
		Description   string       `db:"description"`
		Price         float64      `db:"price"`
		Currency      string       `db:"currency"`
	}

	err = r.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при получении списка товаров")
		return result, fmt.Errorf("ошибка при получении списка товаров: %w", err)
	}

	// Преобразуем результаты запроса в модели
	result.Items = make([]models.Product, 0, len(products))
	for _, p := range products {
		product := models.Product{
			ID:            p.ID,
			CategoryID:    p.CategoryID,
			SubcategoryID: p.SubcategoryID,
			CreatedAt:     p.CreatedAt.Time,
			UpdatedAt:     p.UpdatedAt.Time,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			Currency:      p.Currency,
		}

		// Получаем изображения продукта
		var images []string
		query := `
        SELECT url FROM product_images 
        WHERE product_id = $1 
        ORDER BY is_main DESC, sort_order ASC
        `
		err = r.db.SelectContext(ctx, &images, query, p.ID)
		if err != nil {
			r.logger.WithError(err).Errorf("Ошибка при получении изображений для товара ID=%d", p.ID)
		}

		if len(images) > 0 {
			product.Images = images
		} else {
			product.Images = []string{"/default-product-image.jpg"}
		}

		// Получаем характеристики товара
		var characteristics []struct {
			Key   string `db:"key"`
			Value string `db:"value"`
		}

		query = `
        SELECT key, value 
        FROM product_characteristics 
        WHERE product_id = $1 AND language = $2
        `

		err = r.db.SelectContext(ctx, &characteristics, query, p.ID, filter.Language)
		if err != nil {
			r.logger.WithError(err).Errorf("Ошибка при получении характеристик товара ID=%d", p.ID)
		}

		// Если есть характеристики, добавляем их
		if len(characteristics) > 0 {
			product.Characteristics = make(map[string]string)
			for _, c := range characteristics {
				product.Characteristics[c.Key] = c.Value
			}
		}

		result.Items = append(result.Items, product)
	}

	return result, nil
}

// GetProductByID возвращает детальную информацию о товаре по его ID
func (r *PostgresRepository) GetProductByID(ctx context.Context, id int64, language string) (models.ProductDetail, error) {
	var result models.ProductDetail

	// Получаем основную информацию о товаре
	query := `
    SELECT p.id, p.category_id, p.subcategory_id, p.created_at, p.updated_at,
           pt.name, pt.description, pt.price, pt.currency
    FROM products p
    JOIN product_translations pt ON p.id = pt.product_id
    WHERE p.id = $1 AND pt.language = $2
	`

	var product struct {
		ID            int64        `db:"id"`
		CategoryID    int64        `db:"category_id"`
		SubcategoryID int64        `db:"subcategory_id"`
		CreatedAt     sql.NullTime `db:"created_at"`
		UpdatedAt     sql.NullTime `db:"updated_at"`
		Name          string       `db:"name"`
		Description   string       `db:"description"`
		Price         float64      `db:"price"`
		Currency      string       `db:"currency"`
	}

	err := r.db.GetContext(ctx, &product, query, id, language)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, fmt.Errorf("товар с ID=%d не найден", id)
		}
		r.logger.WithError(err).Errorf("Ошибка при получении товара ID=%d", id)
		return result, fmt.Errorf("ошибка при получении товара: %w", err)
	}

	// Заполняем основные поля
	result.ID = product.ID
	result.CategoryID = product.CategoryID
	result.SubcategoryID = product.SubcategoryID
	result.CreatedAt = product.CreatedAt.Time
	result.UpdatedAt = product.UpdatedAt.Time
	result.Name = product.Name
	result.Description = product.Description
	result.Price = product.Price
	result.Currency = product.Currency

	// Получаем характеристики товара
	var characteristics []struct {
		Key   string `db:"key"`
		Value string `db:"value"`
	}

	query = `
    SELECT key, value 
    FROM product_characteristics 
    WHERE product_id = $1 AND language = $2
    `

	err = r.db.SelectContext(ctx, &characteristics, query, id, language)
	if err != nil {
		r.logger.WithError(err).Errorf("Ошибка при получении характеристик товара ID=%d", id)
	}

	// Если есть характеристики, добавляем их
	if len(characteristics) > 0 {
		result.Characteristics = make(map[string]string)
		for _, c := range characteristics {
			result.Characteristics[c.Key] = c.Value
		}
	}

	// Получаем изображения товара
	var images []string
	query = `
	SELECT url FROM product_images 
	WHERE product_id = $1 
	ORDER BY is_main DESC, sort_order ASC
	`

	err = r.db.SelectContext(ctx, &images, query, id)
	if err != nil {
		r.logger.WithError(err).Errorf("Ошибка при получении изображений товара ID=%d", id)
	}

	if len(images) > 0 {
		result.Images = images
	} else {
		result.Images = []string{"/default-product-image.jpg"}
	}

	// Получаем связанные товары
	relatedProducts, err := r.GetRelatedProducts(ctx, id, 5, language)
	if err != nil {
		r.logger.WithError(err).Errorf("Ошибка при получении связанных товаров для ID=%d", id)
	} else {
		result.RelatedProducts = relatedProducts
	}

	return result, nil
}

// GetRelatedProducts возвращает список связанных товаров
func (r *PostgresRepository) GetRelatedProducts(ctx context.Context, productID int64, limit int, language string) ([]models.Product, error) {
	var result []models.Product

	// Получаем категорию текущего товара
	var categoryID int64
	query := "SELECT category_id FROM products WHERE id = $1"
	err := r.db.GetContext(ctx, &categoryID, query, productID)
	if err != nil {
		return result, fmt.Errorf("ошибка при получении категории товара: %w", err)
	}

	// Получаем товары из той же категории, кроме текущего
	query = `
    SELECT p.id, p.category_id, p.subcategory_id, p.created_at, p.updated_at,
           pt.name, pt.description, pt.price, pt.currency
    FROM products p
    JOIN product_translations pt ON p.id = pt.product_id
    WHERE p.category_id = $1 AND p.id != $2 AND pt.language = $3
    ORDER BY RANDOM()
    LIMIT $4
    `

	var products []struct {
		ID            int64        `db:"id"`
		CategoryID    int64        `db:"category_id"`
		SubcategoryID int64        `db:"subcategory_id"`
		CreatedAt     sql.NullTime `db:"created_at"`
		UpdatedAt     sql.NullTime `db:"updated_at"`
		Name          string       `db:"name"`
		Description   string       `db:"description"`
		Price         float64      `db:"price"`
		Currency      string       `db:"currency"`
	}

	err = r.db.SelectContext(ctx, &products, query, categoryID, productID, language, limit)
	if err != nil {
		return result, fmt.Errorf("ошибка при получении связанных товаров: %w", err)
	}

	// Преобразуем результаты запроса в модели
	for _, p := range products {
		product := models.Product{
			ID:            p.ID,
			CategoryID:    p.CategoryID,
			SubcategoryID: p.SubcategoryID,
			CreatedAt:     p.CreatedAt.Time,
			UpdatedAt:     p.UpdatedAt.Time,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			Currency:      p.Currency,
		}

		// Получаем основное изображение товара
		var image string
		imageQuery := `
		SELECT url FROM product_images 
		WHERE product_id = $1 AND is_main = true
		LIMIT 1
		`

		err = r.db.GetContext(ctx, &image, imageQuery, p.ID)
		if err != nil && err != sql.ErrNoRows {
			r.logger.WithError(err).Errorf("Ошибка при получении изображения для товара ID=%d", p.ID)
		}

		if image != "" {
			product.Images = []string{image}
		} else {
			// Если нет основного изображения, берем первое доступное
			imageQuery = `SELECT url FROM product_images WHERE product_id = $1 LIMIT 1`
			err = r.db.GetContext(ctx, &image, imageQuery, p.ID)
			if err != nil && err != sql.ErrNoRows {
				r.logger.WithError(err).Errorf("Ошибка при получении изображения для товара ID=%d", p.ID)
			}

			if image != "" {
				product.Images = []string{image}
			} else {
				product.Images = []string{"/default-product-image.jpg"}
			}
		}

		result = append(result, product)
	}

	return result, nil
}

// GetCategories возвращает список категорий с подкатегориями
func (r *PostgresRepository) GetCategories(ctx context.Context, language string) ([]models.Category, error) {
	var result []models.Category

	// Получаем все категории верхнего уровня
	query := `
    SELECT c.id, c.parent_id, c.created_at, c.updated_at, ct.name
    FROM categories c
    JOIN category_translations ct ON c.id = ct.category_id
    WHERE c.parent_id IS NULL AND ct.language = $1
    ORDER BY ct.name
    `

	var categories []struct {
		ID        int64         `db:"id"`
		ParentID  sql.NullInt64 `db:"parent_id"`
		CreatedAt sql.NullTime  `db:"created_at"`
		UpdatedAt sql.NullTime  `db:"updated_at"`
		Name      string        `db:"name"`
	}

	err := r.db.SelectContext(ctx, &categories, query, language)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при получении списка категорий")
		return result, fmt.Errorf("ошибка при получении списка категорий: %w", err)
	}

	// Преобразуем результаты запроса в модели
	for _, c := range categories {
		category := models.Category{
			ID:        c.ID,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
			// Прямое присвоение имени
			Name: c.Name,
		}

		// Получаем подкатегории для текущей категории
		subcategories, err := r.getSubcategories(ctx, c.ID, language)
		if err != nil {
			r.logger.WithError(err).Errorf("Ошибка при получении подкатегорий для категории ID=%d", c.ID)
		} else {
			category.Subcategories = subcategories
		}

		result = append(result, category)
	}

	return result, nil
}

// getSubcategories возвращает список подкатегорий для указанной категории
func (r *PostgresRepository) getSubcategories(ctx context.Context, parentID int64, language string) ([]models.Category, error) {
	var result []models.Category

	// Получаем все подкатегории для указанной категории
	query := `
    SELECT c.id, c.parent_id, c.created_at, c.updated_at, ct.name
    FROM categories c
    JOIN category_translations ct ON c.id = ct.category_id
    WHERE c.parent_id = $1 AND ct.language = $2
    ORDER BY ct.name
    `

	var subcategories []struct {
		ID        int64        `db:"id"`
		ParentID  int64        `db:"parent_id"`
		CreatedAt sql.NullTime `db:"created_at"`
		UpdatedAt sql.NullTime `db:"updated_at"`
		Name      string       `db:"name"`
	}

	err := r.db.SelectContext(ctx, &subcategories, query, parentID, language)
	if err != nil {
		return result, fmt.Errorf("ошибка при получении подкатегорий: %w", err)
	}

	// Преобразуем результаты запроса в модели
	for _, c := range subcategories {
		category := models.Category{
			ID:        c.ID,
			ParentID:  c.ParentID,
			CreatedAt: c.CreatedAt.Time,
			UpdatedAt: c.UpdatedAt.Time,
			// Прямое присвоение имени
			Name: c.Name,
		}

		result = append(result, category)
	}

	return result, nil
}

// CreateProduct создает новый товар в базе данных
func (r *PostgresRepository) CreateProduct(ctx context.Context, product *models.Product) (int64, error) {
	// Начинаем транзакцию
	tx, err := r.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при начале транзакции для создания товара")
		return 0, fmt.Errorf("ошибка при начале транзакции: %w", err)
	}

	// Добавляем отложенную функцию для отката транзакции в случае ошибки
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.logger.WithError(rollbackErr).Error("Ошибка при откате транзакции")
			}
		}
	}()

	// Устанавливаем время создания и обновления
	now := time.Now()
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	if product.UpdatedAt.IsZero() {
		product.UpdatedAt = now
	}

	// Вставляем товар (без поля price)
	query := `
    INSERT INTO products (category_id, subcategory_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4)
    RETURNING id
    `

	var subcategoryID *int64
	if product.SubcategoryID != 0 {
		subcategoryID = &product.SubcategoryID
	}

	var productID int64
	err = tx.QueryRowContext(
		ctx,
		query,
		product.CategoryID,
		subcategoryID,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&productID)

	if err != nil {
		r.logger.WithError(err).Error("Ошибка при создании товара")
		return 0, fmt.Errorf("ошибка при создании товара: %w", err)
	}

	// Добавляем переводы товара с ценами
	for lang, translation := range product.Translations {
		query = `
        INSERT INTO product_translations (product_id, language, name, description, price, currency)
        VALUES ($1, $2, $3, $4, $5, $6)
        `

		_, err = tx.ExecContext(
			ctx,
			query,
			productID,
			lang,
			translation.Name,
			translation.Description,
			translation.Price,
			translation.Currency,
		)

		if err != nil {
			r.logger.WithError(err).Errorf("Ошибка при добавлении перевода для языка %s", lang)
			return 0, fmt.Errorf("ошибка при добавлении перевода: %w", err)
		}

		// Добавляем характеристики, если они есть
		if translation.Characteristics != nil && len(translation.Characteristics) > 0 {
			for key, value := range translation.Characteristics {
				query = `
                INSERT INTO product_characteristics (product_id, language, key, value)
                VALUES ($1, $2, $3, $4)
                `

				_, err = tx.ExecContext(
					ctx,
					query,
					productID,
					lang,
					key,
					value,
				)

				if err != nil {
					r.logger.WithError(err).Errorf("Ошибка при добавлении характеристики %s", key)
					return 0, fmt.Errorf("ошибка при добавлении характеристики: %w", err)
				}
			}
		}
	}

	// Добавляем изображения, если они есть
	if len(product.Images) > 0 {
		for i, imageURL := range product.Images {
			isMain := i == 0 // Первое изображение - основное

			query = `
            INSERT INTO product_images (product_id, url, is_main, sort_order, created_at)
            VALUES ($1, $2, $3, $4, $5)
            `

			_, err = tx.ExecContext(
				ctx,
				query,
				productID,
				imageURL,
				isMain,
				i, // sort_order - порядок сортировки
				now,
			)

			if err != nil {
				r.logger.WithError(err).Errorf("Ошибка при добавлении изображения %s", imageURL)
				return 0, fmt.Errorf("ошибка при добавлении изображения: %w", err)
			}
		}
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Ошибка при фиксации транзакции")
		return 0, fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return productID, nil
}
