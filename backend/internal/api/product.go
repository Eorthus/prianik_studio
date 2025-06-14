package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"pryanik_studio/internal/models"
	"pryanik_studio/internal/storage"
)

// ProductHandler обработчик запросов для товаров
type ProductHandler struct {
	repo   storage.ProductRepository
	logger *logrus.Logger
}

// NewProductHandler создает новый экземпляр ProductHandler
func NewProductHandler(repo storage.ProductRepository, logger *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetProducts обработчик для получения списка товаров
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var filter models.ProductFilter

	// Получаем язык из запроса, по умолчанию "ru"
	filter.Language = c.DefaultQuery("language", "ru")

	// Получаем параметры пагинации
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	// Преобразуем строковые параметры в числовые
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}
	filter.Page = pageNum

	pageSizeNum, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeNum < 1 {
		pageSizeNum = 10
	}
	filter.PageSize = pageSizeNum

	// Получаем параметры фильтрации
	categoryIDStr := c.Query("category")
	if categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err == nil {
			filter.CategoryID = &categoryID
		}
	}

	subcategoryIDStr := c.Query("subcategory")
	if subcategoryIDStr != "" {
		subcategoryID, err := strconv.ParseInt(subcategoryIDStr, 10, 64)
		if err == nil {
			filter.SubcategoryID = &subcategoryID
		}
	}

	search := c.Query("search")
	if search != "" {
		filter.Search = &search
	}

	// Получаем параметр сортировки по цене
	sortPrice := c.Query("sort_price")
	if sortPrice != "" && (sortPrice == "asc" || sortPrice == "desc") {
		filter.SortByPrice = &sortPrice
	}

	// Получаем список товаров из репозитория
	products, err := h.repo.GetProducts(c.Request.Context(), filter)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка при получении списка товаров")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при получении списка товаров"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(products))
}

// GetProductByID обработчик для получения товара по ID
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	// Получаем ID товара из URL
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный ID товара"))
		return
	}

	// Получаем язык из запроса, по умолчанию "ru"
	language := c.DefaultQuery("language", "ru")

	// Получаем товар из репозитория
	product, err := h.repo.GetProductByID(c.Request.Context(), id, language)
	if err != nil {
		h.logger.WithError(err).Errorf("Ошибка при получении товара ID=%d", id)
		c.JSON(http.StatusNotFound, models.NewErrorResponse("Товар не найден"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(product))
}

// GetCategories обработчик для получения списка категорий
func (h *ProductHandler) GetCategories(c *gin.Context) {
	// Получаем язык из запроса, по умолчанию "ru"
	language := c.DefaultQuery("language", "ru")

	// Получаем категории из репозитория
	categories, err := h.repo.GetCategories(c.Request.Context(), language)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка при получении списка категорий")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при получении списка категорий"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(categories))
}

// GetRelatedProducts обработчик для получения связанных товаров
func (h *ProductHandler) GetRelatedProducts(c *gin.Context) {
	// Получаем ID товара из URL
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный ID товара"))
		return
	}

	// Получаем язык из запроса, по умолчанию "ru"
	language := c.DefaultQuery("language", "ru")

	// Получаем лимит связанных товаров, по умолчанию 5
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 5
	}

	// Получаем связанные товары из репозитория
	products, err := h.repo.GetRelatedProducts(c.Request.Context(), id, limit, language)
	if err != nil {
		h.logger.WithError(err).Errorf("Ошибка при получении связанных товаров для ID=%d", id)
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при получении связанных товаров"))
		return
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(products))
}

