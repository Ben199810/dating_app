package unit_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"golang_dev_docker/domain/usecase"
)

// Mock MatchRepository
type MockMatchRepository struct {
	mock.Mock
}

func (m *MockMatchRepository) Create(ctx context.Context, match *entity.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockMatchRepository) GetByID(ctx context.Context, id uint) (*entity.Match, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) GetMatches(ctx context.Context, userID uint, status entity.MatchStatus, limit, offset int) ([]*entity.Match, error) {
	args := m.Called(ctx, userID, status, limit, offset)
	return args.Get(0).([]*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) UpdateStatus(ctx context.Context, matchID uint, status entity.MatchStatus) error {
	args := m.Called(ctx, matchID, status)
	return args.Error(0)
}

func (m *MockMatchRepository) GetMatchBetweenUsers(ctx context.Context, user1ID, user2ID uint) (*entity.Match, error) {
	args := m.Called(ctx, user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) GetMutualMatches(ctx context.Context, userID uint, limit, offset int) ([]*entity.Match, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) CreateSwipe(ctx context.Context, match *entity.Match) error {
	args := m.Called(ctx, match)
	return args.Error(0)
}

func (m *MockMatchRepository) GetMatch(ctx context.Context, user1ID, user2ID uint) (*entity.Match, error) {
	args := m.Called(ctx, user1ID, user2ID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) UpdateMatchStatus(ctx context.Context, matchID uint, status entity.MatchStatus) error {
	args := m.Called(ctx, matchID, status)
	return args.Error(0)
}

func (m *MockMatchRepository) ProcessSwipe(ctx context.Context, userID, targetUserID uint, action entity.SwipeAction) (*entity.Match, bool, error) {
	args := m.Called(ctx, userID, targetUserID, action)
	return args.Get(0).(*entity.Match), args.Get(1).(bool), args.Error(2)
}

func (m *MockMatchRepository) GetUserMatches(ctx context.Context, userID uint, status entity.MatchStatus) ([]*entity.Match, error) {
	args := m.Called(ctx, userID, status)
	return args.Get(0).([]*entity.Match), args.Error(1)
}

func (m *MockMatchRepository) GetMatchedUsers(ctx context.Context, userID uint) ([]*entity.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockMatchRepository) HasUserSwiped(ctx context.Context, userID, targetUserID uint) (bool, error) {
	args := m.Called(ctx, userID, targetUserID)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockMatchRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockMatchRepository) GetMatchByID(ctx context.Context, id uint) (*entity.Match, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Match), args.Error(1)
}

// Mock MatchingAlgorithmRepository
type MockMatchingAlgorithmRepository struct {
	mock.Mock
}

func (m *MockMatchingAlgorithmRepository) GetPotentialMatches(ctx context.Context, userID uint, params repository.PotentialMatchParams) ([]*entity.User, error) {
	args := m.Called(ctx, userID, params)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockMatchingAlgorithmRepository) CalculateCompatibilityScore(ctx context.Context, user1ID, user2ID uint) (float64, error) {
	args := m.Called(ctx, user1ID, user2ID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockMatchingAlgorithmRepository) GetMatchingStats(ctx context.Context, userID uint) (*repository.MatchingStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.MatchingStats), args.Error(1)
}

func (m *MockMatchingAlgorithmRepository) GetUsersNearby(ctx context.Context, userID uint, lat, lng float64, maxDistanceKm int, limit int) ([]*entity.User, error) {
	args := m.Called(ctx, userID, lat, lng, maxDistanceKm, limit)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockMatchingAlgorithmRepository) GetUsersByAgeRange(ctx context.Context, userID uint, minAge, maxAge int, limit int) ([]*entity.User, error) {
	args := m.Called(ctx, userID, minAge, maxAge, limit)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockMatchingAlgorithmRepository) GetUsersByCommonInterests(ctx context.Context, userID uint, limit int) ([]*entity.User, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]*entity.User), args.Error(1)
}

// Mock cache interface
type MockMatchingCache struct {
	mock.Mock
}

func (m *MockMatchingCache) CachePotentialMatches(userID uint, matches []*entity.User) error {
	args := m.Called(userID, matches)
	return args.Error(0)
}

func (m *MockMatchingCache) GetPotentialMatches(userID uint) ([]*entity.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockMatchingCache) InvalidatePotentialMatches(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockMatchingCache) CacheUserMatches(userID uint, status entity.MatchStatus, matches []*entity.Match) error {
	args := m.Called(userID, status, matches)
	return args.Error(0)
}

