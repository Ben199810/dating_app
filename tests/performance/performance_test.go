package performance

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/usecase"

	"github.com/stretchr/testify/assert"
)

// MockUserRepository 模擬用戶儲存庫
type MockUserRepository struct {
	users map[uint]*entity.User
	mutex sync.RWMutex
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[uint]*entity.User),
	}
}

func (r *MockUserRepository) Create(user *entity.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	user.ID = uint(len(r.users) + 1)
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) GetByID(id uint) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if user, exists := r.users[id]; exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) GetByEmail(email string) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *MockUserRepository) Update(user *entity.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *MockUserRepository) Delete(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.users, id)
	return nil
}

func (r *MockUserRepository) GetActiveUsers() ([]*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var activeUsers []*entity.User
	for _, user := range r.users {
		if user.IsActive() {
			activeUsers = append(activeUsers, user)
		}
	}
	return activeUsers, nil
}

// MockMatchRepository 模擬配對儲存庫
type MockMatchRepository struct {
	matches map[uint]*entity.Match
	mutex   sync.RWMutex
}

func NewMockMatchRepository() *MockMatchRepository {
	return &MockMatchRepository{
		matches: make(map[uint]*entity.Match),
	}
}

func (r *MockMatchRepository) Create(match *entity.Match) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	match.ID = uint(len(r.matches) + 1)
	r.matches[match.ID] = match
	return nil
}

func (r *MockMatchRepository) GetByID(id uint) (*entity.Match, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if match, exists := r.matches[id]; exists {
		return match, nil
	}
	return nil, fmt.Errorf("match not found")
}

func (r *MockMatchRepository) GetByUsers(userID1, userID2 uint) (*entity.Match, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, match := range r.matches {
		if (match.UserID == userID1 && match.TargetUserID == userID2) ||
			(match.UserID == userID2 && match.TargetUserID == userID1) {
			return match, nil
		}
	}
	return nil, fmt.Errorf("match not found")
}

func (r *MockMatchRepository) GetUserMatches(userID uint, limit int) ([]*entity.Match, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	var matches []*entity.Match
	for _, match := range r.matches {
		if match.UserID == userID || match.TargetUserID == userID {
			matches = append(matches, match)
			if len(matches) >= limit {
				break
			}
		}
	}
	return matches, nil
}

func (r *MockMatchRepository) Update(match *entity.Match) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.matches[match.ID] = match
	return nil
}

func (r *MockMatchRepository) Delete(id uint) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.matches, id)
	return nil
}

// 效能測試：用戶註冊
func BenchmarkUserRegistration(b *testing.B) {
	userRepo := NewMockUserRepository()
	userService := usecase.NewUserService(userRepo)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			counter++
			email := fmt.Sprintf("user%d@example.com", counter)
			_, err := userService.Register(context.Background(), &entity.User{
				Email:     email,
				Password:  "password123",
				FirstName: "Test",
				LastName:  "User",
				BirthDate: time.Now().AddDate(-25, 0, 0),
				Gender:    entity.GenderMale,
			})
			if err != nil {
				b.Error(err)
			}
		}
	})
}

// 效能測試：配對演算法
func BenchmarkMatchingAlgorithm(b *testing.B) {
	userRepo := NewMockUserRepository()
	matchRepo := NewMockMatchRepository()

	// 建立測試數據
	setupTestUsers(userRepo, 1000)

	matchingService := usecase.NewMatchingService(matchRepo, userRepo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		userID := uint(rand.Intn(1000) + 1)
		_, err := matchingService.GetPotentialMatches(context.Background(), userID, 10)
		if err != nil {
			b.Error(err)
		}
	}
}

// 效能測試：並發配對請求
func BenchmarkConcurrentMatching(b *testing.B) {
	userRepo := NewMockUserRepository()
	matchRepo := NewMockMatchRepository()

	// 建立測試數據
	setupTestUsers(userRepo, 500)

	matchingService := usecase.NewMatchingService(matchRepo, userRepo)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			userID := uint(rand.Intn(500) + 1)
			targetID := uint(rand.Intn(500) + 1)
			if userID != targetID {
				_, err := matchingService.SwipeRight(context.Background(), userID, targetID)
				if err != nil {
					b.Logf("Swipe error: %v", err)
				}
			}
		}
	})
}

