package models

import (
	"time"
)

// Product представляет модель товара
type Product struct {
	ID            int64 `json:"id" db:"id"`
	CategoryID    int64 `json:"category_id" db:"category_id"`
	SubcategoryID int64 `json:"subcategory_id,omitempty" db:"subcategory_id"`
	// Удаляем поле price из основной структуры
	Images    []string  `json:"images" db:"-"` // Массив URL изображений
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Переводимые поля прямо в структуре
	Name            string            `json:"name" db:"-"`
	Description     string            `json:"description" db:"-"`
	Price           float64           `json:"price" db:"-"`    // Добавляем price как переводимое поле
	Currency        string            `json:"currency" db:"-"` // Добавляем валюту
	Characteristics map[string]string `json:"characteristics,omitempty" db:"-"`

	// Сохраняем для внутреннего использования, но не возвращаем в API
	Translations map[string]*ProductTranslation `json:"-" db:"-"`
}

// ProductTranslation содержит переводимые поля товара
type ProductTranslation struct {
	ProductID       int64             `json:"-" db:"product_id"`
	Language        string            `json:"-" db:"language"`
	Name            string            `json:"name" db:"name"`
	Description     string            `json:"description" db:"description"`
	Price           float64           `json:"price" db:"price"`       // Добавляем цену в перевод
	Currency        string            `json:"currency" db:"currency"` // Добавляем валюту
	Characteristics map[string]string `json:"characteristics" db:"-"`
}

// Category представляет категорию товаров
type Category struct {
	ID        int64     `json:"id" db:"id"`
	ParentID  int64     `json:"parent_id,omitempty" db:"parent_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	// Переводимое поле прямо в структуре
	Name string `json:"name" db:"-"`

	// Сохраняем для внутреннего использования, но не возвращаем в API
	Translations map[string]*CategoryTranslation `json:"-" db:"-"`

	// Связанные данные
	Subcategories []Category `json:"subcategories,omitempty" db:"-"`
}

// CategoryTranslation содержит переводимые поля категории
type CategoryTranslation struct {
	CategoryID int64  `json:"-" db:"category_id"`
	Language   string `json:"-" db:"language"`
	Name       string `json:"name" db:"name"`
}

// ProductImage представляет изображение товара
type ProductImage struct {
	ID        int64     `json:"id" db:"id"`
	ProductID int64     `json:"-" db:"product_id"`
	URL       string    `json:"url" db:"url"`
	IsMain    bool      `json:"is_main" db:"is_main"`
	SortOrder int       `json:"-" db:"sort_order"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ProductList представляет структуру для возврата списка товаров с пагинацией
type ProductList struct {
	Items      []Product `json:"items"`
	TotalItems int       `json:"total_items"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
	TotalPages int       `json:"total_pages"`
}

// ProductCharacteristic представляет характеристику товара
type ProductCharacteristic struct {
	ProductID int64  `db:"product_id"`
	Language  string `db:"language"`
	Key       string `db:"key"`
	Value     string `db:"value"`
}

// ProductDetail представляет подробную информацию о товаре
type ProductDetail struct {
	Product
	RelatedProducts []Product `json:"related_products,omitempty"`
}

// ProductFilter содержит параметры фильтрации товаров
type ProductFilter struct {
	CategoryID    *int64  `form:"category"`
	SubcategoryID *int64  `form:"subcategory"`
	Search        *string `form:"search"`
	SortByPrice   *string `form:"sort_price"` // "asc", "desc" или пусто
	Page          int     `form:"page,default=1"`
	PageSize      int     `form:"page_size,default=10"`
	Language      string  `form:"language,default=ru"`
}

// ProductCreateRequest представляет запрос на создание товара
type ProductCreateRequest struct {
	CategoryID    int64                                       `json:"category_id" binding:"required"`
	SubcategoryID *int64                                      `json:"subcategory_id"`
	Images        []string                                    `json:"images"`
	Translations  map[string]*ProductTranslationCreateRequest `json:"translations" binding:"required"`
}

// ProductTranslationCreateRequest представляет перевод товара при создании
type ProductTranslationCreateRequest struct {
	Name            string            `json:"name" binding:"required"`
	Description     string            `json:"description"`
	Price           float64           `json:"price" binding:"required,min=0"` // Добавляем цену
	Currency        string            `json:"currency" binding:"required"`    // Добавляем валюту
	Characteristics map[string]string `json:"characteristics"`
}
