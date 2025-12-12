package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimiter создает middleware для ограничения частоты запросов
func RateLimiter(requests int, window int) gin.HandlerFunc {
	// Создаем rate limiter с конфигурацией
	rate := limiter.Rate{
		Period: time.Duration(window) * time.Second,
		Limit:  int64(requests),
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Получаем IP клиента
		ip := c.ClientIP()

		// Проверяем лимит
		context, err := instance.Get(c, ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Ошибка rate limiter",
			})
			c.Abort()
			return
		}

		// Устанавливаем заголовки
		c.Header("X-RateLimit-Limit", string(rune(context.Limit)))
		c.Header("X-RateLimit-Remaining", string(rune(context.Remaining)))
		c.Header("X-RateLimit-Reset", string(rune(context.Reset)))

		// Если лимит превышен
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Превышен лимит запросов. Попробуйте позже.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
