package storage

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер PostgreSQL
	"github.com/sirupsen/logrus"

	"pryanik_studio/internal/config"
)

// NewDatabase создает и инициализирует новое подключение к базе данных
func NewDatabase(cfg config.DatabaseConfig, logger *logrus.Logger) (*sqlx.DB, error) {
	// Сначала пробуем подключиться к серверу PostgreSQL без указания базы данных,
	// чтобы проверить, существует ли наша база данных и создать её при необходимости
	mainDSN := getMainDSN(cfg)

	mainDb, err := sql.Open("postgres", mainDSN)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к серверу PostgreSQL: %w", err)
	}
	defer mainDb.Close()

	// Проверяем соединение с сервером
	if err := mainDb.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с сервером PostgreSQL: %w", err)
	}

	// Проверяем существование базы данных
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = mainDb.QueryRow(query, cfg.DBName).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки существования базы данных: %w", err)
	}

	// Если база данных не существует, создаем её
	if !exists {
		logger.Infof("База данных '%s' не найдена, создаем новую", cfg.DBName)

		// Экранируем имя базы данных для безопасного использования в SQL-запросе
		// Это особенно важно, если имя базы данных может содержать спецсимволы
		escapedDBName := strings.Replace(cfg.DBName, "'", "''", -1)

		// Создаем базу данных
		_, err = mainDb.Exec(fmt.Sprintf("CREATE DATABASE %s", escapedDBName))
		if err != nil {
			return nil, fmt.Errorf("ошибка создания базы данных '%s': %w", cfg.DBName, err)
		}

		logger.Infof("База данных '%s' успешно создана", cfg.DBName)
	}

	// Теперь подключаемся к созданной/существующей базе данных
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Тестирование соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с базой данных: %w", err)
	}

	logger.Info("Успешное подключение к базе данных")
	return db, nil
}

// getMainDSN возвращает строку подключения к серверу PostgreSQL без указания базы данных
func getMainDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.SSLMode)
}

