package repository

import "golang_dev_docker/domain/entity"

type UserRepository interface {
	Create(user *entity.UserInformation) error
	GetByEmail(email string) (*entity.UserInformation, error)
	GetByID(id int) (*entity.UserInformation, error)
	Update(user *entity.UserInformation) error
	Delete(id int) error
}
