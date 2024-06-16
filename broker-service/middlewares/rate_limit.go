package middlewares

import (
	"broker/helpers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

func RateLimitMiddleware(redisClient *redis.Client, rateLimit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := "rate_limit_" + ip

		// Increment request count
		count, err := redisClient.Incr(c, key).Result()
		if err != nil {
			helpers.ErrorJSON(c, err, http.StatusInternalServerError)
			c.Abort()
			return
		}

		// set expiration for the key if this is the first request.
		if count == 1 {
			redisClient.Expire(c, key, duration)
		}

		// check rate limit
		if count > int64(rateLimit) {
			resetTime, _ := redisClient.TTL(c, key).Result()
			helpers.ErrorJSON(c, fmt.Errorf("rate limit exceeded, try again in %.0f seconds", resetTime.Seconds()), http.StatusTooManyRequests)
			c.Abort()
			return
		}
		c.Next()
	}
}
