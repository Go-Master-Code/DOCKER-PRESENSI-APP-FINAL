package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
	memory "github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimiter() gin.HandlerFunc {
	// aturan: 5 request per menit
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  5,
	}

	// gunakan store in-memory (atau Redis untuk produksi)
	store := memory.NewStore()

	// buat middlewarenya
	middleware := ginlimiter.NewMiddleware(limiter.New(store, rate))
	return middleware
}
