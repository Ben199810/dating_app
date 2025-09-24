package entity

import (
	"errors"
	"path/filepath"
	"strings"
	"time"
)

// PhotoType 照片類型枚舉
type PhotoType string

const (
	PhotoTypeProfile PhotoType = "profile" // 個人檔案照片
	PhotoTypeGallery PhotoType = "gallery" // 相簿照片
)

// IsValid 檢查照片類型是否有效
func (pt PhotoType) IsValid() bool {
	return pt == PhotoTypeProfile || pt == PhotoTypeGallery
}

// PhotoStatus 照片狀態枚舉
type PhotoStatus string

const (
	PhotoStatusPending  PhotoStatus = "pending"  // 待審核
	PhotoStatusApproved PhotoStatus = "approved" // 已通過
	PhotoStatusRejected PhotoStatus = "rejected" // 已拒絕
)

// IsValid 檢查照片狀態是否有效
func (ps PhotoStatus) IsValid() bool {
	return ps == PhotoStatusPending || ps == PhotoStatusApproved || ps == PhotoStatusRejected
}

// Photo 用戶照片實體
type Photo struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `gorm:"not null;index" json:"user_id"`
	Type         PhotoType   `gorm:"not null" json:"type"`
	FileName     string      `gorm:"not null;size:255" json:"file_name"`
	FilePath     string      `gorm:"not null;size:500" json:"file_path"`
	FileSize     int64       `gorm:"not null" json:"file_size"`
	MimeType     string      `gorm:"not null;size:100" json:"mime_type"`
	Width        int         `json:"width"`
	Height       int         `json:"height"`
	IsMain       bool        `gorm:"default:false" json:"is_main"`   // 是否為主照片
	DisplayOrder int         `gorm:"default:0" json:"display_order"` // 顯示順序
	Status       PhotoStatus `gorm:"not null;default:'pending'" json:"status"`

	// 審核資訊
	ReviewerID  *uint   `gorm:"index" json:"reviewer_id,omitempty"`
	ReviewNotes *string `gorm:"size:500" json:"review_notes,omitempty"`

	// 時間戳記
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ReviewedAt *time.Time `json:"reviewed_at"`

	// 關聯 - 將在其他實體完成後添加
	// User     User  `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	// Reviewer *User `gorm:"foreignKey:ReviewerID" json:"reviewer,omitempty"`
}

// Validate 驗證照片資料
func (p *Photo) Validate() error {
	if p.UserID == 0 {
		return errors.New("user_id 是必填欄位")
	}

	if !p.Type.IsValid() {
		return errors.New("type 必須是 profile 或 gallery")
	}

	if strings.TrimSpace(p.FileName) == "" {
		return errors.New("file_name 是必填欄位")
	}

	if len(p.FileName) > 255 {
		return errors.New("file_name 不能超過 255 字元")
	}

	if strings.TrimSpace(p.FilePath) == "" {
		return errors.New("file_path 是必填欄位")
	}

	if len(p.FilePath) > 500 {
		return errors.New("file_path 不能超過 500 字元")
	}

	if p.FileSize <= 0 {
		return errors.New("file_size 必須大於 0")
	}

	if strings.TrimSpace(p.MimeType) == "" {
		return errors.New("mime_type 是必填欄位")
	}

	// 檢查是否為有效的圖片格式
	if !p.IsValidImageType() {
		return errors.New("mime_type 必須是有效的圖片格式")
	}

	if p.Width < 0 || p.Height < 0 {
		return errors.New("width 和 height 不能為負數")
	}

	if !p.Status.IsValid() {
		return errors.New("status 必須是 pending、approved 或 rejected")
	}

	if p.DisplayOrder < 0 {
		return errors.New("display_order 不能為負數")
	}

	// 審核備註長度檢查
	if p.ReviewNotes != nil && len(*p.ReviewNotes) > 500 {
		return errors.New("review_notes 不能超過 500 字元")
	}

	return nil
}

