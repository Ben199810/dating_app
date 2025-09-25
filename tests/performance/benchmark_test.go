package performance

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang_dev_docker/domain/entity"
)

// BenchmarkUserValidation 測試用戶驗證效能
func BenchmarkUserValidation(b *testing.B) {
	user := &entity.User{
		Email:     "test@example.com",
		IsActive:  true,
		BirthDate: time.Now().AddDate(-25, 0, 0),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = user.Validate()
	}
}

// BenchmarkUserProfileValidation 測試用戶資料驗證效能
func BenchmarkUserProfileValidation(b *testing.B) {
	profile := &entity.UserProfile{
		UserID:      1,
		DisplayName: "測試用戶",
		Bio:         "這是一個測試用戶的個人簡介",
		Gender:      entity.GenderMale,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = profile.Validate()
	}
}

// BenchmarkChatMessageValidation 測試聊天訊息驗證效能
func BenchmarkChatMessageValidation(b *testing.B) {
	message := &entity.ChatMessage{
		MatchID:    1,
		SenderID:   1,
		ReceiverID: 2,
		Type:       entity.MessageTypeText,
		Content:    "這是一個測試訊息",
		Status:     entity.MessageStatusSent,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = message.Validate()
	}
}

// BenchmarkAgeVerificationValidation 測試年齡驗證效能
func BenchmarkAgeVerificationValidation(b *testing.B) {
	verification := &entity.AgeVerification{
		UserID:            1,
		Method:            entity.VerificationMethodID,
		DocumentNumber:    "A123456789",
		DocumentImagePath: "/path/to/document.jpg",
		Status:            entity.VerificationStatusPending,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = verification.Validate()
	}
}

// TestUserAgeCalculationPerformance 測試用戶年齡計算效能
func TestUserAgeCalculationPerformance(t *testing.T) {
	const numUsers = 10000
	users := make([]*entity.User, numUsers)

	// 建立測試用戶
	for i := 0; i < numUsers; i++ {
		users[i] = &entity.User{
			Email:     fmt.Sprintf("user%d@example.com", i),
			IsActive:  true,
			BirthDate: time.Now().AddDate(-20-i%30, 0, 0), // 20-50歲
		}
	}

	start := time.Now()

	adultCount := 0
	for _, user := range users {
		if user.IsAdult() {
			adultCount++
		}
		_ = user.GetAge()
	}

	duration := time.Since(start)
	t.Logf("Calculated ages for %d users in %v", numUsers, duration)
	t.Logf("Average: %.2f calculations/second", float64(numUsers)/duration.Seconds())
	t.Logf("Adult users: %d/%d", adultCount, numUsers)

	// 效能要求：應該能夠每秒計算至少 100,000 次年齡
	assert.True(t, duration < time.Second, "Age calculation took too long")
	assert.True(t, adultCount > 0, "Should have adult users")
}

// TestChatMessageProcessingPerformance 測試聊天訊息處理效能
func TestChatMessageProcessingPerformance(t *testing.T) {
	const numMessages = 5000
	messages := make([]*entity.ChatMessage, 0, numMessages)

	start := time.Now()

	for i := 0; i < numMessages; i++ {
		message := &entity.ChatMessage{
			MatchID:    uint(i%100 + 1),
			SenderID:   uint(i%50 + 1),
			ReceiverID: uint((i+1)%50 + 1),
			Type:       entity.MessageTypeText,
			Content:    fmt.Sprintf("訊息內容 %d", i),
			Status:     entity.MessageStatusSent,
		}

		err := message.Validate()
		assert.NoError(t, err)

		// 模擬訊息狀態更新
		message.MarkAsDelivered()
		message.MarkAsRead()

		messages = append(messages, message)
	}

	duration := time.Since(start)
	t.Logf("Processed %d messages in %v", numMessages, duration)
	t.Logf("Average: %.2f messages/second", float64(numMessages)/duration.Seconds())

	// 效能要求：應該能夠每秒處理至少 500 條訊息
	assert.True(t, duration < 10*time.Second, "Message processing took too long")
	assert.Equal(t, numMessages, len(messages), "Not all messages were processed")

	// 驗證所有訊息都被標記為已讀
	for _, msg := range messages {
		assert.True(t, msg.IsRead(), "Message should be marked as read")
		assert.True(t, msg.IsDelivered(), "Message should be marked as delivered")
	}
}

