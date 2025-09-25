package mysql

import (
	"context"
	"fmt"
	"time"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"

	"gorm.io/gorm"
)

// MySQLReportRepository MySQL 檢舉儲存庫實作
type MySQLReportRepository struct {
	db *gorm.DB
}

// NewReportRepository 創建新的 MySQL 檢舉儲存庫
func NewReportRepository(db *gorm.DB) repository.ReportRepository {
	return &MySQLReportRepository{db: db}
}

// Create 創建檢舉記錄
func (r *MySQLReportRepository) Create(ctx context.Context, report *entity.Report) error {
	if err := r.db.WithContext(ctx).Create(report).Error; err != nil {
		return fmt.Errorf("創建檢舉記錄失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 獲取檢舉記錄
func (r *MySQLReportRepository) GetByID(ctx context.Context, id uint) (*entity.Report, error) {
	var report entity.Report
	if err := r.db.WithContext(ctx).First(&report, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("檢舉記錄不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢檢舉記錄失敗: %w", err)
	}
	return &report, nil
}

// GetByUserID 獲取用戶的檢舉記錄
func (r *MySQLReportRepository) GetByUserID(ctx context.Context, userID uint, limit int) ([]*entity.Report, error) {
	query := r.db.WithContext(ctx).Where("reporter_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	var reports []*entity.Report
	if err := query.Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("獲取用戶檢舉記錄失敗: %w", err)
	}
	return reports, nil
}

// GetByTargetUserID 獲取針對特定用戶的檢舉記錄
func (r *MySQLReportRepository) GetByTargetUserID(ctx context.Context, targetUserID uint, limit int) ([]*entity.Report, error) {
	query := r.db.WithContext(ctx).Where("reported_user_id = ?", targetUserID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	var reports []*entity.Report
	if err := query.Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("獲取目標用戶檢舉記錄失敗: %w", err)
	}
	return reports, nil
}

// Update 更新檢舉記錄
func (r *MySQLReportRepository) Update(ctx context.Context, report *entity.Report) error {
	if err := r.db.WithContext(ctx).Save(report).Error; err != nil {
		return fmt.Errorf("更新檢舉記錄失敗: %w", err)
	}
	return nil
}

// GetPendingReports 獲取待審核檢舉記錄
func (r *MySQLReportRepository) GetPendingReports(ctx context.Context, limit int) ([]*entity.Report, error) {
	query := r.db.WithContext(ctx).Where("status = ?", entity.ReportStatusPending).Order("created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	var reports []*entity.Report
	if err := query.Find(&reports).Error; err != nil {
		return nil, fmt.Errorf("獲取待審核檢舉記錄失敗: %w", err)
	}
	return reports, nil
}

// SetReportStatus 設定檢舉狀態
func (r *MySQLReportRepository) SetReportStatus(ctx context.Context, reportID uint, status entity.ReportStatus, reviewerID *uint, reviewNotes string) error {
	updates := map[string]interface{}{
		"status":       status,
		"review_notes": reviewNotes,
		"reviewed_at":  time.Now(),
	}

	if reviewerID != nil {
		updates["reviewer_id"] = *reviewerID
	}

	if err := r.db.WithContext(ctx).Model(&entity.Report{}).Where("id = ?", reportID).Updates(updates).Error; err != nil {
		return fmt.Errorf("設定檢舉狀態失敗: %w", err)
	}
	return nil
}

// CheckDuplicateReport 檢查重複檢舉
func (r *MySQLReportRepository) CheckDuplicateReport(ctx context.Context, userID, targetUserID uint, category entity.ReportCategory) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Report{}).
		Where("reporter_id = ? AND reported_user_id = ? AND category = ?", userID, targetUserID, category).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("檢查重複檢舉失敗: %w", err)
	}
	return count > 0, nil
}

// GetReportStats 獲取檢舉統計數據
func (r *MySQLReportRepository) GetReportStats(ctx context.Context, params repository.ReportStatsParams) (*repository.ReportStats, error) {
	stats := &repository.ReportStats{
		CategoryStats: make(map[entity.ReportCategory]int),
		StatusStats:   make(map[entity.ReportStatus]int),
		DailyStats:    make(map[string]int),
	}

	query := r.db.WithContext(ctx).Model(&entity.Report{})

	// 時間範圍篩選
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	// 其他篩選條件
	if params.Category != nil {
		query = query.Where("category = ?", *params.Category)
	}

	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}

	if params.TargetUserID != nil {
		query = query.Where("reported_user_id = ?", *params.TargetUserID)
	}

	if params.ReviewerID != nil {
		query = query.Where("reviewer_id = ?", *params.ReviewerID)
	}

	// 總檢舉數
	var totalReports int64
	if err := query.Count(&totalReports).Error; err != nil {
		return nil, fmt.Errorf("獲取總檢舉數失敗: %w", err)
	}
	stats.TotalReports = int(totalReports)

	// 各狀態統計
	var statusStats []struct {
		Status entity.ReportStatus
		Count  int64
	}

	if err := r.db.WithContext(ctx).Model(&entity.Report{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats).Error; err != nil {
		return nil, fmt.Errorf("獲取狀態統計失敗: %w", err)
	}

	for _, stat := range statusStats {
		stats.StatusStats[stat.Status] = int(stat.Count)
		switch stat.Status {
		case entity.ReportStatusPending:
			stats.PendingReports = int(stat.Count)
		case entity.ReportStatusApproved:
			stats.ApprovedReports = int(stat.Count)
		case entity.ReportStatusRejected:
			stats.RejectedReports = int(stat.Count)
		default:
			stats.ProcessedReports += int(stat.Count)
		}
	}

	// 各類別統計
	var categoryStats []struct {
		Category entity.ReportCategory
		Count    int64
	}

	if err := r.db.WithContext(ctx).Model(&entity.Report{}).
		Select("category, COUNT(*) as count").
		Group("category").
		Scan(&categoryStats).Error; err != nil {
		return nil, fmt.Errorf("獲取類別統計失敗: %w", err)
	}

	for _, stat := range categoryStats {
		stats.CategoryStats[stat.Category] = int(stat.Count)
	}

	// 被檢舉最多的用戶
	var topReported []repository.UserReportSummary
	if err := r.db.WithContext(ctx).
		Table("reports").
		Select("reported_user_id as user_id, COUNT(*) as report_count").
		Group("reported_user_id").
		Order("report_count DESC").
		Limit(10).
		Scan(&topReported).Error; err != nil {
		return nil, fmt.Errorf("獲取被檢舉最多用戶失敗: %w", err)
	}

	// 獲取用戶資訊
	for i := range topReported {
		var user entity.User
		if err := r.db.WithContext(ctx).First(&user, topReported[i].UserID).Error; err == nil {
			topReported[i].User = &user
		}
	}
	stats.TopReportedUsers = topReported

	return stats, nil
}

// Delete 軟刪除檢舉記錄
func (r *MySQLReportRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Report{}, id).Error; err != nil {
		return fmt.Errorf("刪除檢舉記錄失敗: %w", err)
	}
	return nil
}