func (m *MockMatchingCache) GetUserMatches(userID uint, status entity.MatchStatus) ([]*entity.Match, error) {
	args := m.Called(userID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Match), args.Error(1)
}

func (m *MockMatchingCache) InvalidateUserMatches(userID uint, status entity.MatchStatus) error {
	args := m.Called(userID, status)
	return args.Error(0)
}

func (m *MockMatchingCache) InvalidateAllUserMatches(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockMatchingCache) CacheMatchingStats(userID uint, stats *repository.MatchingStats) error {
	args := m.Called(userID, stats)
	return args.Error(0)
}

func (m *MockMatchingCache) GetMatchingStats(userID uint) (*repository.MatchingStats, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.MatchingStats), args.Error(1)
}

func (m *MockMatchingCache) InvalidateMatchingStats(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockMatchingCache) CacheCompatibilityScore(user1ID, user2ID uint, score float64) error {
	args := m.Called(user1ID, user2ID, score)
	return args.Error(0)
}

func (m *MockMatchingCache) GetCompatibilityScore(user1ID, user2ID uint) (float64, error) {
	args := m.Called(user1ID, user2ID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockMatchingCache) CacheNearbyUsers(userID uint, maxDistance int, users []*entity.User) error {
	args := m.Called(userID, maxDistance, users)
	return args.Error(0)
}

func (m *MockMatchingCache) GetNearbyUsers(userID uint, maxDistance int) ([]*entity.User, error) {
	args := m.Called(userID, maxDistance)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockMatchingCache) CacheCommonInterestUsers(userID uint, users []*entity.User) error {
	args := m.Called(userID, users)
	return args.Error(0)
}

func (m *MockMatchingCache) GetCommonInterestUsers(userID uint) ([]*entity.User, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockMatchingCache) InvalidateUserCache(userID uint) error {
	args := m.Called(userID)
	return args.Error(0)
}

// Test setup helper
func setupMatchingService() (*usecase.MatchingService, *MockMatchRepository, *MockMatchingAlgorithmRepository, *MockUserRepository, *MockUserProfileRepository, *MockMatchingCache) {
	matchRepo := &MockMatchRepository{}
	algorithmRepo := &MockMatchingAlgorithmRepository{}
	userRepo := &MockUserRepository{}
	profileRepo := &MockUserProfileRepository{}
	cache := &MockMatchingCache{}

	// 創建配對服務，設置快取
	service := usecase.NewMatchingService(matchRepo, algorithmRepo, userRepo, profileRepo)
	// service.SetCache(cache) // 假設有此方法

	return service, matchRepo, algorithmRepo, userRepo, profileRepo, cache
}

func TestMatchingService_GetPotentialMatches_Success(t *testing.T) {
	service, _, algorithmRepo, userRepo, profileRepo, _ := setupMatchingService()
	ctx := context.Background()

	userID := uint(1)
	user := &entity.User{
		ID:       userID,
		IsActive: true,
	}

	profile := &entity.UserProfile{
		UserID:      userID,
		MaxDistance: 50,
		AgeRangeMin: 20,
		AgeRangeMax: 30,
	}

	potentialUsers := []*entity.User{
		{ID: 2, Email: "user2@example.com", IsActive: true},
		{ID: 3, Email: "user3@example.com", IsActive: true},
	}

	// Mock expectations
	userRepo.On("GetByID", ctx, userID).Return(user, nil)
	profileRepo.On("GetByUserID", ctx, userID).Return(profile, nil)
	algorithmRepo.On("GetPotentialMatches", ctx, userID, 50, 20, 30, 20).Return(potentialUsers, nil)

	// Execute
	result, err := service.GetPotentialMatches(ctx, userID, 20)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, uint(2), result[0].ID)
	assert.Equal(t, uint(3), result[1].ID)

	// Verify expectations
	userRepo.AssertExpectations(t)
	profileRepo.AssertExpectations(t)
	algorithmRepo.AssertExpectations(t)
}

func TestMatchingService_GetPotentialMatches_UserNotFound(t *testing.T) {
	service, _, _, userRepo, _, _ := setupMatchingService()
	ctx := context.Background()

	userID := uint(999)

	// Mock expectations
	userRepo.On("GetByID", ctx, userID).Return((*entity.User)(nil), errors.New("user not found"))

	// Execute
	result, err := service.GetPotentialMatches(ctx, userID, 20)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "獲取用戶資料失敗")

	// Verify expectations
	userRepo.AssertExpectations(t)
}

