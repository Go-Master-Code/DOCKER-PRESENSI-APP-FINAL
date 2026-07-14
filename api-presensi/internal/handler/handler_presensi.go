package handler

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/service"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

// langsung struct implementasi tanpa interface{}
type HandlerPresensi struct {
	service service.ServicePresensi
}

// constructor
func NewHandlerPresensi(service service.ServicePresensi) *HandlerPresensi {
	return &HandlerPresensi{service}
}

// struct method
func (h *HandlerPresensi) GetAllPresensi(c *gin.Context) {
	presensi, err := h.service.GetAllPresensi()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
	}

	helper.SuccessGetData(c, presensi)
}

func (h *HandlerPresensi) CreateOrUpdatePresensi(c *gin.Context) {
	// parsing request body
	var presensi dto.CreatePresensiRequest

	err := c.ShouldBindJSON(&presensi)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	newPresensiDTO, err := h.service.CreateOrUpdatePresensi(presensi)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.SuccessCreateData(c, newPresensiDTO)
}

func (h *HandlerPresensi) GetPresensiPerTanggal(c *gin.Context) {
	tanggal := c.Query("tanggal") // ambil query parameter, bisa lebih dari 1 param
	// cara penulisan URL beda, cek di main.go

	if tanggal == "" {
		helper.ErrorParsingDate(c, errors.New("tanggal harus diisi"))
		return
	}
	presensi, err := h.service.GetPresensiPerTanggal(tanggal)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	// jika tidak ada data yang di return
	if len(presensi) < 1 {
		helper.SuccessGetData(c, presensi) // tetap response success, dan return data kosong
		return
	}

	helper.SuccessGetData(c, presensi)
}

// handler untuk menampilkan presensi hari ini
func (h *HandlerPresensi) GetPresensiHariIni(c *gin.Context) {
	tgl := time.Now().Format("2006-01-02")
	// log.Println(tgl) debug tanggal
	presensi, err := h.service.GetPresensiPerTanggal(tgl)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, presensi)
}

func (h *HandlerPresensi) GetPresensiByIDPerPerPeriode(c *gin.Context) {
	// tangkap 3 query parameter
	idKaryawan := c.Query("id")
	tanggalAwal := c.Query("awal")
	tanggalAkhir := c.Query("akhir")

	presensi, err := h.service.GetPresensiByIDPerPerPeriode(idKaryawan, tanggalAwal, tanggalAkhir)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	// jika tidak ada record yang di return
	if len(presensi) < 1 {
		helper.ErrorDataNotFound(c)
		return
	}

	// jika data ditemukan
	helper.SuccessGetData(c, presensi)
}

func (h *HandlerPresensi) GetPresensiAllKaryawanPerPeriode(c *gin.Context) {
	tglAwal := c.Query("awal")
	tglAkhir := c.Query("akhir")
	presensi, err := h.service.GetPresensiAllKaryawanPerPeriode(tglAwal, tglAkhir)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, presensi)
}
