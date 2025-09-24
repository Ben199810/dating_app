package entity

import (
	"errors"
	"strings"
	"time"
)

// VerificationMethod 年齡驗證方法枚舉
type VerificationMethod string

const (
	VerificationMethodID            VerificationMethod = "id_card"        // 身分證
	VerificationMethodPassport      VerificationMethod = "passport"       // 護照
	VerificationMethodDriverLicense VerificationMethod = "driver_license" // 駕照
	VerificationMethodOther         VerificationMethod = "other"          // 其他
)

// IsValid 檢查驗證方法是否有效
func (vm VerificationMethod) IsValid() bool {
	validMethods := []VerificationMethod{
		VerificationMethodID,
		VerificationMethodPassport,
		VerificationMethodDriverLicense,
		VerificationMethodOther,
	}

	for _, method := range validMethods {
		if vm == method {
			return true
		}
	}
	return false
}

// GetDisplayName 獲取驗證方法的顯示名稱
func (vm VerificationMethod) GetDisplayName() string {
	switch vm {
	case VerificationMethodID:
		return "身分證"
	case VerificationMethodPassport:
		return "護照"
	case VerificationMethodDriverLicense:
		return "駕照"
	case VerificationMethodOther:
		return "其他"
	default:
		return string(vm)
	}
}

// VerificationStatus 驗證狀態枚舉
type VerificationStatus string

const (
	VerificationStatusPending  VerificationStatus = "pending"  // 待審核
	VerificationStatusApproved VerificationStatus = "approved" // 已通過
	VerificationStatusRejected VerificationStatus = "rejected" // 已拒絕
	VerificationStatusExpired  VerificationStatus = "expired"  // 已過期
)

// IsValid 檢查驗證狀態是否有效
func (vs VerificationStatus) IsValid() bool {
	validStatuses := []VerificationStatus{
		VerificationStatusPending,
		VerificationStatusApproved,
		VerificationStatusRejected,
		VerificationStatusExpired,
	}

	for _, status := range validStatuses {
		if vs == status {
			return true
		}
	}
	return false
}

