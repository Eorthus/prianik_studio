package auth

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidToken = errors.New("недействительный токен")
	ErrTokenExpired = errors.New("токен истек")
	ErrNoToken      = errors.New("токен отсутствует")
)

// JWTAuth структура для работы с JWT
type JWTAuth struct {
	secret []byte
	logger *logrus.Logger
}

// Claims структура для JWT claims
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest структура для запроса авторизации
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse структура ответа авторизации
type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserInfo  `json:"user"`
}

// UserInfo информация о пользователе
type UserInfo struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

// NewJWTAuth создает новый экземпляр JWTAuth
func NewJWTAuth(secret string, logger *logrus.Logger) *JWTAuth {
	return &JWTAuth{
		secret: []byte(secret),
		logger: logger,
	}
}

// GenerateToken генерирует JWT токен
func (j *JWTAuth) GenerateToken(username, role string) (string, time.Time, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен на 24 часа

	claims := &Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "prianik-studio",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// ValidateToken проверяет и парсит JWT токен
func (j *JWTAuth) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractTokenFromHeader извлекает токен из заголовка Authorization
func (j *JWTAuth) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", ErrNoToken
	}

	// Формат: "Bearer TOKEN"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", ErrInvalidToken
	}

	return parts[1], nil
}

// Middleware для проверки JWT токена
func (j *JWTAuth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		token, err := j.ExtractTokenFromHeader(authHeader)
		if err != nil {
			j.logger.Warnf("Отсутствует токен авторизации от %s", c.ClientIP())
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Необходима авторизация",
			})
			c.Abort()
			return
		}

		claims, err := j.ValidateToken(token)
		if err != nil {
			j.logger.Warnf("Недействительный токен от %s: %v", c.ClientIP(), err)

			var message string
			switch err {
			case ErrTokenExpired:
				message = "Токен истек"
			default:
				message = "Недействительный токен"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   message,
			})
			c.Abort()
			return
		}

		// Добавляем информацию о пользователе в контекст
		c.Set("user", claims)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// RequireAdmin проверяет, что пользователь - администратор
func (j *JWTAuth) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Необходима авторизация",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			j.logger.Warnf("Попытка доступа без прав администратора от %s", c.ClientIP())
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"error":   "Недостаточно прав доступа",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// CheckCredentials проверяет учетные данные (простая проверка)
func (j *JWTAuth) CheckCredentials(username, password string) bool {
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	return username == adminUsername && password == adminPassword
}
