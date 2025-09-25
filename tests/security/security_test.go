package security

import (
	"strings"
	"testing"
	"time"

	"golang_dev_docker/domain/entity"

	"github.com/stretchr/testify/assert"
)

// TestUserEmailValidation 測試用戶電子郵件安全驗證
func TestUserEmailValidation(t *testing.T) {
	// 測試電子郵件驗證的安全性
	// 注意: User.Validate() 只檢查必填欄位，不檢查格式
	// 格式驗證應該在應用層或資料庫層進行
	tests := []struct {
		name        string
		email       string
		expectError bool
		description string
	}{
		{
			name:        "Valid email",
			email:       "user@example.com",
			expectError: false,
			description: "正常的電子郵件格式",
		},
		{
			name:        "SQL injection attempt",
			email:       "user'; DROP TABLE users; --@example.com",
			expectError: false, // User.Validate() 不檢查格式，只檢查是否為空
			description: "SQL 注入攻擊嘗試",
		},
		{
			name:        "XSS attempt",
			email:       "<script>alert('xss')</script>@example.com",
			expectError: false, // User.Validate() 不檢查格式，只檢查是否為空
			description: "XSS 攻擊嘗試",
		},
		{
			name:        "Long email attack",
			email:       strings.Repeat("a", 1000) + "@example.com",
			expectError: false, // User.Validate() 不檢查長度，只檢查是否為空
			description: "過長電子郵件攻擊",
		},
		{
			name:        "Unicode bypass attempt",
			email:       "user\u202e@example.com",
			expectError: false, // User.Validate() 不檢查格式，只檢查是否為空
			description: "Unicode 字符繞過攻擊",
		},
		{
			name:        "Empty email attack",
			email:       "",
			expectError: true, // User.Validate() 檢查必填欄位
			description: "空電子郵件攻擊",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entity.User{
				Email:        tt.email,
				PasswordHash: "hashedpassword123",
				IsActive:     true,
				BirthDate:    time.Now().AddDate(-25, 0, 0),
			}

			err := user.Validate()

			if tt.expectError {
				assert.Error(t, err, "Should reject %s", tt.description)
				t.Logf("Successfully blocked: %s", tt.description)
			} else {
				assert.NoError(t, err, "Should accept %s", tt.description)
				if err == nil {
					t.Logf("Correctly accepted: %s", tt.description)
				}
			}
		})
	}
}

// TestChatMessageContentSecurity 測試聊天訊息內容安全
func TestChatMessageContentSecurity(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectError bool
		description string
	}{
		{
			name:        "Normal message",
			content:     "Hello, how are you?",
			expectError: false,
			description: "正常訊息內容",
		},
		{
			name:        "Script tag injection",
			content:     "<script>alert('xss')</script>",
			expectError: false, // 內容本身不驗證 HTML，由前端處理
			description: "腳本標籤注入嘗試",
		},
		{
			name:        "SQL injection attempt",
			content:     "'; DROP TABLE messages; --",
			expectError: false, // 內容本身允許特殊字符
			description: "SQL 注入嘗試",
		},
		{
			name:        "Extremely long message",
			content:     strings.Repeat("A", 1001),
			expectError: true,
			description: "超長訊息攻擊",
		},
		{
			name:        "Empty message",
			content:     "",
			expectError: true,
			description: "空訊息",
		},
		{
			name:        "Whitespace only message",
			content:     "   \t\n  ",
			expectError: true,
			description: "僅空白字符的訊息",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := &entity.ChatMessage{
				MatchID:    1,
				SenderID:   1,
				ReceiverID: 2,
				Type:       entity.MessageTypeText,
				Content:    tt.content,
				Status:     entity.MessageStatusSent,
			}

			err := message.Validate()

			if tt.expectError {
				assert.Error(t, err, "Should reject %s", tt.description)
				t.Logf("Successfully blocked: %s", tt.description)
			} else {
				assert.NoError(t, err, "Should accept %s", tt.description)
			}
		})
	}
}

