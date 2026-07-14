package dto

type JenisIjinResponse struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Aktif bool   `json:"aktif"`
}

type CreateJenisIjinRequest struct {
	Nama  string `json:"nama" binding:"required,max=50"`
	Aktif bool   `json:"aktif"`
}

type UpdateJenisIjinRequest struct {
	Nama  *string `json:"nama" binding:"omitempty"`
	Aktif *bool   `json:"aktif"`
}
