package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

// ReportService 檢舉業務邏輯服務
// 負責檢舉處理、內容審核、用戶封鎖等核心業務邏輯
type ReportService struct {
	reportRepo     repository.ReportRepository
	blockRepo      repository.BlockRepository
	moderationRepo repository.ModerationRepository
	userRepo       repository.UserRepository
	matchRepo      repository.MatchRepository
}

// NewReportService 創建新的檢舉服務實例
func NewReportService(
	reportRepo repository.ReportRepository,
	blockRepo repository.BlockRepository,
	moderationRepo repository.ModerationRepository,
	userRepo repository.UserRepository,
	matchRepo repository.MatchRepository,
) *ReportService {
	return &ReportService{
		reportRepo:     reportRepo,
		blockRepo:      blockRepo,
		moderationRepo: moderationRepo,
		userRepo:       userRepo,
		matchRepo:      matchRepo,
	}
}

// SubmitReportRequest 提交檢舉請求
type SubmitReportRequest struct {
	ReporterID     uint                  `json:"reporter_id" validate:"required"`
	ReportedUserID uint                  `json:"reported_user_id" validate:"required"`
	Category       entity.ReportCategory `json:"category" validate:"required"`
	Description    string                `json:"description" validate:"required"`
	Evidence       []string              `json:"evidence,omitempty"`
}

// BlockUserRequest 封鎖用戶請求
type BlockUserRequest struct {
	UserID        uint   `json:"user_id" validate:"required"`
	BlockedUserID uint   `json:"blocked_user_id" validate:"required"`
	Reason        string `json:"reason,omitempty"`
}

// ReportResponse 檢舉回應
type ReportResponse struct {
	Report  *entity.Report `json:"report,omitempty"`
	Success bool           `json:"success"`
	Message string         `json:"message"`
}

// ReviewReportRequest 審核檢舉請求
type ReviewReportRequest struct {
	ReportID    uint                `json:"report_id" validate:"required"`
	ReviewerID  uint                `json:"reviewer_id" validate:"required"`
	Status      entity.ReportStatus `json:"status" validate:"required"`
	ReviewNotes string              `json:"review_notes"`
	Action      *ModerationAction   `json:"action,omitempty"`
}

// ModerationAction 審核行動
type ModerationAction struct {
	Type        string `json:"type"`        // warning, suspend, ban
	Duration    *int   `json:"duration"`    // 暫停天數（如適用）
	Description string `json:"description"` // 行動描述
}

// SubmitReport 提交檢舉
// 處理用戶檢舉其他用戶的不當行為
func (s *ReportService) SubmitReport(ctx context.Context, req *SubmitReportRequest) (*ReportResponse, error) {
	// 驗證請求資料
	if err := s.validateSubmitReportRequest(req); err != nil {
		return &ReportResponse{
			Success: false,
			Message: fmt.Sprintf("檢舉資料驗證失敗: %v", err),
		}, nil
	}

	// 檢查不能檢舉自己
	if req.ReporterID == req.ReportedUserID {
		return &ReportResponse{
			Success: false,
			Message: "不能檢舉自己",
		}, nil
	}

	// 檢查檢舉者是否存在且啟用
	reporter, err := s.userRepo.GetByID(ctx, req.ReporterID)
	if err != nil {
		return &ReportResponse{
			Success: false,
			Message: "檢舉者不存在",
		}, nil
	}

	if !reporter.IsActive {
		return &ReportResponse{
			Success: false,
			Message: "檢舉者帳戶未啟用",
		}, nil
	}

	// 檢查被檢舉用戶是否存在
	_, err = s.userRepo.GetByID(ctx, req.ReportedUserID)
	if err != nil {
		return &ReportResponse{
			Success: false,
			Message: "被檢舉用戶不存在",
		}, nil
	}

	// 檢查重複檢舉
	isDuplicate, err := s.reportRepo.CheckDuplicateReport(ctx, req.ReporterID, req.ReportedUserID, req.Category)
	if err != nil {
		return nil, fmt.Errorf("檢查重複檢舉失敗: %w", err)
	}

	if isDuplicate {
		return &ReportResponse{
			Success: false,
			Message: "您已經對該用戶提交過相同類別的檢舉",
		}, nil
	}

	// 創建檢舉記錄
	report := &entity.Report{
		ReporterID:  req.ReporterID,
		ReportedID:  req.ReportedUserID,
		Category:    req.Category,
		Description: req.Description,
		Status:      entity.ReportStatusPending,
	}

	// 處理證據（如果有）
	// Note: Report entity 沒有 Evidence 字段，可以將證據資訊包含在 Description 中
	if len(req.Evidence) > 0 {
		evidenceStr := strings.Join(req.Evidence, ", ")
		report.Description = fmt.Sprintf("%s\n證據: %s", req.Description, evidenceStr)
	}

	// 保存檢舉記錄
	if err := s.reportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("創建檢舉記錄失敗: %w", err)
	}

	// 創建審核日誌（自動審核預處理）
	moderationLog := &repository.ModerationLog{
		ContentType: "user",
		ContentID:   req.ReportedUserID,
		UserID:      req.ReportedUserID,
		Action:      "report_submitted",
		IsAutomatic: true,
		Confidence:  nil, // 人工檢舉，無AI信心度
		Notes:       fmt.Sprintf("用戶檢舉: %s", req.Description),
	}

	if err := s.moderationRepo.CreateModerationLog(ctx, moderationLog); err != nil {
		// 記錄錯誤但不影響檢舉提交
	}

	return &ReportResponse{
		Report:  report,
		Success: true,
		Message: "檢舉已提交，我們會儘快處理",
	}, nil
}

