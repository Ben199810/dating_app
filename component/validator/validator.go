package validator

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email 不能為空")
	}

	// 更嚴格的 email 格式驗證
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("email 格式不正確")
	}

	return nil
}

func ValidatePassword(password string) error {
	password = strings.TrimSpace(password)
	if password == "" {
		return errors.New("密碼不能為空")
	}

	if len(password) < 6 {
		return errors.New("密碼長度不能少於 6 個字元")
	}

	return nil
}

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return errors.New("名稱不能為空")
	}

	if len(username) < 2 {
		return errors.New("用戶名長度至少需要 2 個字符")
	}

	if len(username) > 10 {
		return errors.New("用戶名長度不能超過 10 個字符")
	}

	// 用戶名只能包含字母、數字和下劃線
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("用戶名只能包含字母、數字和下劃線")
	}

	return nil
}

func ValidateLoginInput(email, password string) error {
	if err := ValidateEmail(email); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	return nil
}

func ValidateRegisterInput(username, email, password string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}

	if err := ValidateEmail(email); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	return nil
}
