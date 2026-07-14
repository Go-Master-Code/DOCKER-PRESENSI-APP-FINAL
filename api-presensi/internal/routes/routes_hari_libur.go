package routes

import (
	"api-presensi/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesHariLibur(rg *gin.RouterGroup, h *handler.HandlerHariLibur) {
	// list handler hari libur
	rg.GET("/hari_libur", h.GetAllHariLibur)
	rg.GET("/hari_libur/:id", h.GetHariLiburByID)
	rg.POST("/hari_libur", h.CreateHariLibur)
	rg.DELETE("/hari_libur/:id", h.DeleteHariLibur)
	rg.PUT("/hari_libur/:id", h.UpdateHariLibur)
	// handler untuk get jumlah hari kerja dalam periode tertentu
	rg.GET("/hari_kerja", h.GetHariKerjaPerPeriode)
}
