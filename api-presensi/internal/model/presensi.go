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
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (Presensi) TableName() string {
	return "presensi_karyawan"
}