// BlockUser 封鎖用戶
// 允許用戶封鎖其他用戶以避免互動
func (s *ReportService) BlockUser(ctx context.Context, req *BlockUserRequest) error {
	// 驗證請求資料
	if err := s.validateBlockUserRequest(req); err != nil {
		return fmt.Errorf("封鎖請求驗證失敗: %w", err)
	}

	// 檢查不能封鎖自己
	if req.UserID == req.BlockedUserID {
		return errors.New("不能封鎖自己")
	}

	// 檢查用戶是否存在且啟用
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return errors.New("用戶不存在")
	}

	if !user.IsActive {
		return errors.New("用戶帳戶未啟用")
	}

	// 檢查被封鎖用戶是否存在
	if _, err := s.userRepo.GetByID(ctx, req.BlockedUserID); err != nil {
		return errors.New("被封鎖用戶不存在")
	}

	// 檢查是否已經封鎖
	isBlocked, err := s.blockRepo.IsBlocked(ctx, req.UserID, req.BlockedUserID)
	if err != nil {
		return fmt.Errorf("檢查封鎖狀態失敗: %w", err)
	}

	if isBlocked {
		return errors.New("用戶已被封鎖")
	}

	// 創建封鎖記錄
	block := &entity.Block{
		BlockerID: req.UserID,
		BlockedID: req.BlockedUserID,
		Reason:    entity.BlockReasonOther, // 預設使用 "其他"，或者根據 req.Reason 轉換
	}

	// 如果有提供原因，設置為備註
	if strings.TrimSpace(req.Reason) != "" {
		block.Notes = &req.Reason
	}

	if err := s.blockRepo.Create(ctx, block); err != nil {
		return fmt.Errorf("創建封鎖記錄失敗: %w", err)
	}

	// 如果有相關配對，更新配對狀態
	match, err := s.matchRepo.GetMatch(ctx, req.UserID, req.BlockedUserID)
	if err == nil && match != nil && match.Status == entity.MatchStatusMatched {
		// 將配對狀態設為未配對
		s.matchRepo.UpdateMatchStatus(ctx, match.ID, entity.MatchStatusUnmatched)
	}

	return nil
}

// UnblockUser 解除封鎖
// 允許用戶解除對其他用戶的封鎖
func (s *ReportService) UnblockUser(ctx context.Context, userID, blockedUserID uint) error {
	// 檢查用戶是否存在且啟用
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("用戶不存在")
	}

	if !user.IsActive {
		return errors.New("用戶帳戶未啟用")
	}

	// 檢查封鎖關係是否存在
	isBlocked, err := s.blockRepo.IsBlocked(ctx, userID, blockedUserID)
	if err != nil {
		return fmt.Errorf("檢查封鎖狀態失敗: %w", err)
	}

	if !isBlocked {
		return errors.New("封鎖關係不存在")
	}

	// 移除封鎖記錄
	return s.blockRepo.DeleteByUsers(ctx, userID, blockedUserID)
}

// GetUserBlockList 獲取用戶封鎖列表
// 返回用戶封鎖的所有用戶
func (s *ReportService) GetUserBlockList(ctx context.Context, userID uint) ([]*entity.Block, error) {
	// 檢查用戶是否存在且啟用
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用戶不存在")
	}

	if !user.IsActive {
		return nil, errors.New("用戶帳戶未啟用")
	}

	return s.blockRepo.GetByUserID(ctx, userID)
}

