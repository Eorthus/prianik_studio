package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pryanik_studio/internal/api"
	"pryanik_studio/internal/config"
	"pryanik_studio/internal/storage"
	"pryanik_studio/internal/utils"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Ошибка загрузки конфигурации: %v\n", err)
		os.Exit(1)
	}

	// Инициализируем логгер
	log := utils.NewLogger(cfg.Logging.Level)
	log.Info("Запуск сервера...")

	// Устанавливаем режим Gin в соответствии с конфигурацией
	if cfg.Server.Mode != "" {
		os.Setenv("GIN_MODE", cfg.Server.Mode)
	}

	// Инициализируем соединение с базой данных
	db, err := storage.NewDatabase(cfg.Database, log)
	if err != nil {
		log.WithError(err).Fatal("Не удалось подключиться к базе данных")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.WithError(err).Error("Ошибка при закрытии соединения с базой данных")
		}
	}()

	// Выполняем миграции базы данных
	if err := storage.MigrateDatabase(db, log); err != nil {
		log.WithError(err).Fatal("Ошибка при выполнении миграций базы данных")
	}

	// Инициализируем репозиторий
	repo := storage.NewPostgresRepository(db, log)

	// Инициализируем отправитель email
	var emailSender utils.Sender

	// Проверяем, какой провайдер email использовать
	switch cfg.Email.Provider {
	case "sendgrid":
		if cfg.Email.SendGridAPIKey != "" {
			emailSender = utils.NewSendGridSender(cfg.Email, log)
		} else {
			log.Warn("SendGrid API ключ не найден, переключаемся на SMTP")
			emailSender = utils.NewGomailSender(cfg.Email, log)
		}
	default:
		emailSender = utils.NewGomailSender(cfg.Email, log)
	}

	// Инициализируем роутер
	router := api.SetupRouter(repo, emailSender, &cfg, log)

	// Создаем HTTP-сервер
	server := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		if cfg.Security.EnableHTTPS {
			log.Info("Запуск HTTPS сервера...")
			if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
				log.WithError(err).Fatal("Ошибка при запуске HTTPS сервера")
			}
		} else {
			log.Info("Запуск HTTP сервера...")
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.WithError(err).Fatal("Ошибка при запуске HTTP сервера")
			}
		}
	}()

	// Канал для перехвата сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал завершения
	<-quit
	log.Info("Завершение работы сервера...")

	// Создаем контекст с таймаутом для корректного завершения
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.WithError(err).Fatal("Ошибка при штатном завершении сервера")
	}

	// Ожидаем завершения обработки запросов
	<-ctx.Done()
	log.Info("Сервер успешно остановлен")
}
