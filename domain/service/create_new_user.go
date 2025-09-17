package user

import (
	"errors"
	"golang_dev_docker/domain/entity"
	"strings"
)

type NewUserInput struct {
	user *entity.UserInformation
}

func NewUser(Username string, Email string, Password string) (NewUserInput, error) {
	// 驗證名稱
	if strings.TrimSpace(Username) == "" {
		return nil, errors.New("名稱不能為空")
	}

	// 驗證 email
	if strings.TrimSpace(Email) == "" {
		return nil, errors.New("email 不能為空")
	}

	// 簡單的 email 格式驗證
	if !strings.Contains(Email, "@") || !strings.Contains(Email, ".") {
		return nil, errors.New("email 格式不正確")
	}

	// 驗證密碼
	if strings.TrimSpace(Password) == "" {
		return nil, errors.New("密碼不能為空")
	}

	// 密碼長度驗證
	if len(Password) < 6 {
		return nil, errors.New("密碼長度至少需要 6 個字符")
	}

	// 建立新用戶
	user := &entity.UserInformation{
		Username: Username,
		Email:    Email,
		Password: Password, // 注意：實際應用中應該對密碼進行哈希處理
	}

	return NewUserInput{user: user}, nil

}
