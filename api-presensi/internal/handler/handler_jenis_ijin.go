package handler

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// no interface

// struct implementasi
type HandlerJenisIjin struct {
	service service.ServiceJenisIjin
}

// constructor
func NewHandlerJenisIjin(service service.ServiceJenisIjin) *HandlerJenisIjin {
	return &HandlerJenisIjin{service}
}

// struct method
func (h *HandlerJenisIjin) GetAllJenisIjin(c *gin.Context) {
	jenisIjin, err := h.service.GetAllJenisIjin()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, jenisIjin)
}

func (h *HandlerJenisIjin) GetJenisIjinByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	jenisIjin, err := h.service.GetJenisIjinByID(idInt)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, jenisIjin)
}

func (h *HandlerJenisIjin) CreateJenisIjin(c *gin.Context) {
	var jenisIjin dto.CreateJenisIjinRequest

	err := c.ShouldBindJSON(&jenisIjin)

	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	newJenisIjinDTO, err := h.service.CreateJenisIjin(jenisIjin)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.SuccessCreateData(c, newJenisIjinDTO)
}

func (h *HandlerJenisIjin) UpdateJenisIjin(c *gin.Context) {
	// ambil param dari URL
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	// parsing request body
	var updateJenisIjin dto.UpdateJenisIjinRequest
	err = c.ShouldBindJSON(&updateJenisIjin)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	updateJenisIjinDTO, err := h.service.UpdateJenisIjin(idInt, updateJenisIjin)
	if err != nil {
		helper.ErrorUpdateData(c, err)
		return
	}

	helper.SuccessUpdateData(c, updateJenisIjinDTO)
}

func (h *HandlerJenisIjin) DeleteJenisIjin(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	jenisIjin, err := h.service.DeleteJenisIjin(idInt)

	// jika data dengan param id tidak ditemukan
	if jenisIjin.ID == 0 {
		helper.ErrorDataNotFound(c)
		return
	}

	// error lainnya
	if err != nil {
		helper.ErrorDeleteData(c, err)
		return
	}

	helper.SuccessDeleteData(c, jenisIjin)
}
