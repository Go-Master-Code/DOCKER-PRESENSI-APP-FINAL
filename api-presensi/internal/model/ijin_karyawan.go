package model

import (
	"time"

	"gorm.io/gorm"
)

type IjinKaryawan struct {
	ID         int       `json:"id"`
	Tanggal    time.Time `json:"tanggal"`
	KaryawanID string    `json:"karyawan_id"`
	// Karyawan    Karyawan       `json:"karyawan" gorm:"foreignKey:KaryawanID;references:ID"`
	Karyawan    Karyawan `json:"karyawan"`
	JenisIjinID int      `json:"jenis_ijin_id"`
	//JenisIjin   JenisIjin      `json:"jenis_ijin" gorm:"foreignKey:JenisIjinID;references:ID"`
	JenisIjin  JenisIjin      `json:"jenis_ijin"`
	Keterangan string         `json:"keterangan"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (IjinKaryawan) TableName() string {
	return "ijin_karyawan"
}
