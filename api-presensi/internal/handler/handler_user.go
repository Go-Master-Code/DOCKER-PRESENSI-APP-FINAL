package handler

import (
	"api-presensi/auth"
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/service"
	"api-presensi/internal/utils/crypto"
	"strconv"

	"github.com/gin-gonic/gin"
)

// no interface, langsung struct implementasi
type HandlerUser struct {
	service service.ServiceUser
}

// constructor
func NewHandlerUser(service service.ServiceUser) *HandlerUser {
	return &HandlerUser{service}
}

// struct method
func (h *HandlerUser) GetAllUser(c *gin.Context) {
	user, err := h.service.GetAllUser()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, user)
}

func (h *HandlerUser) GetUserByID(c *gin.Context) {
	// ambil id dari param URL
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	user, err := h.service.GetUserByID(idInt)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, user)
}

func (h *HandlerUser) CreateUser(c *gin.Context) {
	// parsing request body
	var user dto.CreateUserRequest

	err := c.ShouldBindJSON(&user)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	// hash password dengan bcrypt
	hashPassword, err := crypto.HashPassword(user.Password)
	if err != nil {
		helper.ErrorHashingPassword(c, err)
		return
	}

	// ganti password dengan hasil hash
	user.Password = hashPassword

	newUser, err := h.service.CreateUser(user)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.SuccessCreateData(c, newUser)
}

func (h *HandlerUser) UpdateUser(c *gin.Context) {
	// ambil id dari param URL
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}
	// parsing request body
	var req dto.UpdateUserRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	// hash password dengan bcrypt
	hashPassword, err := crypto.HashPassword(*req.Password)
	if err != nil {
		helper.ErrorHashingPassword(c, err)
		return
	}

	// ganti password dengan hasil hash
	req.Password = &hashPassword

	// eksekusi method
	updateUser, err := h.service.UpdateUser(idInt, req)
	if err != nil {
		helper.ErrorUpdateData(c, err)
		return
	}

	// response success
	helper.SuccessUpdateData(c, updateUser)
}

func (h *HandlerUser) DeleteUserByID(c *gin.Context) {
	// ambil param id
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		helper.ErrorParsingAtoi(c, err)
		return
	}

	user, err := h.service.DeleteUserByID(idInt)
	if err != nil {
		helper.ErrorDeleteData(c, err)
		return
	}

	helper.SuccessDeleteData(c, user)
}

func (h *HandlerUser) Login(c *gin.Context) {
	// parsing request body
	var loginReq dto.LoginRequest
	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	user, err := h.service.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		helper.ErrorLoginInvalid(c, err)
		return
	}

	// jika username ada dan password benar
	token, err := auth.GenerateToken(user.Username, user.RoleNama)

	if err != nil {
		helper.ErrorGenerateJWTToken(c, err)
		return
	}

	// login berhasil => kirim data user sekaligus token
	helper.SuccessLogin(c, user, token)
}
