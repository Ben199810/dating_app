package mysql

import (
	"context"
	"fmt"
	"math"
	"time"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"

	"gorm.io/gorm"
)

// MySQLMatchRepository MySQL 配對儲存庫實作
type MySQLMatchRepository struct {
	db *gorm.DB
}

// NewMatchRepository 創建新的 MySQL 配對儲存庫
func NewMatchRepository(db *gorm.DB) repository.MatchRepository {
	return &MySQLMatchRepository{db: db}
}

// CreateSwipe 記錄滑動動作
func (r *MySQLMatchRepository) CreateSwipe(ctx context.Context, match *entity.Match) error {
	if err := r.db.WithContext(ctx).Create(match).Error; err != nil {
		return fmt.Errorf("記錄滑動動作失敗: %w", err)
	}
	return nil
}

// GetMatch 獲取配對記錄
func (r *MySQLMatchRepository) GetMatch(ctx context.Context, user1ID, user2ID uint) (*entity.Match, error) {
	var match entity.Match

	// 查找任一方向的配對記錄
	if err := r.db.WithContext(ctx).Where(
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		user1ID, user2ID, user2ID, user1ID,
	).First(&match).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("配對記錄不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢配對記錄失敗: %w", err)
	}
	return &match, nil
}

// GetMatchByID 根據 ID 獲取配對記錄
func (r *MySQLMatchRepository) GetMatchByID(ctx context.Context, id uint) (*entity.Match, error) {
	var match entity.Match
	if err := r.db.WithContext(ctx).First(&match, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("配對記錄不存在: %w", err)
		}
		return nil, fmt.Errorf("查詢配對記錄失敗: %w", err)
	}
	return &match, nil
}

// UpdateMatchStatus 更新配對狀態
func (r *MySQLMatchRepository) UpdateMatchStatus(ctx context.Context, matchID uint, status entity.MatchStatus) error {
	updates := map[string]interface{}{
		"status": status,
	}

	// 如果是配對成功，記錄配對時間
	if status == entity.MatchStatusMatched {
		updates["matched_at"] = "NOW()"
	}

	if err := r.db.WithContext(ctx).Model(&entity.Match{}).Where("id = ?", matchID).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新配對狀態失敗: %w", err)
	}
	return nil
}

// ProcessSwipe 處理滑動動作並檢查是否配對成功
func (r *MySQLMatchRepository) ProcessSwipe(ctx context.Context, userID, targetUserID uint, action entity.SwipeAction) (*entity.Match, bool, error) {
	var match *entity.Match
	var isMatched bool

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 檢查是否已存在反向滑動記錄
		var existingMatch entity.Match
		err := tx.Where("user1_id = ? AND user2_id = ?", targetUserID, userID).First(&existingMatch).Error

		if err != nil && err != gorm.ErrRecordNotFound {
			return fmt.Errorf("查詢現有配對記錄失敗: %w", err)
		}

		// 創建新的滑動記錄
		newMatch := &entity.Match{
			User1ID:     userID,
			User2ID:     targetUserID,
			User1Action: action,
			Status:      entity.MatchStatusPending,
		}

		if err := tx.Create(newMatch).Error; err != nil {
			return fmt.Errorf("創建滑動記錄失敗: %w", err)
		}

		match = newMatch

		// 如果存在反向記錄且都是 like，則配對成功
		if err == nil && existingMatch.User1Action == entity.SwipeActionLike && action == entity.SwipeActionLike {
			// 更新雙方記錄為配對成功
			if err := tx.Model(&existingMatch).Updates(map[string]interface{}{
				"user2_action": action,
				"status":       entity.MatchStatusMatched,
				"matched_at":   "NOW()",
			}).Error; err != nil {
				return fmt.Errorf("更新原配對記錄失敗: %w", err)
			}

			if err := tx.Model(newMatch).Updates(map[string]interface{}{
				"status":     entity.MatchStatusMatched,
				"matched_at": "NOW()",
			}).Error; err != nil {
				return fmt.Errorf("更新新配對記錄失敗: %w", err)
			}

			isMatched = true
			match.Status = entity.MatchStatusMatched
		}

		return nil
	})

	if err != nil {
		return nil, false, err
	}

	return match, isMatched, nil
}