// IsUserBlocked 檢查用戶是否被封鎖
// 用於配對推薦和聊天權限檢查
func (s *ReportService) IsUserBlocked(ctx context.Context, userID, targetUserID uint) (bool, error) {
	return s.blockRepo.IsBlocked(ctx, userID, targetUserID)
}

// ReviewReport 審核檢舉
// 管理員審核檢舉並採取相應行動
func (s *ReportService) ReviewReport(ctx context.Context, req *ReviewReportRequest) error {
	// 驗證請求資料
	if err := s.validateReviewReportRequest(req); err != nil {
		return fmt.Errorf("審核請求驗證失敗: %w", err)
	}

	// 檢查審核者權限（簡化處理，實際應用中需要更複雜的權限檢查）
	reviewer, err := s.userRepo.GetByID(ctx, req.ReviewerID)
	if err != nil {
		return errors.New("審核者不存在")
	}

	if !reviewer.IsActive {
		return errors.New("審核者帳戶未啟用")
	}

	// 獲取檢舉記錄
	report, err := s.reportRepo.GetByID(ctx, req.ReportID)
	if err != nil {
		return errors.New("檢舉記錄不存在")
	}

	// 檢查檢舉狀態（只有待處理或審查中的可以審核）
	if report.Status != entity.ReportStatusPending && report.Status != entity.ReportStatusReviewing {
		return errors.New("該檢舉已被處理，無法再次審核")
	}

	// 更新檢舉狀態
	if err := s.reportRepo.SetReportStatus(ctx, req.ReportID, req.Status, &req.ReviewerID, req.ReviewNotes); err != nil {
		return fmt.Errorf("更新檢舉狀態失敗: %w", err)
	}

	// 如果審核通過且有指定行動，執行相應的行動
	if req.Status == entity.ReportStatusApproved && req.Action != nil {
		if err := s.executeModerationAction(ctx, report.ReportedID, req.Action, req.ReviewerID); err != nil {
			// 記錄錯誤但不回滾審核結果
			// 實際應用中可能需要更複雜的錯誤處理
		}
	}

	// 創建審核日誌
	moderationLog := &repository.ModerationLog{
		ContentType: "report",
		ContentID:   req.ReportID,
		UserID:      report.ReportedID,
		ModeratorID: &req.ReviewerID,
		Action:      string(req.Status),
		IsAutomatic: false,
		Confidence:  nil,
		Notes:       req.ReviewNotes,
	}

	if err := s.moderationRepo.CreateModerationLog(ctx, moderationLog); err != nil {
		// 記錄錯誤但不影響審核流程
	}

	return nil
}

// GetPendingReports 獲取待審核檢舉
// 供管理員審核面板使用
func (s *ReportService) GetPendingReports(ctx context.Context, limit int) ([]*entity.Report, error) {
	if limit <= 0 {
		limit = 20 // 預設返回20條
	}
	if limit > 100 {
		limit = 100 // 最大限制100條
	}

	return s.reportRepo.GetPendingReports(ctx, limit)
}

// GetReportStats 獲取檢舉統計數據
// 用於系統監控和管理分析
func (s *ReportService) GetReportStats(ctx context.Context, params repository.ReportStatsParams) (*repository.ReportStats, error) {
	return s.reportRepo.GetReportStats(ctx, params)
}

// GetUserReports 獲取用戶檢舉記錄
// 返回用戶提交的檢舉記錄
func (s *ReportService) GetUserReports(ctx context.Context, userID uint, limit int) ([]*entity.Report, error) {
	// 檢查用戶是否存在
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用戶不存在")
	}

	if !user.IsActive {
		return nil, errors.New("用戶帳戶未啟用")
	}

	if limit <= 0 {
		limit = 10 // 預設返回10條
	}
	if limit > 50 {
		limit = 50 // 最大限制50條
	}

	return s.reportRepo.GetByUserID(ctx, userID, limit)
}

// GetReportsAgainstUser 獲取針對用戶的檢舉記錄
// 返回針對特定用戶的檢舉記錄
func (s *ReportService) GetReportsAgainstUser(ctx context.Context, targetUserID uint, limit int) ([]*entity.Report, error) {
	if limit <= 0 {
		limit = 10 // 預設返回10條
	}
	if limit > 50 {
		limit = 50 // 最大限制50條
	}

	return s.reportRepo.GetByTargetUserID(ctx, targetUserID, limit)
}

