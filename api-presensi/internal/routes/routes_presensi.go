package routes

import (
	"api-presensi/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesPresensi(rg *gin.RouterGroup, h *handler.HandlerPresensi) {
	// endpoint presensi
	rg.GET("/presensi", h.GetAllPresensi)
	rg.GET("/presensi/today", h.GetPresensiHariIni)
	rg.POST("/presensi", h.CreateOrUpdatePresensi)
	rg.GET("/presensi/harian", h.GetPresensiPerTanggal)
	rg.GET("/presensi/karyawan/periode", h.GetPresensiByIDPerPerPeriode) // presensi 1 orang karyawan di periode tertentu
	rg.GET("/presensi/all/periode", h.GetPresensiAllKaryawanPerPeriode)  // presensi seluruh karyawan di periode tertentu (hanya tampil jumlah kehadiran)

}