// GetUserMatches 獲取用戶的所有配對記錄
func (r *MySQLMatchRepository) GetUserMatches(ctx context.Context, userID uint, status entity.MatchStatus) ([]*entity.Match, error) {
	var matches []*entity.Match
	query := r.db.WithContext(ctx).Where("user1_id = ? OR user2_id = ?", userID, userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at DESC").Find(&matches).Error; err != nil {
		return nil, fmt.Errorf("獲取用戶配對記錄失敗: %w", err)
	}
	return matches, nil
}

// GetMatchedUsers 獲取用戶配對成功的用戶列表
func (r *MySQLMatchRepository) GetMatchedUsers(ctx context.Context, userID uint) ([]*entity.User, error) {
	var users []*entity.User

	if err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*").
		Joins(`
			INNER JOIN matches ON 
			(matches.user1_id = ? AND matches.user2_id = users.id) OR 
			(matches.user2_id = ? AND matches.user1_id = users.id)
		`, userID, userID).
		Where("matches.status = ?", entity.MatchStatusMatched).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("獲取配對用戶列表失敗: %w", err)
	}
	return users, nil
}

// HasUserSwiped 檢查用戶是否已經滑動過目標用戶
func (r *MySQLMatchRepository) HasUserSwiped(ctx context.Context, userID, targetUserID uint) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Match{}).
		Where("user1_id = ? AND user2_id = ?", userID, targetUserID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("檢查滑動記錄失敗: %w", err)
	}
	return count > 0, nil
}

// Delete 刪除配對記錄
func (r *MySQLMatchRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Match{}, id).Error; err != nil {
		return fmt.Errorf("刪除配對記錄失敗: %w", err)
	}
	return nil
}

// MySQLMatchingAlgorithmRepository MySQL 配對演算法儲存庫實作
type MySQLMatchingAlgorithmRepository struct {
	db *gorm.DB
}

// NewMatchingAlgorithmRepository 創建新的 MySQL 配對演算法儲存庫
func NewMatchingAlgorithmRepository(db *gorm.DB) repository.MatchingAlgorithmRepository {
	return &MySQLMatchingAlgorithmRepository{db: db}
}

// GetPotentialMatches 獲取潛在配對對象
func (r *MySQLMatchingAlgorithmRepository) GetPotentialMatches(ctx context.Context, userID uint, params repository.PotentialMatchParams) ([]*entity.User, error) {
	query := r.db.WithContext(ctx).
		Table("users").
		Select("DISTINCT users.*").
		Joins("INNER JOIN user_profiles ON users.id = user_profiles.user_id").
		Where("users.id != ? AND users.is_active = ? AND users.is_verified = ?", userID, true, true)

	// 排除已滑動過的用戶
	if params.ExcludeSwipedUsers {
		query = query.Where(`
			users.id NOT IN (
				SELECT DISTINCT 
					CASE 
						WHEN user1_id = ? THEN user2_id 
						ELSE user1_id 
					END
				FROM matches 
				WHERE user1_id = ? OR user2_id = ?
			)
		`, userID, userID, userID)
	}

	// 排除已封鎖的用戶
	if params.ExcludeBlockedUsers {
		query = query.Where(`
			users.id NOT IN (
				SELECT DISTINCT 
					CASE 
						WHEN blocker_id = ? THEN blocked_id 
						ELSE blocker_id 
					END
				FROM blocks 
				WHERE blocker_id = ? OR blocked_id = ?
			)
		`, userID, userID, userID)
	}

	// 地理位置篩選
	if params.Latitude != nil && params.Longitude != nil && params.MaxDistance != nil {
		query = query.Where(`
			ST_Distance_Sphere(
				POINT(user_profiles.location_lng, user_profiles.location_lat),
				POINT(?, ?)
			) <= ? * 1000
		`, *params.Longitude, *params.Latitude, *params.MaxDistance)
	}

	// 年齡篩選
	if params.MinAge != nil {
		query = query.Where("YEAR(NOW()) - YEAR(users.birth_date) >= ?", *params.MinAge)
	}

	if params.MaxAge != nil {
		query = query.Where("YEAR(NOW()) - YEAR(users.birth_date) <= ?", *params.MaxAge)
	}

	// 性別偏好
	if params.PreferredGender != nil {
		query = query.Where("user_profiles.gender = ?", *params.PreferredGender)
	}

	// 共同興趣篩選
	if params.RequireCommonInterests {
		minCommon := 1
		if params.MinCommonInterests != nil {
			minCommon = *params.MinCommonInterests
		}

		query = query.Having(`
			(
				SELECT COUNT(DISTINCT ui2.interest_id)
				FROM user_interests ui1
				INNER JOIN user_interests ui2 ON ui1.interest_id = ui2.interest_id
				WHERE ui1.user_id = ? AND ui2.user_id = users.id
			) >= ?
		`, userID, minCommon)
	}

	// 分頁
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	var users []*entity.User
	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("獲取潛在配對對象失敗: %w", err)
	}

	return users, nil
}

