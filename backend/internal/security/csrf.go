package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// CSRFHeaderName имя HTTP-заголовка для CSRF-токена
	CSRFHeaderName = "X-CSRF-Token"

	// CSRFCookieName имя cookie для CSRF-токена
	CSRFCookieName = "csrf_token"

	// TokenExpiration время жизни токена в секундах (1 час)
	TokenExpiration = 3600
)

var (
	// ErrInvalidToken ошибка при недействительном токене
	ErrInvalidToken = errors.New("недействительный CSRF-токен")

	// ErrTokenExpired ошибка при истекшем токене
	ErrTokenExpired = errors.New("срок действия CSRF-токена истек")
)

// CSRFProtection защита от CSRF-атак
type CSRFProtection struct {
	secret []byte
	logger interface {
		Warnf(format string, args ...interface{})
	}
}

// NewCSRFProtection создает новый экземпляр CSRFProtection
func NewCSRFProtection(secret string, logger interface {
	Warnf(format string, args ...interface{})
}) *CSRFProtection {
	return &CSRFProtection{
		secret: []byte(secret),
		logger: logger,
	}
}

// GenerateToken генерирует новый CSRF-токен
func (c *CSRFProtection) GenerateToken() string {
	// Формат токена: [timestamp].[hmac]
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Создаем HMAC для timestamp
	h := hmac.New(sha256.New, c.secret)
	h.Write([]byte(timestamp))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// Объединяем timestamp и подпись
	token := fmt.Sprintf("%s.%s", timestamp, signature)
	return token
}

// ValidateToken проверяет действительность CSRF-токена
func (c *CSRFProtection) ValidateToken(token string) error {
	// Разделяем токен на timestamp и подпись
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return ErrInvalidToken
	}

	timestamp, signature := parts[0], parts[1]

	// Проверяем, не истек ли срок действия токена
	timestampInt, err := parseTimestamp(timestamp)
	if err != nil {
		return ErrInvalidToken
	}

	// Проверяем, не истек ли токен
	if time.Now().Unix()-timestampInt > TokenExpiration {
		return ErrTokenExpired
	}

	// Проверяем подпись
	h := hmac.New(sha256.New, c.secret)
	h.Write([]byte(timestamp))
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if signature != expectedSignature {
		return ErrInvalidToken
	}

	return nil
}

// SetCSRFCookie устанавливает CSRF-токен в cookie
func (c *CSRFProtection) SetCSRFCookie(ctx *gin.Context) {
	token := c.GenerateToken()

	// Устанавливаем cookie с SameSite=Strict для дополнительной защиты
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.SetCookie(
		CSRFCookieName,
		token,
		TokenExpiration,
		"/",
		"",
		ctx.Request.TLS != nil, // Secure flag только для HTTPS
		true,                   // HttpOnly для защиты от XSS
	)
}

// CSRFMiddleware создает middleware для защиты от CSRF-атак
func CSRFMiddleware(csrf *CSRFProtection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем GET, HEAD, OPTIONS и TRACE запросы
		if c.Request.Method == "GET" ||
			c.Request.Method == "HEAD" ||
			c.Request.Method == "OPTIONS" ||
			c.Request.Method == "TRACE" {
			// Устанавливаем новый CSRF-токен для последующих запросов
			csrf.SetCSRFCookie(c)
			c.Next()
			return
		}

		// Для POST, PUT, DELETE и PATCH проверяем CSRF-токен
		token := c.GetHeader(CSRFHeaderName)
		if token == "" {
			csrf.logger.Warnf("CSRF-токен отсутствует в запросе от %s", c.ClientIP())
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "CSRF-токен отсутствует",
			})
			c.Abort()
			return
		}

		// Проверяем токен
		if err := csrf.ValidateToken(token); err != nil {
			csrf.logger.Warnf("Недействительный CSRF-токен от %s: %v", c.ClientIP(), err)

			if err == ErrTokenExpired {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"error":   "Срок действия CSRF-токена истек",
				})
			} else {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"error":   "Недействительный CSRF-токен",
				})
			}

			c.Abort()
			return
		}

		// Устанавливаем новый CSRF-токен для следующего запроса
		csrf.SetCSRFCookie(c)
		c.Next()
	}
}

// Вспомогательная функция для анализа временной метки
func parseTimestamp(timestamp string) (int64, error) {
	var ts int64
	if _, err := fmt.Sscanf(timestamp, "%d", &ts); err != nil {
		return 0, err
	}
	return ts, nil
}
