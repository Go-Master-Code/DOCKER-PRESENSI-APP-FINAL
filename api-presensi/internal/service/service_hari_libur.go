package service

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
	"api-presensi/internal/repository"
	"errors"
	"log"
	"time"
)

// interface
type ServiceHariLibur interface {
	GetAllHariLibur() ([]dto.HariLiburResponse, error)
	GetHariLiburByID(id int) (dto.HariLiburResponse, error)
	CreateHariLibur(hl dto.CreateHariLiburRequest) (dto.HariLiburResponse, error)
	DeleteHariLibur(id int) (dto.HariLiburResponse, error)
	UpdateHariLibur(id int, req dto.UpdateHariLiburRequest) (dto.HariLiburResponse, error)
	GetHariKerjaPerPeriode(tglAwal, tglAkhir time.Time) (int, error)
}

// struct implementasi
type serviceHariLibur struct {
	repo repository.RepositoryHariLibur
}

// constructor
func NewServiceHariLibur(repo repository.RepositoryHariLibur) ServiceHariLibur {
	return &serviceHariLibur{repo}
}

// struct method
func (s *serviceHariLibur) GetAllHariLibur() ([]dto.HariLiburResponse, error) {
	hariLibur, err := s.repo.GetAllHariLibur()
	if err != nil {
		return []dto.HariLiburResponse{}, err
	}

	// convert model to dto
	hariLiburDTO := helper.ConvertToDTOHariLiburPlural(hariLibur)
	return hariLiburDTO, nil
}

func (s *serviceHariLibur) GetHariLiburByID(id int) (dto.HariLiburResponse, error) {
	hariLibur, err := s.repo.GetHariLiburByID(id)
	if err != nil {
		return dto.HariLiburResponse{}, err
	}

	// convert model to dto
	hariLiburDTO := helper.ConvertToDTOHariLiburSingle(hariLibur)
	return hariLiburDTO, nil
}

func (s *serviceHariLibur) CreateHariLibur(hl dto.CreateHariLiburRequest) (dto.HariLiburResponse, error) {
	// cek apakah tanggal libur sudah ada
	exist := s.repo.ExistByDate(hl.Tanggal)

	// jika data sudah ada
	if exist {
		return dto.HariLiburResponse{}, errors.New("hari libur sudah ada")
	}

	// parsing string dari dto (request body) ke tipe data date
	tgl, err := time.Parse("2006-01-02", hl.Tanggal)
	if err != nil {
		return dto.HariLiburResponse{}, err
	}

	// buat var model untuk eksekusi repo
	req := model.HariLibur{
		Tanggal:    tgl, // tgl yang sudah di convert (tipe data time.Time)
		Keterangan: hl.Keterangan,
	}

	newHariLibur, err := s.repo.CreateHariLibur(req)
	if err != nil {
		return dto.HariLiburResponse{}, err
	}

	// convert model to dto
	newHariLiburDTO := helper.ConvertToDTOHariLiburSingle(newHariLibur)
	return newHariLiburDTO, nil
}

func (s *serviceHariLibur) DeleteHariLibur(id int) (dto.HariLiburResponse, error) {
	hariLibur, err := s.repo.DeleteHariLibur(id)
	if err != nil {
		return dto.HariLiburResponse{}, err
	}

	// convert model to dto
	hariLiburDTO := helper.ConvertToDTOHariLiburSingle(hariLibur)
	return hariLiburDTO, nil
}

func (s *serviceHariLibur) UpdateHariLibur(id int, req dto.UpdateHariLiburRequest) (dto.HariLiburResponse, error) {
	// map untuk update data repository
	var updateMap = map[string]any{}

	if (req.Tanggal) != nil {
		updateMap["tanggal"] = req.Tanggal
	}
	if (req.Keterangan) != nil {
		updateMap["keterangan"] = req.Keterangan
	}

	updatedHariLibur, err := s.repo.UpdateHariLibur(id, updateMap)
	if err != nil {
		return dto.HariLiburResponse{}, err
	}

	// convert model to dto
	hariLiburDTO := helper.ConvertToDTOHariLiburSingle(updatedHariLibur)
	return hariLiburDTO, nil
}

// func untuk menghitung jumlah hari kerja pada suatu periode
func (s *serviceHariLibur) GetHariKerjaPerPeriode(tglAwal, tglAkhir time.Time) (int, error) {
	// ambil hari libur per bulan dari repo
	// log.Println("Tgl Awal dari handler: ", tglAwal)
	//log.Println("Tgl akhir dari handler :", tglAkhir)
	hariLibur, err := s.repo.GetHariKerjaPerPeriode(tglAwal, tglAkhir)
	if err != nil {
		return 0, err
	}

	log.Println("Daftar hari libur:", hariLibur)

	// masukkan hari libur ke dalam map
	liburMap := make(map[string]bool)
	for _, l := range hariLibur {
		// iterasi hariLibur, masukkan ke dalam map (key=tanggal, value=true)
		liburMap[l.Tanggal.Format("2006-01-02")] = true
	}

	// hitung hari kerja: Senin - Jumat, dan bukan hari libur
	totalHariKerja := 0
	for hk := tglAwal; !hk.After(tglAkhir); hk = hk.AddDate(0, 0, 1) {
		// hk = tgl awal, hk <= end, hk+=1
		weekday := hk.Weekday() // returns the day of the week
		// iterasi hari kerja hanya dari senin - jumat
		if weekday >= time.Monday && weekday <= time.Friday {
			if !liburMap[hk.Format("2006-01-02")] { // jika value hari tsb = false di map (artinya bukan hari libur nasional)
				// tambah jumlah hari kerja
				totalHariKerja++
			}
		}
	}

	return totalHariKerja, nil
}
