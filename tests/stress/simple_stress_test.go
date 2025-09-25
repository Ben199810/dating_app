package stress_test

import (
	"fmt"
	"golang_dev_docker/domain/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEntityValidationStress 測試實體驗證的壓力
func TestEntityValidationStress(t *testing.T) {
	operations := 100000

	startTime := time.Now()

	// 測試User實體驗證壓力
	t.Run("User_Validation_Stress", func(t *testing.T) {
		successCount := 0

		for i := 0; i < operations; i++ {
			user := &entity.User{
				Email:        fmt.Sprintf("stress%d@example.com", i),
				PasswordHash: "hashedpassword123",
				BirthDate:    time.Now().AddDate(-25, 0, 0),
				IsActive:     true,
			}

			err := user.Validate()
			if err == nil {
				successCount++
			}
		}

		duration := time.Since(startTime)

		t.Logf("User驗證壓力測試結果:")
		t.Logf("- 操作數: %d", operations)
		t.Logf("- 成功數: %d", successCount)
		t.Logf("- 總時間: %v", duration)
		t.Logf("- 平均每秒操作數: %.2f", float64(operations)/duration.Seconds())

		assert.Equal(t, operations, successCount, "所有用戶驗證都應該成功")
	})

	// 測試ChatMessage實體驗證壓力
	t.Run("ChatMessage_Validation_Stress", func(t *testing.T) {
		successCount := 0

		for i := 0; i < operations; i++ {
			message := &entity.ChatMessage{
				MatchID:    1,
				SenderID:   1,
				ReceiverID: 2,
				Content:    fmt.Sprintf("壓力測試訊息 #%d", i),
				Type:       entity.MessageTypeText,
				Status:     entity.MessageStatusSent,
			}

			err := message.Validate()
			if err == nil {
				successCount++
			}
		}

		duration := time.Since(startTime)

		t.Logf("ChatMessage驗證壓力測試結果:")
		t.Logf("- 操作數: %d", operations)
		t.Logf("- 成功數: %d", successCount)
		t.Logf("- 總時間: %v", duration)
		t.Logf("- 平均每秒操作數: %.2f", float64(operations)/duration.Seconds())

		assert.Equal(t, operations, successCount, "所有訊息驗證都應該成功")
	})
}

// TestConcurrentEntityOperations 測試併發實體操作
func TestConcurrentEntityOperations(t *testing.T) {
	concurrency := 100
	operationsPerWorker := 1000

	// 模擬併發用戶驗證
	t.Run("Concurrent_User_Validation", func(t *testing.T) {
		results := make(chan bool, concurrency)

		for i := 0; i < concurrency; i++ {
			go func(workerID int) {
				defer func() {
					results <- true
				}()

				for j := 0; j < operationsPerWorker; j++ {
					user := &entity.User{
						Email:        fmt.Sprintf("worker%d_user%d@example.com", workerID, j),
						PasswordHash: "hashedpassword123",
						BirthDate:    time.Now().AddDate(-25, 0, 0),
						IsActive:     true,
						IsVerified:   true,
					}

					err := user.Validate()
					assert.NoError(t, err, "用戶驗證應該成功")

					// 測試業務邏輯方法
					assert.True(t, user.IsAdult(), "用戶應該是成年人")
					assert.Greater(t, user.GetAge(), 17, "用戶年齡應該大於17")
					assert.True(t, user.IsEligible(), "用戶應該符合條件")
				}
			}(i)
		}

		// 等待所有worker完成
		for i := 0; i < concurrency; i++ {
			<-results
		}

		t.Logf("併發用戶操作完成:")
		t.Logf("- 工作者數量: %d", concurrency)
		t.Logf("- 每個工作者操作數: %d", operationsPerWorker)
		t.Logf("- 總操作數: %d", concurrency*operationsPerWorker)
	})
}

// TestAgeVerificationStress 測試年齡驗證壓力
func TestAgeVerificationStress(t *testing.T) {
	operations := 50000

	t.Run("AgeVerification_Stress", func(t *testing.T) {
		successCount := 0

		for i := 0; i < operations; i++ {
			verification := &entity.AgeVerification{
				UserID:            uint(i + 1),
				Method:            entity.VerificationMethodID,
				DocumentNumber:    fmt.Sprintf("ID%08d", i),
				DocumentImagePath: fmt.Sprintf("/uploads/doc_%d.jpg", i),
				Status:            entity.VerificationStatusPending,
			}

			err := verification.Validate()
			if err == nil {
				successCount++

				// 測試業務邏輯
				assert.True(t, verification.IsPending(), "驗證應該是待審核狀態")
				assert.False(t, verification.IsApproved(), "驗證不應該是已通過狀態")
				assert.False(t, verification.IsRejected(), "驗證不應該是已拒絕狀態")
				assert.False(t, verification.IsExpired(), "驗證不應該是已過期狀態")
			}
		}

		t.Logf("年齡驗證壓力測試結果:")
		t.Logf("- 操作數: %d", operations)
		t.Logf("- 成功數: %d", successCount)
		t.Logf("- 成功率: %.2f%%", float64(successCount)/float64(operations)*100)

		assert.Equal(t, operations, successCount, "所有年齡驗證都應該成功")
	})
}

// TestUserProfileStress 測試用戶資料壓力
func TestUserProfileStress(t *testing.T) {
	operations := 50000

	t.Run("UserProfile_Stress", func(t *testing.T) {
		successCount := 0

		for i := 0; i < operations; i++ {
			profile := &entity.UserProfile{
				UserID:      uint(i + 1),
				DisplayName: fmt.Sprintf("用戶%d", i),
				Bio:         fmt.Sprintf("這是用戶%d的個人簡介", i),
				Gender:      entity.GenderMale,
				MaxDistance: 50,
				AgeRangeMin: 18,
				AgeRangeMax: 99,
			}

			err := profile.Validate()
			if err == nil {
				successCount++

				// 測試業務邏輯
				assert.True(t, profile.Gender.IsValid(), "性別應該有效")
				assert.False(t, profile.HasLocation(), "沒有設定位置資訊")
			}
		}

		t.Logf("用戶資料壓力測試結果:")
		t.Logf("- 操作數: %d", operations)
		t.Logf("- 成功數: %d", successCount)
		t.Logf("- 成功率: %.2f%%", float64(successCount)/float64(operations)*100)

		assert.Equal(t, operations, successCount, "所有用戶資料驗證都應該成功")
	})
}

// BenchmarkUserValidation 基準測試用戶驗證
func BenchmarkUserValidation(b *testing.B) {
	user := &entity.User{
		Email:        "benchmark@example.com",
		PasswordHash: "hashedpassword123",
		BirthDate:    time.Now().AddDate(-25, 0, 0),
		IsActive:     true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Validate()
	}
}

// BenchmarkChatMessageValidation 基準測試聊天訊息驗證
func BenchmarkChatMessageValidation(b *testing.B) {
	message := &entity.ChatMessage{
		MatchID:  1,
		SenderID: 1,
		Content:  "基準測試訊息",
		Type:     entity.MessageTypeText,
		Status:   entity.MessageStatusSent,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = message.Validate()
	}
}

// BenchmarkAgeCalculation 基準測試年齡計算
func BenchmarkAgeCalculation(b *testing.B) {
	user := &entity.User{
		Email:        "benchmark@example.com",
		PasswordHash: "hashedpassword123",
		BirthDate:    time.Now().AddDate(-25, 0, 0),
		IsActive:     true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.GetAge()
		_ = user.IsAdult()
		_ = user.IsEligible()
	}
}