// AgeVerification 年齡驗證實體
type AgeVerification struct {
	ID                uint               `gorm:"primaryKey" json:"id"`
	UserID            uint               `gorm:"uniqueIndex;not null" json:"user_id"`
	Method            VerificationMethod `gorm:"not null" json:"method"`
	DocumentNumber    string             `gorm:"not null;size:100" json:"document_number"`
	DocumentImagePath string             `gorm:"not null;size:500" json:"document_image_path"`
	Status            VerificationStatus `gorm:"not null;default:'pending'" json:"status"`

	// 從文件中提取的資訊
	ExtractedBirthDate *time.Time `json:"extracted_birth_date"`
	ExtractedAge       *int       `json:"extracted_age"`
	ExtractedName      *string    `gorm:"size:100" json:"extracted_name,omitempty"`

	// 審核資訊
	ReviewerID      *uint   `gorm:"index" json:"reviewer_id,omitempty"`
	ReviewNotes     *string `gorm:"size:1000" json:"review_notes,omitempty"`
	RejectionReason *string `gorm:"size:500" json:"rejection_reason,omitempty"`

	// 時間戳記
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ReviewedAt *time.Time `json:"reviewed_at"`
	ApprovedAt *time.Time `json:"approved_at"`
	ExpiresAt  *time.Time `json:"expires_at"` // 驗證過期時間

	// 關聯 - 將在其他實體完成後添加
	// User     User  `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	// Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// Validate 驗證年齡驗證資料
func (av *AgeVerification) Validate() error {
	if av.UserID == 0 {
		return errors.New("user_id 是必填欄位")
	}

	if !av.Method.IsValid() {
		return errors.New("method 必須是有效的驗證方法")
	}

	if strings.TrimSpace(av.DocumentNumber) == "" {
		return errors.New("document_number 是必填欄位")
	}

	if len(av.DocumentNumber) > 100 {
		return errors.New("document_number 不能超過 100 字元")
	}

	if strings.TrimSpace(av.DocumentImagePath) == "" {
		return errors.New("document_image_path 是必填欄位")
	}

	if len(av.DocumentImagePath) > 500 {
		return errors.New("document_image_path 不能超過 500 字元")
	}

	if !av.Status.IsValid() {
		return errors.New("status 必須是有效的驗證狀態")
	}

	// 檢查提取的年齡
	if av.ExtractedAge != nil && (*av.ExtractedAge < 0 || *av.ExtractedAge > 150) {
		return errors.New("extracted_age 必須在合理範圍內")
	}

	// 檢查備註長度
	if av.ReviewNotes != nil && len(*av.ReviewNotes) > 1000 {
		return errors.New("review_notes 不能超過 1000 字元")
	}

	if av.RejectionReason != nil && len(*av.RejectionReason) > 500 {
		return errors.New("rejection_reason 不能超過 500 字元")
	}

	if av.ExtractedName != nil && len(*av.ExtractedName) > 100 {
		return errors.New("extracted_name 不能超過 100 字元")
	}

	return nil
}

// IsPending 檢查是否待審核
func (av *AgeVerification) IsPending() bool {
	return av.Status == VerificationStatusPending
}

// IsApproved 檢查是否已通過
func (av *AgeVerification) IsApproved() bool {
	return av.Status == VerificationStatusApproved && !av.IsExpired()
}

// IsRejected 檢查是否已拒絕
func (av *AgeVerification) IsRejected() bool {
	return av.Status == VerificationStatusRejected
}

// IsExpired 檢查是否已過期
func (av *AgeVerification) IsExpired() bool {
	if av.Status == VerificationStatusExpired {
		return true
	}

	if av.ExpiresAt != nil && time.Now().After(*av.ExpiresAt) {
		return true
	}

	return false
}

// IsValidAge 檢查提取的年齡是否滿足要求（18歲以上）
func (av *AgeVerification) IsValidAge() bool {
	if av.ExtractedAge != nil {
		return *av.ExtractedAge >= 18
	}

	if av.ExtractedBirthDate != nil {
		age := av.calculateAge(*av.ExtractedBirthDate)
		return age >= 18
	}

	return false
}

// calculateAge 計算年齡
func (av *AgeVerification) calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

// Approve 通過驗證
func (av *AgeVerification) Approve(reviewerID uint, notes string) error {
	if av.IsApproved() {
		return errors.New("驗證已通過")
	}

	if !av.IsValidAge() {
		return errors.New("年齡不符合要求，無法通過驗證")
	}

	av.Status = VerificationStatusApproved
	av.ReviewerID = &reviewerID

	cleanNotes := strings.TrimSpace(notes)
	if cleanNotes != "" {
		av.ReviewNotes = &cleanNotes
	}

	now := time.Now()
	av.ReviewedAt = &now
	av.ApprovedAt = &now
	av.UpdatedAt = now

	// 設定過期時間（2年後）
	expiresAt := now.AddDate(2, 0, 0)
	av.ExpiresAt = &expiresAt

	return nil
}

// Reject 拒絕驗證
func (av *AgeVerification) Reject(reviewerID uint, reason, notes string) error {
	if av.IsRejected() {
		return errors.New("驗證已被拒絕")
	}

	cleanReason := strings.TrimSpace(reason)
	if cleanReason == "" {
		return errors.New("拒絕驗證必須提供原因")
	}

	av.Status = VerificationStatusRejected
	av.ReviewerID = &reviewerID
	av.RejectionReason = &cleanReason

	cleanNotes := strings.TrimSpace(notes)
	if cleanNotes != "" {
		av.ReviewNotes = &cleanNotes
	}

	now := time.Now()
	av.ReviewedAt = &now
	av.UpdatedAt = now

	return nil
}

// MarkAsExpired 標記為過期
func (av *AgeVerification) MarkAsExpired() {
	if av.Status == VerificationStatusApproved {
		av.Status = VerificationStatusExpired
		av.UpdatedAt = time.Now()
	}
}

// SetExtractedInfo 設定從文件中提取的資訊
func (av *AgeVerification) SetExtractedInfo(birthDate *time.Time, name *string) {
	av.ExtractedBirthDate = birthDate
	av.ExtractedName = name

	if birthDate != nil {
		age := av.calculateAge(*birthDate)
		av.ExtractedAge = &age
	}

	av.UpdatedAt = time.Now()
}

// GetMethodDisplayName 獲取驗證方法的顯示名稱
func (av *AgeVerification) GetMethodDisplayName() string {
	return av.Method.GetDisplayName()
}

// GetRemainingDays 獲取驗證剩餘有效天數
func (av *AgeVerification) GetRemainingDays() int {
	if !av.IsApproved() || av.ExpiresAt == nil {
		return 0
	}

	days := int(time.Until(*av.ExpiresAt).Hours() / 24)
	if days < 0 {
		return 0
	}

	return days
}

// NeedsRenewal 檢查是否需要續約（剩餘30天內）
func (av *AgeVerification) NeedsRenewal() bool {
	return av.IsApproved() && av.GetRemainingDays() <= 30
}

// GetExtractedAge 獲取提取的年齡
func (av *AgeVerification) GetExtractedAge() int {
	if av.ExtractedAge != nil {
		return *av.ExtractedAge
	}

	if av.ExtractedBirthDate != nil {
		return av.calculateAge(*av.ExtractedBirthDate)
	}

	return 0
}
