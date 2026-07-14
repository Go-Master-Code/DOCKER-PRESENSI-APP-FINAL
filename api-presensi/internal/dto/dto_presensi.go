package dto

type PresensiResponse struct {
	ID           int    `json:"id"`
	KaryawanID   string `json:"karyawan_id"`
	KaryawanNama string `json:"karyawan_nama"`
	Tanggal      string `json:"tanggal"`
	WaktuMasuk   string `json:"waktu_masuk"`
	Terlambat    bool   `json:"terlambat"`
	WaktuPulang  string `json:"waktu_pulang"`
}

type CreatePresensiRequest struct {
	KaryawanID  string `json:"karyawan_id" binding:"required"`
	Tanggal     string `json:"tanggal" binding:"required"`
	WaktuMasuk  string `json:"waktu_masuk" binding:"required"`
	WaktuPulang string `json:"waktu_pulang"` // tidak required karena bisa hanya dipakai untuk absen masuk, waktu pulang nya default null
}

type UpdatePresensiRequest struct {
	KaryawanID  *string `json:"karyawan_id" binding:"omitempty"`
	Tanggal     *string `json:"tanggal" binding:"omitempty"`
	WaktuMasuk  *string `json:"waktu_masuk" binding:"omitempty"`
	WaktuPulang *string `json:"waktu_pulang" binding:"omitempty"` // tidak required karena bisa hanya dipakai untuk absen masuk, waktu keluar nya default null
}

type KehadiranResult struct {
	KaryawanID string `json:"karyawan_id"`
	Nama       string `json:"nama"`
	Kehadiran  int    `json:"kehadiran"`
	JumlahIjin int    `json:"jumlah_ijin"` // untuk laporan jumlah ijin per row karyawan
}
