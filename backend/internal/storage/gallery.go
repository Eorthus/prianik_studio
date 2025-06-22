package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"pryanik_studio/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// GetGalleryItems возвращает список элементов галереи с фильтрацией без пагинации
// GetGalleryItems возвращает список элементов галереи с фильтрацией без пагинации
func (r *PostgresRepository) GetGalleryItems(ctx context.Context, filter models.GalleryFilter) (models.GalleryList, error) {
	var result models.GalleryList

	// Базовый запрос для выборки элементов
	query := `
	SELECT gi.id, gi.category_id, gi.thumbnail, gi.full_image, gi.created_at, gi.updated_at,
	       git.title, git.description
	FROM gallery_items gi
	JOIN gallery_item_translations git ON gi.id = git.gallery_item_id
	WHERE git.language = $1
	`

	// Добавляем условия фильтрации
	args := []interface{}{filter.Language}

	if filter.CategoryID != nil {
		query += " AND gi.category_id = $2"
		args = append(args, *filter.CategoryID)
	}

	// Добавляем сортировку
	query += " ORDER BY gi.created_at DESC"

	// Логируем выполняемый запрос
	r.logger.WithFields(logrus.Fields{
		"query": query,
		"args":  args,
	}).Debug("Выполнение SQL запроса для галереи")

	// Запрос элементов галереи
	var items []struct {
		ID          int64        `db:"id"`
		CategoryID  int64        `db:"category_id"`
		Thumbnail   string       `db:"thumbnail"`
		FullImage   string       `db:"full_image"`
		CreatedAt   sql.NullTime `db:"created_at"`
		UpdatedAt   sql.NullTime `db:"updated_at"`
		Title       string       `db:"title"`
		Description string       `db:"description"`
	}

	err := r.db.SelectContext(ctx, &items, query, args...)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при получении списка элементов галереи")
		return result, fmt.Errorf("ошибка при получении списка элементов галереи: %w", err)
	}

	// Логируем количество найденных записей
	r.logger.WithField("count", len(items)).Debug("Получены записи из таблицы gallery_items")

	// Если записей нет, возвращаем пустой результат
	if len(items) == 0 {
		r.logger.Info("Записи для галереи не найдены")
		return result, nil
	}

	// Преобразуем результаты запроса в модели
	result.Items = make([]models.GalleryItem, 0, len(items))
	for _, item := range items {
		galleryItem := models.GalleryItem{
			ID:         item.ID,
			CategoryID: item.CategoryID,
			Thumbnail:  item.Thumbnail,
			FullImage:  item.FullImage,
			CreatedAt:  item.CreatedAt.Time,
			UpdatedAt:  item.UpdatedAt.Time,
			// Прямой доступ к полям перевода
			Title:       item.Title,
			Description: item.Description,
		}

		result.Items = append(result.Items, galleryItem)
	}

	return result, nil
}

// CreateGalleryItem создает новый элемент галереи в базе данных
func (r *PostgresRepository) CreateGalleryItem(ctx context.Context, item *models.GalleryItem) (int64, error) {
	// Начинаем транзакцию
	tx, err := r.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		r.logger.WithError(err).Error("Ошибка при начале транзакции для создания элемента галереи")
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
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	if item.UpdatedAt.IsZero() {
		item.UpdatedAt = now
	}

	// Вставляем элемент галереи
	query := `
    INSERT INTO gallery_items (category_id, thumbnail, full_image, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id
    `

	var itemID int64
	err = tx.QueryRowContext(
		ctx,
		query,
		item.CategoryID,
		item.Thumbnail,
		item.FullImage,
		item.CreatedAt,
		item.UpdatedAt,
	).Scan(&itemID)

	if err != nil {
		r.logger.WithError(err).Error("Ошибка при создании элемента галереи")
		return 0, fmt.Errorf("ошибка при создании элемента галереи: %w", err)
	}

	// Добавляем переводы элемента галереи
	for lang, translation := range item.Translations {
		query = `
        INSERT INTO gallery_item_translations (gallery_item_id, language, title, description)
        VALUES ($1, $2, $3, $4)
        `

		_, err = tx.ExecContext(
			ctx,
			query,
			itemID,
			lang,
			translation.Title,
			translation.Description,
		)

		if err != nil {
			r.logger.WithError(err).Errorf("Ошибка при добавлении перевода для языка %s", lang)
			return 0, fmt.Errorf("ошибка при добавлении перевода: %w", err)
		}
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Ошибка при фиксации транзакции")
		return 0, fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return itemID, nil
}

func (r *PostgresRepository) DeleteGalleryItem(ctx context.Context, id int64) error {
	// Начинаем транзакцию
	tx, err := r.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("ошибка при начале транзакции: %w", err)
	}

	// Добавляем отложенную функцию для отката транзакции в случае ошибки
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				r.logger.WithError(rollbackErr).Error("Ошибка при откате транзакции")
			}
		}
	}()

	// Сначала проверяем, существует ли элемент
	var exists bool
	checkQuery := "SELECT EXISTS(SELECT 1 FROM gallery_items WHERE id = $1)"
	err = tx.QueryRowContext(ctx, checkQuery, id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования элемента: %w", err)
	}

	if !exists {
		return fmt.Errorf("gallery item not found")
	}

	// Удаляем переводы элемента галереи
	deleteTranslationsQuery := "DELETE FROM gallery_item_translations WHERE gallery_item_id = $1"
	_, err = tx.ExecContext(ctx, deleteTranslationsQuery, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении переводов элемента галереи: %w", err)
	}

	// Удаляем сам элемент галереи
	deleteItemQuery := "DELETE FROM gallery_items WHERE id = $1"
	_, err = tx.ExecContext(ctx, deleteItemQuery, id)
	if err != nil {
		return fmt.Errorf("ошибка при удалении элемента галереи: %w", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		r.logger.WithError(err).Error("Ошибка при фиксации транзакции")
		return fmt.Errorf("ошибка при фиксации транзакции: %w", err)
	}

	return nil
}
