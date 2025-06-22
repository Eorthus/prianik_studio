package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"pryanik_studio/internal/models"
	"pryanik_studio/internal/storage"
)

// GalleryHandler обработчик запросов для галереи
type GalleryHandler struct {
	repo   storage.GalleryRepository
	logger *logrus.Logger
}

// NewGalleryHandler создает новый экземпляр GalleryHandler
func NewGalleryHandler(repo storage.GalleryRepository, logger *logrus.Logger) *GalleryHandler {
	return &GalleryHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetGalleryItems обработчик для получения элементов галереи
func (h *GalleryHandler) GetGalleryItems(c *gin.Context) {
	var filter models.GalleryFilter

	// Получаем язык из запроса, по умолчанию "ru"
	filter.Language = c.DefaultQuery("language", "ru")

	// Получаем параметры фильтрации по категории
	categoryIDStr := c.Query("category")
	if categoryIDStr != "" {
		categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
		if err == nil {
			filter.CategoryID = &categoryID
		}
	}

	// Логируем параметры запроса
	h.logger.WithFields(logrus.Fields{
		"language":   filter.Language,
		"categoryID": filter.CategoryID,
	}).Info("Получение элементов галереи")

	// Получаем список элементов галереи из репозитория
	galleryList, err := h.repo.GetGalleryItems(c.Request.Context(), filter)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка при получении элементов галереи")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при получении элементов галереи"))
		return
	}

	// Логируем количество найденных элементов
	h.logger.WithField("count", len(galleryList.Items)).Info("Найдены элементы галереи")

	// Проверка на пустой список
	if len(galleryList.Items) == 0 {
		// Возвращаем пустой массив вместо null
		c.JSON(http.StatusOK, models.NewSuccessResponse([]models.GalleryItem{}))
		return
	}

	// Возвращаем список элементов
	c.JSON(http.StatusOK, models.NewSuccessResponse(galleryList.Items))
}

// CreateGalleryItem обработчик для создания нового элемента галереи
func (h *GalleryHandler) CreateGalleryItem(c *gin.Context) {
	var request models.GalleryItemCreateRequest

	// Парсим JSON из тела запроса
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.WithError(err).Error("Ошибка при разборе JSON запроса на создание элемента галереи")
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный формат запроса"))
		return
	}

	// Проверка наличия переводов для русского языка (обязательно)
	if _, ok := request.Translations["ru"]; !ok {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Отсутствует обязательный перевод для русского языка"))
		return
	}

	// Создаем модель элемента галереи из запроса
	galleryItem := &models.GalleryItem{
		CategoryID:   request.CategoryID,
		Thumbnail:    request.Thumbnail,
		FullImage:    request.FullImage,
		Translations: make(map[string]*models.GalleryItemTranslation),
	}

	// Преобразуем переводы
	for lang, translation := range request.Translations {
		galleryItem.Translations[lang] = &models.GalleryItemTranslation{
			Title:       translation.Title,
			Description: translation.Description,
		}
	}

	// Добавляем элемент галереи в базу данных
	itemID, err := h.repo.CreateGalleryItem(c.Request.Context(), galleryItem)
	if err != nil {
		h.logger.WithError(err).Error("Ошибка при создании элемента галереи")
		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при создании элемента галереи"))
		return
	}

	// Возвращаем успешный ответ с ID созданного элемента
	response := map[string]interface{}{
		"id":      itemID,
		"message": "Элемент галереи успешно создан",
	}

	c.JSON(http.StatusCreated, models.NewSuccessResponse(response))
}

func (h *GalleryHandler) DeleteGalleryItem(c *gin.Context) {
	// Получаем ID элемента из URL
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewErrorResponse("Некорректный ID элемента"))
		return
	}

	// Удаляем элемент из базы данных
	err = h.repo.DeleteGalleryItem(c.Request.Context(), id)
	if err != nil {
		// Проверяем, не найден ли элемент
		if err.Error() == "gallery item not found" {
			c.JSON(http.StatusNotFound, models.NewErrorResponse("Элемент галереи не найден"))
			return
		}

		c.JSON(http.StatusInternalServerError, models.NewErrorResponse("Ошибка при удалении элемента галереи"))
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, models.NewSuccessResponse(map[string]interface{}{
		"message": "Элемент галереи успешно удален",
		"id":      id,
	}))
}
