package dto

type HariLiburResponse struct {
	ID         int    `json:"id"`
	Tanggal    string `json:"tanggal"`
	Hari       string `json:"hari"` // hari hanya ada di dto, di model tidak perlu
	Keterangan string `json:"keterangan"`
}

type CreateHariLiburRequest struct {
	Tanggal    string `json:"tanggal" binding:"required"`
	Keterangan string `json:"keterangan" binding:"required,min=3,max=100"`
}

type UpdateHariLiburRequest struct {
	Tanggal    *string `json:"tanggal" binding:"omitempty"`
	Keterangan *string `json:"keterangan" binding:"omitempty"`
}
