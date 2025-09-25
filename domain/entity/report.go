package entity

import (
	"errors"
	"strings"
	"time"
)

// ReportCategory 檢舉類別枚舉
type ReportCategory string

const (
	ReportCategoryInappropriateBehavior ReportCategory = "inappropriate_behavior"
	ReportCategoryHarassment            ReportCategory = "harassment"
	ReportCategorySpam                  ReportCategory = "spam"
	ReportCategoryFakeProfile           ReportCategory = "fake_profile"
	ReportCategoryUnderage              ReportCategory = "underage"
	ReportCategoryViolenceThreat        ReportCategory = "violence_threat"
	ReportCategoryInappropriateContent  ReportCategory = "inappropriate_content"
	ReportCategoryOther                 ReportCategory = "other"
)

// IsValid 檢查檢舉類別是否有效
func (rc ReportCategory) IsValid() bool {
	validCategories := []ReportCategory{
		ReportCategoryInappropriateBehavior,
		ReportCategoryHarassment,
		ReportCategorySpam,
		ReportCategoryFakeProfile,
		ReportCategoryUnderage,
		ReportCategoryViolenceThreat,
		ReportCategoryInappropriateContent,
		ReportCategoryOther,
	}

	for _, category := range validCategories {
		if rc == category {
			return true
		}
	}
	return false
}

// GetDisplayName 獲取檢舉類別的顯示名稱
func (rc ReportCategory) GetDisplayName() string {
	switch rc {
	case ReportCategoryInappropriateBehavior:
		return "不當行為"
	case ReportCategoryHarassment:
		return "騷擾"
	case ReportCategorySpam:
		return "垃圾訊息"
	case ReportCategoryFakeProfile:
		return "虛假檔案"
	case ReportCategoryUnderage:
		return "未成年"
	case ReportCategoryViolenceThreat:
		return "暴力威脅"
	case ReportCategoryInappropriateContent:
		return "不當內容"
	case ReportCategoryOther:
		return "其他"
	default:
		return string(rc)
	}
}

// ReportStatus 檢舉狀態枚舉
type ReportStatus string

const (
	ReportStatusPending   ReportStatus = "pending"   // 待處理
	ReportStatusReviewing ReportStatus = "reviewing" // 審查中
	ReportStatusApproved  ReportStatus = "approved"  // 已通過
	ReportStatusRejected  ReportStatus = "rejected"  // 已拒絕
	ReportStatusResolved  ReportStatus = "resolved"  // 已解決
)

// IsValid 檢查檢舉狀態是否有效
func (rs ReportStatus) IsValid() bool {
	validStatuses := []ReportStatus{
		ReportStatusPending,
		ReportStatusReviewing,
		ReportStatusApproved,
		ReportStatusRejected,
		ReportStatusResolved,
	}

	for _, status := range validStatuses {
		if rs == status {
			return true
		}
	}
	return false
}

// Report 檢舉記錄實體
type Report struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ReporterID  uint           `gorm:"not null;index" json:"reporter_id"`
	ReportedID  uint           `gorm:"not null;index" json:"reported_id"`
	Category    ReportCategory `gorm:"not null" json:"category"`
	Description string         `gorm:"not null;size:1000" json:"description"`
	Status      ReportStatus   `gorm:"not null;default:'pending'" json:"status"`

	// 管理員處理資訊
	ReviewerID  *uint   `gorm:"index" json:"reviewer_id,omitempty"`
	ReviewNotes *string `gorm:"size:1000" json:"review_notes,omitempty"`

	// 時間戳記
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ReviewedAt *time.Time `json:"reviewed_at"`
	ResolvedAt *time.Time `json:"resolved_at"`

	// 關聯 - 將在其他實體完成後添加
	// Reporter User  `gorm:"foreignKey:ReporterID;constraint:OnDelete:CASCADE" json:"reporter"`
	// Reported User  `gorm:"foreignKey:ReportedID;constraint:OnDelete:CASCADE" json:"reported"`
	// Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// Validate 驗證檢舉記錄資料
func (r *Report) Validate() error {
	if r.ReporterID == 0 {
		return errors.New("reporter_id 是必填欄位")
	}

	if r.ReportedID == 0 {
		return errors.New("reported_id 是必填欄位")
	}

	if r.ReporterID == r.ReportedID {
		return errors.New("不能檢舉自己")
	}

	if !r.Category.IsValid() {
		return errors.New("category 必須是有效的檢舉類別")
	}

	if strings.TrimSpace(r.Description) == "" {
		return errors.New("description 是必填欄位")
	}

	if len(r.Description) < 10 {
		return errors.New("description 至少需要 10 個字元")
	}

	if len(r.Description) > 1000 {
		return errors.New("description 不能超過 1000 字元")
	}

	if !r.Status.IsValid() {
		return errors.New("status 必須是有效的檢舉狀態")
	}

	return nil
}

// IsPending 檢查檢舉是否待處理
func (r *Report) IsPending() bool {
	return r.Status == ReportStatusPending
}

// IsReviewing 檢查檢舉是否審查中
func (r *Report) IsReviewing() bool {
	return r.Status == ReportStatusReviewing
}

// IsResolved 檢查檢舉是否已解決
func (r *Report) IsResolved() bool {
	return r.Status == ReportStatusApproved ||
		r.Status == ReportStatusRejected ||
		r.Status == ReportStatusResolved
}

// StartReview 開始審查檢舉
func (r *Report) StartReview(reviewerID uint) error {
	if r.IsResolved() {
		return errors.New("檢舉已解決，無法重新審查")
	}

	r.Status = ReportStatusReviewing
	r.ReviewerID = &reviewerID
	now := time.Now()
	r.ReviewedAt = &now
	r.UpdatedAt = now

	return nil
}

// Approve 通過檢舉
func (r *Report) Approve(reviewerID uint, notes string) error {
	if !r.IsReviewing() && !r.IsPending() {
		return errors.New("只有待處理或審查中的檢舉可以通過")
	}

	r.Status = ReportStatusApproved
	r.ReviewerID = &reviewerID
	r.ReviewNotes = &notes
	now := time.Now()
	r.ReviewedAt = &now
	r.ResolvedAt = &now
	r.UpdatedAt = now

	return nil
}

// Reject 拒絕檢舉
func (r *Report) Reject(reviewerID uint, notes string) error {
	if !r.IsReviewing() && !r.IsPending() {
		return errors.New("只有待處理或審查中的檢舉可以拒絕")
	}

	r.Status = ReportStatusRejected
	r.ReviewerID = &reviewerID
	r.ReviewNotes = &notes
	now := time.Now()
	r.ReviewedAt = &now
	r.ResolvedAt = &now
	r.UpdatedAt = now

	return nil
}

// Resolve 解決檢舉
func (r *Report) Resolve(reviewerID uint, notes string) error {
	if r.IsResolved() {
		return errors.New("檢舉已解決")
	}

	r.Status = ReportStatusResolved
	r.ReviewerID = &reviewerID
	r.ReviewNotes = &notes
	now := time.Now()
	r.ReviewedAt = &now
	r.ResolvedAt = &now
	r.UpdatedAt = now

	return nil
}

// IsReportedBy 檢查是否由指定用戶檢舉
func (r *Report) IsReportedBy(userID uint) bool {
	return r.ReporterID == userID
}

// IsAbout 檢查是否關於指定用戶的檢舉
func (r *Report) IsAbout(userID uint) bool {
	return r.ReportedID == userID
}

// GetCategoryDisplayName 獲取檢舉類別的顯示名稱
func (r *Report) GetCategoryDisplayName() string {
	return r.Category.GetDisplayName()
}
