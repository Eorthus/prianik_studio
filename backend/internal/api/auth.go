package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"pryanik_studio/internal/auth"
)

// AuthHandler обработчик аутентификации
type AuthHandler struct {
	jwtAuth *auth.JWTAuth
	logger  *logrus.Logger
}

// NewAuthHandler создает новый AuthHandler
func NewAuthHandler(jwtAuth *auth.JWTAuth, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		jwtAuth: jwtAuth,
		logger:  logger,
	}
}

// Login обрабатывает запрос на авторизацию
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warnf("Некорректные данные авторизации от %s: %v", c.ClientIP(), err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Некорректные данные запроса",
		})
		return
	}

	// Проверяем учетные данные
	if !h.jwtAuth.CheckCredentials(req.Username, req.Password) {
		h.logger.Warnf("Неуспешная попытка входа от %s: неверные учетные данные", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Неверные учетные данные",
		})
		return
	}

	// Генерируем токен
	token, expiresAt, err := h.jwtAuth.GenerateToken(req.Username, "admin")
	if err != nil {
		h.logger.Errorf("Ошибка генерации токена: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Ошибка сервера",
		})
		return
	}

	h.logger.Infof("Успешная авторизация пользователя %s от %s", req.Username, c.ClientIP())

	// Возвращаем токен
	response := auth.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: auth.UserInfo{
			Username: req.Username,
			Role:     "admin",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    response,
	})
}
