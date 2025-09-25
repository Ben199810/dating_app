package entity

import (
	"errors"
	"strings"
	"time"
)

// BlockReason 封鎖原因枚舉
type BlockReason string

const (
	BlockReasonInappropriateBehavior BlockReason = "inappropriate_behavior"
	BlockReasonHarassment            BlockReason = "harassment"
	BlockReasonSpam                  BlockReason = "spam"
	BlockReasonNotInterested         BlockReason = "not_interested"
	BlockReasonFakeProfile           BlockReason = "fake_profile"
	BlockReasonOther                 BlockReason = "other"
)

// IsValid 檢查封鎖原因是否有效
func (br BlockReason) IsValid() bool {
	validReasons := []BlockReason{
		BlockReasonInappropriateBehavior,
		BlockReasonHarassment,
		BlockReasonSpam,
		BlockReasonNotInterested,
		BlockReasonFakeProfile,
		BlockReasonOther,
	}

	for _, reason := range validReasons {
		if br == reason {
			return true
		}
	}
	return false
}

// GetDisplayName 獲取封鎖原因的顯示名稱
func (br BlockReason) GetDisplayName() string {
	switch br {
	case BlockReasonInappropriateBehavior:
		return "不當行為"
	case BlockReasonHarassment:
		return "騷擾"
	case BlockReasonSpam:
		return "垃圾訊息"
	case BlockReasonNotInterested:
		return "不感興趣"
	case BlockReasonFakeProfile:
		return "虛假檔案"
	case BlockReasonOther:
		return "其他"
	default:
		return string(br)
	}
}

// Block 封鎖記錄實體
type Block struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	BlockerID uint        `gorm:"not null;index" json:"blocker_id"`
	BlockedID uint        `gorm:"not null;index" json:"blocked_id"`
	Reason    BlockReason `gorm:"not null" json:"reason"`
	Notes     *string     `gorm:"size:500" json:"notes,omitempty"`
	IsActive  bool        `gorm:"default:true" json:"is_active"`

	// 時間戳記
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	UnblockedAt *time.Time `json:"unblocked_at"`

	// 關聯 - 將在其他實體完成後添加
	// Blocker User `gorm:"foreignKey:BlockerID;constraint:OnDelete:CASCADE" json:"blocker"`
	// Blocked User `gorm:"foreignKey:BlockedID;constraint:OnDelete:CASCADE" json:"blocked"`
}

// Validate 驗證封鎖記錄資料
func (b *Block) Validate() error {
	if b.BlockerID == 0 {
		return errors.New("blocker_id 是必填欄位")
	}

	if b.BlockedID == 0 {
		return errors.New("blocked_id 是必填欄位")
	}

	if b.BlockerID == b.BlockedID {
		return errors.New("不能封鎖自己")
	}

	if !b.Reason.IsValid() {
		return errors.New("reason 必須是有效的封鎖原因")
	}

	// 檢查備註長度（如果有提供）
	if b.Notes != nil && len(*b.Notes) > 500 {
		return errors.New("notes 不能超過 500 字元")
	}

	return nil
}

// IsActiveBlock 檢查封鎖是否仍然有效
func (b *Block) IsActiveBlock() bool {
	return b.IsActive && b.UnblockedAt == nil
}

// Unblock 解除封鎖
func (b *Block) Unblock() {
	b.IsActive = false
	now := time.Now()
	b.UnblockedAt = &now
	b.UpdatedAt = now
}

// Reblock 重新封鎖
func (b *Block) Reblock() error {
	if b.IsActiveBlock() {
		return errors.New("用戶已被封鎖")
	}

	b.IsActive = true
	b.UnblockedAt = nil
	b.UpdatedAt = time.Now()

	return nil
}

// IsBlockedBy 檢查是否被指定用戶封鎖
func (b *Block) IsBlockedBy(blockerID uint) bool {
	return b.BlockerID == blockerID && b.IsActiveBlock()
}

// IsBlocking 檢查是否封鎖了指定用戶
func (b *Block) IsBlocking(blockedID uint) bool {
	return b.BlockedID == blockedID && b.IsActiveBlock()
}

// GetReasonDisplayName 獲取封鎖原因的顯示名稱
func (b *Block) GetReasonDisplayName() string {
	return b.Reason.GetDisplayName()
}

// GetDuration 計算封鎖持續時間
func (b *Block) GetDuration() time.Duration {
	if b.IsActiveBlock() {
		return time.Since(b.CreatedAt)
	}

	if b.UnblockedAt != nil {
		return b.UnblockedAt.Sub(b.CreatedAt)
	}

	return 0
}

// AddNotes 添加或更新備註
func (b *Block) AddNotes(notes string) {
	cleanNotes := strings.TrimSpace(notes)
	if cleanNotes == "" {
		b.Notes = nil
	} else {
		b.Notes = &cleanNotes
	}
	b.UpdatedAt = time.Now()
}

// HasNotes 檢查是否有備註
func (b *Block) HasNotes() bool {
	return b.Notes != nil && strings.TrimSpace(*b.Notes) != ""
}

// GetNotes 獲取備註（如果有）
func (b *Block) GetNotes() string {
	if b.Notes != nil {
		return *b.Notes
	}
	return ""
}