// TestUserProfileSecurityValidation 測試用戶資料安全驗證
func TestUserProfileSecurityValidation(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		bio         string
		expectError bool
		description string
	}{
		{
			name:        "Normal profile",
			displayName: "John Doe",
			bio:         "I love traveling and reading books.",
			expectError: false,
			description: "正常的用戶資料",
		},
		{
			name:        "XSS in display name",
			displayName: "<script>alert('xss')</script>",
			bio:         "Normal bio",
			expectError: false, // UserProfile.Validate() 不檢查內容格式，只檢查長度
			description: "顯示名稱中的 XSS 攻擊",
		},
		{
			name:        "HTML injection in bio",
			displayName: "John Doe",
			bio:         "<iframe src='javascript:alert(\"xss\")'></iframe>",
			expectError: false, // Bio 允許一些特殊字符，由前端處理
			description: "個人簡介中的 HTML 注入",
		},
		{
			name:        "Extremely long display name",
			displayName: strings.Repeat("A", 51),
			bio:         "Normal bio",
			expectError: true,
			description: "超長顯示名稱",
		},
		{
			name:        "Extremely long bio",
			displayName: "John Doe",
			bio:         strings.Repeat("A", 501),
			expectError: true,
			description: "超長個人簡介",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := &entity.UserProfile{
				UserID:      1,
				DisplayName: tt.displayName,
				Bio:         tt.bio,
				Gender:      entity.GenderMale,
				MaxDistance: 50,
				AgeRangeMin: 18,
				AgeRangeMax: 99,
			}

			err := profile.Validate()

			if tt.expectError {
				assert.Error(t, err, "Should reject %s", tt.description)
				t.Logf("Successfully blocked: %s", tt.description)
			} else {
				assert.NoError(t, err, "Should accept %s", tt.description)
				if err == nil {
					t.Logf("Correctly handled: %s", tt.description)
				}
			}
		})
	}
}

// TestAgeVerificationSecurity 測試年齡驗證安全性
func TestAgeVerificationSecurity(t *testing.T) {
	tests := []struct {
		name           string
		documentNumber string
		imagePath      string
		expectError    bool
		description    string
	}{
		{
			name:           "Normal verification",
			documentNumber: "A123456789",
			imagePath:      "/uploads/doc_123.jpg",
			expectError:    false,
			description:    "正常的年齡驗證",
		},
		{
			name:           "Path traversal attack",
			documentNumber: "A123456789",
			imagePath:      "../../../etc/passwd",
			expectError:    false, // 路徑驗證應該在檔案處理時進行
			description:    "路徑遍歷攻擊",
		},
		{
			name:           "Document number injection",
			documentNumber: "'; DROP TABLE age_verifications; --",
			imagePath:      "/uploads/doc.jpg",
			expectError:    false, // 內容本身允許特殊字符
			description:    "證件號碼注入攻擊",
		},
		{
			name:           "Long document number",
			documentNumber: strings.Repeat("A", 101),
			imagePath:      "/uploads/doc.jpg",
			expectError:    true,
			description:    "超長證件號碼",
		},
		{
			name:           "Long file path",
			documentNumber: "A123456789",
			imagePath:      strings.Repeat("/path", 200),
			expectError:    true,
			description:    "超長檔案路徑",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			verification := &entity.AgeVerification{
				UserID:            1,
				Method:            entity.VerificationMethodID,
				DocumentNumber:    tt.documentNumber,
				DocumentImagePath: tt.imagePath,
				Status:            entity.VerificationStatusPending,
			}

			err := verification.Validate()

			if tt.expectError {
				assert.Error(t, err, "Should reject %s", tt.description)
				t.Logf("Successfully blocked: %s", tt.description)
			} else {
				assert.NoError(t, err, "Should accept %s", tt.description)
			}
		})
	}
}

