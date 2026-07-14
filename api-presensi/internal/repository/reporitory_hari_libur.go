package repository

import (
	"api-presensi/internal/model"
	"time"

	"gorm.io/gorm"
)

// interface
type RepositoryHariLibur interface {
	GetAllHariLibur() ([]model.HariLibur, error)
	GetHariKerjaPerPeriode(tglAwal, tglAkhir time.Time) ([]model.HariLibur, error)
	GetHariLiburByID(id int) (model.HariLibur, error)
	CreateHariLibur(hl model.HariLibur) (model.HariLibur, error)
	UpdateHariLibur(id int, updateMap map[string]interface{}) (model.HariLibur, error)
	DeleteHariLibur(id int) (model.HariLibur, error)
	ExistByDate(tanggal string) bool
}

// struct implementasi
type repositoryHariLibur struct {
	db *gorm.DB
}

// constructor
func NewRepositoryHariLibur(db *gorm.DB) RepositoryHariLibur {
	return &repositoryHariLibur{db}
}

// struct method
func (r *repositoryHariLibur) GetAllHariLibur() ([]model.HariLibur, error) {
	var hariLibur []model.HariLibur
	err := r.db.Find(&hariLibur).Error
	return hariLibur, err
}

func (r *repositoryHariLibur) GetHariKerjaPerPeriode(tglAwal, tglAkhir time.Time) ([]model.HariLibur, error) {
	var hariLibur []model.HariLibur
	err := r.db.Where("tanggal BETWEEN ? AND ?", tglAwal, tglAkhir).Find(&hariLibur).Error
	return hariLibur, err
}

func (r *repositoryHariLibur) GetHariLiburByID(id int) (model.HariLibur, error) {
	var hariLibur model.HariLibur
	err := r.db.First(&hariLibur, id).Error
	return hariLibur, err
}

func (r *repositoryHariLibur) CreateHariLibur(hl model.HariLibur) (model.HariLibur, error) {
	err := r.db.Create(&hl).Error
	if err != nil {
		return model.HariLibur{}, err
	}

	// preload semua relasi untuk ditampilkan di response, get data ulang pakai .First()
	return hl, err
}

func (r *repositoryHariLibur) UpdateHariLibur(id int, updateMap map[string]interface{}) (model.HariLibur, error) {
	// get data dulu untuk update dan delete
	var hariLibur model.HariLibur
	err := r.db.First(&hariLibur, id).Error
	if err != nil {
		return model.HariLibur{}, err
	}

	// hindari error relasi tabel, kosongkan data struct yang berelasi

	err = r.db.Model(&hariLibur).Updates(updateMap).Error
	if err != nil {
		return model.HariLibur{}, err
	}

	// preload lagi relasi sambil get data
	err = r.db.First(&hariLibur).Error
	if err != nil {
		return model.HariLibur{}, err
	}

	return hariLibur, nil
}

func (r *repositoryHariLibur) DeleteHariLibur(id int) (model.HariLibur, error) {
	// get data dulu untuk update dan delete
	var hariLibur model.HariLibur
	err := r.db.First(&hariLibur, id).Error
	if err != nil {
		return model.HariLibur{}, err
	}

	err = r.db.Delete(&hariLibur).Error
	if err != nil {
		return model.HariLibur{}, err
	}

	return hariLibur, nil
}

func (r *repositoryHariLibur) ExistByDate(tanggal string) bool {
	var hl model.HariLibur
	var count int64
	r.db.Where("tanggal = ?", tanggal).First(&hl).Count(&count)
	return count > 0
}
