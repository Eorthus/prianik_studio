package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggingMiddleware логирует информацию о запросах
func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Выполняем запрос
		c.Next()

		// Логируем информацию после выполнения запроса
		status := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()

		// Добавляем цветную подсветку в зависимости от статуса
		var statusColor, resetColor string
		if gin.Mode() != gin.ReleaseMode {
			resetColor = "\033[0m"
			if status >= 200 && status < 300 {
				statusColor = "\033[32m" // Зеленый для 2xx
			} else if status >= 300 && status < 400 {
				statusColor = "\033[33m" // Желтый для 3xx
			} else if status >= 400 && status < 500 {
				statusColor = "\033[31m" // Красный для 4xx
			} else {
				statusColor = "\033[31m\033[1m" // Жирный красный для 5xx
			}
		}

		// Логируем с соответствующим уровнем
		if status >= 500 {
			logger.WithFields(logrus.Fields{
				"status":    status,
				"path":      path,
				"method":    method,
				"client_ip": clientIP,
				"latency":   c.Writer.Size(),
			}).Error(statusColor, "Ошибка сервера", resetColor)
		} else if status >= 400 {
			logger.WithFields(logrus.Fields{
				"status":    status,
				"path":      path,
				"method":    method,
				"client_ip": clientIP,
				"size":      c.Writer.Size(),
			}).Warn(statusColor, "Ошибка клиента", resetColor)
		} else {
			logger.WithFields(logrus.Fields{
				"status":    status,
				"path":      path,
				"method":    method,
				"client_ip": clientIP,
				"size":      c.Writer.Size(),
			}).Info(statusColor, "Запрос обработан", resetColor)
		}
	}
}

// SecureHeadersMiddleware добавляет заголовки безопасности
func SecureHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Устанавливаем заголовки безопасности
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data: https:; script-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self'; connect-src 'self'")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Feature-Policy", "camera 'none'; microphone 'none'; geolocation 'none'")

		c.Next()
	}
}

// SanitizeInputMiddleware санитизирует входные данные
func SanitizeInputMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Для POST, PUT, PATCH методов
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			// Проверяем Content-Type
			contentType := c.GetHeader("Content-Type")

			// Для JSON запросов валидация будет выполнена при привязке
			if strings.Contains(contentType, "application/json") {
				c.Next()
				return
			}

			// Для form-data или x-www-form-urlencoded
			if strings.Contains(contentType, "multipart/form-data") ||
				strings.Contains(contentType, "application/x-www-form-urlencoded") {
				// Получаем form-data
				if err := c.Request.ParseForm(); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"success": false,
						"error":   "Невозможно разобрать данные формы",
					})
					c.Abort()
					return
				}

				// Здесь можно добавить валидацию и санитизацию полей формы
				// ...
			}
		}

		c.Next()
	}
}
