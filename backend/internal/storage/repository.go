package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"pryanik_studio/internal/models"
)

// DatabaseConnection интерфейс для работы с базой данных
type DatabaseConnection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// Repository интерфейс для работы с хранилищем данных
type Repository interface {
	// Интерфейсы для работы с товарами
	ProductRepository

	// Интерфейсы для работы с галереей
	GalleryRepository

	// Интерфейсы для работы с заказами
	OrderRepository
}

// ProductRepository интерфейс для работы с товарами
type ProductRepository interface {
	GetProducts(ctx context.Context, filter models.ProductFilter) (models.ProductList, error)
	GetProductByID(ctx context.Context, id int64, language string) (models.ProductDetail, error)
	GetRelatedProducts(ctx context.Context, productID int64, limit int, language string) ([]models.Product, error)
	GetCategories(ctx context.Context, language string) ([]models.Category, error)
	CreateProduct(ctx context.Context, product *models.Product) (int64, error)
	UpdateProduct(ctx context.Context, product *models.Product) error
}

// GalleryRepository интерфейс для работы с галереей
type GalleryRepository interface {
	GetGalleryItems(ctx context.Context, filter models.GalleryFilter) (models.GalleryList, error)
	CreateGalleryItem(ctx context.Context, item *models.GalleryItem) (int64, error)
}

// OrderRepository интерфейс для работы с заказами
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *models.Order) (int64, error)
	GetOrderByID(ctx context.Context, id int64) (models.Order, error)
}

// PostgresRepository реализация Repository для PostgreSQL
type PostgresRepository struct {
	db     DatabaseConnection
	logger *logrus.Logger
}

// NewPostgresRepository создает новый PostgresRepository
func NewPostgresRepository(db DatabaseConnection, logger *logrus.Logger) *PostgresRepository {
	return &PostgresRepository{
		db:     db,
		logger: logger,
	}
}