// MySQLBlockRepository MySQL 封鎖儲存庫實作
type MySQLBlockRepository struct {
	db *gorm.DB
}

// NewBlockRepository 創建新的 MySQL 封鎖儲存庫
func NewBlockRepository(db *gorm.DB) repository.BlockRepository {
	return &MySQLBlockRepository{db: db}
}

// Create 創建封鎖記錄
func (r *MySQLBlockRepository) Create(ctx context.Context, block *entity.Block) error {
	if err := r.db.WithContext(ctx).Create(block).Error; err != nil {
		return fmt.Errorf("創建封鎖記錄失敗: %w", err)
	}
	return nil
}

// GetByID 根據 ID 獲取封鎖記錄
func (r *MySQLBlockRepository) GetByID(ctx context.Context, id uint) (*entity.Block, error) {
	var block entity.Block
	if err := r.db.WithContext(ctx).First(&block, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("封鎖記錄不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢封鎖記錄失敗: %w", err)
	}
	return &block, nil
}

// GetByUserID 獲取用戶的封鎖列表
func (r *MySQLBlockRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Block, error) {
	var blocks []*entity.Block
	if err := r.db.WithContext(ctx).Where("blocker_id = ?", userID).Order("created_at DESC").Find(&blocks).Error; err != nil {
		return nil, fmt.Errorf("獲取用戶封鎖列表失敗: %w", err)
	}
	return blocks, nil
}

