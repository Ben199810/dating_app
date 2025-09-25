package unit_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang_dev_docker/domain/entity"
)

func TestVerificationMethod_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		method   entity.VerificationMethod
		expected bool
	}{
		{
			name:     "Valid ID card verification",
			method:   entity.VerificationMethodID,
			expected: true,
		},
		{
			name:     "Valid passport verification", 
			method:   entity.VerificationMethodPassport,
			expected: true,
		},
		{
			name:     "Valid driver license verification",
			method:   entity.VerificationMethodDriverLicense,
			expected: true,
		},
		{
			name:     "Valid other verification",
			method:   entity.VerificationMethodOther,
			expected: true,
		},
		{
			name:     "Invalid verification method",
			method:   "invalid_method",
			expected: false,
		},
		{
			name:     "Empty verification method",
			method:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.method.IsValid())
		})
	}
}

func TestVerificationMethod_GetDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		method   entity.VerificationMethod
		expected string
	}{
		{
			name:     "ID card display name",
			method:   entity.VerificationMethodID,
			expected: "身分證",
		},
		{
			name:     "Passport display name",
			method:   entity.VerificationMethodPassport,
			expected: "護照",
		},
		{
			name:     "Driver license display name",
			method:   entity.VerificationMethodDriverLicense,
			expected: "駕照",
		},
		{
			name:     "Other display name",
			method:   entity.VerificationMethodOther,
			expected: "其他",
		},
		{
			name:     "Invalid method display name",
			method:   "invalid",
			expected: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.method.GetDisplayName())
		})
	}
}

func TestVerificationStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   entity.VerificationStatus
		expected bool
	}{
		{
			name:     "Valid pending status",
			status:   entity.VerificationStatusPending,
			expected: true,
		},
		{
			name:     "Valid approved status",
			status:   entity.VerificationStatusApproved,
			expected: true,
		},
		{
			name:     "Valid rejected status",
			status:   entity.VerificationStatusRejected,
			expected: true,
		},
		{
			name:     "Valid expired status",
			status:   entity.VerificationStatusExpired,
			expected: true,
		},
		{
			name:     "Invalid status",
			status:   "invalid",
			expected: false,
		},
		{
			name:     "Empty status",
			status:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestAgeVerification_Validate(t *testing.T) {
	tests := []struct {
		name         string
		verification entity.AgeVerification
		expectError  bool
		expectedError string
	}{
		{
			name: "Valid age verification",
			verification: entity.AgeVerification{
				UserID:            1,
				Method:            entity.VerificationMethodID,
				DocumentNumber:    "A123456789",
				DocumentImagePath: "/path/to/document.jpg",
				Status:            entity.VerificationStatusPending,
			},
			expectError: false,
		},
		{
			name: "Invalid - zero UserID",
			verification: entity.AgeVerification{
				UserID:            0,
				Method:            entity.VerificationMethodID,
				DocumentNumber:    "A123456789",
				DocumentImagePath: "/path/to/document.jpg",
				Status:            entity.VerificationStatusPending,
			},
			expectError:   true,
			expectedError: "user_id 是必填欄位",
		},
		{
			name: "Invalid - invalid method",
			verification: entity.AgeVerification{
				UserID:            1,
				Method:            "invalid",
				DocumentNumber:    "A123456789",
				DocumentImagePath: "/path/to/document.jpg",
				Status:            entity.VerificationStatusPending,
			},
			expectError:   true,
			expectedError: "method 必須是有效的驗證方法",
		},
		{
			name: "Invalid - empty document number",
			verification: entity.AgeVerification{
				UserID:           1,
				Method:           entity.VerificationMethodID,
				DocumentNumber:   "",
				DocumentImagePath: "/path/to/document.jpg",
				Status:           entity.VerificationStatusPending,
			},
			expectError:   true,
			expectedError: "document_number 是必填欄位",
		},
		{
			name: "Invalid - empty document image path",
			verification: entity.AgeVerification{
				UserID:           1,
				Method:           entity.VerificationMethodID,
				DocumentNumber:   "A123456789",
				DocumentImagePath: "",
				Status:           entity.VerificationStatusPending,
			},
			expectError:   true,
			expectedError: "document_image_path 是必填欄位",
		},
		{
			name: "Invalid - invalid status",
			verification: entity.AgeVerification{
				UserID:           1,
				Method:           entity.VerificationMethodID,
				DocumentNumber:   "A123456789",
				DocumentImagePath: "/path/to/document.jpg",
				Status:           "invalid_status",
			},
			expectError:   true,
			expectedError: "status 必須是有效的驗證狀態",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.verification.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAgeVerification_IsPending(t *testing.T) {
	tests := []struct {
		name         string
		verification entity.AgeVerification
		expected     bool
	}{
		{
			name: "Is pending - pending status",
			verification: entity.AgeVerification{
				Status: entity.VerificationStatusPending,
			},
			expected: true,
		},
		{
			name: "Not pending - approved status",
			verification: entity.AgeVerification{
				Status: entity.VerificationStatusApproved,
			},
			expected: false,
		},
		{
			name: "Not pending - rejected status",
			verification: entity.AgeVerification{
				Status: entity.VerificationStatusRejected,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.verification.IsPending())
		})
	}
}

func TestAgeVerification_IsApproved(t *testing.T) {
	futureTime := time.Now().Add(30 * 24 * time.Hour)
	
	tests := []struct {
		name         string
		verification entity.AgeVerification
		expected     bool
	}{
		{
			name: "Not approved - pending status",
			verification: entity.AgeVerification{
				Status: entity.VerificationStatusPending,
			},
			expected: false,
		},
		{
			name: "Is approved - approved status not expired",
			verification: entity.AgeVerification{
				Status:    entity.VerificationStatusApproved,
				ExpiresAt: &futureTime,
			},
			expected: true,
		},
		{
			name: "Not approved - rejected status",
			verification: entity.AgeVerification{
				Status: entity.VerificationStatusRejected,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.verification.IsApproved())
		})
	}
}

func TestAgeVerification_Approve(t *testing.T) {
	validAge := 25
	verification := &entity.AgeVerification{
		UserID:            1,
		Method:            entity.VerificationMethodID,
		DocumentNumber:    "A123456789",
		DocumentImagePath: "/path/to/document.jpg",
		Status:            entity.VerificationStatusPending,
		ExtractedAge:      &validAge,
		UpdatedAt:         time.Now().Add(-1 * time.Hour),
	}

	reviewerID := uint(123)
	err := verification.Approve(reviewerID, "文件已審核通過")

	assert.NoError(t, err)
	assert.Equal(t, entity.VerificationStatusApproved, verification.Status)
	assert.Equal(t, &reviewerID, verification.ReviewerID)
	assert.NotNil(t, verification.ReviewedAt)
	assert.NotNil(t, verification.ApprovedAt)
	assert.NotNil(t, verification.ExpiresAt)
}

func TestAgeVerification_Reject(t *testing.T) {
	verification := &entity.AgeVerification{
		Status:    entity.VerificationStatusPending,
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	reviewerID := uint(123)
	reason := "文件模糊無法辨識"
	notes := "請重新上傳清晰的文件"

	err := verification.Reject(reviewerID, reason, notes)

	assert.NoError(t, err)
	assert.Equal(t, entity.VerificationStatusRejected, verification.Status)
	assert.Equal(t, &reviewerID, verification.ReviewerID)
	assert.Equal(t, &reason, verification.RejectionReason)
	assert.Equal(t, &notes, verification.ReviewNotes)
	assert.NotNil(t, verification.ReviewedAt)
}

func TestAgeVerification_IsExpired(t *testing.T) {
	now := time.Now()
	pastTime := now.Add(-31 * 24 * time.Hour) // 31 days ago
	futureTime := now.Add(30 * 24 * time.Hour) // 30 days from now
	
	tests := []struct {
		name         string
		verification entity.AgeVerification
		expected     bool
	}{
		{
			name: "Not expired - recent submission",
			verification: entity.AgeVerification{
				Status:    entity.VerificationStatusApproved,
				ExpiresAt: &futureTime,
			},
			expected: false,
		},
		{
			name: "Expired - old expiry date",
			verification: entity.AgeVerification{
				Status:    entity.VerificationStatusApproved,
				ExpiresAt: &pastTime,
			},
			expected: true,
		},
		{
			name: "Edge case - expired status",
			verification: entity.AgeVerification{
				Status: entity.VerificationStatusExpired,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.verification.IsExpired())
		})
	}
}

func TestAgeVerification_IsValidAge(t *testing.T) {
	validAge := 25
	invalidAge := 16
	validBirthDate := time.Now().AddDate(-20, 0, 0) // 20 years ago
	invalidBirthDate := time.Now().AddDate(-16, 0, 0) // 16 years ago

	tests := []struct {
		name         string
		verification entity.AgeVerification
		expected     bool
	}{
		{
			name: "Valid age - extracted age 25",
			verification: entity.AgeVerification{
				ExtractedAge: &validAge,
			},
			expected: true,
		},
		{
			name: "Invalid age - extracted age 16", 
			verification: entity.AgeVerification{
				ExtractedAge: &invalidAge,
			},
			expected: false,
		},
		{
			name: "Valid age - birth date 20 years ago",
			verification: entity.AgeVerification{
				ExtractedBirthDate: &validBirthDate,
			},
			expected: true,
		},
		{
			name: "Invalid age - birth date 16 years ago",
			verification: entity.AgeVerification{
				ExtractedBirthDate: &invalidBirthDate,
			},
			expected: false,
		},
		{
			name: "No age info - should be false",
			verification: entity.AgeVerification{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.verification.IsValidAge())
		})
	}
}

func TestAgeVerification_BusinessRules(t *testing.T) {
	t.Run("Document number validation rules", func(t *testing.T) {
		// Test document number length validation
		longDocNumber := make([]byte, 101)
		for i := range longDocNumber {
			longDocNumber[i] = 'A'
		}

		verification := entity.AgeVerification{
			UserID:            1,
			Method:            entity.VerificationMethodID,
			DocumentNumber:    string(longDocNumber),
			DocumentImagePath: "/path/to/document.jpg",
			Status:            entity.VerificationStatusPending,
		}

		err := verification.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "document_number 不能超過 100 字元")
	})

	t.Run("Method and document type compatibility", func(t *testing.T) {
		verification := entity.AgeVerification{
			UserID:            1,
			Method:            entity.VerificationMethodPassport,
			DocumentNumber:    "P123456789",
			DocumentImagePath: "/path/to/passport.jpg",
			Status:            entity.VerificationStatusPending,
		}

		err := verification.Validate()
		assert.NoError(t, err)
	})
}