// MigrateDatabase выполняет необходимые миграции базы данных
func MigrateDatabase(db *sqlx.DB, logger *logrus.Logger) error {
	// Здесь создаем таблицы напрямую через SQL-запросы

	// Создание схемы для таблиц, если она не существует
	schema := `
	-- Категории
	CREATE TABLE IF NOT EXISTS categories (
		id SERIAL PRIMARY KEY,
		parent_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	-- Переводы категорий
	CREATE TABLE IF NOT EXISTS category_translations (
		category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
		language VARCHAR(5) NOT NULL,
		name VARCHAR(255) NOT NULL,
		PRIMARY KEY (category_id, language)
	);

	-- Товары
	CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
		subcategory_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	-- Переводы товаров
	CREATE TABLE IF NOT EXISTS product_translations (
		product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		language VARCHAR(5) NOT NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		PRIMARY KEY (product_id, language)
	);

	-- Характеристики товаров
	CREATE TABLE IF NOT EXISTS product_characteristics (
		product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		language VARCHAR(5) NOT NULL,
		key VARCHAR(100) NOT NULL,
		value TEXT NOT NULL,
		PRIMARY KEY (product_id, language, key)
	);

	-- Изображения товаров
	CREATE TABLE IF NOT EXISTS product_images (
		id SERIAL PRIMARY KEY,
		product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		url VARCHAR(255) NOT NULL,
		is_main BOOLEAN NOT NULL DEFAULT false,
		sort_order INTEGER NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	-- Элементы галереи
	CREATE TABLE IF NOT EXISTS gallery_items (
		id SERIAL PRIMARY KEY,
		category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
		thumbnail VARCHAR(255) NOT NULL,
		full_image VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	-- Переводы элементов галереи
	CREATE TABLE IF NOT EXISTS gallery_item_translations (
		gallery_item_id INTEGER NOT NULL REFERENCES gallery_items(id) ON DELETE CASCADE,
		language VARCHAR(5) NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		PRIMARY KEY (gallery_item_id, language)
	);

	-- Заказы
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		phone VARCHAR(50) NOT NULL,
		comment TEXT,
		language VARCHAR(5) DEFAULT 'ru',
		status VARCHAR(50) NOT NULL DEFAULT 'new',
		total_cost DECIMAL(10, 2) NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	-- Товары в заказе
	CREATE TABLE IF NOT EXISTS order_items (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
		quantity INTEGER NOT NULL DEFAULT 1,
		price DECIMAL(10, 2) NOT NULL
	);
	`

	// Выполняем SQL запрос для создания таблиц
	_, err := db.Exec(schema)
	if err != nil {
		logger.WithError(err).Error("Ошибка при выполнении миграций")
		return fmt.Errorf("ошибка при выполнении миграций: %w", err)
	}

	logger.Info("Миграции успешно выполнены")

	// Проверка наличия колонки language в таблице orders
	// Выполняем проверку и добавление колонки, если она не существует
	var languageExists bool
	err = db.Get(&languageExists, `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'orders' AND column_name = 'language'
		)
	`)

	if err != nil {
		logger.WithError(err).Error("Ошибка при проверке наличия колонки language в таблице orders")
		return fmt.Errorf("ошибка при проверке наличия колонки language: %w", err)
	}

	// Если колонка не существует, добавляем её
	if !languageExists {
		_, err = db.Exec("ALTER TABLE orders ADD COLUMN language VARCHAR(5) DEFAULT 'ru'")
		if err != nil {
			logger.WithError(err).Error("Ошибка при добавлении колонки language в таблицу orders")
			return fmt.Errorf("ошибка при добавлении колонки language: %w", err)
		}
		logger.Info("Колонка language добавлена в таблицу orders")
	}

	// Проверяем, есть ли уже категории в базе данных
	var count int
	err = db.Get(&count, "SELECT COUNT(*) FROM categories")
	if err != nil {
		logger.WithError(err).Error("Ошибка при проверке наличия категорий")
		return fmt.Errorf("ошибка при проверке наличия категорий: %w", err)
	}

	// Добавляем начальные данные, если таблица категорий пуста
	if count == 0 {
		if err := seedInitialData(db, logger); err != nil {
			logger.WithError(err).Error("Ошибка при добавлении начальных данных")
			return fmt.Errorf("ошибка при добавлении начальных данных: %w", err)
		}
	}

	// Проверяем, есть ли уже колонки price и currency в таблице product_translations
	var hasColumns bool
	err = db.Get(&hasColumns, `
    SELECT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'product_translations' 
        AND column_name = 'price'
    )
`)
	if err != nil {
		logger.WithError(err).Error("Ошибка при проверке наличия колонки price")
		return fmt.Errorf("ошибка при проверке наличия колонки: %w", err)
	}

	// Если колонки нет, то добавляем
	if !hasColumns {
		_, err = db.Exec(`
        -- Добавляем колонки price и currency в таблицу product_translations
        ALTER TABLE product_translations ADD COLUMN price DECIMAL(10, 2);
        ALTER TABLE product_translations ADD COLUMN currency VARCHAR(3);
        
        -- Копируем значение price из таблицы products в таблицу product_translations
        UPDATE product_translations pt
        SET price = (SELECT price FROM products WHERE id = pt.product_id),
            currency = CASE 
                WHEN pt.language = 'ru' THEN 'RUB'
                WHEN pt.language = 'en' THEN 'USD'
                WHEN pt.language = 'es' THEN 'EUR'
                ELSE 'USD'
            END;
        
        -- Делаем колонки NOT NULL
        ALTER TABLE product_translations ALTER COLUMN price SET NOT NULL;
        ALTER TABLE product_translations ALTER COLUMN currency SET NOT NULL;
    `)
		if err != nil {
			logger.WithError(err).Error("Ошибка при добавлении колонок price и currency")
			return fmt.Errorf("ошибка при добавлении колонок: %w", err)
		}

		logger.Info("Колонки price и currency успешно добавлены в таблицу product_translations")
	}

	return nil
}

