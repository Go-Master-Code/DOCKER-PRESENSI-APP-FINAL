package model

import (
	"time"

	"gorm.io/gorm"
)

type Karyawan struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Nama         string         `json:"nama"`
	Aktif        bool           `json:"aktif"`
	IjinKaryawan []IjinKaryawan // reverse relation: 1 karyawan bisa beberapa kali ijin
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (Karyawan) TableName() string {
	return "karyawan"
}
