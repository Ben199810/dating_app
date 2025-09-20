package service

import (
	"errors"
	"golang_dev_docker/component/validator"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) *AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
	User  *entity.UserInformation
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	// 使用通用驗證器驗證輸入
	if err := validator.ValidateLoginInput(req.Email, req.Password); err != nil {
		return nil, err
	}

	// 根據 email 獲取用戶
	user, err := s.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用戶不存在或密碼錯誤")
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用戶不存在或密碼錯誤")
	}

	// 更新最後登入時間
	if err := s.authRepo.UpdateLastLoginTime(user.ID); err != nil {
		// 記錄日誌但不影響登入流程
		// log.Printf("更新最後登入時間失敗: %v", err)
	}

	// 生成 token
	token := s.generateToken(user.ID)

	// 清空密碼再回傳
	user.Password = ""

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// generateToken 生成簡單的 token（實際應使用 JWT）
func (s *AuthService) generateToken(userID int) string {
	// 這裡簡化實作，實際應使用 JWT
	return strconv.FormatInt(time.Now().Unix(), 10) + "_" + strconv.Itoa(userID)
}
