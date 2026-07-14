package routes

import (
	"api-presensi/internal/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutesUser(rg *gin.RouterGroup, h *handler.HandlerUser) {
	// endpoint user
	rg.GET("/user", h.GetAllUser)
	rg.GET("/user/:id", h.GetUserByID)
	rg.POST("/user", h.CreateUser)
	rg.PUT("/user/:id", h.UpdateUser)
	rg.DELETE("/user/:id", h.DeleteUserByID)
}
