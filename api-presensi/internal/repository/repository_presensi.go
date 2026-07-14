package repository

import (
	"api-presensi/internal/dto"
	"api-presensi/internal/model"

	"gorm.io/gorm"
)

// interface
type RepositoryPresensi interface {
	GetAllPresensi() ([]model.Presensi, error)
	GetPresensiPerTanggal(tanggal string) ([]model.Presensi, error)
	GetPresensiByIDPerPerPeriode(idKaryawan string, tanggalAwal string, tanggalAkhir string) ([]model.Presensi, error) // report kehadiran per karyawan per periode
	CreatePresensi(presensi model.Presensi) (model.Presensi, error)
	CekPresensiMasuk(idKaryawan string, tanggal string) (model.Presensi, error)
	UpdateWaktuPulang(idKaryawan string, tanggal string, waktuPulang string) (model.Presensi, error)
	GetPresensiAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]dto.KehadiranResult, error) // rekap jumlah kehadiran karyawan dalam suatu periode -> untuk report juga
	GetTerlambatPerKaryawanPerPeriode(idKaryawan, tglAwal, tglAkhir string) ([]model.Presensi, error)
}

// struct implementasi
type repositoryPresensi struct {
	db *gorm.DB
}

// constructor
func NewRepositoryPresensi(db *gorm.DB) RepositoryPresensi {
	return &repositoryPresensi{db}
}

// struct method
func (r *repositoryPresensi) GetAllPresensi() ([]model.Presensi, error) {
	var presensi []model.Presensi
	err := r.db.Preload("Karyawan").Find(&presensi).Error
	return presensi, err
}

func (r *repositoryPresensi) GetPresensiPerTanggal(tanggal string) ([]model.Presensi, error) {
	var presensi []model.Presensi
	err := r.db.Preload("Karyawan").Where("tanggal = ?", tanggal).Find(&presensi).Error
	return presensi, err
}

func (r *repositoryPresensi) GetPresensiByIDPerPerPeriode(idKaryawan string, tanggalAwal string, tanggalAkhir string) ([]model.Presensi, error) {
	var presensi []model.Presensi
	err := r.db.Preload("Karyawan").Where("karyawan_id = ? and tanggal between ? and ?", idKaryawan, tanggalAwal, tanggalAkhir).Find(&presensi).Error
	return presensi, err
}

func (r *repositoryPresensi) CreatePresensi(presensi model.Presensi) (model.Presensi, error) {
	err := r.db.Create(&presensi).Error
	if err != nil {
		return model.Presensi{}, err
	}

	// load ulang data karena harus preload dengan tabel relasi (Karyawan)
	err = r.db.Preload("Karyawan").First(&presensi).Error
	if err != nil {
		return model.Presensi{}, err
	}

	return presensi, nil
}

func (r *repositoryPresensi) CekPresensiMasuk(idKaryawan string, tanggal string) (model.Presensi, error) {
	var presensi model.Presensi
	err := r.db.Where("karyawan_id = ? and tanggal = ?", idKaryawan, tanggal).First(&presensi).Error
	return presensi, err
}

func (r *repositoryPresensi) UpdateWaktuPulang(idKaryawan string, tanggal string, waktuPulang string) (model.Presensi, error) {
	// cari dulu recordnya
	var presensi model.Presensi
	err := r.db.Where("karyawan_id = ? and tanggal = ?", idKaryawan, tanggal).First(&presensi).Error
	if err != nil {
		return model.Presensi{}, err
	}

	// untuk menghindari conflict relasi, kosongkan struct relasi
	presensi.Karyawan = model.Karyawan{}

	// ubah isi var presensi (atribut waktu pulang)
	presensi.WaktuPulang = waktuPulang

	// update data presensi
	err = r.db.Model(&presensi).Update("waktu_pulang", waktuPulang).Error
	if err != nil {
		return model.Presensi{}, err
	}

	// get data setelah diupdate pakai preload terhadap struct relasi lain
	err = r.db.Preload("Karyawan").First(&presensi).Error
	if err != nil {
		return model.Presensi{}, err
	}

	return presensi, nil
}

func (r *repositoryPresensi) GetPresensiAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]dto.KehadiranResult, error) {
	var results []dto.KehadiranResult

	err := r.db.Table("karyawan k").
		Select(`
			k.id as karyawan_id,
			k.nama,
			COALESCE(p.kehadiran, 0) as kehadiran,
			COALESCE(i.jumlah_ijin, 0) as jumlah_ijin
		`).

		// 🔥 SUBQUERY KEHADIRAN
		Joins(`
			LEFT JOIN (
				SELECT 
					p.karyawan_id,
					COUNT(*) as kehadiran
				FROM presensi_karyawan p
				LEFT JOIN hari_libur l ON p.tanggal = l.tanggal
				WHERE 
					p.tanggal BETWEEN ? AND ?
					AND l.tanggal IS NULL
					AND DAYOFWEEK(p.tanggal) BETWEEN 2 AND 6
				GROUP BY p.karyawan_id
			) p ON p.karyawan_id = k.id
		`, tglAwal, tglAkhir).

		// 🔥 SUBQUERY IJIN
		Joins(`
			LEFT JOIN (
				SELECT 
					i.karyawan_id,
					COUNT(*) as jumlah_ijin
				FROM ijin_karyawan i
				WHERE 
					i.tanggal BETWEEN ? AND ?
				GROUP BY i.karyawan_id
			) i ON i.karyawan_id = k.id
		`, tglAwal, tglAkhir).
		// kondisi WHERE hanya untuk karyawan yang masih ada (belum di delete)
		Where("k.deleted_at IS NULL").
		Scan(&results).Error

	return results, err
}

func (r *repositoryPresensi) GetTerlambatPerKaryawanPerPeriode(idKaryawan, tglAwal, tglAkhir string) ([]model.Presensi, error) {
	var presensi []model.Presensi
	batasWaktuTerlambat := "07:00:00"
	err := r.db. // jika ingin mengambil properti dari tabel yang di join (karyawan) maka wajib join dulu agar kolom dari tabel karyawan bisa masuk kriteria where
			Joins("JOIN karyawan ON karyawan.id = presensi.karyawan_id").
			Preload("Karyawan").
			Where("karyawan.deleted_at IS NULL AND karyawan_id = ? and waktu_masuk > ? and tanggal between ? and ?", idKaryawan, batasWaktuTerlambat, tglAwal, tglAkhir).
			Find(&presensi).Error
	return presensi, err
}