func TestMatchingService_SwipeRight_Success_CreateMatch(t *testing.T) {
	service, matchRepo, _, _, _, _ := setupMatchingService()
	ctx := context.Background()

	userID := uint(1)
	targetUserID := uint(2)

	// Mock expectations - 沒有現有配對
	matchRepo.On("GetMatchBetweenUsers", ctx, userID, targetUserID).Return((*entity.Match)(nil), errors.New("match not found"))
	matchRepo.On("GetMatchBetweenUsers", ctx, targetUserID, userID).Return((*entity.Match)(nil), errors.New("match not found"))
	matchRepo.On("Create", ctx, mock.AnythingOfType("*entity.Match")).Return(nil)

	// Execute
	result, err := service.SwipeRight(ctx, userID, targetUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, targetUserID, result.TargetUserID)
	assert.Equal(t, entity.MatchStatusPending, result.Status)

	// Verify expectations
	matchRepo.AssertExpectations(t)
}

func TestMatchingService_SwipeRight_Success_MutualMatch(t *testing.T) {
	service, matchRepo, _, _, _, _ := setupMatchingService()
	ctx := context.Background()

	userID := uint(1)
	targetUserID := uint(2)

	existingMatch := &entity.Match{
		ID:           1,
		UserID:       targetUserID,
		TargetUserID: userID,
		Status:       entity.MatchStatusPending,
	}

	// Mock expectations - 目標用戶已經向我滑右
	matchRepo.On("GetMatchBetweenUsers", ctx, userID, targetUserID).Return((*entity.Match)(nil), errors.New("match not found"))
	matchRepo.On("GetMatchBetweenUsers", ctx, targetUserID, userID).Return(existingMatch, nil)
	matchRepo.On("Create", ctx, mock.AnythingOfType("*entity.Match")).Return(nil)
	matchRepo.On("UpdateStatus", ctx, uint(1), entity.MatchStatusMatched).Return(nil)

	// Execute
	result, err := service.SwipeRight(ctx, userID, targetUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.MatchStatusMatched, result.Status)

	// Verify expectations
	matchRepo.AssertExpectations(t)
}

func TestMatchingService_SwipeLeft_Success(t *testing.T) {
	service, matchRepo, _, _, _, _ := setupMatchingService()
	ctx := context.Background()

	userID := uint(1)
	targetUserID := uint(2)

	// Mock expectations - 沒有現有配對
	matchRepo.On("GetMatchBetweenUsers", ctx, userID, targetUserID).Return((*entity.Match)(nil), errors.New("match not found"))
	matchRepo.On("Create", ctx, mock.AnythingOfType("*entity.Match")).Return(nil)

	// Execute
	result, err := service.SwipeLeft(ctx, userID, targetUserID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, targetUserID, result.TargetUserID)
	assert.Equal(t, entity.MatchStatusRejected, result.Status)

	// Verify expectations
	matchRepo.AssertExpectations(t)
}

func TestMatchingService_GetUserMatches_Success(t *testing.T) {
	service, matchRepo, _, _, _, _ := setupMatchingService()
	ctx := context.Background()

	userID := uint(1)
	expectedMatches := []*entity.Match{
		{ID: 1, UserID: userID, TargetUserID: 2, Status: entity.MatchStatusMatched},
		{ID: 2, UserID: userID, TargetUserID: 3, Status: entity.MatchStatusMatched},
	}

	// Mock expectations
	matchRepo.On("GetMutualMatches", ctx, userID, 20, 0).Return(expectedMatches, nil)

	// Execute
	result, err := service.GetUserMatches(ctx, userID, 20, 0)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, uint(1), result[0].ID)
	assert.Equal(t, uint(2), result[1].ID)

	// Verify expectations
	matchRepo.AssertExpectations(t)
}

func TestMatchingService_CalculateCompatibility_Success(t *testing.T) {
	service, _, algorithmRepo, _, _, _ := setupMatchingService()
	ctx := context.Background()

	user1ID := uint(1)
	user2ID := uint(2)
	expectedScore := 0.85

	// Mock expectations
	algorithmRepo.On("CalculateCompatibilityScore", ctx, user1ID, user2ID).Return(expectedScore, nil)

	// Execute
	result, err := service.CalculateCompatibility(ctx, user1ID, user2ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedScore, result)

	// Verify expectations
	algorithmRepo.AssertExpectations(t)
}