// IsValidImageType 檢查是否為有效的圖片類型
func (p *Photo) IsValidImageType() bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if p.MimeType == validType {
			return true
		}
	}
	return false
}

// GetFileExtension 獲取檔案副檔名
func (p *Photo) GetFileExtension() string {
	return strings.ToLower(filepath.Ext(p.FileName))
}

// IsProfilePhoto 檢查是否為個人檔案照片
func (p *Photo) IsProfilePhoto() bool {
	return p.Type == PhotoTypeProfile
}

// IsGalleryPhoto 檢查是否為相簿照片
func (p *Photo) IsGalleryPhoto() bool {
	return p.Type == PhotoTypeGallery
}

// IsMainPhoto 檢查是否為主照片
func (p *Photo) IsMainPhoto() bool {
	return p.IsMain
}

// SetAsMain 設定為主照片
func (p *Photo) SetAsMain() {
	p.IsMain = true
	p.UpdatedAt = time.Now()
}

// UnsetAsMain 取消主照片設定
func (p *Photo) UnsetAsMain() {
	p.IsMain = false
	p.UpdatedAt = time.Now()
}

// IsPending 檢查是否待審核
func (p *Photo) IsPending() bool {
	return p.Status == PhotoStatusPending
}

// IsApproved 檢查是否已通過審核
func (p *Photo) IsApproved() bool {
	return p.Status == PhotoStatusApproved
}

// IsRejected 檢查是否已被拒絕
func (p *Photo) IsRejected() bool {
	return p.Status == PhotoStatusRejected
}

// Approve 通過照片審核
func (p *Photo) Approve(reviewerID uint, notes string) error {
	if p.IsApproved() {
		return errors.New("照片已通過審核")
	}

	p.Status = PhotoStatusApproved
	p.ReviewerID = &reviewerID

	cleanNotes := strings.TrimSpace(notes)
	if cleanNotes != "" {
		p.ReviewNotes = &cleanNotes
	}

	now := time.Now()
	p.ReviewedAt = &now
	p.UpdatedAt = now

	return nil
}

// Reject 拒絕照片審核
func (p *Photo) Reject(reviewerID uint, notes string) error {
	if p.IsRejected() {
		return errors.New("照片已被拒絕")
	}

	p.Status = PhotoStatusRejected
	p.ReviewerID = &reviewerID

	cleanNotes := strings.TrimSpace(notes)
	if cleanNotes == "" {
		return errors.New("拒絕照片必須提供原因")
	}
	p.ReviewNotes = &cleanNotes

	now := time.Now()
	p.ReviewedAt = &now
	p.UpdatedAt = now

	return nil
}

// UpdateDisplayOrder 更新顯示順序
func (p *Photo) UpdateDisplayOrder(order int) error {
	if order < 0 {
		return errors.New("display_order 不能為負數")
	}

	p.DisplayOrder = order
	p.UpdatedAt = time.Now()

	return nil
}

// GetFileSizeInMB 獲取檔案大小（MB）
func (p *Photo) GetFileSizeInMB() float64 {
	return float64(p.FileSize) / (1024 * 1024)
}

// GetDimensions 獲取照片尺寸
func (p *Photo) GetDimensions() (width, height int) {
	return p.Width, p.Height
}

// GetAspectRatio 計算長寬比
func (p *Photo) GetAspectRatio() float64 {
	if p.Height == 0 {
		return 0
	}
	return float64(p.Width) / float64(p.Height)
}

// IsSquare 檢查是否為正方形照片
func (p *Photo) IsSquare() bool {
	return p.Width == p.Height && p.Width > 0
}

// IsLandscape 檢查是否為橫向照片
func (p *Photo) IsLandscape() bool {
	return p.Width > p.Height
}

// IsPortrait 檢查是否為直向照片
func (p *Photo) IsPortrait() bool {
	return p.Height > p.Width
}
