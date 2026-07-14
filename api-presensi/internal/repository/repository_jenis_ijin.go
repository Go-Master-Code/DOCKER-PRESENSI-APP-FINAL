package repository

import (
	"api-presensi/internal/model"

	"gorm.io/gorm"
)

type RepositoryJenisIjin interface {
	GetAllJenisIjin() ([]model.JenisIjin, error)
	GetJenisIjinByID(id int) (model.JenisIjin, error)
	CreateJenisIjin(jenisIjin model.JenisIjin) (model.JenisIjin, error)
	UpdateJenisIjin(id int, updateMap map[string]any) (model.JenisIjin, error)
	DeleteJenisIjin(id int) (model.JenisIjin, error)
	ExistByNama(nama string) bool
}

// struct implemetasi
type repositoryJenisIjin struct {
	db *gorm.DB
}

// constructor
func NewRepositoryJenisIjin(db *gorm.DB) RepositoryJenisIjin {
	return &repositoryJenisIjin{db}
}

// struct method
func (r *repositoryJenisIjin) GetAllJenisIjin() ([]model.JenisIjin, error) {
	var jenisIjin []model.JenisIjin
	err := r.db.Find(&jenisIjin).Error
	return jenisIjin, err
}

func (r *repositoryJenisIjin) GetJenisIjinByID(id int) (model.JenisIjin, error) {
	var jenisIjin model.JenisIjin
	err := r.db.First(&jenisIjin, id).Error // klausa where untuk db bisa pakai seperti ini pada method First asalkan field primary key (id) bertipe int
	return jenisIjin, err
}

func (r *repositoryJenisIjin) CreateJenisIjin(jenisIjin model.JenisIjin) (model.JenisIjin, error) {
	err := r.db.Create(&jenisIjin).Error
	return jenisIjin, err
}

func (r *repositoryJenisIjin) UpdateJenisIjin(id int, updateMap map[string]any) (model.JenisIjin, error) {
	// untuk update dan delete data wajib select data dulu
	var jenisIjin model.JenisIjin

	err := r.db.First(&jenisIjin, id).Error
	if err != nil {
		return model.JenisIjin{}, err
	}

	// jika ada relasi dengan tabel lain, kosongkan struct yang berelasi dengan struct ini, contoh
	// jenisIjin.Ijin = model.Ijin{}

	err = r.db.Model(&jenisIjin).Updates(updateMap).Error
	if err != nil {
		return model.JenisIjin{}, err
	}

	// jika tidak ada error, tampilkan kembali data terbaru
	err = r.db.First(&jenisIjin).Error // + preload relasi
	if err != nil {
		return model.JenisIjin{}, err
	}

	return jenisIjin, nil
}

func (r *repositoryJenisIjin) DeleteJenisIjin(id int) (model.JenisIjin, error) {
	// untuk update dan delete wajib get data dulu pakai .First() untuk ditampilkan di response
	var jenisIjin model.JenisIjin
	err := r.db.First(&jenisIjin, id).Error
	if err != nil {
		return model.JenisIjin{}, err
	}

	err = r.db.Delete(&jenisIjin).Error
	if err != nil {
		return model.JenisIjin{}, err
	}

	return jenisIjin, nil
}

func (r *repositoryJenisIjin) ExistByNama(nama string) bool {
	var jenisIjin model.JenisIjin
	var count int64

	r.db.Where("nama = ?", nama).First(&jenisIjin).Count(&count)
	return count > 0 // menghasilkan true jika count > 0 (artinya data sudah ada)
}
