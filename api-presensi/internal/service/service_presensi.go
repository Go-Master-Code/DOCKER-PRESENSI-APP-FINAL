package service

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
	"api-presensi/internal/repository"
	"errors"
	"time"
)

// interface
type ServicePresensi interface {
	GetAllPresensi() ([]dto.PresensiResponse, error)
	GetPresensiPerTanggal(tanggal string) ([]dto.PresensiResponse, error)
	GetPresensiByIDPerPerPeriode(idKaryawan string, tanggalAwal string, tanggalAkhir string) ([]dto.PresensiResponse, error) // untuk report presensi per karyawan per periode
	CreateOrUpdatePresensi(presensi dto.CreatePresensiRequest) (dto.PresensiResponse, error)
	GetPresensiAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]dto.KehadiranResult, error)
	GetTerlambatPerKaryawanPerPeriode(idKaryawan, tglAwal, tglAkhir string) ([]dto.PresensiResponse, error)
}

// struct implementasi
type servicePresensi struct {
	repo         repository.RepositoryPresensi
	repoKaryawan repository.RepositoryKaryawan
}

// constructor
func NewServicePresensi(repo repository.RepositoryPresensi, repoKaryawan repository.RepositoryKaryawan) ServicePresensi {
	return &servicePresensi{repo, repoKaryawan}
}

// struct method
func (s *servicePresensi) GetAllPresensi() ([]dto.PresensiResponse, error) {
	presensi, err := s.repo.GetAllPresensi()
	if err != nil {
		return []dto.PresensiResponse{}, err
	}

	// convert model to dto
	presensiDTO := helper.ConvertToDTOPresensiPlural(presensi)
	return presensiDTO, nil
}

func (s *servicePresensi) GetPresensiPerTanggal(tanggal string) ([]dto.PresensiResponse, error) {
	presensi, err := s.repo.GetPresensiPerTanggal(tanggal)
	if err != nil {
		return []dto.PresensiResponse{}, err
	}

	// convert model to dto
	presensiDTO := helper.ConvertToDTOPresensiPlural(presensi)
	return presensiDTO, nil
}

func (s *servicePresensi) CreateOrUpdatePresensi(presensi dto.CreatePresensiRequest) (dto.PresensiResponse, error) {
	// cek apakah ID karyawan terdaftar di sistem?
	exist := s.repoKaryawan.ExistByID(presensi.KaryawanID)
	if !exist {
		return dto.PresensiResponse{}, errors.New("id tidak terdaftar")
	}

	// convert var tanggal di dto (string) ke tipe data date (atribut model)
	tglDate, err := time.Parse("2006-01-02", presensi.Tanggal)
	if err != nil {
		return dto.PresensiResponse{}, err
	}

	// buat var model sebagai parameter ke repo
	var req = model.Presensi{
		Tanggal:    tglDate,
		KaryawanID: presensi.KaryawanID,
		// waktu masuk dan pulang ditentukan berdasarkan hasil method CekPresensiMasuk
		// Keterangan: presensi.Keterangan,
	}

	// cek apakah presensi masuk sudah dilakukan
	data, _ := s.repo.CekPresensiMasuk(presensi.KaryawanID, presensi.Tanggal) // ignore error karena jika data presensi masuk belum ditemukan akan muncul error record not found

	// cek apakah karyawan tertentu (data) sudah melakukan presensi masuk
	if data.KaryawanID == "" { // jika KaryawanID nya empty string, artinya karyawan itu belum absen masuk, eksekusi repo->create data
		req.WaktuMasuk = presensi.WaktuMasuk
		req.WaktuPulang = presensi.WaktuMasuk

		newPresensi, err := s.repo.CreatePresensi(req)
		if err != nil {
			return dto.PresensiResponse{}, err
		}

		// convert model to dto
		presensiDTO := helper.ConvertToDTOPresensiSingle(newPresensi)
		return presensiDTO, nil
	} else { // jika karyawanID sudah ada (artinya sudah melakukan absen masuk)
		req.WaktuPulang = presensi.WaktuPulang                                                             // ambil data waktu pulang
		updatePresensi, err := s.repo.UpdateWaktuPulang(req.KaryawanID, presensi.Tanggal, req.WaktuPulang) // param tanggal string, jadi ambil dari presensi, bukan dari req (tipe data di req time.Time)
		if err != nil {
			return dto.PresensiResponse{}, err
		}

		// convert model to dto
		updatePresensiDTO := helper.ConvertToDTOPresensiSingle(updatePresensi)
		return updatePresensiDTO, nil
	}
}

func (s *servicePresensi) GetPresensiByIDPerPerPeriode(idKaryawan string, tanggalAwal string, tanggalAkhir string) ([]dto.PresensiResponse, error) {
	presensi, err := s.repo.GetPresensiByIDPerPerPeriode(idKaryawan, tanggalAwal, tanggalAkhir)
	if err != nil {
		return []dto.PresensiResponse{}, err
	}

	// convert model to dto
	presensiDTO := helper.ConvertToDTOPresensiPlural(presensi)
	return presensiDTO, nil
}

func (s *servicePresensi) GetPresensiAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]dto.KehadiranResult, error) {
	presensi, err := s.repo.GetPresensiAllKaryawanPerPeriode(tglAwal, tglAkhir)
	if err != nil {
		return []dto.KehadiranResult{}, err
	}

	return presensi, nil
}

func (s *servicePresensi) GetTerlambatPerKaryawanPerPeriode(idKaryawan, tglAwal, tglAkhir string) ([]dto.PresensiResponse, error) {
	terlambat, err := s.repo.GetTerlambatPerKaryawanPerPeriode(idKaryawan, tglAwal, tglAkhir)
	if err != nil {
		return []dto.PresensiResponse{}, err
	}

	// convert model to dto
	terlambatDTO := helper.ConvertToDTOPresensiPlural(terlambat)
	return terlambatDTO, nil
}
