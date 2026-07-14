package model

import (
	"time"

	"gorm.io/gorm"
)

type HariLibur struct {
	ID         int            `json:"id" gorm:"primaryKey"`
	Tanggal    time.Time      `json:"tanggal"`
	Keterangan string         `json:"keterangan"`
	CreatedAt  time.Time      `json:"created_at;autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

func (HariLibur) TableName() string {
	return "hari_libur"
}
