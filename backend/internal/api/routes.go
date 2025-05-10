package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

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

	// csrf := security.NewCSRFProtection(cfg.Security.JWTSecret, logger)
	// router.Use(security.CSRFMiddleware(csrf))
	// Создаем обработчики
	productHandler := NewProductHandler(repo, logger)
	galleryHandler := NewGalleryHandler(repo, logger)
	orderHandler := NewOrderHandler(repo, repo, emailSender, logger)

	// Группа API
	api := router.Group("/api")
	{
		// Товары и категории
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProductByID)
		api.GET("/products/:id/related", productHandler.GetRelatedProducts)
		api.GET("/categories", productHandler.GetCategories)
		api.POST("/products", productHandler.CreateProduct)
		api.PATCH("/products/:id", productHandler.UpdateProduct)

		// Галерея
		api.GET("/gallery", galleryHandler.GetGalleryItems)
		api.POST("/gallery", galleryHandler.CreateGalleryItem)

		// Заказы и формы
		api.POST("/orders", orderHandler.CreateOrder)
		api.POST("/contact", orderHandler.SubmitContactForm)
	}

	return router
}
