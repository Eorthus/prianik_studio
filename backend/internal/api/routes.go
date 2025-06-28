package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"pryanik_studio/internal/auth"
	"pryanik_studio/internal/config"
	"pryanik_studio/internal/security"
	"pryanik_studio/internal/storage"
	"pryanik_studio/internal/utils"
)

// SetupRouter настраивает маршруты для API
func SetupRouter(
	repo storage.Repository,
	emailSender utils.Sender,
	cfg *config.Config,
	logger *logrus.Logger,
) *gin.Engine {
	// Создаем роутер
	router := gin.Default()

	// Настройка CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Настройка ограничения скорости запросов
	limiter := security.NewIPRateLimiter(rate.Limit(cfg.Security.APIRateLimit), 20, logger)
	router.Use(security.RateLimitMiddleware(limiter))

	// Инициализируем JWT аутентификацию
	jwtAuth := auth.NewJWTAuth(cfg.Security.JWTSecret, logger)

	// Создаем обработчики
	authHandler := NewAuthHandler(jwtAuth, logger)
	productHandler := NewProductHandler(repo, logger)
	galleryHandler := NewGalleryHandler(repo, logger)
	orderHandler := NewOrderHandler(repo, repo, emailSender, logger)

	// Группа API
	api := router.Group("/api")
	{
		// Аутентификация (открытые эндпоинты)
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
		}

		// Публичные эндпоинты (без авторизации)
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProductByID)
		api.GET("/products/:id/related", productHandler.GetRelatedProducts)
		api.GET("/categories", productHandler.GetCategories)
		api.GET("/gallery", galleryHandler.GetGalleryItems)

		// Публичные формы
		api.POST("/orders", orderHandler.CreateOrder)
		api.POST("/contact", orderHandler.SubmitContactForm)

		// Админские эндпоинты (требуют авторизации и роли admin)
		admin := api.Group("/admin")
		admin.Use(jwtAuth.Middleware(), jwtAuth.RequireAdmin())
		{
			// Управление товарами
			admin.POST("/products", productHandler.CreateProduct)
			admin.PATCH("/products/:id", productHandler.UpdateProduct)

			// Управление галереей
			admin.POST("/gallery", galleryHandler.CreateGalleryItem)
			admin.DELETE("/gallery/:id", galleryHandler.DeleteGalleryItem)
		}
	}

	return router
}
