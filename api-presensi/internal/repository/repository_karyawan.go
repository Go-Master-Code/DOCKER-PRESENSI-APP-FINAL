package repository

import (
	"api-presensi/internal/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

// interface (kontrak)
type RepositoryKaryawan interface {
	// berisi list method, param, dan return value dari semua func
	GetAllKaryawan() ([]model.Karyawan, error)
	GetKaryawanByID(id string) (model.Karyawan, error)
	GetKaryawanBelumAbsen(tanggal string) ([]model.Karyawan, error)
	GetKaryawanBelumIjin(tanggal string) ([]model.Karyawan, error)
	CreateKaryawan(karyawan model.Karyawan) (model.Karyawan, error)
	UpdateKaryawan(id string, updateMap map[string]any) (model.Karyawan, error)
	DeleteKaryawan(id string) (model.Karyawan, error)
	ExistByID(id string) bool
	ExistByNama(nama string) bool
	CekUpdateExistByNama(id, nama string) bool
	ImportKaryawan(newKaryawan []model.Karyawan) ([]model.Karyawan, error)
}

// struct implementasi
type repositoryKaryawan struct {
	// 1 field : gorm.DB (ambil dari package database)
	db *gorm.DB
}

// constructor
func NewRepositoryKaryawan(db *gorm.DB) RepositoryKaryawan {
	return &repositoryKaryawan{db}
}

// struct method
func (r *repositoryKaryawan) GetAllKaryawan() ([]model.Karyawan, error) {
	var karyawan []model.Karyawan
	err := r.db.Find(&karyawan).Error
	return karyawan, err
}

func (r *repositoryKaryawan) GetKaryawanByID(id string) (model.Karyawan, error) {
	var karyawan model.Karyawan
	// err := r.db.First(&karyawan, id).Error // method First mendukung param id untuk klausa where
	// method di atas harus hati-hati jika primary key bukan berupa angka (int). misal string, query yang dihasilkan akan salah
	err := r.db.Where("id = ?", id).First(&karyawan).Error
	return karyawan, err
}

func (r *repositoryKaryawan) GetKaryawanBelumAbsen(tanggal string) ([]model.Karyawan, error) {
	var karyawan []model.Karyawan
	err := r.db.Model(&model.Karyawan{}).
		Where(`
        NOT EXISTS (
            SELECT 1 
            FROM presensi_karyawan p
            WHERE p.karyawan_id = karyawan.id
            AND DATE(p.tanggal) = ?
        )`, tanggal).Find(&karyawan).Error
	return karyawan, err
}

func (r *repositoryKaryawan) GetKaryawanBelumIjin(tanggal string) ([]model.Karyawan, error) {
	var karyawan []model.Karyawan
	err := r.db.Model(&model.Karyawan{}).
		Where(`
        NOT EXISTS (
            SELECT 1 
            FROM ijin_karyawan i
            WHERE i.karyawan_id = karyawan.id
            AND DATE(i.tanggal) = ?
        )`, tanggal).Find(&karyawan).Error
	return karyawan, err
}

func (r *repositoryKaryawan) CreateKaryawan(karyawan model.Karyawan) (model.Karyawan, error) {
	err := r.db.Create(&karyawan).Error
	if err != nil {
		return model.Karyawan{}, err
	}

	// preload relation untuk ditampilkan ke response
	return karyawan, err
}

func (r *repositoryKaryawan) UpdateKaryawan(id string, updateMap map[string]any) (model.Karyawan, error) {
	var karyawan model.Karyawan

	// cari dulu data karyawan dengan id pada param
	err := r.db.Where("id = ?", id).First(&karyawan).Error

	if err != nil {
		return model.Karyawan{}, err
	}

	// jika ada relasi dengan tabel lain, kosongkan data struct yang terelasi

	// proses update data
	err = r.db.Model(&karyawan).Updates(updateMap).Error

	if err != nil {
		return model.Karyawan{}, err
	}

	// jika ada relasi preload dengan tabel lain, wajib select data lagi pakai .Preload dan .First

	return karyawan, nil
}

func (r *repositoryKaryawan) DeleteKaryawan(id string) (model.Karyawan, error) {
	// get data untuk ditampilkan di response json
	var karyawan model.Karyawan

	err := r.db.Where("id = ?", id).First(&karyawan).Error
	if err != nil {
		return model.Karyawan{}, err
	}

	err = r.db.Delete(&karyawan).Error
	if err != nil {
		return model.Karyawan{}, err
	}

	return karyawan, nil
}

func (r *repositoryKaryawan) ExistByID(id string) bool {
	var karyawan model.Karyawan
	var count int64

	r.db.Where("id = ?", id).First(&karyawan).Count(&count)
	return count > 0 // menghasilkan true jika count > 0 (artinya data sudah ada)
}

func (r *repositoryKaryawan) ExistByNama(nama string) bool {
	var karyawan model.Karyawan
	var count int64

	r.db.Where("nama = ?", nama).First(&karyawan).Count(&count)
	return count > 0 // menghasilkan true jika count > 0 (artinya data sudah ada)
}

func (r *repositoryKaryawan) CekUpdateExistByNama(id, nama string) bool {
	var karyawan model.Karyawan
	var count int64

	r.db.Where("nama = ? and id <> ?", nama, id).First(&karyawan).Count(&count)
	return count > 0
}

func (r *repositoryKaryawan) ImportKaryawan(newKaryawan []model.Karyawan) ([]model.Karyawan, error) {
	// transaction - bisa commit dan rollback
	tx := r.db.Begin()

	// ambil semua id dari input
	var ids []string
	var names []string

	for _, k := range newKaryawan {
		ids = append(ids, k.ID)
		names = append(names, k.Nama)
	}

	// 🔥 cek data yang sudah ada di DB (hindari duplicate insert)
	var existingIDs []string
	var existingNames []string

	// cek semua id ke db
	err := tx.Model(&model.Karyawan{}).Where("id IN ?", ids).Pluck("id", &existingIDs).Error
	if err != nil {
		tx.Rollback()
		return []model.Karyawan{}, err
	}

	// cek semua nama ke db
	err = tx.Model(model.Karyawan{}).Where("nama IN ?", names).Pluck("nama", &existingNames).Error
	if err != nil {
		tx.Rollback()
		return []model.Karyawan{}, err
	}

	// setelah proses di atas, existing akan berisi data id yang sudah ada di db

	// ubah ke map biar cepat lookup
	existingIDMap := make(map[string]bool) // loop existing, masukan data ke map dan statusnya true
	for _, id := range existingIDs {
		existingIDMap[id] = true
	}

	existingNamaMap := make(map[string]bool)
	for _, nama := range existingNames {
		existingNamaMap[nama] = true // ubah value map[nama] jadi true jika sudah ada di db
	}

	// filter data yang benar-benar valid
	var validData []model.Karyawan
	for _, k := range newKaryawan {
		if existingIDMap[k.ID] { // cek jika existing map dengan id dari newKaryawan bernilai true (sudah ada di db)
			continue // skip ID duplicate
		}

		if existingNamaMap[k.Nama] { // cek jika existing map dengan nama dari newKaryawan bernilai true (sudah ada di db)
			continue // skip nama duplicate
		}

		// jika id dan nama belum ada, append
		validData = append(validData, k)
		log.Println(k)
	}

	// 🔥 kalau tidak ada data valid
	if len(validData) == 0 {
		tx.Rollback()
		return []model.Karyawan{}, errors.New("tidak ada data baru untuk diimport")
	}

	// insert bulk (hanya insert data yang valid saja, jangan newKaryawan -> data masih belum tervalidasi)
	err = tx.Create(&validData).Error
	if err != nil {
		tx.Rollback()
		return []model.Karyawan{}, err
	}

	// commit transaction
	err = tx.Commit().Error
	if err != nil {
		return []model.Karyawan{}, err
	}

	return validData, nil // return data yang valid tanpa duplikat
}
