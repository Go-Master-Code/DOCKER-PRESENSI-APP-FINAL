package repository

import (
	"api-presensi/internal/model"

	"gorm.io/gorm"
)

type RepositoryLog interface {
	CreateLog(log model.Log) (model.Log, error)
}

// struct implementasi
type repositoryLog struct {
	db *gorm.DB
}

// constructor
func NewRepositoryLog(db *gorm.DB) RepositoryLog {
	return &repositoryLog{db}
}

// struct method
func (r *repositoryLog) CreateLog(log model.Log) (model.Log, error) {
	err := r.db.Create(&log).Error
	return log, err
}
