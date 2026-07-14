package service

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
	"api-presensi/internal/repository"
	"time"
)

// interface
type ServiceIjinKaryawan interface {
	GetAllIjinKaryawan() ([]dto.IjinKaryawanResponse, error)
	GetIjinKaryawanPerTanggal(tanggal string) ([]dto.IjinKaryawanResponse, error)
	GetIjinKaryawanByID(id int) (dto.IjinKaryawanResponse, error)
	CreateIjinKaryawan(ijin dto.CreateIjinKaryawanRequest) (dto.IjinKaryawanResponse, error)
	UpdateIjinKaryawan(id int, req dto.UpdateIjinKaryawanRequest) (dto.IjinKaryawanResponse, error)
	GetIjinAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]dto.IjinKaryawanResponse, error)
	DeleteIjinKaryawan(id int) (dto.IjinKaryawanResponse, error)
}

// struct implementasi
type serviceIjinKaryawan struct {
	repo repository.RepositoryIjinKaryawan
}

// constructor
func NewServiceIjinKaryawan(repo repository.RepositoryIjinKaryawan) ServiceIjinKaryawan {
	return &serviceIjinKaryawan{repo}
}

// struct method
func (s *serviceIjinKaryawan) GetAllIjinKaryawan() ([]dto.IjinKaryawanResponse, error) {
	ijin, err := s.repo.GetAllIjinKaryawan()
	if err != nil {
		return []dto.IjinKaryawanResponse{}, err
	}

	// convert model to dto
	ijinDTO := helper.ConvertToDTOIjinKaryawanPlural(ijin)
	return ijinDTO, nil
}

func (s *serviceIjinKaryawan) GetIjinKaryawanPerTanggal(tanggal string) ([]dto.IjinKaryawanResponse, error) {
	ijin, err := s.repo.GetIjinKaryawanPerTanggal(tanggal)
	if err != nil {
		return []dto.IjinKaryawanResponse{}, err
	}

	// convert model to dto
	ijinDTO := helper.ConvertToDTOIjinKaryawanPlural(ijin)
	return ijinDTO, nil
}

func (s *serviceIjinKaryawan) GetIjinKaryawanByID(id int) (dto.IjinKaryawanResponse, error) {
	ijin, err := s.repo.GetIjinKaryawanByID(id)
	if err != nil {
		return dto.IjinKaryawanResponse{}, err
	}

	// convert model to dto
	ijinDTO := helper.ConvertToDTOIjinKaryawanSingle(ijin)
	return ijinDTO, nil
}

func (s *serviceIjinKaryawan) CreateIjinKaryawan(ijin dto.CreateIjinKaryawanRequest) (dto.IjinKaryawanResponse, error) {
	// convert var tanggal di dto (string) ke tipe data date (atribut model)
	tglDate, err := time.Parse("2006-01-02", ijin.Tanggal)
	if err != nil {
		return dto.IjinKaryawanResponse{}, err
	}

	// buat var model (convert dto ke model) untuk eksekusi repo
	req := model.IjinKaryawan{
		Tanggal:     tglDate,
		KaryawanID:  ijin.KaryawanID,
		JenisIjinID: ijin.JenisIjinID,
		Keterangan:  ijin.Keterangan,
	}

	newIjin, err := s.repo.CreateIjinKaryawan(req)
	if err != nil {
		return dto.IjinKaryawanResponse{}, err
	}

	// convert model to dto
	ijinKaryawanDTO := helper.ConvertToDTOIjinKaryawanSingle(newIjin)

	return ijinKaryawanDTO, nil
}

func (s *serviceIjinKaryawan) UpdateIjinKaryawan(id int, req dto.UpdateIjinKaryawanRequest) (dto.IjinKaryawanResponse, error) {
	// buat var map untuk tampung data update
	var updateMap = map[string]any{}

	if req.Tanggal != nil {
		updateMap["tanggal"] = req.Tanggal
	}
	if req.KaryawanID != nil {
		updateMap["karyawan_id"] = req.KaryawanID
	}
	if req.JenisIjinID != nil {
		updateMap["jenis_ijin_id"] = req.JenisIjinID
	}
	if req.Keterangan != nil {
		updateMap["keterangan"] = req.Keterangan
	}

	updatedIjin, err := s.repo.UpdateIjinKaryawan(id, updateMap)
	if err != nil {
		return dto.IjinKaryawanResponse{}, err
	}

	// convert model to dto
	ijinDTO := helper.ConvertToDTOIjinKaryawanSingle(updatedIjin)

	return ijinDTO, nil
}

func (s *serviceIjinKaryawan) GetIjinAllKaryawanPerPeriode(tglAwal string, tglAkhir string) ([]dto.IjinKaryawanResponse, error) {
	ijin, err := s.repo.GetIjinAllKaryawanPerPeriode(tglAwal, tglAkhir)
	if err != nil {
		return []dto.IjinKaryawanResponse{}, err
	}

	// convert model to dto
	ijinDTO := helper.ConvertToDTOIjinKaryawanPlural(ijin)
	return ijinDTO, nil
}

func (s *serviceIjinKaryawan) DeleteIjinKaryawan(id int) (dto.IjinKaryawanResponse, error) {
	ijin, err := s.repo.DeleteIjinKaryawan(id)
	if err != nil {
		return dto.IjinKaryawanResponse{}, err
	}

	// convert model to dto
	ijinDTO := helper.ConvertToDTOIjinKaryawanSingle(ijin)

	return ijinDTO, nil
}
