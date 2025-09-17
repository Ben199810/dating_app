package service

import (
	"errors"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"strings"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *UserService) CreateUser(req *CreateUserRequest) (*entity.UserInformation, error) {
	// 驗證輸入
	if err := s.validateInput(req.Username, req.Email, req.Password); err != nil {
		return nil, err
	}

	// 檢查 email 是否已存在
	existingUser, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email 已被使用")
	}

	// 密碼加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 建立用戶
	user := &entity.UserInformation{
		Username: strings.TrimSpace(req.Username),
		Email:    strings.TrimSpace(req.Email),
		Password: string(hashedPassword),
	}

	// 儲存到資料庫
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 不返回密碼
	user.Password = ""
	return user, nil
}

func (s *UserService) validateInput(username, email, password string) error {
	// 驗證名稱
	if strings.TrimSpace(username) == "" {
		return errors.New("名稱不能為空")
	}

	// 驗證 email
	if strings.TrimSpace(email) == "" {
		return errors.New("email 不能為空")
	}

	// 簡單的 email 格式驗證
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("email 格式不正確")
	}

	// 驗證密碼
	if strings.TrimSpace(password) == "" {
		return errors.New("密碼不能為空")
	}

	// 密碼長度驗證
	if len(password) < 6 {
		return errors.New("密碼長度至少需要 6 個字符")
	}

	return nil
}