// IsBlocked 檢查是否封鎖關係
func (r *MySQLBlockRepository) IsBlocked(ctx context.Context, userID, targetUserID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Block{}).
		Where("blocker_id = ? AND blocked_id = ?", userID, targetUserID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("檢查封鎖關係失敗: %w", err)
	}
	return count > 0, nil
}

// IsMutuallyBlocked 檢查雙向封鎖關係
func (r *MySQLBlockRepository) IsMutuallyBlocked(ctx context.Context, user1ID, user2ID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Block{}).
		Where("(blocker_id = ? AND blocked_id = ?) OR (blocker_id = ? AND blocked_id = ?)",
			user1ID, user2ID, user2ID, user1ID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("檢查雙向封鎖關係失敗: %w", err)
	}
	return count >= 2, nil
}

// Delete 移除封鎖記錄
func (r *MySQLBlockRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Block{}, id).Error; err != nil {
		return fmt.Errorf("移除封鎖記錄失敗: %w", err)
	}
	return nil
}

// DeleteByUsers 根據用戶關係移除封鎖
func (r *MySQLBlockRepository) DeleteByUsers(ctx context.Context, userID, blockedUserID uint) error {
	if err := r.db.WithContext(ctx).Where("blocker_id = ? AND blocked_id = ?", userID, blockedUserID).Delete(&entity.Block{}).Error; err != nil {
		return fmt.Errorf("根據用戶關係移除封鎖失敗: %w", err)
	}
	return nil
}

// GetBlockedUsers 獲取被封鎖的用戶列表
func (r *MySQLBlockRepository) GetBlockedUsers(ctx context.Context, userID uint) ([]uint, error) {
	var blockedUserIDs []uint
	if err := r.db.WithContext(ctx).Model(&entity.Block{}).
		Where("blocker_id = ?", userID).
		Pluck("blocked_id", &blockedUserIDs).Error; err != nil {
		return nil, fmt.Errorf("獲取被封鎖用戶列表失敗: %w", err)
	}
	return blockedUserIDs, nil
}

// GetBlockingUsers 獲取封鎖該用戶的用戶列表
func (r *MySQLBlockRepository) GetBlockingUsers(ctx context.Context, userID uint) ([]uint, error) {
	var blockingUserIDs []uint
	if err := r.db.WithContext(ctx).Model(&entity.Block{}).
		Where("blocked_id = ?", userID).
		Pluck("blocker_id", &blockingUserIDs).Error; err != nil {
		return nil, fmt.Errorf("獲取封鎖該用戶的用戶列表失敗: %w", err)
	}
	return blockingUserIDs, nil
}

// MySQLModerationRepository MySQL 內容審核儲存庫實作
type MySQLModerationRepository struct {
	db *gorm.DB
}

// NewModerationRepository 創建新的 MySQL 內容審核儲存庫
func NewModerationRepository(db *gorm.DB) repository.ModerationRepository {
	return &MySQLModerationRepository{db: db}
}

// CreateModerationLog 創建審核日誌
func (r *MySQLModerationRepository) CreateModerationLog(ctx context.Context, log *repository.ModerationLog) error {
	if err := r.db.WithContext(ctx).Table("moderation_logs").Create(log).Error; err != nil {
		return fmt.Errorf("創建審核日誌失敗: %w", err)
	}
	return nil
}

// GetModerationHistory 獲取內容審核歷史
func (r *MySQLModerationRepository) GetModerationHistory(ctx context.Context, contentType string, contentID uint) ([]*repository.ModerationLog, error) {
	var logs []*repository.ModerationLog
	if err := r.db.WithContext(ctx).Table("moderation_logs").
		Where("content_type = ? AND content_id = ?", contentType, contentID).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("獲取內容審核歷史失敗: %w", err)
	}
	return logs, nil
}

// GetUserModerationHistory 獲取用戶審核歷史
func (r *MySQLModerationRepository) GetUserModerationHistory(ctx context.Context, userID uint, limit int) ([]*repository.ModerationLog, error) {
	query := r.db.WithContext(ctx).Table("moderation_logs").Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	var logs []*repository.ModerationLog
	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("獲取用戶審核歷史失敗: %w", err)
	}
	return logs, nil
}

