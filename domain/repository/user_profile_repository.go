package repository

import "golang_dev_docker/domain/entity"

// UserProfileRepository 用戶詳細資料儲存庫
type UserProfileRepository interface {
	CreateProfile(profile *entity.UserProfile) error
	GetProfileByUserID(userID int) (*entity.UserProfile, error)
	UpdateProfile(profile *entity.UserProfile) error
	DeleteProfile(userID int) error
}

// UserPhotoRepository 用戶照片儲存庫
type UserPhotoRepository interface {
	CreatePhoto(photo *entity.UserPhoto) error
	GetPhotosByUserID(userID int) ([]*entity.UserPhoto, error)
	GetPhotoByID(id int) (*entity.UserPhoto, error)
	UpdatePhoto(photo *entity.UserPhoto) error
	DeletePhoto(id int) error
	SetPrimaryPhoto(userID, photoID int) error
	GetPrimaryPhoto(userID int) (*entity.UserPhoto, error)
}

// UserPreferenceRepository 用戶偏好設定儲存庫
type UserPreferenceRepository interface {
	CreatePreference(preference *entity.UserPreference) error
	GetPreferenceByUserID(userID int) (*entity.UserPreference, error)
	UpdatePreference(preference *entity.UserPreference) error
	DeletePreference(userID int) error
}