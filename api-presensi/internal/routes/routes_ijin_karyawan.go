package routes

import (
	"api-presensi/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesIjinKaryawan(rg *gin.RouterGroup, h *handler.HandlerIjinKaryawan) {
	// endpoint ijin karyawan
	rg.GET("/ijin_karyawan", h.GetAllIjinKaryawan)
	rg.GET("/ijin_karyawan/today", h.GetIjinKaryawanHariIni)
	rg.GET("/ijin_karyawan/harian", h.GetIjinKaryawanPerTanggal)
	rg.GET("/ijin_karyawan/:id", h.GetIjinKaryawanByID)
	rg.GET("/ijin_karyawan/all/periode", h.GetIjinAllKaryawanPerPeriode)
	rg.POST("/ijin_karyawan", h.CreateIjinKaryawan)
	rg.PUT("/ijin_karyawan/:id", h.UpdateIjinKaryawan)
	rg.DELETE("/ijin_karyawan/:id", h.DeleteIjinKaryawan)
}
