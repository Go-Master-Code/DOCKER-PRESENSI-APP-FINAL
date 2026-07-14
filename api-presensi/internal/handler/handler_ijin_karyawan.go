package handler

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/service"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// struct impementasi
type HandlerIjinKaryawan struct {
	service service.ServiceIjinKaryawan
}

// constructor
func NewHandlerIjinKaryawan(service service.ServiceIjinKaryawan) *HandlerIjinKaryawan {
	return &HandlerIjinKaryawan{service}
}

// struct method
func (h *HandlerIjinKaryawan) GetAllIjinKaryawan(c *gin.Context) {
	ijin, err := h.service.GetAllIjinKaryawan()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, ijin)
}

func (h *HandlerIjinKaryawan) GetIjinKaryawanPerTanggal(c *gin.Context) {
	// ambil query parameter tanggal
	tanggal := c.Query("tanggal") // string

	if tanggal == "" {
		helper.ErrorParsingDate(c, errors.New("tanggal harus diisi"))
		return
	}

	ijin, err := h.service.GetIjinKaryawanPerTanggal(tanggal)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, ijin)
}

func (h *HandlerIjinKaryawan) GetIjinKaryawanHariIni(c *gin.Context) {
	// tanggal hari ini
	tgl := time.Now().Format("2006-01-02")
	log.Println(tgl)
	ijin, err := h.service.GetIjinKaryawanPerTanggal(tgl)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, ijin)
}

func (h *HandlerIjinKaryawan) GetIjinKaryawanByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	ijin, err := h.service.GetIjinKaryawanByID(idInt)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, ijin)
}

func (h *HandlerIjinKaryawan) CreateIjinKaryawan(c *gin.Context) {
	// parsing request body
	var ijinKaryawan dto.CreateIjinKaryawanRequest
	err := c.ShouldBindJSON(&ijinKaryawan)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	newIjinKaryawan, err := h.service.CreateIjinKaryawan(ijinKaryawan)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.SuccessCreateData(c, newIjinKaryawan)
}

func (h *HandlerIjinKaryawan) UpdateIjinKaryawan(c *gin.Context) {
	// parsing request body
	var ijin dto.UpdateIjinKaryawanRequest

	err := c.ShouldBindJSON(&ijin)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	// get param from URL
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	updatedIjin, err := h.service.UpdateIjinKaryawan(idInt, ijin)
	if err != nil {
		helper.ErrorUpdateData(c, err)
		return
	}

	helper.SuccessUpdateData(c, updatedIjin)
}

func (h *HandlerIjinKaryawan) GetIjinAllKaryawanPerPeriode(c *gin.Context) {
	tglAwal := c.Query("awal")
	tglAkhir := c.Query("akhir")

	ijin, err := h.service.GetIjinAllKaryawanPerPeriode(tglAwal, tglAkhir)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, ijin)
}

func (h *HandlerIjinKaryawan) DeleteIjinKaryawan(c *gin.Context) {
	// ambil param dari URL
	id := c.Param("id")
	// convert to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	deletedIjin, err := h.service.DeleteIjinKaryawan(idInt)
	if err != nil {
		helper.ErrorDeleteData(c, err)
		return
	}

	helper.SuccessDeleteData(c, deletedIjin)
}
