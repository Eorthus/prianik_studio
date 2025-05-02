package security

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// IPRateLimiter ограничивает скорость запросов по IP-адресу
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	logger *logrus.Logger
}

// NewIPRateLimiter создает новый экземпляр IPRateLimiter
func NewIPRateLimiter(r rate.Limit, burst int, logger *logrus.Logger) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  burst,
		logger: logger,
	}
}

// AddIP добавляет IP-адрес в лимитер
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.ips[ip] = limiter
	return limiter
}

// GetLimiter возвращает лимитер для IP-адреса
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}
	return limiter
}

// RateLimitMiddleware создает middleware для ограничения скорости запросов
func RateLimitMiddleware(rateLimiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем IP клиента
		clientIP := c.ClientIP()

		// Получаем лимитер для этого IP
		limiter := rateLimiter.GetLimiter(clientIP)

		// Проверяем, не превышен ли лимит
		if !limiter.Allow() {
			rateLimiter.logger.Warnf("Превышено ограничение скорости запросов для IP %s", clientIP)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Превышено ограничение скорости запросов. Пожалуйста, повторите позже.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
