package routes

import (
	"api-presensi/internal/handler"
	"api-presensi/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesKaryawan(rg *gin.RouterGroup, h *handler.HandlerKaryawan) {
	rg.GET("/karyawan", h.GetAllKaryawan) // endpoint ini hanya bisa diakses oleh admin
	rg.GET("/karyawan/:id", h.GetKaryawanByID)
	rg.GET("/karyawan/absen/:tanggal", h.GetKaryawanBelumAbsen)
	rg.GET("/karyawan/ijin/:tanggal", h.GetKaryawanBelumIjin)
	rg.POST("/karyawan", h.CreateKaryawan) // endpoint ini hanya bisa diakses oleh admin
	rg.POST("/karyawan/import", h.ImportKaryawan)
	rg.PUT("/karyawan/:id", middleware.AuthRole("admin"), h.UpdateKaryawan)
	rg.DELETE("/karyawan/:id", middleware.AuthRole("admin"), h.DeleteKaryawan) // endpoint ini hanya bisa diakses oleh admin
}