// UpdateModerationAction 更新審核動作結果
func (r *MySQLModerationRepository) UpdateModerationAction(ctx context.Context, logID uint, action repository.ModerationAction, notes string) error {
	updates := map[string]interface{}{
		"action":     action,
		"notes":      notes,
		"updated_at": time.Now(),
	}

	if err := r.db.WithContext(ctx).Table("moderation_logs").Where("id = ?", logID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新審核動作失敗: %w", err)
	}
	return nil
}

// GetPendingModerations 獲取待審核內容
func (r *MySQLModerationRepository) GetPendingModerations(ctx context.Context, contentType string, limit int) ([]*repository.ModerationLog, error) {
	query := r.db.WithContext(ctx).Table("moderation_logs").
		Where("action = ?", repository.ModerationActionFlagged).
		Order("created_at ASC")

	if contentType != "" {
		query = query.Where("content_type = ?", contentType)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	var logs []*repository.ModerationLog
	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("獲取待審核內容失敗: %w", err)
	}
	return logs, nil
}

// GetModerationStats 獲取審核統計數據
func (r *MySQLModerationRepository) GetModerationStats(ctx context.Context, params repository.ModerationStatsParams) (*repository.ModerationStats, error) {
	stats := &repository.ModerationStats{
		ActionStats:      make(map[repository.ModerationAction]int),
		ContentTypeStats: make(map[string]int),
	}

	query := r.db.WithContext(ctx).Table("moderation_logs")

	// 時間範圍篩選
	if params.StartDate != nil {
		query = query.Where("created_at >= ?", *params.StartDate)
	}

	if params.EndDate != nil {
		query = query.Where("created_at <= ?", *params.EndDate)
	}

	// 其他篩選條件
	if params.ContentType != nil {
		query = query.Where("content_type = ?", *params.ContentType)
	}

	if params.ModeratorID != nil {
		query = query.Where("moderator_id = ?", *params.ModeratorID)
	}

	if params.IsAutomatic != nil {
		query = query.Where("is_automatic = ?", *params.IsAutomatic)
	}

	// 總審核數
	var totalModerations int64
	if err := query.Count(&totalModerations).Error; err != nil {
		return nil, fmt.Errorf("獲取總審核數失敗: %w", err)
	}
	stats.TotalModerations = int(totalModerations)

	// 自動/人工審核統計
	var autoCount int64
	if err := r.db.WithContext(ctx).Table("moderation_logs").Where("is_automatic = ?", true).Count(&autoCount).Error; err != nil {
		return nil, fmt.Errorf("獲取自動審核數失敗: %w", err)
	}
	stats.AutoModerations = int(autoCount)
	stats.ManualModerations = stats.TotalModerations - stats.AutoModerations

	// 各動作統計
	var actionStats []struct {
		Action repository.ModerationAction
		Count  int64
	}

	if err := r.db.WithContext(ctx).Table("moderation_logs").
		Select("action, COUNT(*) as count").
		Group("action").
		Scan(&actionStats).Error; err != nil {
		return nil, fmt.Errorf("獲取動作統計失敗: %w", err)
	}

	for _, stat := range actionStats {
		stats.ActionStats[stat.Action] = int(stat.Count)
		if stat.Action == repository.ModerationActionApproved {
			stats.ApprovedCount = int(stat.Count)
		} else if stat.Action == repository.ModerationActionRejected {
			stats.RejectedCount = int(stat.Count)
		}
	}

	// 各內容類型統計
	var contentTypeStats []struct {
		ContentType string
		Count       int64
	}

	if err := r.db.WithContext(ctx).Table("moderation_logs").
		Select("content_type, COUNT(*) as count").
		Group("content_type").
		Scan(&contentTypeStats).Error; err != nil {
		return nil, fmt.Errorf("獲取內容類型統計失敗: %w", err)
	}

	for _, stat := range contentTypeStats {
		stats.ContentTypeStats[stat.ContentType] = int(stat.Count)
	}

	return stats, nil
}
