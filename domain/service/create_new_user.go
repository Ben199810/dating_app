package user

import (
	"errors"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// NewUserInput 創建新用戶的輸入結構
type NewUserInput struct {
	Username string
	Email    string
	Password string
}

// UserService 用戶服務結構
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService 創建新的用戶服務
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateNewUser 創建新用戶，包含完整驗證和唯一性檢查
func (s *UserService) CreateNewUser(input NewUserInput) (*entity.User, error) {
	// 驗證用戶名
	if err := validateUsername(input.Username); err != nil {
		return nil, err
	}

	// 驗證電子郵件
	if err := validateEmail(input.Email); err != nil {
		return nil, err
	}

	// 驗證密碼
	if err := validatePassword(input.Password); err != nil {
		return nil, err
	}

	// 檢查用戶名唯一性
	if err := s.checkUsernameUniqueness(input.Username); err != nil {
		return nil, err
	}

	// 檢查電子郵件唯一性
	if err := s.checkEmailUniqueness(input.Email); err != nil {
		return nil, err
	}

	// 創建並返回新的 User 實例
	user := &entity.User{
		ID:          uuid.New(),
		Username:    strings.TrimSpace(input.Username),
		Email:       strings.ToLower(strings.TrimSpace(input.Email)),
		Password:    input.Password, // 注意：這裡應該是已經哈希過的密碼
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastLoginAt: nil, // 新用戶尚未登入
	}

	// 保存用戶到儲存庫
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// checkUsernameUniqueness 檢查用戶名唯一性
func (s *UserService) checkUsernameUniqueness(username string) error {
	username = strings.TrimSpace(username)
	exists, err := s.userRepo.ExistsByUsername(username)
	if err != nil {
		return errors.New("檢查用戶名唯一性時發生錯誤")
	}
	if exists {
		return errors.New("用戶名已被使用")
	}
	return nil
}

// checkEmailUniqueness 檢查電子郵件唯一性
func (s *UserService) checkEmailUniqueness(email string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return errors.New("檢查電子郵件唯一性時發生錯誤")
	}
	if exists {
		return errors.New("電子郵件已被使用")
	}
	return nil
}

// validateUsername 驗證用戶名
func validateUsername(username string) error {
	username = strings.TrimSpace(username)

	// 1. 不能為空
	if username == "" {
		return errors.New("用戶名不能為空")
	}

	// 2. 長度限制 (3-20 字符)
	if len(username) < 3 || len(username) > 20 {
		return errors.New("用戶名長度必須在3-20個字符之間")
	}

	// 3. 只能包含字母、數字、底線和連字號
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", username)
	if !matched {
		return errors.New("用戶名只能包含字母、數字、底線和連字號")
	}

	// 4. 不能以數字開頭
	if unicode.IsDigit(rune(username[0])) {
		return errors.New("用戶名不能以數字開頭")
	}

	// 5. 不能包含敏感詞
	forbiddenWords := []string{"admin", "root", "system", "guest", "anonymous"}
	usernameLower := strings.ToLower(username)
	for _, word := range forbiddenWords {
		if strings.Contains(usernameLower, word) {
			return errors.New("用戶名包含禁用詞彙")
		}
	}

	return nil
}

// validateEmail 驗證電子郵件
func validateEmail(email string) error {
	email = strings.TrimSpace(email)

	// 1. 不能為空
	if email == "" {
		return errors.New("電子郵件不能為空")
	}

	// 2. 基本格式驗證
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	if !matched {
		return errors.New("電子郵件格式不正確")
	}

	// 3. 長度限制
	if len(email) > 254 {
		return errors.New("電子郵件地址過長")
	}

	return nil
}

// validatePassword 驗證密碼
func validatePassword(password string) error {
	// 1. 不能為空
	if password == "" {
		return errors.New("密碼不能為空")
	}

	// 2. 長度限制 (8-50 字符)
	if len(password) < 8 || len(password) > 50 {
		return errors.New("密碼長度必須在8-50個字符之間")
	}

	// 3. 必須包含至少一個小寫字母
	hasLower := false
	// 4. 必須包含至少一個大寫字母
	hasUpper := false
	// 5. 必須包含至少一個數字
	hasDigit := false
	// 6. 必須包含至少一個特殊字符
	hasSpecial := false

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	if !hasLower {
		return errors.New("密碼必須包含至少一個小寫字母")
	}
	if !hasUpper {
		return errors.New("密碼必須包含至少一個大寫字母")
	}
	if !hasDigit {
		return errors.New("密碼必須包含至少一個數字")
	}
	if !hasSpecial {
		return errors.New("密碼必須包含至少一個特殊字符 (!@#$%^&*()_+-=[]{}|;:,.<>?)")
	}

	// 7. 不能包含常見弱密碼
	weakPasswords := []string{"password", "123456", "qwerty", "abc123"}
	passwordLower := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if strings.Contains(passwordLower, weak) {
			return errors.New("密碼不能包含常見的弱密碼模式")
		}
	}

	return nil
}
