package dto

type KaryawanResponse struct { // response json untuk entitas karyawan
	ID    string `json:"id"`
	Nama  string `json:"nama"`
	Aktif bool   `json:"aktif"`
}

type CreateKaryawanRequest struct {
	ID    string `json:"id" binding:"required,min=10,max=10"`  // ID user 10 digit pas
	Nama  string `json:"nama" binding:"required,min=3,max=50"` // required, ada min dan max char
	Aktif bool   `json:"aktif"`                                // ga usah pake binding required untuk type bool
}

type UpdateKaryawanRequest struct {
	Nama  *string `json:"nama" binding:"omitempty"`
	Aktif *bool   `json:"aktif" binding:"omitempty"`
}
