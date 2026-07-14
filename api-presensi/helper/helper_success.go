package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// struct AllSuccess berisi field untuk setiap jenis pesan sukses
type AllSuccess struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type LoginSuccess struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Token   string `json:"token"`
}

type SuccessImportExcel struct {
	Code         int           `json:"code"`
	InsertedData int           `json:"inserted_data"`
	ErrorCount   int           `json:"error_count"`
	Errors       []ImportError `json:"errors"`
}

// --> Helper sukses untuk berbagai macam kondisi
func SuccessGetData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, AllSuccess{
		Code:    http.StatusOK,
		Message: "sukses tarik data dari database",
		Data:    data,
	})
}

func SuccessCreateData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, AllSuccess{
		Code:    http.StatusOK,
		Message: "sukses create data",
		Data:    data,
	})
}

func SuccessImportData(c *gin.Context, data any, inserted int, errorCount int, err []ImportError) {
	c.JSON(http.StatusOK, SuccessImportExcel{
		Code:         http.StatusOK,
		InsertedData: inserted,
		ErrorCount:   errorCount,
		Errors:       err,
	})
}

func SuccessUpdateData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, AllSuccess{
		Code:    http.StatusOK,
		Message: "sukses update data",
		Data:    data,
	})
}

func SuccessDeleteData(c *gin.Context, data any) {
	c.JSON(http.StatusOK, AllSuccess{
		Code:    http.StatusOK,
		Message: "sukses delete data",
		Data:    data,
	})
}

func SuccessLogin(c *gin.Context, data any, token string) {
	c.JSON(http.StatusOK, LoginSuccess{
		Code:    http.StatusOK,
		Message: "login berhasil",
		Data:    data,
		Token:   token,
	})
}
