package helper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// struct error -> definisikan field untuk segala jenis error
type AllErrors struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   any    `json:"error"`
}

// struct error untuk import from file
type ImportError struct {
	Row     int    `json:"row"`
	Message string `json:"message"`
}

// definisikan jenis-jenis error melalui func
func ErrorFetchDataFromDB(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, AllErrors{
		Code:    http.StatusInternalServerError,
		Message: "gagal tarik data dari database",
		Error:   err.Error(),
	})
}

func ErrorParsingAtoi(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "gagal parsing atoi",
		Error:   err.Error(),
	})
}

func ErrorParsingBoolean(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "gagal parsing string ke bool",
		Error:   err.Error(),
	})
}

func ErrorDataNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, AllErrors{
		Code:    http.StatusNotFound,
		Message: "data tidak ditemukan",
		Error:   errors.New("record not found").Error(), // wajib pakai .Error() di belakang errors.New
	})
}

func ErrorCreateData(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "gagal create data",
		Error:   err.Error(),
	})
}

func ErrorUpdateData(c *gin.Context, err error) {
	c.JSON(http.StatusConflict, AllErrors{
		Code:    http.StatusConflict,
		Message: "gagal update data",
		Error:   err.Error(),
	})
}

func ErrorDeleteData(c *gin.Context, err error) {
	c.JSON(http.StatusConflict, AllErrors{
		Code:    http.StatusConflict,
		Message: "gagal delete data",
		Error:   err.Error(),
	})
}

func ErrorParsingRequestBody(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "gagal parsing request body",
		Error:   err.Error(),
	})
}

func ErrorParsingDate(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "gagal parsing string ke date",
		Error:   err.Error(),
	})
}

func ErrorParsingFile(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "file tidak ditemukan",
		Error:   err.Error(),
	})
}

func ErrorSavingFile(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, AllErrors{
		Code:    http.StatusInternalServerError,
		Message: "gagal simpan file",
		Error:   err.Error(),
	})
}

func ErrorReadingFile(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, AllErrors{
		Code:    http.StatusInternalServerError,
		Message: "gagal baca file",
		Error:   err.Error(),
	})
}

func ErrorExcelSheetNotFound(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, AllErrors{
		Code:    http.StatusInternalServerError,
		Message: "nama sheet excel tidak ditemukan",
		Error:   err.Error(),
	})
}

func ErrorImportData(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, AllErrors{
		Code:    http.StatusInternalServerError,
		Message: "gagal import data dari excel",
		Error:   err.Error(),
	})
}

func ErrorDataImportKosong(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "tidak ada data yang ditarik dari file excel",
		Error:   err.Error(),
	})
}

func ErrorHashingPassword(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, AllErrors{
		Code:    http.StatusInternalServerError,
		Message: "gagal hash password",
		Error:   err.Error(),
	})
}

func ErrorLoginInvalid(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, AllErrors{
		Code:    http.StatusBadRequest,
		Message: "gagal login",
		Error:   err.Error(),
	})
}

func ErrorGenerateJWTToken(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, AllErrors{
		Code:    http.StatusInternalServerError,
		Message: "gagal generate token",
		Error:   err.Error(),
	})
}
