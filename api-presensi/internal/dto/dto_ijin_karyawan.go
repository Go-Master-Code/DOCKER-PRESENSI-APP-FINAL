package dto

type IjinKaryawanResponse struct {
	ID            int    `json:"id"`
	Tanggal       string `json:"tanggal"` // response tipe data yang seharusnya time.Time jadikan string saja
	KaryawanID    string `json:"karyawan_id"`
	KaryawanNama  string `json:"karyawan_nama"`
	JenisIjinID   int    `json:"jenis_ijin_id"`
	JenisIjinNama string `json:"jenis_ijin_nama"`
	Keterangan    string `json:"keterangan"`
}

type CreateIjinKaryawanRequest struct {
	Tanggal     string `json:"tanggal" binding:"required"` // tipe data yang seharusnya time.Time dijadikan string juga pada create request
	KaryawanID  string `json:"karyawan_id" binding:"required"`
	JenisIjinID int    `json:"jenis_ijin_id" binding:"required"`
	Keterangan  string `json:"keterangan"`
}

type UpdateIjinKaryawanRequest struct {
	Tanggal     *string `json:"tanggal" binding:"omitempty"`
	KaryawanID  *string `json:"karyawan_id" binding:"omitempty"`
	JenisIjinID *int    `json:"jenis_ijin_id" binding:"omitempty"`
	Keterangan  *string `json:"keterangan" binding:"omitempty"`
}