// GetUsersNearby 獲取附近的用戶
func (r *MySQLMatchingAlgorithmRepository) GetUsersNearby(ctx context.Context, userID uint, lat, lng float64, maxDistanceKm int, limit int) ([]*entity.User, error) {
	var users []*entity.User

	query := r.db.WithContext(ctx).
		Table("users").
		Select("users.*, ST_Distance_Sphere(POINT(user_profiles.location_lng, user_profiles.location_lat), POINT(?, ?)) as distance", lng, lat).
		Joins("INNER JOIN user_profiles ON users.id = user_profiles.user_id").
		Where("users.id != ? AND users.is_active = ? AND users.is_verified = ?", userID, true, true).
		Where("user_profiles.location_lat IS NOT NULL AND user_profiles.location_lng IS NOT NULL").
		Where("ST_Distance_Sphere(POINT(user_profiles.location_lng, user_profiles.location_lat), POINT(?, ?)) <= ?", lng, lat, maxDistanceKm*1000).
		Order("distance").
		Limit(limit)

	if err := query.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("獲取附近用戶失敗: %w", err)
	}

	return users, nil
}

// GetUsersByAgeRange 根據年齡範圍獲取用戶
func (r *MySQLMatchingAlgorithmRepository) GetUsersByAgeRange(ctx context.Context, userID uint, minAge, maxAge int, limit int) ([]*entity.User, error) {
	var users []*entity.User

	if err := r.db.WithContext(ctx).
		Where("id != ? AND is_active = ? AND is_verified = ?", userID, true, true).
		Where("YEAR(NOW()) - YEAR(birth_date) BETWEEN ? AND ?", minAge, maxAge).
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("根據年齡獲取用戶失敗: %w", err)
	}

	return users, nil
}

// GetUsersByCommonInterests 根據共同興趣獲取用戶
func (r *MySQLMatchingAlgorithmRepository) GetUsersByCommonInterests(ctx context.Context, userID uint, limit int) ([]*entity.User, error) {
	var users []*entity.User

	if err := r.db.WithContext(ctx).
		Table("users").
		Select("users.*, COUNT(ui2.interest_id) as common_interests").
		Joins(`
			INNER JOIN user_interests ui2 ON users.id = ui2.user_id
			INNER JOIN user_interests ui1 ON ui1.interest_id = ui2.interest_id AND ui1.user_id = ?
		`, userID).
		Where("users.id != ? AND users.is_active = ? AND users.is_verified = ?", userID, true, true).
		Group("users.id").
		Order("common_interests DESC").
		Limit(limit).
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("根據共同興趣獲取用戶失敗: %w", err)
	}

	return users, nil
}

// CalculateCompatibilityScore 計算用戶相容性分數
func (r *MySQLMatchingAlgorithmRepository) CalculateCompatibilityScore(ctx context.Context, user1ID, user2ID uint) (float64, error) {
	// 獲取用戶資料
	var user1, user2 entity.User
	var profile1, profile2 entity.UserProfile

	if err := r.db.WithContext(ctx).First(&user1, user1ID).Error; err != nil {
		return 0, fmt.Errorf("獲取用戶1資料失敗: %w", err)
	}

	if err := r.db.WithContext(ctx).First(&user2, user2ID).Error; err != nil {
		return 0, fmt.Errorf("獲取用戶2資料失敗: %w", err)
	}

	if err := r.db.WithContext(ctx).Where("user_id = ?", user1ID).First(&profile1).Error; err != nil {
		return 0, fmt.Errorf("獲取用戶1檔案失敗: %w", err)
	}

	if err := r.db.WithContext(ctx).Where("user_id = ?", user2ID).First(&profile2).Error; err != nil {
		return 0, fmt.Errorf("獲取用戶2檔案失敗: %w", err)
	}

	var score float64

	// 年齡相容性 (40%)
	age1 := calculateAge(user1.BirthDate)
	age2 := calculateAge(user2.BirthDate)
	ageDiff := math.Abs(float64(age1 - age2))
	ageScore := math.Max(0, 1.0-ageDiff/20.0) // 年齡差20歲以內得分
	score += ageScore * 0.4

	// 地理位置相容性 (30%)
	if profile1.LocationLat != nil && profile1.LocationLng != nil &&
		profile2.LocationLat != nil && profile2.LocationLng != nil {
		distance := calculateDistance(*profile1.LocationLat, *profile1.LocationLng,
			*profile2.LocationLat, *profile2.LocationLng)
		distanceScore := math.Max(0, 1.0-distance/100.0) // 100km以內得分
		score += distanceScore * 0.3
	}

	// 共同興趣 (30%)
	var commonInterests int64
	r.db.WithContext(ctx).
		Table("user_interests ui1").
		Joins("INNER JOIN user_interests ui2 ON ui1.interest_id = ui2.interest_id").
		Where("ui1.user_id = ? AND ui2.user_id = ?", user1ID, user2ID).
		Count(&commonInterests)

	interestScore := math.Min(1.0, float64(commonInterests)/5.0) // 最多5個共同興趣
	score += interestScore * 0.3

	return math.Min(1.0, score), nil
}