// TestAgeVerificationPerformance 測試年齡驗證效能
func TestAgeVerificationPerformance(t *testing.T) {
	const numVerifications = 1000
	verifications := make([]*entity.AgeVerification, 0, numVerifications)

	start := time.Now()

	for i := 0; i < numVerifications; i++ {
		age := 25
		verification := &entity.AgeVerification{
			UserID:            uint(i + 1),
			Method:            entity.VerificationMethodID,
			DocumentNumber:    fmt.Sprintf("A%d56789", i%10),
			DocumentImagePath: fmt.Sprintf("/path/to/doc_%d.jpg", i),
			Status:            entity.VerificationStatusPending,
			ExtractedAge:      &age,
		}

		err := verification.Validate()
		assert.NoError(t, err)

		// 模擬審核流程
		if i%2 == 0 {
			err = verification.Approve(1, "文件清晰，驗證通過")
			assert.NoError(t, err)
		} else {
			err = verification.Reject(1, "文件模糊", "請重新上傳清晰的文件")
			assert.NoError(t, err)
		}

		verifications = append(verifications, verification)
	}

	duration := time.Since(start)
	t.Logf("Processed %d age verifications in %v", numVerifications, duration)
	t.Logf("Average: %.2f verifications/second", float64(numVerifications)/duration.Seconds())

	// 效能要求：應該能夠每秒處理至少 100 個驗證
	assert.True(t, duration < 10*time.Second, "Age verification processing took too long")
	assert.Equal(t, numVerifications, len(verifications), "Not all verifications were processed")

	// 統計審核結果
	approvedCount := 0
	rejectedCount := 0
	for _, verification := range verifications {
		if verification.IsApproved() {
			approvedCount++
		} else if verification.IsRejected() {
			rejectedCount++
		}
	}

	t.Logf("Approved: %d, Rejected: %d", approvedCount, rejectedCount)
	assert.Equal(t, numVerifications, approvedCount+rejectedCount, "All verifications should be processed")
}

// TestUserProfileValidationPerformance 測試用戶資料驗證效能
func TestUserProfileValidationPerformance(t *testing.T) {
	const numProfiles = 5000
	profiles := make([]*entity.UserProfile, 0, numProfiles)

	start := time.Now()

	for i := 0; i < numProfiles; i++ {
		profile := &entity.UserProfile{
			UserID:      uint(i + 1),
			DisplayName: fmt.Sprintf("用戶%d", i),
			Bio:         fmt.Sprintf("這是用戶%d的個人簡介，包含一些基本資訊。", i),
			Gender:      entity.Gender([]entity.Gender{entity.GenderMale, entity.GenderFemale}[i%2]),
			ShowAge:     i%3 == 0,
			MaxDistance: 50 + i%50,
			AgeRangeMin: 18 + i%5,
			AgeRangeMax: 35 + i%15,
		}

		err := profile.Validate()
		assert.NoError(t, err)

		profiles = append(profiles, profile)
	}

	duration := time.Since(start)
	t.Logf("Validated %d user profiles in %v", numProfiles, duration)
	t.Logf("Average: %.2f profiles/second", float64(numProfiles)/duration.Seconds())

	// 效能要求：應該能夠每秒驗證至少 1000 個資料
	assert.True(t, duration < 5*time.Second, "Profile validation took too long")
	assert.Equal(t, numProfiles, len(profiles), "Not all profiles were processed")
}