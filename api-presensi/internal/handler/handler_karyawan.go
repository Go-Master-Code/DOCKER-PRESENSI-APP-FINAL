package handler

import (
	"api-presensi/helper"
	"api-presensi/internal/config"
	"api-presensi/internal/dto"
	"api-presensi/internal/service"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// tidak pakai interface kosong

// langsung struct implementasi
type HandlerKaryawan struct {
	service service.ServiceKaryawan
}

// constructor
func NewHandlerKaryawan(service service.ServiceKaryawan) *HandlerKaryawan {
	return &HandlerKaryawan{service}
}

// struct method
func (h *HandlerKaryawan) GetAllKaryawan(c *gin.Context) {
	karyawan, err := h.service.GetAllKaryawan()
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, karyawan)
}

func (h *HandlerKaryawan) GetKaryawanByID(c *gin.Context) {
	// ambil param id dari URL
	id := c.Param("id")

	karyawan, err := h.service.GetKaryawanByID(id)
	if err != nil {
		helper.ErrorFetchDataFromDB(c, err)
		return
	}

	helper.SuccessGetData(c, karyawan)
}

func (h *HandlerKaryawan) GetKaryawanBelumAbsen(c *gin.Context) {
	// get param tanggal dari URL
	tanggal := c.Param("tanggal")

	karyawan, err := h.service.GetKaryawanBelumAbsen(tanggal)
	if err != nil {
		helper.ErrorDataNotFound(c)
		return
	}

	helper.SuccessGetData(c, karyawan)
}

func (h *HandlerKaryawan) GetKaryawanBelumIjin(c *gin.Context) {
	// get param tanggal dari URL
	tanggal := c.Param("tanggal")

	karyawan, err := h.service.GetKaryawanBelumIjin(tanggal)
	if err != nil {
		helper.ErrorDataNotFound(c)
		return
	}

	helper.SuccessGetData(c, karyawan)
}

func (h *HandlerKaryawan) CreateKaryawan(c *gin.Context) {
	// parsing request body
	var newKaryawan dto.CreateKaryawanRequest
	err := c.ShouldBindJSON(&newKaryawan)
	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	newKaryawanDTO, err := h.service.CreateKaryawan(newKaryawan)
	if err != nil {
		helper.ErrorCreateData(c, err)
		return
	}

	helper.SuccessCreateData(c, newKaryawanDTO)
}

func (h *HandlerKaryawan) UpdateKaryawan(c *gin.Context) {
	var req dto.UpdateKaryawanRequest

	// parsing request body
	err := c.ShouldBindJSON(&req)

	if err != nil {
		helper.ErrorParsingRequestBody(c, err)
		return
	}

	// tangkap parameter id dari URL
	id := c.Param("id")

	karyawanDTO, err := h.service.UpdateKaryawan(id, req)
	if err != nil {
		helper.ErrorUpdateData(c, err)
		return
	}

	// jika sukses
	helper.SuccessUpdateData(c, karyawanDTO)
}

func (h *HandlerKaryawan) DeleteKaryawan(c *gin.Context) {
	// ambil param dari URL
	id := c.Param("id")

	karyawanDTO, err := h.service.DeleteKaryawan(id)

	// jika data tidak ditemukan
	if karyawanDTO.ID == "" {
		helper.ErrorDataNotFound(c)
		return
	}

	// jika terjadi error lainnya
	if err != nil {
		helper.ErrorDeleteData(c, err)
		return
	}

	helper.SuccessDeleteData(c, karyawanDTO)
}

func (h *HandlerKaryawan) ImportKaryawan(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		helper.ErrorParsingFile(c, err)
		return
	}

	// simpan sementara
	path := config.App.UploadPath + "/" + file.Filename
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		helper.ErrorSavingFile(c, err)
		return
	}

	f, err := excelize.OpenFile(path)
	if err != nil {
		helper.ErrorReadingFile(c, err)
		return
	}

	defer os.Remove(path) // hapus file setelah diproses

	/* alur import data from excel:
	Upload Excel
	↓
	Baca Excel
	↓
	Import DB
	↓
	File otomatis dihapus
	*/

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		helper.ErrorExcelSheetNotFound(c, err)
		return
	}

	// buat var plural untuk menampung data dtoCreateKaryawanRequest
	var newKaryawan []dto.CreateKaryawanRequest

	// var error dari struct ImportError untuk tampung data error
	var importError []helper.ImportError

	// counter jumlah data yang diinsert
	var inserted int

	// map untuk menampung ID dan nama yang duplikat di file excel
	seenID := make(map[string]bool)   // 🔥 cek duplikat di file
	seenNama := make(map[string]bool) // optional

	// masukkan data rows ke dalam var dtoCreateKaryawanRequest
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		rowNumber := i + 1 // biar sesuai dengan nomor row excel

		// contoh kolom:
		// [0]=id, [1]=nama, dst
		// validasi jumlah kolom pada satu baris excel cukup sebelum diproses
		// len(row) WAJIB 2
		// Excel		row				len(row)
		// K001, Budi	["K001","Budi"]	2
		if len(row) < 2 { // angka 2 dipengaruhi oleh jumlah row yang ada di excel
			importError = append(importError, helper.ImportError{
				Row:     rowNumber,
				Message: "Minimal 2 kolom: ID, Nama",
			})
			continue
		}

		// hilangkan spasi
		id := strings.TrimSpace(row[0])
		nama := strings.TrimSpace(row[1])

		// validasi isi id dan nama
		if id == "" || nama == "" {
			importError = append(importError, helper.ImportError{
				Row:     rowNumber,
				Message: "ID atau nama kosong",
			})
			continue
		}

		// cek duplikat id
		if seenID[id] {
			importError = append(importError, helper.ImportError{
				Row:     rowNumber,
				Message: "ID duplikat di file",
			})
			continue
		}

		// 🔥 DUPLIKAT DI FILE (NAMA - opsional)
		if seenNama[nama] {
			importError = append(importError, helper.ImportError{
				Row:     rowNumber,
				Message: "Nama duplikat di file",
			})
			continue
		}

		seenID[id] = true
		seenNama[nama] = true

		// append data ke dalam struct dtoCreateRequest
		newKaryawan = append(newKaryawan, dto.CreateKaryawanRequest{
			ID:    id,
			Nama:  nama,
			Aktif: true,
		})

		// increment counter
		inserted++
	}

	// cek apakah struct sudah terisi
	if len(newKaryawan) == 0 {
		helper.ErrorDataImportKosong(c, err)
		return
	}

	// kalo lolos validasi error dan struct sudah terisi, insert data
	karyawan, err := h.service.ImportKaryawan(newKaryawan)
	if err != nil {
		helper.ErrorImportData(c, err)
		return
	}

	helper.SuccessImportData(c, karyawan, len(newKaryawan), len(importError), importError)
}
