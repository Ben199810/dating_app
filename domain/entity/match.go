package entity

import (
	"errors"
	"time"
)

// SwipeAction 滑動動作枚舉
type SwipeAction string

const (
	SwipeActionLike SwipeAction = "like"
	SwipeActionPass SwipeAction = "pass"
)

// IsValid 檢查滑動動作是否有效
func (sa SwipeAction) IsValid() bool {
	return sa == SwipeActionLike || sa == SwipeActionPass
}

// MatchStatus 配對狀態枚舉
type MatchStatus string

const (
	MatchStatusPending   MatchStatus = "pending"   // 等待對方回應
	MatchStatusMatched   MatchStatus = "matched"   // 雙向配對成功
	MatchStatusUnmatched MatchStatus = "unmatched" // 配對失敗
)

// IsValid 檢查配對狀態是否有效
func (ms MatchStatus) IsValid() bool {
	return ms == MatchStatusPending || ms == MatchStatusMatched || ms == MatchStatusUnmatched
}

// Match 配對記錄實體
type Match struct {
	ID      uint        `gorm:"primaryKey" json:"id"`
	User1ID uint        `gorm:"not null;index" json:"user1_id"`
	User2ID uint        `gorm:"not null;index" json:"user2_id"`
	Status  MatchStatus `gorm:"not null;default:'pending'" json:"status"`

	// 滑動記錄
	User1Action SwipeAction  `gorm:"not null" json:"user1_action"`
	User2Action *SwipeAction `json:"user2_action"` // null 表示尚未回應

	// 時間戳記
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	MatchedAt *time.Time `json:"matched_at"` // 配對成功時間

	// 關聯 - 將在User實體完成後添加
	// User1 User `gorm:"foreignKey:User1ID;constraint:OnDelete:CASCADE" json:"user1"`
	// User2 User `gorm:"foreignKey:User2ID;constraint:OnDelete:CASCADE" json:"user2"`
}

// Validate 驗證配對記錄資料
func (m *Match) Validate() error {
	if m.User1ID == 0 {
		return errors.New("user1_id 是必填欄位")
	}

	if m.User2ID == 0 {
		return errors.New("user2_id 是必填欄位")
	}

	if m.User1ID == m.User2ID {
		return errors.New("用戶不能配對自己")
	}

	if !m.User1Action.IsValid() {
		return errors.New("user1_action 必須是 like 或 pass")
	}

	if m.User2Action != nil && !m.User2Action.IsValid() {
		return errors.New("user2_action 必須是 like 或 pass")
	}

	if !m.Status.IsValid() {
		return errors.New("status 必須是 pending、matched 或 unmatched")
	}

	return nil
}

// IsMutualLike 檢查是否雙向喜歡
func (m *Match) IsMutualLike() bool {
	return m.User1Action == SwipeActionLike &&
		m.User2Action != nil &&
		*m.User2Action == SwipeActionLike
}

// IsCompleted 檢查配對流程是否完成
func (m *Match) IsCompleted() bool {
	return m.User2Action != nil
}

// ProcessSwipe 處理第二個用戶的滑動動作
func (m *Match) ProcessSwipe(action SwipeAction) error {
	if !action.IsValid() {
		return errors.New("無效的滑動動作")
	}

	if m.IsCompleted() {
		return errors.New("配對已完成，無法再次滑動")
	}

	m.User2Action = &action
	m.UpdatedAt = time.Now()

	// 判斷最終狀態
	if m.IsMutualLike() {
		m.Status = MatchStatusMatched
		now := time.Now()
		m.MatchedAt = &now
	} else {
		m.Status = MatchStatusUnmatched
	}

	return nil
}

// GetPartnerID 獲取配對夥伴的用戶ID
func (m *Match) GetPartnerID(userID uint) (uint, error) {
	if userID == m.User1ID {
		return m.User2ID, nil
	} else if userID == m.User2ID {
		return m.User1ID, nil
	}

	return 0, errors.New("用戶不在此配對記錄中")
}

// IsUserInMatch 檢查指定用戶是否參與此配對
func (m *Match) IsUserInMatch(userID uint) bool {
	return userID == m.User1ID || userID == m.User2ID
}

// GetUserAction 獲取指定用戶的滑動動作
func (m *Match) GetUserAction(userID uint) (*SwipeAction, error) {
	if userID == m.User1ID {
		return &m.User1Action, nil
	} else if userID == m.User2ID {
		return m.User2Action, nil
	}

	return nil, errors.New("用戶不在此配對記錄中")
}

// IsWaitingForResponse 檢查是否等待指定用戶回應
func (m *Match) IsWaitingForResponse(userID uint) bool {
	return userID == m.User2ID && m.User2Action == nil
}

// CanChat 檢查用戶是否可以開始聊天
func (m *Match) CanChat() bool {
	return m.Status == MatchStatusMatched
}
