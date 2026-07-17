package handler

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/service"

	"github.com/gin-gonic/gin"
)

// tidak ada interface{}, langsung ke struct implementasi
type HandlerLog struct {
	service service.ServiceLog
}

// constructor
func NewHandlerLog(service service.ServiceLog) *HandlerLog {
	return &HandlerLog{service}
}

// struct method
func (h *HandlerLog) CreateLog(c *gin.Context) {
	var req dto.LogRequestAndResponse
	// parsing json request
	err := c.ShouldBindJSON(&req)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	logDTO, err := h.service.CreateLog(req)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.SuccessCreateData(c, logDTO)
}

func (h *HandlerLog) GetAllLogs(c *gin.Context) {
	logsDTO, err := h.service.GetAllLogs()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, logsDTO)
}
