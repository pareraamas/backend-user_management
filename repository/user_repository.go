package repository

import (
	"github.com/pareraamas/user_mgmt/model"
	"gorm.io/gorm"
)

// Repository untuk akses data user
// Komentar dalam Bahasa Indonesia

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update user di database
func (r *UserRepository) Update(user *model.User) error {
	return r.DB.Save(user).Error
}