// TestBusinessLogicSecurity 測試業務邏輯安全性
func TestBusinessLogicSecurity(t *testing.T) {
	t.Run("Age verification bypass attempt", func(t *testing.T) {
		// 嘗試用未成年用戶進行年齡驗證
		underageUser := &entity.User{
			Email:     "underage@example.com",
			IsActive:  true,
			BirthDate: time.Now().AddDate(-16, 0, 0), // 16歲
		}

		assert.False(t, underageUser.IsAdult(), "User should not be adult")
		assert.Equal(t, 16, underageUser.GetAge(), "Age should be 16")

		// 嘗試創建年齡驗證
		minorAge := 16
		verification := &entity.AgeVerification{
			UserID:            1,
			Method:            entity.VerificationMethodID,
			DocumentNumber:    "A123456789",
			DocumentImagePath: "/path/to/doc.jpg",
			Status:            entity.VerificationStatusPending,
			ExtractedAge:      &minorAge,
		}

		// 嘗試通過驗證應該失敗
		err := verification.Approve(1, "嘗試通過未成年驗證")
		assert.Error(t, err, "Should not approve underage verification")
		assert.Contains(t, err.Error(), "年齡不符合要求", "Error should mention age requirement")
	})

	t.Run("Self message attempt", func(t *testing.T) {
		// 嘗試自己發送訊息給自己
		selfMessage := &entity.ChatMessage{
			MatchID:    1,
			SenderID:   1,
			ReceiverID: 1, // 相同的用戶ID
			Type:       entity.MessageTypeText,
			Content:    "Self message",
			Status:     entity.MessageStatusSent,
		}

		err := selfMessage.Validate()
		assert.Error(t, err, "Should not allow self messages")
		assert.Contains(t, err.Error(), "發送者和接收者不能是同一人", "Error should mention sender/receiver validation")
	})

	t.Run("Invalid gender bypass", func(t *testing.T) {
		profile := &entity.UserProfile{
			UserID:      1,
			DisplayName: "Test User",
			Bio:         "Test bio",
			Gender:      "invalid_gender", // 無效的性別
		}

		err := profile.Validate()
		assert.Error(t, err, "Should reject invalid gender")
	})
}

// TestDataIntegrityValidation 測試資料完整性驗證
func TestDataIntegrityValidation(t *testing.T) {
	t.Run("Required field validation", func(t *testing.T) {
		// 測試所有必填欄位
		tests := []struct {
			name   string
			entity interface {
				Validate() error
			}
			description string
		}{
			{
				name: "User with empty email",
				entity: &entity.User{
					Email:     "", // 空的必填欄位
					IsActive:  true,
					BirthDate: time.Now().AddDate(-25, 0, 0),
				},
				description: "空電子郵件",
			},
			{
				name: "ChatMessage with zero MatchID",
				entity: &entity.ChatMessage{
					MatchID:    0, // 零值的必填欄位
					SenderID:   1,
					ReceiverID: 2,
					Type:       entity.MessageTypeText,
					Content:    "Test message",
					Status:     entity.MessageStatusSent,
				},
				description: "零值 MatchID",
			},
			{
				name: "AgeVerification with zero UserID",
				entity: &entity.AgeVerification{
					UserID:            0, // 零值的必填欄位
					Method:            entity.VerificationMethodID,
					DocumentNumber:    "A123456789",
					DocumentImagePath: "/path/to/doc.jpg",
					Status:            entity.VerificationStatusPending,
				},
				description: "零值 UserID",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.entity.Validate()
				assert.Error(t, err, "Should reject %s", tt.description)
				t.Logf("Successfully validated required field: %s", tt.description)
			})
		}
	})

	t.Run("Data type validation", func(t *testing.T) {
		// 測試枚舉值驗證
		invalidMessageType := entity.MessageType("invalid")
		assert.False(t, invalidMessageType.IsValid(), "Invalid message type should be rejected")

		invalidMessageStatus := entity.MessageStatus("invalid")
		assert.False(t, invalidMessageStatus.IsValid(), "Invalid message status should be rejected")

		invalidGender := entity.Gender("invalid")
		assert.False(t, invalidGender.IsValid(), "Invalid gender should be rejected")

		invalidVerificationMethod := entity.VerificationMethod("invalid")
		assert.False(t, invalidVerificationMethod.IsValid(), "Invalid verification method should be rejected")

		invalidVerificationStatus := entity.VerificationStatus("invalid")
		assert.False(t, invalidVerificationStatus.IsValid(), "Invalid verification status should be rejected")
	})
}