// GetMatchingStats 獲取配對統計數據
func (r *MySQLMatchingAlgorithmRepository) GetMatchingStats(ctx context.Context, userID uint) (*repository.MatchingStats, error) {
	stats := &repository.MatchingStats{}

	// 總滑動次數
	var totalSwipes int64
	if err := r.db.WithContext(ctx).Model(&entity.Match{}).
		Where("user1_id = ?", userID).
		Count(&totalSwipes).Error; err != nil {
		return nil, fmt.Errorf("獲取總滑動次數失敗: %w", err)
	}
	stats.TotalSwipes = int(totalSwipes)

	// 給出的 like 數
	var likesGiven int64
	if err := r.db.WithContext(ctx).Model(&entity.Match{}).
		Where("user1_id = ? AND user1_action = ?", userID, entity.SwipeActionLike).
		Count(&likesGiven).Error; err != nil {
		return nil, fmt.Errorf("獲取給出like數失敗: %w", err)
	}
	stats.LikesGiven = int(likesGiven)

	// 收到的 like 數
	var likesReceived int64
	if err := r.db.WithContext(ctx).Model(&entity.Match{}).
		Where("user2_id = ? AND user1_action = ?", userID, entity.SwipeActionLike).
		Count(&likesReceived).Error; err != nil {
		return nil, fmt.Errorf("獲取收到like數失敗: %w", err)
	}
	stats.LikesReceived = int(likesReceived)

	// 總配對數
	var totalMatches int64
	if err := r.db.WithContext(ctx).Model(&entity.Match{}).
		Where("(user1_id = ? OR user2_id = ?) AND status = ?", userID, userID, entity.MatchStatusMatched).
		Count(&totalMatches).Error; err != nil {
		return nil, fmt.Errorf("獲取總配對數失敗: %w", err)
	}
	stats.TotalMatches = int(totalMatches)

	// 活躍配對數（有聊天記錄）
	var activeMatches int64
	if err := r.db.WithContext(ctx).
		Table("matches").
		Joins("INNER JOIN chat_messages ON matches.id = chat_messages.match_id").
		Where("(matches.user1_id = ? OR matches.user2_id = ?) AND matches.status = ?", userID, userID, entity.MatchStatusMatched).
		Count(&activeMatches).Error; err != nil {
		return nil, fmt.Errorf("獲取活躍配對數失敗: %w", err)
	}
	stats.ActiveMatches = int(activeMatches)

	// 計算比率
	if stats.TotalSwipes > 0 {
		stats.MatchRate = float64(stats.TotalMatches) / float64(stats.TotalSwipes)
	}

	// 計算受歡迎度（收到的like / 被滑動的總數）
	var totalSwipedBy int64
	if err := r.db.WithContext(ctx).Model(&entity.Match{}).
		Where("user2_id = ?", userID).
		Count(&totalSwipedBy).Error; err == nil && totalSwipedBy > 0 {
		stats.PopularityRate = float64(stats.LikesReceived) / float64(totalSwipedBy)
	}

	return stats, nil
}

// calculateAge 計算年齡的輔助函數
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()

	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age
}

// calculateDistance 計算兩點間距離的輔助函數 (Haversine formula)
func calculateDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadius = 6371 // 地球半徑 (公里)

	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
