package routes

import (
	"api-presensi/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesJenisIjin(rg *gin.RouterGroup, h *handler.HandlerJenisIjin) {
	rg.GET("/jenis_ijin", h.GetAllJenisIjin)
	rg.GET("/jenis_ijin/:id", h.GetJenisIjinByID)
	rg.POST("/jenis_ijin", h.CreateJenisIjin)
	rg.DELETE("/jenis_ijin/:id", h.DeleteJenisIjin)
	rg.PUT("/jenis_ijin/:id", h.UpdateJenisIjin)
}
