package routes

import (
	"api-presensi/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesRole(rg *gin.RouterGroup, h *handler.HandlerRole) {
	// endpoint role
	rg.GET("/role", h.GetAllRole)
}
