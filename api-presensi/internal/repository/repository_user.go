package repository

import (
	"api-presensi/internal/model"

	"gorm.io/gorm"
)

// interface
type RepositoryUser interface {
	GetAllUser() ([]model.User, error)
	GetUserByID(id int) (model.User, error)
	GetUserByUsername(username string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	ExistByUsername(username string) bool
	CekUpdateExistByUsername(id int, username string) bool
	UpdateUser(id int, updateMap map[string]any) (model.User, error)
	DeleteUserByID(id int) (model.User, error)
}

// struct
type repositoryUser struct {
	db *gorm.DB
}

// constructor
func NewRepositoryUser(db *gorm.DB) RepositoryUser {
	return &repositoryUser{db}
}

// struct method
func (r *repositoryUser) GetAllUser() ([]model.User, error) {
	var user []model.User
	err := r.db.Preload("Role").Find(&user).Error
	return user, err
}

func (r *repositoryUser) GetUserByID(id int) (model.User, error) {
	var user model.User
	err := r.db.Preload("Role").First(&user, id).Error
	return user, err
}

func (r *repositoryUser) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Where("username = ?", username).First(&user).Error
	return user, err
}

func (r *repositoryUser) ExistByUsername(username string) bool {
	var user model.User
	var count int64
	r.db.Where("username = ?", username).First(&user).Count(&count)
	return count > 0 // bool, jika count > 0 menghasilkan true
}

func (r *repositoryUser) CreateUser(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return model.User{}, err
	}

	// preload relasi untuk ditampilkan di output
	err = r.db.Preload("Role").First(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repositoryUser) UpdateUser(id int, updateMap map[string]any) (model.User, error) {
	// update dan delete wajib select data dulu (cari dulu datanya)
	var user model.User

	err := r.db.First(&user, id).Error
	if err != nil {
		return model.User{}, err
	}

	// kosongkan relasi struct yang terhubung dengan struct user
	user.Role = model.Role{}

	// lakukan update data
	err = r.db.Model(&user).Updates(updateMap).Error
	if err != nil {
		return model.User{}, err
	}

	// preload data untuk refresh output
	err = r.db.Preload("Role").First(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repositoryUser) DeleteUserByID(id int) (model.User, error) {
	// get data dulu untuk update dan delete agar bisa ditampilkan di response (return value)
	var user model.User
	err := r.db.Preload("Role").First(&user, id).Error // jika ada relasi dengan struct lain, langsung preload disini
	if err != nil {
		return model.User{}, err
	}

	err = r.db.Delete(&user).Error
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repositoryUser) CekUpdateExistByUsername(id int, username string) bool {
	var user model.User
	var count int64
	r.db.Where("id <> ? and username = ?", id, username).Find(&user).Count(&count)
	return count > 0
}
