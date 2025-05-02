package models

import (
	"time"
)

type GalleryItem struct {
	ID         int64     `json:"id" db:"id"`
	CategoryID int64     `json:"category_id" db:"category_id"`
	Thumbnail  string    `json:"thumbnail" db:"thumbnail"`
	FullImage  string    `json:"full" db:"full_image"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`

	// Переводимые поля прямо в структуре
	Title       string `json:"title" db:"-"`
	Description string `json:"description" db:"-"`

	// Сохраняем для внутреннего использования, но не возвращаем в API
	Translations map[string]*GalleryItemTranslation `json:"-" db:"-"`
}

// GalleryItemTranslation содержит переводимые поля элемента галереи
type GalleryItemTranslation struct {
	GalleryItemID int64  `json:"-" db:"gallery_item_id"`
	Language      string `json:"-" db:"language"`
	Title         string `json:"title" db:"title"`
	Description   string `json:"description" db:"description"`
}

// GalleryFilter содержит параметры фильтрации галереи
type GalleryFilter struct {
	CategoryID *int64 `form:"category"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=15"`
	Language   string `form:"language,default=ru"`
}

// GalleryList представляет структуру для возврата списка элементов галереи с пагинацией
type GalleryList struct {
	Items      []GalleryItem `json:"items"`
	TotalItems int           `json:"total_items"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// GalleryItemCreateRequest представляет запрос на создание элемента галереи
type GalleryItemCreateRequest struct {
	CategoryID   int64                                     `json:"category_id" binding:"required"`
	Thumbnail    string                                    `json:"thumbnail" binding:"required"`
	FullImage    string                                    `json:"full_image" binding:"required"`
	Translations map[string]*GalleryItemTranslationRequest `json:"translations" binding:"required"`
}

// GalleryItemTranslationRequest представляет перевод элемента галереи при создании
type GalleryItemTranslationRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}
