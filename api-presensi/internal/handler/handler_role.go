package handler

import (
	"api-presensi/helper"
	"api-presensi/internal/service"

	"github.com/gin-gonic/gin"
)

// no interface, struct implementasi
type HandlerRole struct {
	service service.ServiceRole
}

// constructor
func NewHandlerRole(service service.ServiceRole) *HandlerRole {
	return &HandlerRole{service}
}

// struct method
func (h *HandlerRole) GetAllRole(c *gin.Context) {
	roles, err := h.service.GetAllRole()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, roles)
}
