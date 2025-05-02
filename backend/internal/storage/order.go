package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"pryanik_studio/internal/models"
)

// CreateOrder создает новый заказ в базе данных
func (r *PostgresRepository) CreateOrder(ctx context.Context, order *models.Order) (int64, error) {
	// Начинаем транзакцию
	tx, err := r.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при начале транзакции для создания заказа")
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

	// Устанавливаем язык по умолчанию, если он не задан
	if order.Language == "" {
		order.Language = "ru"
	}

	// Устанавливаем статус, если он не задан
	if order.Status == "" {
		order.Status = "new"
	}

	// Устанавливаем время создания и обновления
	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	if order.UpdatedAt.IsZero() {
		order.UpdatedAt = now
	}

	// Вставляем заказ
	query := `
	INSERT INTO orders (name, email, phone, comment, status, total_cost, language, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id
	`

	var orderID int64
	err = tx.QueryRowContext(
		ctx,
		query,
		order.Name,
		order.Email,
		order.Phone,
		order.Comment,
		order.Status,
		order.TotalCost,
		order.Language,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&orderID)

	if err != nil {
		r.logger.WithError(err).Error("Ошибка при создании заказа")
		return 0, fmt.Errorf("ошибка при создании заказа: %w", err)
	}

	// Вставляем товары заказа, если они есть
	if len(order.Items) > 0 {
		// SQL-запрос для вставки товаров заказа
		itemsQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)
		`

		for _, item := range order.Items {
			_, err = tx.ExecContext(
				ctx,
				itemsQuery,
				orderID,
				item.ProductID,
				item.Quantity,
				item.Price,
			)

			if err != nil {
				r.logger.WithError(err).Errorf("Ошибка при добавлении товара (ID=%d) в заказ", item.ProductID)
				return 0, fmt.Errorf("ошибка при добавлении товара в заказ: %w", err)
			}
		}
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Ошибка при фиксации транзакции")
		return 0, fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return orderID, nil
}

// GetOrderByID возвращает заказ по его ID
func (r *PostgresRepository) GetOrderByID(ctx context.Context, id int64) (models.Order, error) {
	var order models.Order

	// Получаем основную информацию о заказе
	query := `
	SELECT id, name, email, phone, comment, status, total_cost, language, created_at, updated_at
	FROM orders
	WHERE id = $1
	`

	var result struct {
		ID        int64          `db:"id"`
		Name      string         `db:"name"`
		Email     string         `db:"email"`
		Phone     string         `db:"phone"`
		Comment   string         `db:"comment"`
		Status    string         `db:"status"`
		TotalCost float64        `db:"total_cost"`
		Language  sql.NullString `db:"language"`
		CreatedAt sql.NullTime   `db:"created_at"`
		UpdatedAt sql.NullTime   `db:"updated_at"`
	}

	err := r.db.GetContext(ctx, &result, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, fmt.Errorf("заказ с ID=%d не найден", id)
		}
		r.logger.WithError(err).Errorf("Ошибка при получении заказа ID=%d", id)
		return order, fmt.Errorf("ошибка при получении заказа: %w", err)
	}

	// Заполняем основные поля заказа
	order.ID = result.ID
	order.Name = result.Name
	order.Email = result.Email
	order.Phone = result.Phone
	order.Comment = result.Comment
	order.Status = result.Status
	order.TotalCost = result.TotalCost

	// Устанавливаем язык, учитывая возможность NULL значения
	if result.Language.Valid {
		order.Language = result.Language.String
	} else {
		order.Language = "ru" // Язык по умолчанию
	}

	// Преобразуем NullTime в Time, учитывая возможность NULL значений
	if result.CreatedAt.Valid {
		order.CreatedAt = result.CreatedAt.Time
	}
	if result.UpdatedAt.Valid {
		order.UpdatedAt = result.UpdatedAt.Time
	}

	// Получаем товары заказа
	itemsQuery := `
	SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price,
	       pt.name as product_name
	FROM order_items oi
	LEFT JOIN product_translations pt ON oi.product_id = pt.product_id AND pt.language = $2
	WHERE oi.order_id = $1
	ORDER BY oi.id
	`

	var items []struct {
		ID          int64          `db:"id"`
		OrderID     int64          `db:"order_id"`
		ProductID   int64          `db:"product_id"`
		Quantity    int            `db:"quantity"`
		Price       float64        `db:"price"`
		ProductName sql.NullString `db:"product_name"`
	}

	err = r.db.SelectContext(ctx, &items, itemsQuery, id, order.Language)
	if err != nil {
		r.logger.WithError(err).Errorf("Ошибка при получении товаров заказа ID=%d", id)
		return order, fmt.Errorf("ошибка при получении товаров заказа: %w", err)
	}

	// Преобразуем результаты запроса в модели
	order.Items = make([]models.OrderItem, 0, len(items))
	for _, item := range items {
		orderItem := models.OrderItem{
			ID:        item.ID,
			OrderID:   item.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}

		// Устанавливаем наименование товара, если оно доступно
		if item.ProductName.Valid {
			orderItem.ProductName = item.ProductName.String
		}

		// Получаем изображение товара
		var image string
		imageQuery := `
		SELECT url FROM product_images 
		WHERE product_id = $1 AND is_main = true
		LIMIT 1
		`

		err = r.db.GetContext(ctx, &image, imageQuery, item.ProductID)
		if err != nil && err != sql.ErrNoRows {
			r.logger.WithError(err).Errorf("Ошибка при получении изображения для товара ID=%d", item.ProductID)
		}

		if image != "" {
			orderItem.ProductImage = image
		} else {
			// Если нет основного изображения, берем первое доступное
			imageQuery = `SELECT url FROM product_images WHERE product_id = $1 LIMIT 1`
			err = r.db.GetContext(ctx, &image, imageQuery, item.ProductID)
			if err != nil && err != sql.ErrNoRows {
				r.logger.WithError(err).Errorf("Ошибка при получении изображения для товара ID=%d", item.ProductID)
			}

			if image != "" {
				orderItem.ProductImage = image
			} else {
				// Если изображение не найдено, используем изображение по умолчанию
				orderItem.ProductImage = "/default-product-image.jpg"
			}
		}

		// Добавляем товар в список
		order.Items = append(order.Items, orderItem)
	}

	return order, nil
}
