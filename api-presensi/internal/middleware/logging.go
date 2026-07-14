package middleware

import (
	"api-presensi/internal/dto"
	"api-presensi/internal/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// custom logger, sebenarnya sudah ada default Gin
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// lanjutkan ke middleware atau handler berikutnya
		c.Next()

		// setelah handler berikutnya selesai dieksekusi
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()
		clientIP := c.ClientIP()

		log.Printf("Log backend: [%d] %s %s | IP=%s | Agent=%s| Duration=%v",
			status,
			method,
			path,
			clientIP,
			userAgent,
			duration,
		)
	}
}

func RequestLoggerDB(logService service.ServiceLog) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // proses request

		// Filter hanya method tertentu
		switch c.Request.Method {
		case "POST", "PUT", "DELETE": // simpan hanya ketiga method ini agar method lain (yang tidak diperlukan) tidak ikut tersimpan di db
			username := getUserIDFromContext(c)

			logDTO := dto.LogRequestAndResponse{
				UserID:    username,
				Method:    c.Request.Method,
				Endpoint:  c.FullPath(),
				IPAddress: c.ClientIP(),
				UserAgent: c.Request.UserAgent(),
			}

			// simpan log secara asynchronus agar tidak blocking response
			go logService.CreateLog(logDTO)
		}
	}
}

// Mengambil userID dari context (di-set oleh middleware auth)
func getUserIDFromContext(c *gin.Context) string {
	username, exists := c.Get("username")
	if !exists {
		return ""
	}

	user, ok := username.(string) // coba parsing dulu ke string
	if !ok {
		return ""
	}

	return user
}
