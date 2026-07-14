package service

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
	"api-presensi/internal/repository"
)

type ServiceLog interface {
	CreateLog(log dto.LogRequestAndResponse) (dto.LogRequestAndResponse, error)
}

// struct implementasi
type serviceLog struct {
	repo repository.RepositoryLog
}

// constructor
func NewServiceLog(repo repository.RepositoryLog) ServiceLog {
	return &serviceLog{repo}
}

// struct method
func (s *serviceLog) CreateLog(log dto.LogRequestAndResponse) (dto.LogRequestAndResponse, error) {
	// convert dto to modal
	req := model.Log{
		UserID:    log.UserID,
		Method:    log.Method,
		Endpoint:  log.Endpoint,
		IPAddress: log.IPAddress,
		UserAgent: log.UserAgent,
	}

	newLog, err := s.repo.CreateLog(req)
	if err != nil {
		return dto.LogRequestAndResponse{}, err
	}

	// jika sukses, convert back to dto
	logDTO := helper.ConvertToDTOLogSingle(newLog)
	return logDTO, nil
}
