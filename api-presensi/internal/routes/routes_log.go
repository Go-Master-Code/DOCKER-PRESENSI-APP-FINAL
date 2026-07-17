package routes

import (
	"api-presensi/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesLog(rg *gin.RouterGroup, h *handler.HandlerLog) {
	// list handler hari libur
	rg.GET("/logs", h.GetAllLogs)
}
