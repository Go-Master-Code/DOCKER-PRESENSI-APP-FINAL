package handler

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// no interface

// struct implementasi
type HandlerHariLibur struct {
	service service.ServiceHariLibur
}

// constructor
func NewHandlerHariLibur(service service.ServiceHariLibur) *HandlerHariLibur {
	return &HandlerHariLibur{service}
}

// struct method
func (h *HandlerHariLibur) GetAllHariLibur(c *gin.Context) {
	hariLibur, err := h.service.GetAllHariLibur()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, hariLibur)
}

func (h *HandlerHariLibur) GetHariLiburByID(c *gin.Context) {
	// get param by URL
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	hariLibur, err := h.service.GetHariLiburByID(idInt)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, hariLibur)
}

func (h *HandlerHariLibur) CreateHariLibur(c *gin.Context) {
	// parsing request body
	var newHariLibur dto.CreateHariLiburRequest

	err := c.ShouldBindJSON(&newHariLibur)
	// log.Println(newHariLibur) debug daftar hari libur

	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	hariLiburDTO, err := h.service.CreateHariLibur(newHariLibur)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.SuccessCreateData(c, hariLiburDTO)
}

func (h *HandlerHariLibur) UpdateHariLibur(c *gin.Context) {
	// parsing request body
	var hariLibur dto.UpdateHariLiburRequest
	err := c.ShouldBindJSON(&hariLibur)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	// ambil param id dari URL
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	hariLiburDTO, err := h.service.UpdateHariLibur(idInt, hariLibur)
	if err != nil {
		helper.ErrorUpdateData(c, err)
		return
	}

	helper.SuccessUpdateData(c, hariLiburDTO)
}

func (h *HandlerHariLibur) DeleteHariLibur(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorDeleteData(c, err)
		return
	}

	hariLibur, err := h.service.DeleteHariLibur(idInt)

	// jika data dengan param id tidak ditemukan
	if hariLibur.Hari == "" {
		helper.ErrorDataNotFound(c)
		return
	}

	// error lainnya
	if err != nil {
		helper.ErrorDeleteData(c, err)
		return
	}

	helper.SuccessDeleteData(c, hariLibur)
}

func (h *HandlerHariLibur) GetHariKerjaPerPeriode(c *gin.Context) {
	// loc penting untuk parsing tanggal ke string ke format date yyyy-mm-dd
	// tanpa loc, hasil convert ke date akan menghasilkan 2026-04-27 00:00:00 +0000 UTC

	// loc := time.FixedZone("WIB", 7*3600) // 7 jam dari UTC
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic("gagal load time Asia/Jakarta")
	}

	// ambil query param tglAwal dan tglAkhir
	tglAwal := c.Query("awal")
	// parsing string ke date dengan method ParseInLocation agar muncul output WIB
	tglAwalDate, err := time.ParseInLocation("2006-01-02", tglAwal, loc)
	if err != nil {
		helper.ErrorParsingDate(c, err)
		return
	}

	tglAkhir := c.Query("akhir")
	// parsing string ke date dengan method ParseInLocation agar muncul output WIB
	tglAkhirDate, err := time.ParseInLocation("2006-01-02", tglAkhir, loc)
	if err != nil {
		helper.ErrorParsingDate(c, err)
		return
	}

	hariKerja, err := h.service.GetHariKerjaPerPeriode(tglAwalDate, tglAkhirDate)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, hariKerja)
}
