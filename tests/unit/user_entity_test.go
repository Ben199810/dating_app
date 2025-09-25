package unit_test

import (
	"testing"
	"time"

	"golang_dev_docker/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestUser_IsAdult(t *testing.T) {
	tests := []struct {
		name      string
		birthDate time.Time
		expected  bool
	}{
		{
			name:      "Adult user - exactly 18 years old today",
			birthDate: time.Now().AddDate(-18, 0, 0),
			expected:  true,
		},
		{
			name:      "Adult user - over 18 years old",
			birthDate: time.Now().AddDate(-25, 0, 0),
			expected:  true,
		},
		{
			name:      "Minor user - under 18",
			birthDate: time.Now().AddDate(-17, 0, 0),
			expected:  false,
		},
		{
			name:      "Minor user - 17 years 364 days",
			birthDate: time.Now().AddDate(-18, 0, 1), // 1 day from now in the past, but 18 years ago
			expected:  false,
		},
		{
			name:      "Adult user - just turned 18",
			birthDate: time.Now().AddDate(-18, 0, -1), // 1 day ago, 18 years ago
			expected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entity.User{BirthDate: tt.birthDate}
			assert.Equal(t, tt.expected, user.IsAdult())
		})
	}
}

func TestUser_GetAge(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name        string
		birthDate   time.Time
		expectedAge int
	}{
		{
			name:        "User exactly 18 years old (birthday passed)",
			birthDate:   time.Date(now.Year()-18, now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			expectedAge: 18,
		},
		{
			name:        "User 25 years old (birthday passed)",
			birthDate:   time.Date(now.Year()-25, now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			expectedAge: 25,
		},
		{
			name:        "User 17 years old (birthday passed)",
			birthDate:   time.Date(now.Year()-17, now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			expectedAge: 17,
		},
		{
			name:        "User born today (0 years old)",
			birthDate:   time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			expectedAge: 0,
		},
		{
			name:        "User almost 18 (birthday tomorrow)",
			birthDate:   time.Date(now.Year()-18, now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC), // birthday is tomorrow, so still 17
			expectedAge: 17,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &entity.User{BirthDate: tt.birthDate}
			assert.Equal(t, tt.expectedAge, user.GetAge())
		})
	}
}

func TestUser_IsEligible(t *testing.T) {
	tests := []struct {
		name     string
		user     entity.User
		expected bool
	}{
		{
			name: "Eligible user - active, verified, adult",
			user: entity.User{
				BirthDate:  time.Now().AddDate(-20, 0, 0),
				IsActive:   true,
				IsVerified: true,
			},
			expected: true,
		},
		{
			name: "Ineligible user - not active",
			user: entity.User{
				BirthDate:  time.Now().AddDate(-20, 0, 0),
				IsActive:   false,
				IsVerified: true,
			},
			expected: false,
		},
		{
			name: "Ineligible user - not verified",
			user: entity.User{
				BirthDate:  time.Now().AddDate(-20, 0, 0),
				IsActive:   true,
				IsVerified: false,
			},
			expected: false,
		},
		{
			name: "Ineligible user - minor",
			user: entity.User{
				BirthDate:  time.Now().AddDate(-16, 0, 0),
				IsActive:   true,
				IsVerified: true,
			},
			expected: false,
		},
		{
			name: "Ineligible user - all conditions false",
			user: entity.User{
				BirthDate:  time.Now().AddDate(-16, 0, 0),
				IsActive:   false,
				IsVerified: false,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.user.IsEligible())
		})
	}
}

func TestUser_Validate(t *testing.T) {
	validBirthDate := time.Now().AddDate(-20, 0, 0)

	tests := []struct {
		name        string
		user        entity.User
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid user",
			user: entity.User{
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
				BirthDate:    validBirthDate,
			},
			expectError: false,
		},
		{
			name: "Invalid - empty email",
			user: entity.User{
				Email:        "",
				PasswordHash: "hashedpassword",
				BirthDate:    validBirthDate,
			},
			expectError: true,
			errorMsg:    "email 是必填欄位",
		},
		{
			name: "Invalid - empty password",
			user: entity.User{
				Email:        "test@example.com",
				PasswordHash: "",
				BirthDate:    validBirthDate,
			},
			expectError: true,
			errorMsg:    "password 是必填欄位",
		},
		{
			name: "Invalid - zero birth date",
			user: entity.User{
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
				BirthDate:    time.Time{},
			},
			expectError: true,
			errorMsg:    "birth_date 是必填欄位",
		},
		{
			name: "Invalid - user under 18",
			user: entity.User{
				Email:        "test@example.com",
				PasswordHash: "hashedpassword",
				BirthDate:    time.Now().AddDate(-16, 0, 0),
			},
			expectError: true,
			errorMsg:    "用戶必須年滿 18 歲",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_Deactivate(t *testing.T) {
	user := &entity.User{
		IsActive:  true,
		UpdatedAt: time.Now().Add(-1 * time.Hour), // Set to 1 hour ago
	}

	beforeDeactivate := user.UpdatedAt
	user.Deactivate()

	assert.False(t, user.IsActive)
	assert.True(t, user.UpdatedAt.After(beforeDeactivate))
}

func TestUser_Activate(t *testing.T) {
	user := &entity.User{
		IsActive:  false,
		UpdatedAt: time.Now().Add(-1 * time.Hour), // Set to 1 hour ago
	}

	beforeActivate := user.UpdatedAt
	user.Activate()

	assert.True(t, user.IsActive)
	assert.True(t, user.UpdatedAt.After(beforeActivate))
}

func TestUser_MarkAsVerified(t *testing.T) {
	user := &entity.User{
		IsVerified: false,
		UpdatedAt:  time.Now().Add(-1 * time.Hour), // Set to 1 hour ago
	}

	beforeVerification := user.UpdatedAt
	user.MarkAsVerified()

	assert.True(t, user.IsVerified)
	assert.True(t, user.UpdatedAt.After(beforeVerification))
}

// Edge case tests for age calculation around year boundaries
func TestUser_AgeCalculation_EdgeCases(t *testing.T) {
	// Test leap year scenarios
	t.Run("Leap year birthday", func(t *testing.T) {
		// Born on Feb 29, 2000 (leap year)
		leapYearBirthday := time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)
		user := &entity.User{BirthDate: leapYearBirthday}

		// Calculate age - this should not panic and should handle leap year correctly
		age := user.GetAge()
		assert.True(t, age >= 0) // Basic sanity check

		// Check if adult (should be true since 2000 was more than 18 years ago)
		assert.True(t, user.IsAdult())
	})

	// Test year day calculation edge case
	t.Run("Year day calculation accuracy", func(t *testing.T) {
		now := time.Now()

		// Create a birthday exactly 18 years ago but 1 day in the future
		almostEighteen := time.Date(now.Year()-18, now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		user := &entity.User{BirthDate: almostEighteen}

		// Should be 17 since birthday hasn't occurred yet this year
		if almostEighteen.Month() == now.Month() && almostEighteen.Day() == now.Day()+1 {
			assert.Equal(t, 17, user.GetAge())
			assert.False(t, user.IsAdult())
		}
	})
}
