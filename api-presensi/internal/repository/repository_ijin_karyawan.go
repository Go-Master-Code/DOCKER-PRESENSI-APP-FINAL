package repository

import (
	"api-presensi/internal/model"

	"gorm.io/gorm"
)

// interface
type RepositoryIjinKaryawan interface {
	GetAllIjinKaryawan() ([]model.IjinKaryawan, error)
	GetIjinKaryawanPerTanggal(tanggal string) ([]model.IjinKaryawan, error)
	GetIjinKaryawanByID(id int) (model.IjinKaryawan, error)
	CreateIjinKaryawan(ijin model.IjinKaryawan) (model.IjinKaryawan, error)
	UpdateIjinKaryawan(id int, updateMap map[string]any) (model.IjinKaryawan, error)
	GetIjinAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]model.IjinKaryawan, error)
	DeleteIjinKaryawan(id int) (model.IjinKaryawan, error)
}

// struct implementasi
type repositoryIjinKaryawan struct {
	db *gorm.DB
}

// constructor
func NewRepositoryIjinKaryawan(db *gorm.DB) RepositoryIjinKaryawan {
	return &repositoryIjinKaryawan{db}
}

// struct method
func (r *repositoryIjinKaryawan) GetAllIjinKaryawan() ([]model.IjinKaryawan, error) {
	var ijin []model.IjinKaryawan
	err := r.db.Preload("Karyawan").Preload("JenisIjin").Find(&ijin).Error
	return ijin, err
}

func (r *repositoryIjinKaryawan) GetIjinKaryawanPerTanggal(tanggal string) ([]model.IjinKaryawan, error) {
	var ijin []model.IjinKaryawan
	err := r.db.Preload("Karyawan").Preload("JenisIjin").Where("tanggal = ?", tanggal).Find(&ijin).Error
	return ijin, err
}

func (r *repositoryIjinKaryawan) GetIjinKaryawanByID(id int) (model.IjinKaryawan, error) {
	var ijin model.IjinKaryawan
	err := r.db.Preload("Karyawan").Preload("JenisIjin").First(&ijin, id).Error
	return ijin, err
}

func (r *repositoryIjinKaryawan) CreateIjinKaryawan(ijin model.IjinKaryawan) (model.IjinKaryawan, error) {
	err := r.db.Create(&ijin).Error
	if err != nil {
		return model.IjinKaryawan{}, err
	}

	// preload semua relasi untuk response dto
	err = r.db.Preload("Karyawan").Preload("JenisIjin").First(&ijin).Error
	if err != nil {
		return model.IjinKaryawan{}, err
	}

	// jika sukses
	return ijin, nil
}

func (r *repositoryIjinKaryawan) UpdateIjinKaryawan(id int, updateMap map[string]any) (model.IjinKaryawan, error) {
	// untuk update dan delete, wajib get data dulu
	var ijin model.IjinKaryawan
	err := r.db.First(&ijin, id).Error
	if err != nil {
		return model.IjinKaryawan{}, err
	}

	// untuk menghindari error relasi tabel, kosongkan data struct JenisIjin dan Karyawan yang terhubung dengan struct IjinKaryawan ini
	ijin.JenisIjin = model.JenisIjin{}
	ijin.Karyawan = model.Karyawan{}

	// lakukan update data berdasarkan data map
	err = r.db.Model(&ijin).Updates(updateMap).Error
	if err != nil {
		return model.IjinKaryawan{}, err
	}

	// reload data + preload relasi untuk mendapatkan data terbaru
	err = r.db.Preload("Karyawan").Preload("JenisIjin").First(&ijin).Error
	if err != nil {
		return model.IjinKaryawan{}, err
	}

	return ijin, nil
}

func (r *repositoryIjinKaryawan) GetIjinAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]model.IjinKaryawan, error) {
	var ijin []model.IjinKaryawan
	err := r.db.Preload("JenisIjin").Preload("Karyawan").Where("tanggal between ? and ?", tglAwal, tglAkhir).Find(&ijin).Error
	return ijin, err
}

func (r *repositoryIjinKaryawan) DeleteIjinKaryawan(id int) (model.IjinKaryawan, error) {
	// wajib get data + preloadnya
	var ijin model.IjinKaryawan
	err := r.db.Preload("Karyawan").Preload("JenisIjin").Where("id = ?", id).First(&ijin).Error
	if err != nil {
		return model.IjinKaryawan{}, err
	}

	err = r.db.Delete(&ijin).Error
	if err != nil {
		return model.IjinKaryawan{}, err
	}

	return ijin, nil
}
