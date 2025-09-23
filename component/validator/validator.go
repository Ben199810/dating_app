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

// ValidateAge 驗證年齡
func ValidateAge(age *int) error {
	if age == nil {
		return nil // 年齡可以為空
	}

	if *age < 18 {
		return errors.New("年齡必須大於等於18歲")
	}

	if *age > 120 {
		return errors.New("年齡不能超過120歲")
	}

	return nil
}

// ValidateGender 驗證性別
func ValidateGender(gender *string) error {
	if gender == nil {
		return nil // 性別可以為空
	}

	validGenders := []string{"male", "female", "other"}
	for _, validGender := range validGenders {
		if *gender == validGender {
			return nil
		}
	}

	return errors.New("性別必須是 male、female 或 other")
}

// ValidateBio 驗證自我介紹
func ValidateBio(bio *string) error {
	if bio == nil {
		return nil // 自我介紹可以為空
	}

	bioText := strings.TrimSpace(*bio)
	if len(bioText) > 500 {
		return errors.New("自我介紹不能超過500個字符")
	}

	return nil
}

// ValidateInterests 驗證興趣列表
func ValidateInterests(interests []string) error {
	if len(interests) > 10 {
		return errors.New("興趣不能超過10個")
	}

	for _, interest := range interests {
		interest = strings.TrimSpace(interest)
		if len(interest) == 0 {
			return errors.New("興趣不能為空")
		}
		if len(interest) > 50 {
			return errors.New("每個興趣不能超過50個字符")
		}
	}

	return nil
}

// ValidateLocation 驗證地理位置
func ValidateLocation(lat, lng *float64) error {
	if (lat == nil && lng != nil) || (lat != nil && lng == nil) {
		return errors.New("經緯度必須同時提供或同時為空")
	}

	if lat != nil && lng != nil {
		if *lat < -90 || *lat > 90 {
			return errors.New("緯度必須在-90到90之間")
		}
		if *lng < -180 || *lng > 180 {
			return errors.New("經度必須在-180到180之間")
		}
	}

	return nil
}

// ValidateHeight 驗證身高
func ValidateHeight(height *int) error {
	if height == nil {
		return nil
	}

	if *height < 100 || *height > 250 {
		return errors.New("身高必須在100-250公分之間")
	}

	return nil
}

// ValidateWeight 驗證體重
func ValidateWeight(weight *int) error {
	if weight == nil {
		return nil
	}

	if *weight < 30 || *weight > 300 {
		return errors.New("體重必須在30-300公斤之間")
	}

	return nil
}

// ValidatePhotoURL 驗證照片URL
func ValidatePhotoURL(url string) error {
	url = strings.TrimSpace(url)
	if url == "" {
		return errors.New("照片URL不能為空")
	}

	// 基本URL格式驗證
	urlRegex := regexp.MustCompile(`^https?://.*\.(jpg|jpeg|png|gif|webp)$`)
	if !urlRegex.MatchString(strings.ToLower(url)) {
		return errors.New("照片URL格式不正確，必須是有效的圖片連結")
	}

	return nil
}

// ValidateProfileData 驗證完整的個人資料
func ValidateProfileData(age *int, gender *string, bio *string, interests []string, lat, lng *float64) error {
	if err := ValidateAge(age); err != nil {
		return err
	}

	if err := ValidateGender(gender); err != nil {
		return err
	}

	if err := ValidateBio(bio); err != nil {
		return err
	}

	if err := ValidateInterests(interests); err != nil {
		return err
	}

	if err := ValidateLocation(lat, lng); err != nil {
		return err
	}

	return nil
}

// ValidateAgeRange 驗證年齡範圍偏好
func ValidateAgeRange(minAge, maxAge *int) error {
	if minAge != nil && maxAge != nil {
		if *minAge > *maxAge {
			return errors.New("最小年齡不能大於最大年齡")
		}
		if *minAge < 18 {
			return errors.New("最小年齡不能少於18歲")
		}
		if *maxAge > 120 {
			return errors.New("最大年齡不能超過120歲")
		}
	}

	return nil
}

// ValidateDistance 驗證距離範圍
func ValidateDistance(distance *int) error {
	if distance == nil {
		return nil
	}

	if *distance < 1 || *distance > 500 {
		return errors.New("距離範圍必須在1-500公里之間")
	}

	return nil
}