// 私有輔助方法

// validateSubmitReportRequest 驗證提交檢舉請求
func (s *ReportService) validateSubmitReportRequest(req *SubmitReportRequest) error {
	if req.ReporterID == 0 {
		return errors.New("檢舉者ID不能為空")
	}

	if req.ReportedUserID == 0 {
		return errors.New("被檢舉用戶ID不能為空")
	}

	if !req.Category.IsValid() {
		return errors.New("無效的檢舉類別")
	}

	if strings.TrimSpace(req.Description) == "" {
		return errors.New("檢舉描述不能為空")
	}

	if len(req.Description) > 1000 {
		return errors.New("檢舉描述不能超過1000個字符")
	}

	return nil
}

// validateBlockUserRequest 驗證封鎖用戶請求
func (s *ReportService) validateBlockUserRequest(req *BlockUserRequest) error {
	if req.UserID == 0 {
		return errors.New("用戶ID不能為空")
	}

	if req.BlockedUserID == 0 {
		return errors.New("被封鎖用戶ID不能為空")
	}

	if len(req.Reason) > 500 {
		return errors.New("封鎖理由不能超過500個字符")
	}

	return nil
}

// validateReviewReportRequest 驗證審核檢舉請求
func (s *ReportService) validateReviewReportRequest(req *ReviewReportRequest) error {
	if req.ReportID == 0 {
		return errors.New("檢舉ID不能為空")
	}

	if req.ReviewerID == 0 {
		return errors.New("審核者ID不能為空")
	}

	if !req.Status.IsValid() {
		return errors.New("無效的審核狀態")
	}

	if len(req.ReviewNotes) > 1000 {
		return errors.New("審核備註不能超過1000個字符")
	}

	return nil
}

// executeModerationAction 執行審核行動
func (s *ReportService) executeModerationAction(ctx context.Context, userID uint, action *ModerationAction, moderatorID uint) error {
	switch action.Type {
	case "warning":
		// 發送警告（可以發送系統訊息或郵件）
		return s.sendWarningToUser(ctx, userID, action.Description)

	case "suspend":
		// 暫停用戶帳戶
		if action.Duration != nil {
			return s.suspendUser(ctx, userID, *action.Duration, action.Description)
		}
		return errors.New("暫停行動必須指定時間")

	case "ban":
		// 永久封禁用戶
		return s.banUser(ctx, userID, action.Description)

	default:
		return fmt.Errorf("不支援的審核行動類型: %s", action.Type)
	}
}

// sendWarningToUser 發送警告給用戶
func (s *ReportService) sendWarningToUser(ctx context.Context, userID uint, message string) error {
	// 簡化實現：更新用戶狀態或發送通知
	// 實際應用中可能需要通知系統

	moderationLog := &repository.ModerationLog{
		ContentType: "user",
		ContentID:   userID,
		UserID:      userID,
		Action:      "warning_sent",
		IsAutomatic: false,
		Notes:       message,
	}

	return s.moderationRepo.CreateModerationLog(ctx, moderationLog)
}

// suspendUser 暫停用戶
func (s *ReportService) suspendUser(ctx context.Context, userID uint, days int, reason string) error {
	// 簡化實現：設置用戶為非啟用狀態
	// 實際應用中需要更複雜的暫停機制

	if err := s.userRepo.SetActive(ctx, userID, false); err != nil {
		return err
	}

	moderationLog := &repository.ModerationLog{
		ContentType: "user",
		ContentID:   userID,
		UserID:      userID,
		Action:      "suspended",
		IsAutomatic: false,
		Notes:       fmt.Sprintf("暫停 %d 天: %s", days, reason),
	}

	return s.moderationRepo.CreateModerationLog(ctx, moderationLog)
}

// banUser 封禁用戶
func (s *ReportService) banUser(ctx context.Context, userID uint, reason string) error {
	// 永久停用用戶帳戶
	if err := s.userRepo.SetActive(ctx, userID, false); err != nil {
		return err
	}

	moderationLog := &repository.ModerationLog{
		ContentType: "user",
		ContentID:   userID,
		UserID:      userID,
		Action:      "banned",
		IsAutomatic: false,
		Notes:       fmt.Sprintf("永久封禁: %s", reason),
	}

	return s.moderationRepo.CreateModerationLog(ctx, moderationLog)
}