// CreateProduct обработчик для создания нового товара
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var request models.ProductCreateRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Ошибка при разборе JSON запроса на создание товара")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный формат запроса"))
		return
	}

	// Проверка наличия переводов для русского языка (обязательно)
	if _, ok := request.Translations["ru"]; !ok {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Отсутствует обязательный перевод для русского языка"))
		return
	}

	// Создаем модель товара из запроса
	product := &models.Product{
		CategoryID:   request.CategoryID,
		Images:       request.Images,
		Translations: make(map[string]*models.ProductTranslation),
	}

	// Устанавливаем подкатегорию, если она указана
	if request.SubcategoryID != nil {
		product.SubcategoryID = request.SubcategoryID
	}

	// Преобразуем переводы
	for lang, translation := range request.Translations {
		product.Translations[lang] = &models.ProductTranslation{
			Name:            translation.Name,
			Description:     translation.Description,
			Price:           translation.Price,
			Currency:        translation.Currency,
			Characteristics: translation.Characteristics,
		}
	}

	// Добавляем товар в базу данных
	productID, err := h.repo.CreateProduct(c.Request.Context(), product)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка при создании товара")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при создании товара"))
		return
	}

	// Получаем созданный товар для возврата полной информации
	createdProduct, err := h.repo.GetProductByID(c.Request.Context(), productID, "ru")
	if err != nil {
		h.logger.WithError(err).Errorf("Ошибка при получении созданного товара ID=%d", productID)
		// Продолжаем выполнение, так как товар уже создан
	}

	// Возвращаем успешный ответ с ID созданного товара и, если доступно, с полной информацией
	response := map[string]interface{}{
		"id":      productID,
		"message": "Товар успешно создан",
	}

	if err == nil {
		response["product"] = createdProduct
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(response))
}

// UpdateProduct обработчик для обновления существующего товара
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	// Получаем ID товара из URL
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный ID товара"))
		return
	}

	// Проверяем существование товара в базе
	language := c.DefaultQuery("language", "ru") // Получаем язык из запроса
	_, err = h.repo.GetProductByID(c.Request.Context(), id, language)
	if err != nil {
		h.logger.WithError(err).Errorf("Товар с ID=%d не найден", id)
		c.JSON(http.StatusNotFound, models.NewErrorResponse("Товар не найден"))
		return
	}

	// Парсим JSON из тела запроса
	var request models.ProductUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Ошибка при разборе JSON запроса на обновление товара")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный формат запроса"))
		return
	}

	// Создаем модель товара для обновления
	product := &models.Product{
		ID:           id,
		Images:       request.Images,
		Translations: make(map[string]*models.ProductTranslation),
	}

	// Устанавливаем CategoryID, если он предоставлен
	if request.CategoryID != nil {
		product.CategoryID = *request.CategoryID
	}

	// Устанавливаем SubcategoryID, если он предоставлен
	if request.SubcategoryID != nil {
		product.SubcategoryID = request.SubcategoryID
	}

	// Преобразуем переводы, если они есть
	for lang, translation := range request.Translations {
		// Создаем перевод с значениями по умолчанию
		productTranslation := &models.ProductTranslation{
			Characteristics: translation.Characteristics,
		}

		// Устанавливаем поля, только если они предоставлены
		if translation.Name != nil {
			productTranslation.Name = *translation.Name
		}

		if translation.Description != nil {
			productTranslation.Description = *translation.Description
		}

		if translation.Price != nil {
			productTranslation.Price = *translation.Price
		}

		if translation.Currency != nil {
			productTranslation.Currency = *translation.Currency
		}

		// Добавляем перевод в карту переводов
		product.Translations[lang] = productTranslation
	}

	// Обновляем товар в базе данных
	err = h.repo.UpdateProduct(c.Request.Context(), product)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка при обновлении товара")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при обновлении товара"))
		return
	}

	// Получаем обновленный товар для возврата
	updatedProduct, err := h.repo.GetProductByID(c.Request.Context(), id, language)
	if err != nil {
		h.logger.WithError(err).Errorf("Ошибка при получении обновленного товара ID=%d", id)
		// Продолжаем выполнение, так как товар уже обновлен
	}

	// Возвращаем успешный ответ
	response := map[string]interface{}{
		"id":      id,
		"message": "Товар успешно обновлен",
	}

	if err == nil {
		response["product"] = updatedProduct
	}

	c.JSON(http.StatusOK, models.NewSuccessResponse(response))
}
