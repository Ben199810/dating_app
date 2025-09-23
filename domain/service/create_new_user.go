package service

import (
	"errors"
	"golang_dev_docker/component/validator"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository // 用於檢查用戶是否已存在
}

func NewUserService(userRepo repository.UserRepository, authRepo repository.AuthRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		authRepo: authRepo,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateUser 處理用戶註冊
func (s *UserService) CreateUser(req *CreateUserRequest) (*entity.UserInformation, error) {
	// 使用通用驗證器驗證輸入
	if err := validator.ValidateRegisterInput(req.Username, req.Email, req.Password); err != nil {
		return nil, err
	}

	// 檢查 email 是否已存在
	existingUser, err := s.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email 已被使用")
	}

	// 檢查用戶名是否已存在
	usernameExists, err := s.authRepo.UserExists(req.Email, req.Username)
	if err != nil {
		return nil, err
	}
	if usernameExists {
		return nil, errors.New("用戶名已被使用")
	}

	// 密碼加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 建立用戶
	user := &entity.UserInformation{
		Username:   strings.TrimSpace(req.Username),
		Email:      strings.TrimSpace(req.Email),
		Password:   string(hashedPassword),
		IsVerified: false,
		Status:     entity.UserStatusActive,
		// 移除 Interests，因為它現在在 UserProfile 中
	}

	// 儲存到資料庫
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// 不返回密碼
	user.Password = ""
	return user, nil
}