// 壓力測試：大量並發用戶註冊
func TestConcurrentUserRegistration(t *testing.T) {
	userRepo := NewMockUserRepository()
	userService := usecase.NewUserService(userRepo)

	const numGoroutines = 100
	const usersPerGoroutine = 10

	var wg sync.WaitGroup
	errorsChan := make(chan error, numGoroutines*usersPerGoroutine)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(routineID int) {
			defer wg.Done()
			for j := 0; j < usersPerGoroutine; j++ {
				email := fmt.Sprintf("stress_user_%d_%d@example.com", routineID, j)
				_, err := userService.Register(context.Background(), &entity.User{
					Email:     email,
					Password:  "password123",
					FirstName: "Stress",
					LastName:  "Test",
					BirthDate: time.Now().AddDate(-25, 0, 0),
					Gender:    entity.GenderFemale,
				})
				if err != nil {
					errorsChan <- err
				}
			}
		}(i)
	}

	wg.Wait()
	close(errorsChan)

	duration := time.Since(start)
	totalUsers := numGoroutines * usersPerGoroutine

	// 檢查錯誤
	errorCount := 0
	for err := range errorsChan {
		errorCount++
		t.Logf("Registration error: %v", err)
	}

	t.Logf("Registered %d users in %v", totalUsers-errorCount, duration)
	t.Logf("Registration rate: %.2f users/second", float64(totalUsers-errorCount)/duration.Seconds())

	assert.True(t, errorCount < totalUsers/10, "Error rate should be less than 10%")
	assert.True(t, duration < 10*time.Second, "Registration should complete within 10 seconds")
}

// 壓力測試：大量並發配對請求
func TestConcurrentMatchingStress(t *testing.T) {
	userRepo := NewMockUserRepository()
	matchRepo := NewMockMatchRepository()

	// 建立測試數據
	setupTestUsers(userRepo, 200)

	matchingService := usecase.NewMatchingService(matchRepo, userRepo)

	const numGoroutines = 50
	const requestsPerGoroutine = 20

	var wg sync.WaitGroup
	errorsChan := make(chan error, numGoroutines*requestsPerGoroutine)
	matchesChan := make(chan bool, numGoroutines*requestsPerGoroutine)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				userID := uint(rand.Intn(200) + 1)
				targetID := uint(rand.Intn(200) + 1)
				if userID != targetID {
					result, err := matchingService.SwipeRight(context.Background(), userID, targetID)
					if err != nil {
						errorsChan <- err
					} else {
						matchesChan <- result.IsMatch
					}
				}
			}
		}()
	}

	wg.Wait()
	close(errorsChan)
	close(matchesChan)

	duration := time.Since(start)
	totalRequests := numGoroutines * requestsPerGoroutine

	// 統計結果
	errorCount := 0
	for err := range errorsChan {
		errorCount++
		t.Logf("Matching error: %v", err)
	}

	matchCount := 0
	for isMatch := range matchesChan {
		if isMatch {
			matchCount++
		}
	}

	t.Logf("Processed %d matching requests in %v", totalRequests-errorCount, duration)
	t.Logf("Matching rate: %.2f requests/second", float64(totalRequests-errorCount)/duration.Seconds())
	t.Logf("Match success rate: %.2f%%", float64(matchCount)/float64(totalRequests-errorCount)*100)

	assert.True(t, errorCount < totalRequests/10, "Error rate should be less than 10%")
	assert.True(t, duration < 30*time.Second, "Matching should complete within 30 seconds")
}

// setupTestUsers 建立測試用戶數據
func setupTestUsers(userRepo *MockUserRepository, count int) {
	genders := []entity.Gender{entity.GenderMale, entity.GenderFemale}
	interests := []string{"音樂", "電影", "運動", "旅行", "讀書", "美食", "攝影", "舞蹈"}

	for i := 1; i <= count; i++ {
		gender := genders[i%2]

		// 隨機選擇興趣
		userInterests := make([]string, rand.Intn(4)+1)
		for j := range userInterests {
			userInterests[j] = interests[rand.Intn(len(interests))]
		}

		user := &entity.User{
			Email:     fmt.Sprintf("testuser%d@example.com", i),
			Password:  "password123",
			FirstName: fmt.Sprintf("Test%d", i),
			LastName:  "User",
			BirthDate: time.Now().AddDate(-rand.Intn(20)-20, 0, 0), // 20-40歲
			Gender:    gender,
			Status:    entity.UserStatusActive,
			Profile: &entity.UserProfile{
				Bio:       fmt.Sprintf("Test user %d profile", i),
				Location:  "台北市",
				Interests: entity.StringArray(userInterests),
			},
		}

		err := userRepo.Create(user)
		if err != nil {
			panic(fmt.Sprintf("Failed to create test user %d: %v", i, err))
		}
	}
}