// seedInitialData добавляет начальные данные в базу данных
func seedInitialData(db *sqlx.DB, logger *logrus.Logger) error {
	// Начинаем транзакцию для атомарной вставки всех данных
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Добавляем основные категории
	engravingCategoryID := int64(0)
	printingCategoryID := int64(0)

	// Категория "Выжигание"
	err = tx.QueryRow(`
		INSERT INTO categories (created_at, updated_at) 
		VALUES (NOW(), NOW()) 
		RETURNING id
	`).Scan(&engravingCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении категории 'Выжигание': %w", err)
	}

	// Добавляем переводы для категории "Выжигание"
	_, err = tx.Exec(`
		INSERT INTO category_translations (category_id, language, name) 
		VALUES ($1, 'ru', 'Выжигание'),
		       ($1, 'en', 'Engraving'),
		       ($1, 'es', 'Grabado')
	`, engravingCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении переводов для категории 'Выжигание': %w", err)
	}

	// Категория "3D Печать"
	err = tx.QueryRow(`
		INSERT INTO categories (created_at, updated_at) 
		VALUES (NOW(), NOW()) 
		RETURNING id
	`).Scan(&printingCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении категории '3D Печать': %w", err)
	}

	// Добавляем переводы для категории "3D Печать"
	_, err = tx.Exec(`
		INSERT INTO category_translations (category_id, language, name) 
		VALUES ($1, 'ru', '3D Печать'),
		       ($1, 'en', '3D Printing'),
		       ($1, 'es', 'Impresión 3D')
	`, printingCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении переводов для категории '3D Печать': %w", err)
	}

	// Добавляем подкатегории для "Выжигание"

	// Подкатегория "Дерево"
	woodCategoryID := int64(0)
	err = tx.QueryRow(`
		INSERT INTO categories (parent_id, created_at, updated_at) 
		VALUES ($1, NOW(), NOW()) 
		RETURNING id
	`, engravingCategoryID).Scan(&woodCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении подкатегории 'Дерево': %w", err)
	}

	// Добавляем переводы для подкатегории "Дерево"
	_, err = tx.Exec(`
		INSERT INTO category_translations (category_id, language, name) 
		VALUES ($1, 'ru', 'Дерево'),
		       ($1, 'en', 'Wood'),
		       ($1, 'es', 'Madera')
	`, woodCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении переводов для подкатегории 'Дерево': %w", err)
	}

	// Подкатегория "Металл"
	metalCategoryID := int64(0)
	err = tx.QueryRow(`
		INSERT INTO categories (parent_id, created_at, updated_at) 
		VALUES ($1, NOW(), NOW()) 
		RETURNING id
	`, engravingCategoryID).Scan(&metalCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении подкатегории 'Металл': %w", err)
	}

	// Добавляем переводы для подкатегории "Металл"
	_, err = tx.Exec(`
		INSERT INTO category_translations (category_id, language, name) 
		VALUES ($1, 'ru', 'Металл'),
		       ($1, 'en', 'Metal'),
		       ($1, 'es', 'Metal')
	`, metalCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении переводов для подкатегории 'Металл': %w", err)
	}

	// Подкатегория "Другое"
	otherCategoryID := int64(0)
	err = tx.QueryRow(`
		INSERT INTO categories (parent_id, created_at, updated_at) 
		VALUES ($1, NOW(), NOW()) 
		RETURNING id
	`, engravingCategoryID).Scan(&otherCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении подкатегории 'Другое': %w", err)
	}

	// Добавляем переводы для подкатегории "Другое"
	_, err = tx.Exec(`
		INSERT INTO category_translations (category_id, language, name) 
		VALUES ($1, 'ru', 'Другое'),
		       ($1, 'en', 'Other'),
		       ($1, 'es', 'Otro')
	`, otherCategoryID)

	if err != nil {
		return fmt.Errorf("ошибка при добавлении переводов для подкатегории 'Другое': %w", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	logger.Info("Начальные данные успешно добавлены")
	return nil
}
