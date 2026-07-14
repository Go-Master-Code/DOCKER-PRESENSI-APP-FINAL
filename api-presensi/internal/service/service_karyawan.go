package service

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
	"api-presensi/internal/repository"
	"errors"
	"log"
)

// interface kosong
type ServiceKaryawan interface {
	GetAllKaryawan() ([]dto.KaryawanResponse, error)
	GetKaryawanByID(id string) (dto.KaryawanResponse, error)
	GetKaryawanBelumAbsen(tanggal string) ([]dto.KaryawanResponse, error)
	GetKaryawanBelumIjin(tanggal string) ([]dto.KaryawanResponse, error)
	CreateKaryawan(karyawan dto.CreateKaryawanRequest) (dto.KaryawanResponse, error)
	UpdateKaryawan(id string, req dto.UpdateKaryawanRequest) (dto.KaryawanResponse, error)
	DeleteKaryawan(id string) (dto.KaryawanResponse, error)
	ImportKaryawan(req []dto.CreateKaryawanRequest) ([]dto.KaryawanResponse, error)
}

// struct implementasi
type serviceKaryawan struct {
	repo repository.RepositoryKaryawan
}

// constructor
func NewServiceKaryawan(repo repository.RepositoryKaryawan) ServiceKaryawan {
	return &serviceKaryawan{repo}
}

// struct method
func (s *serviceKaryawan) GetAllKaryawan() ([]dto.KaryawanResponse, error) {
	karyawan, err := s.repo.GetAllKaryawan() // karyawan adalah var model (dari repo)
	if err != nil {
		return []dto.KaryawanResponse{}, err
	}

	// convert model to dto
	karyawanDTO := helper.ConvertToDTOKaryawanPlural(karyawan)
	return karyawanDTO, nil
}

func (s *serviceKaryawan) GetKaryawanByID(id string) (dto.KaryawanResponse, error) {
	karyawan, err := s.repo.GetKaryawanByID(id)
	if err != nil {
		return dto.KaryawanResponse{}, err
	}

	// convert model to dto
	karyawanDTO := helper.ConvertToDTOKaryawanSingle(karyawan)
	return karyawanDTO, nil
}

func (s *serviceKaryawan) GetKaryawanBelumAbsen(tanggal string) ([]dto.KaryawanResponse, error) {
	karyawan, err := s.repo.GetKaryawanBelumAbsen(tanggal)
	if err != nil {
		return []dto.KaryawanResponse{}, err
	}

	// convert model to dto
	karyawanDTO := helper.ConvertToDTOKaryawanPlural(karyawan)
	return karyawanDTO, nil
}

func (s *serviceKaryawan) GetKaryawanBelumIjin(tanggal string) ([]dto.KaryawanResponse, error) {
	karyawan, err := s.repo.GetKaryawanBelumIjin(tanggal)
	if err != nil {
		return []dto.KaryawanResponse{}, err
	}

	// convert model to dto
	karyawanDTO := helper.ConvertToDTOKaryawanPlural(karyawan)
	return karyawanDTO, nil
}

func (s *serviceKaryawan) CreateKaryawan(karyawan dto.CreateKaryawanRequest) (dto.KaryawanResponse, error) {
	// cek apakah ID sudah ada di db
	exist := s.repo.ExistByID(karyawan.ID)

	// cek apakah exist = true (jika count > 0)
	if exist {
		log.Println("ID karyawan sudah ada!")
		return dto.KaryawanResponse{}, errors.New("id karyawan sudah ada")
	}

	// cek apakah nama karyawan sudah ada di db
	exist = s.repo.ExistByNama(karyawan.Nama)

	// cek apakah exist = true (jika count > 0)
	if exist {
		log.Println("Nama karyawan sudah ada!")
		return dto.KaryawanResponse{}, errors.New("nama karyawan sudah ada")
	}

	// buat dulu request dalam bentuk model
	req := model.Karyawan{
		ID:    karyawan.ID,
		Nama:  karyawan.Nama,
		Aktif: karyawan.Aktif,
	}

	newKaryawan, err := s.repo.CreateKaryawan(req)
	if err != nil {
		return dto.KaryawanResponse{}, err
	}

	// jika tidak ada error, convert model ke dto
	newKaryawanDTO := helper.ConvertToDTOKaryawanSingle(newKaryawan)
	return newKaryawanDTO, nil
}

func (s *serviceKaryawan) UpdateKaryawan(id string, req dto.UpdateKaryawanRequest) (dto.KaryawanResponse, error) {
	// cek dulu apakah nama karyawan sudah exist
	exist := s.repo.CekUpdateExistByNama(id, *req.Nama)
	if exist {
		return dto.KaryawanResponse{}, errors.New("nama karyawan sudah ada")
	}

	// buat dulu map untuk dijadikan param func repo
	var updateMap = map[string]any{}

	if req.Nama != nil {
		updateMap["nama"] = req.Nama
	}
	if req.Aktif != nil {
		updateMap["aktif"] = req.Aktif
	}

	updatedKaryawan, err := s.repo.UpdateKaryawan(id, updateMap)

	if err != nil {
		return dto.KaryawanResponse{}, err
	}

	// convert model to dto
	karyawanDTO := helper.ConvertToDTOKaryawanSingle(updatedKaryawan)

	return karyawanDTO, nil
}

func (s *serviceKaryawan) DeleteKaryawan(id string) (dto.KaryawanResponse, error) {
	karyawan, err := s.repo.DeleteKaryawan(id)
	if err != nil {
		return dto.KaryawanResponse{}, err
	}

	// convert model to dto
	karyawanDTO := helper.ConvertToDTOKaryawanSingle(karyawan)
	return karyawanDTO, nil
}

func (s *serviceKaryawan) ImportKaryawan(req []dto.CreateKaryawanRequest) ([]dto.KaryawanResponse, error) {
	// convert dan append dtoCreateRequest menjadi model
	var newKaryawan []model.Karyawan
	for _, k := range req {
		newKaryawan = append(newKaryawan, model.Karyawan{
			ID:    k.ID,
			Nama:  k.Nama,
			Aktif: k.Aktif,
		})
	}

	// eksekusi repo
	karyawan, err := s.repo.ImportKaryawan(newKaryawan)
	if err != nil {
		return []dto.KaryawanResponse{}, err
	}

	// convert model to dto
	karyawanDTO := helper.ConvertToDTOKaryawanPlural(karyawan)
	return karyawanDTO, nil
}
