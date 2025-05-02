package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"pryanik_studio/internal/models"
	"pryanik_studio/internal/storage"
	"pryanik_studio/internal/utils"
)

// OrderHandler обработчик запросов для заказов и форм обратной связи
type OrderHandler struct {
	repo        storage.OrderRepository
	productRepo storage.ProductRepository
	emailSender utils.Sender
	validator   *validator.Validate
	logger      *logrus.Logger
}

// NewOrderHandler создает новый экземпляр OrderHandler
func NewOrderHandler(
	repo storage.OrderRepository,
	productRepo storage.ProductRepository,
	emailSender utils.Sender,
	logger *logrus.Logger,
) *OrderHandler {
	return &OrderHandler{
		repo:        repo,
		productRepo: productRepo,
		emailSender: emailSender,
		validator:   validator.New(),
		logger:      logger,
	}
}

// CreateOrder обработчик для создания заказа
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request models.OrderRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Ошибка при разборе JSON запроса на создание заказа")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный формат запроса"))
		return
	}

	// Валидируем данные формы
	if err := h.validator.Struct(request); err != nil {
		h.logger.WithError(err).Error("Ошибка валидации данных заказа")

		// Формируем детальные ошибки валидации
		var validationErrors []models.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, models.ValidationError{
				Field:   err.Field(),
				Message: getValidationErrorMessage(err),
			})
		}

		c.JSON(http.StatusBadRequest, models.NewValidationErrorResponse(validationErrors))
		return
	}

	// Если язык не указан, используем язык из заголовка Accept-Language или "ru" по умолчанию
	if request.Language == "" {
		request.Language = getPreferredLanguage(c)
	}

	// Создаем объект заказа
	order := &models.Order{
		Name:      request.Name,
		Email:     request.Email,
		Phone:     request.Phone,
		Comment:   request.Comment,
		Language:  request.Language,
		Status:    "new",
		TotalCost: 0,
		Items:     request.Items,
	}

	// Если в заказе есть товары, проверяем их и рассчитываем общую стоимость
	if len(request.Items) > 0 {
		// Для хранения обработанных товаров
		var processedItems []models.OrderItem
		totalCost := 0.0

		// Получаем информацию о каждом товаре и рассчитываем стоимость
		for _, item := range request.Items {
			// Получаем товар из базы данных для проверки наличия и цены
			product, err := h.productRepo.GetProductByID(c.Request.Context(), item.ProductID, request.Language)
			if err != nil {
				h.logger.WithError(err).Errorf("Товар с ID=%d не найден", item.ProductID)
				c.JSON(http.StatusBadRequest, models.NewErrorResponse("Один из товаров не найден"))
				return
			}

			// Заполняем дополнительные данные о товаре
			item.Price = product.Price
			item.ProductName = product.Name
			if len(product.Images) > 0 {
				item.ProductImage = product.Images[0]
			}

			// Рассчитываем стоимость позиции
			itemCost := product.Price * float64(item.Quantity)
			totalCost += itemCost

			processedItems = append(processedItems, item)
		}

		// Обновляем заказ с обработанными товарами и общей стоимостью
		order.Items = processedItems
		order.TotalCost = totalCost
	}

	// Сохраняем заказ в базе данных
	orderID, err := h.repo.CreateOrder(c.Request.Context(), order)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка при создании заказа")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при создании заказа"))
		return
	}

	// Получаем полную информацию о заказе для отправки email
	createdOrder, err := h.repo.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		h.logger.WithError(err).Errorf("Ошибка при получении созданного заказа ID=%d", orderID)
		// Продолжаем выполнение, так как заказ уже создан
	} else {
		// Отправляем уведомление о заказе по email
		if err := h.emailSender.SendOrderConfirmation(&createdOrder); err != nil {
			h.logger.WithError(err).Errorf("Ошибка при отправке уведомления о заказе ID=%d", orderID)
			// Продолжаем выполнение, так как основная операция создания заказа уже выполнена успешно
		}
	}

	// Возвращаем успешный ответ с сообщением на соответствующем языке
	message := getOrderSuccessMessage(request.Language)
	c.JSON(http.StatusOK, models.OrderResponse{
		Success: true,
		OrderID: orderID,
		Message: message,
	})
}

// SubmitContactForm обработчик для отправки формы обратной связи
func (h *OrderHandler) SubmitContactForm(c *gin.Context) {
	var request models.ContactFormRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Ошибка при разборе JSON запроса формы обратной связи")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный формат запроса"))
		return
	}

	// Валидируем данные формы
	if err := h.validator.Struct(request); err != nil {
		h.logger.WithError(err).Error("Ошибка валидации данных формы обратной связи")

		// Формируем детальные ошибки валидации
		var validationErrors []models.ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, models.ValidationError{
				Field:   err.Field(),
				Message: getValidationErrorMessage(err),
			})
		}

		c.JSON(http.StatusBadRequest, models.NewValidationErrorResponse(validationErrors))
		return
	}

	// Если язык не указан, используем язык из заголовка Accept-Language или "ru" по умолчанию
	if request.Language == "" {
		request.Language = getPreferredLanguage(c)
	}

	// Отправляем уведомление по email
	if err := h.emailSender.SendContactForm(&request); err != nil {
		h.logger.WithError(err).Error("Ошибка при отправке уведомления о сообщении")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при отправке сообщения"))
		return
	}

	// Возвращаем успешный ответ с сообщением на соответствующем языке
	message := getContactSuccessMessage(request.Language)
	c.JSON(http.StatusOK, models.NewSuccessResponse(map[string]interface{}{
		"message": message,
	}))
}

// getPreferredLanguage определяет предпочтительный язык пользователя
func getPreferredLanguage(c *gin.Context) string {
	// Получаем заголовок Accept-Language
	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		return "ru" // По умолчанию русский
	}

	// Разбираем заголовок (простая реализация)
	// Формат: ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7
	languages := []string{"ru", "en", "es"} // Поддерживаемые языки

	// Проверяем наличие поддерживаемых языков в заголовке
	for _, lang := range languages {
		if len(acceptLanguage) >= len(lang) && acceptLanguage[:len(lang)] == lang {
			return lang
		}
	}

	return "ru" // По умолчанию русский
}

// getOrderSuccessMessage возвращает сообщение об успешном создании заказа на нужном языке
func getOrderSuccessMessage(lang string) string {
	switch lang {
	case "en":
		return "Order successfully created"
	case "es":
		return "Pedido creado exitosamente"
	default:
		return "Заказ успешно создан"
	}
}

// getContactSuccessMessage возвращает сообщение об успешной отправке формы на нужном языке
func getContactSuccessMessage(lang string) string {
	switch lang {
	case "en":
		return "Message sent successfully"
	case "es":
		return "Mensaje enviado exitosamente"
	default:
		return "Сообщение успешно отправлено"
	}
}

// getValidationErrorMessage возвращает сообщение об ошибке валидации на нужном языке
func getValidationErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Поле обязательно для заполнения"
	case "email":
		return "Некорректный формат email"
	case "min":
		return "Значение должно быть больше или равно " + err.Param()
	case "max":
		return "Значение должно быть меньше или равно " + err.Param()
	default:
		return "Некорректное значение"
	}
}
