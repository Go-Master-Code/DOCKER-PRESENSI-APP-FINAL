package model

import (
	"time"

	"gorm.io/gorm"
)

type Presensi struct {
	ID          int            `json:"id" gorm:"primaryKey"`
	Tanggal     time.Time      `json:"tanggal"`
	KaryawanID  string         `json:"karyawan_id"`
	Karyawan    Karyawan       `json:"karyawan"`
	WaktuMasuk  string         `json:"waktu_masuk"`
	WaktuPulang string         `json:"waktu_pulang"`
	CreatedAt   time.Time      `json:"created_at;autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (Presensi) TableName() string {
	return "presensi_karyawan"
}
