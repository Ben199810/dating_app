package repository

import (
	"context"
	"time"

	"golang_dev_docker/domain/entity"
)

// ReportRepository 檢舉數據儲存庫介面
// 提供檢舉系統的持久化操作，包括檢舉創建、審核、統計等功能
type ReportRepository interface {
	// Create 創建檢舉記錄
	// 用於用戶檢舉其他用戶的不當行為
	Create(ctx context.Context, report *entity.Report) error

	// GetByID 根據 ID 獲取檢舉記錄
	// 用於檢舉詳情查詢和審核
	GetByID(ctx context.Context, id uint) (*entity.Report, error)

	// GetByUserID 獲取用戶的檢舉記錄
	// 用於查看用戶的檢舉歷史（作為檢舉者）
	GetByUserID(ctx context.Context, userID uint, limit int) ([]*entity.Report, error)

	// GetByTargetUserID 獲取針對特定用戶的檢舉記錄
	// 用於查看用戶被檢舉的記錄
	GetByTargetUserID(ctx context.Context, targetUserID uint, limit int) ([]*entity.Report, error)

	// Update 更新檢舉記錄
	// 用於審核狀態更新和處理結果記錄
	Update(ctx context.Context, report *entity.Report) error

	// GetPendingReports 獲取待審核檢舉記錄
	// 用於管理員審核面板
	GetPendingReports(ctx context.Context, limit int) ([]*entity.Report, error)

	// SetReportStatus 設定檢舉狀態
	// 用於審核通過或拒絕操作
	SetReportStatus(ctx context.Context, reportID uint, status entity.ReportStatus, reviewerID *uint, reviewNotes string) error

	// CheckDuplicateReport 檢查重複檢舉
	// 防止同一用戶重複檢舉同一目標
	CheckDuplicateReport(ctx context.Context, userID, targetUserID uint, category entity.ReportCategory) (bool, error)

	// GetReportStats 獲取檢舉統計數據
	// 用於系統監控和用戶行為分析
	GetReportStats(ctx context.Context, params ReportStatsParams) (*ReportStats, error)

	// Delete 軟刪除檢舉記錄
	// 用於管理員刪除無效檢舉
	Delete(ctx context.Context, id uint) error
}

// BlockRepository 封鎖數據儲存庫介面
// 提供用戶封鎖功能的持久化操作，保護用戶免受騷擾
type BlockRepository interface {
	// Create 創建封鎖記錄
	// 用於用戶封鎖其他用戶
	Create(ctx context.Context, block *entity.Block) error

	// GetByID 根據 ID 獲取封鎖記錄
	// 用於封鎖記錄查詢
	GetByID(ctx context.Context, id uint) (*entity.Block, error)

	// GetByUserID 獲取用戶的封鎖列表
	// 用於封鎖管理頁面展示
	GetByUserID(ctx context.Context, userID uint) ([]*entity.Block, error)

	// IsBlocked 檢查是否封鎖關係
	// 用於配對推薦和聊天權限檢查
	IsBlocked(ctx context.Context, userID, targetUserID uint) (bool, error)

	// IsMutuallyBlocked 檢查雙向封鎖關係
	// 用於完全隔離兩個用戶
	IsMutuallyBlocked(ctx context.Context, user1ID, user2ID uint) (bool, error)

	// Delete 移除封鎖記錄
	// 用於解除封鎖功能
	Delete(ctx context.Context, id uint) error

	// DeleteByUsers 根據用戶關係移除封鎖
	// 用於直接解除兩個用戶間的封鎖關係
	DeleteByUsers(ctx context.Context, userID, blockedUserID uint) error

	// GetBlockedUsers 獲取被封鎖的用戶列表
	// 用於封鎖管理和推薦系統過濾
	GetBlockedUsers(ctx context.Context, userID uint) ([]uint, error)

	// GetBlockingUsers 獲取封鎖該用戶的用戶列表
	// 用於系統分析和反向查詢
	GetBlockingUsers(ctx context.Context, userID uint) ([]uint, error)
}

// ModerationRepository 內容審核數據儲存庫介面
// 提供內容審核功能的持久化操作，自動和人工審核結合
type ModerationRepository interface {
	// CreateModerationLog 創建審核日誌
	// 用於記錄自動或人工審核結果
	CreateModerationLog(ctx context.Context, log *ModerationLog) error

	// GetModerationHistory 獲取內容審核歷史
	// 用於查看特定內容的審核記錄
	GetModerationHistory(ctx context.Context, contentType string, contentID uint) ([]*ModerationLog, error)

	// GetUserModerationHistory 獲取用戶審核歷史
	// 用於分析用戶行為模式
	GetUserModerationHistory(ctx context.Context, userID uint, limit int) ([]*ModerationLog, error)

	// UpdateModerationAction 更新審核動作結果
	// 用於記錄審核後採取的行動（警告、暫停等）
	UpdateModerationAction(ctx context.Context, logID uint, action ModerationAction, notes string) error

	// GetPendingModerations 獲取待審核內容
	// 用於人工審核隊列
	GetPendingModerations(ctx context.Context, contentType string, limit int) ([]*ModerationLog, error)

	// GetModerationStats 獲取審核統計數據
	// 用於系統監控和審核效率分析
	GetModerationStats(ctx context.Context, params ModerationStatsParams) (*ModerationStats, error)
}

