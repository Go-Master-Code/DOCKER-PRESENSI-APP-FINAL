package service

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
	"api-presensi/internal/repository"
	"errors"
)

type ServiceJenisIjin interface {
	GetAllJenisIjin() ([]dto.JenisIjinResponse, error)
	GetJenisIjinByID(id int) (dto.JenisIjinResponse, error)
	CreateJenisIjin(jenisIjin dto.CreateJenisIjinRequest) (dto.JenisIjinResponse, error)
	DeleteJenisIjin(id int) (dto.JenisIjinResponse, error)
	UpdateJenisIjin(id int, req dto.UpdateJenisIjinRequest) (dto.JenisIjinResponse, error)
}

// struct implementasi
type serviceJenisIjin struct {
	repo repository.RepositoryJenisIjin
}

// constructor
func NewServiceJenisIjin(repo repository.RepositoryJenisIjin) ServiceJenisIjin {
	return &serviceJenisIjin{repo}
}

// struct method
func (s *serviceJenisIjin) GetAllJenisIjin() ([]dto.JenisIjinResponse, error) {
	jenisIjin, err := s.repo.GetAllJenisIjin()
	if err != nil {
		return []dto.JenisIjinResponse{}, err
	}

	// convert model to dto
	jenisIjinDTO := helper.ConvertToDTOJenisIjinPlural(jenisIjin)
	return jenisIjinDTO, nil
}

func (s *serviceJenisIjin) GetJenisIjinByID(id int) (dto.JenisIjinResponse, error) {
	jenisIjin, err := s.repo.GetJenisIjinByID(id)
	if err != nil {
		return dto.JenisIjinResponse{}, err
	}

	// convert model to dto
	jenisIjinDTO := helper.ConvertToDTOJenisIjinSingle(jenisIjin)
	return jenisIjinDTO, nil
}

func (s *serviceJenisIjin) CreateJenisIjin(jenisIjin dto.CreateJenisIjinRequest) (dto.JenisIjinResponse, error) {
	// cek dulu apakah data jenis ijin sudah ada
	exist := s.repo.ExistByNama(jenisIjin.Nama)

	if exist {
		return dto.JenisIjinResponse{}, errors.New("nama jenis ijin sudah ada")
	}

	var req model.JenisIjin
	// convert dto.CreateJenisIjinRequest ke dalam model untuk execute repo
	req.Nama = jenisIjin.Nama
	req.Aktif = jenisIjin.Aktif

	newJenisIjin, err := s.repo.CreateJenisIjin(req)
	if err != nil {
		return dto.JenisIjinResponse{}, err
	}

	// convert model to dto
	newJenisIjinDTO := helper.ConvertToDTOJenisIjinSingle(newJenisIjin)

	return newJenisIjinDTO, nil
}

func (s *serviceJenisIjin) DeleteJenisIjin(id int) (dto.JenisIjinResponse, error) {
	// select data dulu
	jenisIjin, err := s.repo.DeleteJenisIjin(id)
	if err != nil {
		return dto.JenisIjinResponse{}, err
	}

	// convert model to dto
	jenisIjinDTO := helper.ConvertToDTOJenisIjinSingle(jenisIjin)

	return jenisIjinDTO, nil
}

func (s *serviceJenisIjin) UpdateJenisIjin(id int, req dto.UpdateJenisIjinRequest) (dto.JenisIjinResponse, error) {
	// buat request dalam format map
	var updateMap = map[string]any{}

	if req.Nama != nil {
		updateMap["nama"] = req.Nama
	}

	if req.Aktif != nil {
		updateMap["aktif"] = req.Aktif
	}

	updatedJenisIjin, err := s.repo.UpdateJenisIjin(id, updateMap)
	if err != nil {
		return dto.JenisIjinResponse{}, err
	}

	// convert model to dto
	updatedJenisIjinDTO := helper.ConvertToDTOJenisIjinSingle(updatedJenisIjin)

	return updatedJenisIjinDTO, nil
}
