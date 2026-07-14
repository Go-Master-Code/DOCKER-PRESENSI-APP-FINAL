package service

import (
	"api-presensi/helper"
	"api-presensi/internal/dto"
	"api-presensi/internal/model"
	"api-presensi/internal/repository"
	"api-presensi/internal/utils/crypto"
	"errors"
	"log"
)

// interface
type ServiceUser interface {
	GetAllUser() ([]dto.UserResponse, error)
	GetUserByID(id int) (dto.UserResponse, error)
	CreateUser(user dto.CreateUserRequest) (dto.UserResponse, error)
	UpdateUser(id int, user dto.UpdateUserRequest) (dto.UserResponse, error)
	DeleteUserByID(id int) (dto.UserResponse, error)
	Login(username, password string) (dto.UserResponse, error)
}

// struct implementasi
type serviceUser struct {
	repo repository.RepositoryUser
}

// constructor
func NewServiceUser(repo repository.RepositoryUser) ServiceUser {
	return &serviceUser{repo}
}

// struct method
func (s *serviceUser) GetAllUser() ([]dto.UserResponse, error) {
	user, err := s.repo.GetAllUser()
	if err != nil {
		return []dto.UserResponse{}, err
	}

	// convert model to dto
	userDTO := helper.ConvertToDTOUserPlural(user)

	return userDTO, nil
}

func (s *serviceUser) GetUserByID(id int) (dto.UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// convert model to dto
	userDTO := helper.ConvertToDTOUserSingle(user)

	return userDTO, nil
}

func (s *serviceUser) CreateUser(user dto.CreateUserRequest) (dto.UserResponse, error) {
	// cek dulu apakah username sudah ada (ingat username punya key UNIQUE di db)
	exist := s.repo.ExistByUsername(user.Username)

	if exist { // jika exist bernilai true (username sudah ada)
		log.Println("username sudah ada")
		return dto.UserResponse{}, errors.New("username sudah terdaftar")
	}

	// convert ke model untuk param repo
	req := model.User{
		//Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
		RoleID:   user.RoleID,
	}

	newUser, err := s.repo.CreateUser(req)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// convert model to dto
	newUserDTO := helper.ConvertToDTOUserSingle(newUser)

	return newUserDTO, nil
}

func (s *serviceUser) UpdateUser(id int, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	// cek dulu apakah nama karyawan sudah exist
	exist := s.repo.CekUpdateExistByUsername(id, *req.Username)
	if exist {
		return dto.UserResponse{}, errors.New("username sudah ada")
	}

	// convert dto ke map untuk parameter repo
	var updateMap = map[string]any{}

	// if req.Email != nil {
	// 	updateMap["email"] = req.Email
	// }
	if req.Username != nil {
		updateMap["username"] = req.Username
	}
	if req.Password != nil {
		updateMap["password"] = req.Password
	}
	if req.RoleID != nil {
		updateMap["role_id"] = req.RoleID
	}

	updatedUser, err := s.repo.UpdateUser(id, updateMap)
	if err != nil {
		return dto.UserResponse{}, nil
	}

	// convert model to dto
	updatedUserDTO := helper.ConvertToDTOUserSingle(updatedUser)

	return updatedUserDTO, nil

}

func (s *serviceUser) DeleteUserByID(id int) (dto.UserResponse, error) {
	user, err := s.repo.DeleteUserByID(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// convert model to dto
	userDTO := helper.ConvertToDTOUserSingle(user)

	return userDTO, nil
}

func (s *serviceUser) Login(username, password string) (dto.UserResponse, error) {
	// get data user dengan username = username di param
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return dto.UserResponse{}, errors.New("username tidak ditemukan")
	}

	// setelah dapat data user, compare password di param dengan password yang di hash di db
	valid := crypto.CheckPassword(password, user.Password)
	if !valid {
		return dto.UserResponse{}, errors.New("username atau password salah")
	}

	// jika password dan hash sama, convert model to dto, return
	userDTO := helper.ConvertToDTOUserSingle(user)
	return userDTO, nil
}