// ReportStatsParams 檢舉統計查詢參數
type ReportStatsParams struct {
	// 時間範圍
	StartDate *time.Time
	EndDate   *time.Time

	// 篩選條件
	Category     *entity.ReportCategory // 檢舉類別
	Status       *entity.ReportStatus   // 處理狀態
	TargetUserID *uint                  // 特定被檢舉用戶
	ReviewerID   *uint                  // 特定審核員

	// 統計維度
	GroupByCategory bool // 按類別分組
	GroupByStatus   bool // 按狀態分組
	GroupByDate     bool // 按日期分組
}

// ReportStats 檢舉統計資料
type ReportStats struct {
	TotalReports     int                           `json:"total_reports"`           // 總檢舉數
	PendingReports   int                           `json:"pending_reports"`         // 待處理檢舉
	ProcessedReports int                           `json:"processed_reports"`       // 已處理檢舉
	ApprovedReports  int                           `json:"approved_reports"`        // 審核通過檢舉
	RejectedReports  int                           `json:"rejected_reports"`        // 審核拒絕檢舉
	CategoryStats    map[entity.ReportCategory]int `json:"category_stats"`          // 各類別統計
	StatusStats      map[entity.ReportStatus]int   `json:"status_stats"`            // 各狀態統計
	DailyStats       map[string]int                `json:"daily_stats"`             // 每日統計
	TopReportedUsers []UserReportSummary           `json:"top_reported_users"`      // 被檢舉最多的用戶
	ProcessingTime   float64                       `json:"average_processing_time"` // 平均處理時間（小時）
}

// UserReportSummary 用戶檢舉摘要
type UserReportSummary struct {
	UserID      uint         `json:"user_id"`      // 用戶 ID
	ReportCount int          `json:"report_count"` // 被檢舉次數
	User        *entity.User `json:"user"`         // 用戶資訊
}

// ModerationLog 審核日誌
type ModerationLog struct {
	ID          uint      `json:"id"`
	ContentType string    `json:"content_type"` // "profile", "photo", "message", "report"
	ContentID   uint      `json:"content_id"`   // 內容 ID
	UserID      uint      `json:"user_id"`      // 內容所屬用戶
	ModeratorID *uint     `json:"moderator_id"` // 審核員 ID（自動審核時為 null）
	Action      string    `json:"action"`       // "approved", "rejected", "flagged", "deleted"
	Reason      string    `json:"reason"`       // 審核原因
	IsAutomatic bool      `json:"is_automatic"` // 是否自動審核
	Confidence  *float64  `json:"confidence"`   // 自動審核信心度（0-1）
	Notes       string    `json:"notes"`        // 審核備註
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ModerationAction 審核動作枚舉
type ModerationAction string

const (
	ModerationActionApproved ModerationAction = "approved" // 審核通過
	ModerationActionRejected ModerationAction = "rejected" // 審核拒絕
	ModerationActionFlagged  ModerationAction = "flagged"  // 標記問題
	ModerationActionDeleted  ModerationAction = "deleted"  // 刪除內容
	ModerationActionWarning  ModerationAction = "warning"  // 發出警告
	ModerationActionSuspend  ModerationAction = "suspend"  // 暫停帳戶
	ModerationActionBan      ModerationAction = "ban"      // 禁用帳戶
)

// ModerationStatsParams 審核統計查詢參數
type ModerationStatsParams struct {
	StartDate   *time.Time
	EndDate     *time.Time
	ContentType *string // 內容類型篩選
	ModeratorID *uint   // 特定審核員
	IsAutomatic *bool   // 是否自動審核
}

// ModerationStats 審核統計資料
type ModerationStats struct {
	TotalModerations   int                      `json:"total_moderations"`    // 總審核數
	AutoModerations    int                      `json:"auto_moderations"`     // 自動審核數
	ManualModerations  int                      `json:"manual_moderations"`   // 人工審核數
	ApprovedCount      int                      `json:"approved_count"`       // 通過數
	RejectedCount      int                      `json:"rejected_count"`       // 拒絕數
	ActionStats        map[ModerationAction]int `json:"action_stats"`         // 各動作統計
	ContentTypeStats   map[string]int           `json:"content_type_stats"`   // 各內容類型統計
	AverageProcessTime float64                  `json:"average_process_time"` // 平均處理時間
	AutoAccuracyRate   float64                  `json:"auto_accuracy_rate"`   // 自動審核準確率
}
