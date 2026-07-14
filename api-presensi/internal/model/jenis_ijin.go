package model

import (
	"time"

	"gorm.io/gorm"
)

type JenisIjin struct {
	ID           int            `json:"id" gorm:"primaryKey"`
	Nama         string         `json:"nama"`
	Aktif        bool           `json:"aktif"`
	IjinKaryawan []IjinKaryawan // reverse relation: 1 jenis ijin bisa terdapat beberapa kali di ijin karyawan
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (JenisIjin) TableName() string {
	return "jenis_ijin"
}
