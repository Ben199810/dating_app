package unit_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/usecase"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) SetVerified(ctx context.Context, id uint, verified bool) error {
	args := m.Called(ctx, id, verified)
	return args.Error(0)
}

func (m *MockUserRepository) SetActive(ctx context.Context, id uint, active bool) error {
	args := m.Called(ctx, id, active)
	return args.Error(0)
}

// Mock UserProfileRepository
type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserProfileRepository) GetByUserID(ctx context.Context, userID uint) (*entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserProfileRepository) Delete(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserProfileRepository) UpdateLocation(ctx context.Context, userID uint, lat, lng *float64) error {
	args := m.Called(ctx, userID, lat, lng)
	return args.Error(0)
}

func (m *MockUserProfileRepository) UpdateMatchingPreferences(ctx context.Context, userID uint, maxDistance, ageMin, ageMax int) error {
	args := m.Called(ctx, userID, maxDistance, ageMin, ageMax)
	return args.Error(0)
}

// Mock PhotoRepository
type MockPhotoRepository struct {
	mock.Mock
}

func (m *MockPhotoRepository) Create(ctx context.Context, photo *entity.Photo) error {
	args := m.Called(ctx, photo)
	return args.Error(0)
}

func (m *MockPhotoRepository) GetByID(ctx context.Context, id uint) (*entity.Photo, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Photo), args.Error(1)
}

func (m *MockPhotoRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Photo, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.Photo), args.Error(1)
}

func (m *MockPhotoRepository) Update(ctx context.Context, photo *entity.Photo) error {
	args := m.Called(ctx, photo)
	return args.Error(0)
}

func (m *MockPhotoRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPhotoRepository) SetAllNonPrimary(ctx context.Context, userID uint) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockPhotoRepository) SetPrimary(ctx context.Context, userID, photoID uint) error {
	args := m.Called(ctx, userID, photoID)
	return args.Error(0)
}

// Mock InterestRepository
type MockInterestRepository struct {
	mock.Mock
}

func (m *MockInterestRepository) Create(ctx context.Context, interest *entity.Interest) error {
	args := m.Called(ctx, interest)
	return args.Error(0)
}

func (m *MockInterestRepository) GetAll(ctx context.Context) ([]*entity.Interest, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Interest), args.Error(1)
}

func (m *MockInterestRepository) GetByUserID(ctx context.Context, userID uint) ([]*entity.Interest, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*entity.Interest), args.Error(1)
}

func (m *MockInterestRepository) AddUserInterest(ctx context.Context, userID uint, interestID uint) error {
	args := m.Called(ctx, userID, interestID)
	return args.Error(0)
}

func (m *MockInterestRepository) RemoveUserInterest(ctx context.Context, userID uint, interestID uint) error {
	args := m.Called(ctx, userID, interestID)
	return args.Error(0)
}

// Mock AgeVerificationRepository
type MockAgeVerificationRepository struct {
	mock.Mock
}

func (m *MockAgeVerificationRepository) Create(ctx context.Context, verification *entity.AgeVerification) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

func (m *MockAgeVerificationRepository) GetByUserID(ctx context.Context, userID uint) (*entity.AgeVerification, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.AgeVerification), args.Error(1)
}

func (m *MockAgeVerificationRepository) Update(ctx context.Context, verification *entity.AgeVerification) error {
	args := m.Called(ctx, verification)
	return args.Error(0)
}

func (m *MockAgeVerificationRepository) GetPendingVerifications(ctx context.Context, limit int) ([]*entity.AgeVerification, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*entity.AgeVerification), args.Error(1)
}

// Test setup helper
func setupUserService() (*usecase.UserService, *MockUserRepository, *MockUserProfileRepository, *MockPhotoRepository, *MockInterestRepository, *MockAgeVerificationRepository) {
	userRepo := &MockUserRepository{}
	userProfileRepo := &MockUserProfileRepository{}
	photoRepo := &MockPhotoRepository{}
	interestRepo := &MockInterestRepository{}
	ageVerificationRepo := &MockAgeVerificationRepository{}

	service := usecase.NewUserService(userRepo, userProfileRepo, photoRepo, interestRepo, ageVerificationRepo)

	return service, userRepo, userProfileRepo, photoRepo, interestRepo, ageVerificationRepo
}

func TestUserService_Register_Success(t *testing.T) {
	service, userRepo, userProfileRepo, _, _, _ := setupUserService()
	ctx := context.Background()

	req := &usecase.RegisterRequest{
		Email:       "test@example.com",
		Password:    "password123",
		BirthDate:   time.Now().AddDate(-20, 0, 0),
		DisplayName: "Test User",
		Gender:      "男",
		Biography:   "Test biography",
	}

	// Mock expectations
	userRepo.On("GetByEmail", ctx, "test@example.com").Return((*entity.User)(nil), errors.New("user not found"))
	userRepo.On("Create", ctx, mock.AnythingOfType("*entity.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*entity.User)
		user.ID = 1 // Simulate database ID assignment
	})
	userProfileRepo.On("Create", ctx, mock.AnythingOfType("*entity.UserProfile")).Return(nil)

	// Execute
	result, err := service.Register(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.DisplayName)
	assert.True(t, result.IsActive)

	// Verify all expectations
	userRepo.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestUserService_Register_EmailAlreadyExists(t *testing.T) {
	service, userRepo, _, _, _, _ := setupUserService()
	ctx := context.Background()

	req := &usecase.RegisterRequest{
		Email:       "existing@example.com",
		Password:    "password123",
		BirthDate:   time.Now().AddDate(-20, 0, 0),
		DisplayName: "Test User",
		Gender:      "男",
	}

	existingUser := &entity.User{
		ID:    1,
		Email: "existing@example.com",
	}

	// Mock expectations
	userRepo.On("GetByEmail", ctx, "existing@example.com").Return(existingUser, nil)

	// Execute
	result, err := service.Register(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "電子郵件已被使用")

	// Verify expectations
	userRepo.AssertExpectations(t)
}

func TestUserService_Register_UnderageUser(t *testing.T) {
	service, _, _, _, _, _ := setupUserService()
	ctx := context.Background()

	req := &usecase.RegisterRequest{
		Email:       "minor@example.com",
		Password:    "password123",
		BirthDate:   time.Now().AddDate(-16, 0, 0), // 16 years old
		DisplayName: "Minor User",
		Gender:      "女",
	}

	// Execute
	result, err := service.Register(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "必須年滿 18 歲")
}

func TestUserService_Login_Success(t *testing.T) {
	service, userRepo, userProfileRepo, _, _, _ := setupUserService()
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &entity.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		IsVerified:   true,
		BirthDate:    time.Now().AddDate(-20, 0, 0),
	}

	profile := &entity.UserProfile{
		UserID:      1,
		DisplayName: "Test User",
		Gender:      entity.GenderMale,
	}

	req := &usecase.LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	// Mock expectations
	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)
	userProfileRepo.On("GetByUserID", ctx, uint(1)).Return(profile, nil)

	// Execute
	result, err := service.Login(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.DisplayName)

	// Verify expectations
	userRepo.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	service, userRepo, _, _, _, _ := setupUserService()
	ctx := context.Background()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	user := &entity.User{
		ID:           1,
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		IsVerified:   true,
	}

	req := &usecase.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	// Mock expectations
	userRepo.On("GetByEmail", ctx, "test@example.com").Return(user, nil)

	// Execute
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "電子郵件或密碼錯誤")

	// Verify expectations
	userRepo.AssertExpectations(t)
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	service, userRepo, _, _, _, _ := setupUserService()
	ctx := context.Background()

	req := &usecase.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	// Mock expectations
	userRepo.On("GetByEmail", ctx, "notfound@example.com").Return((*entity.User)(nil), errors.New("user not found"))

	// Execute
	result, err := service.Login(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "電子郵件或密碼錯誤")

	// Verify expectations
	userRepo.AssertExpectations(t)
}

func TestUserService_GetProfile_Success(t *testing.T) {
	service, userRepo, userProfileRepo, _, _, _ := setupUserService()
	ctx := context.Background()

	user := &entity.User{
		ID:        1,
		Email:     "test@example.com",
		IsActive:  true,
		BirthDate: time.Now().AddDate(-25, 0, 0),
	}

	profile := &entity.UserProfile{
		UserID:      1,
		DisplayName: "Test User",
		Gender:      entity.GenderMale,
		Biography:   "Test bio",
	}

	// Mock expectations
	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	userProfileRepo.On("GetByUserID", ctx, uint(1)).Return(profile, nil)

	// Execute
	result, err := service.GetProfile(ctx, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "test@example.com", result.Email)
	assert.Equal(t, "Test User", result.DisplayName)
	assert.Equal(t, 25, result.Age)

	// Verify expectations
	userRepo.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestUserService_UpdateProfile_Success(t *testing.T) {
	service, userRepo, userProfileRepo, _, _, _ := setupUserService()
	ctx := context.Background()

	user := &entity.User{
		ID:        1,
		Email:     "test@example.com",
		IsActive:  true,
		BirthDate: time.Now().AddDate(-25, 0, 0),
	}

	profile := &entity.UserProfile{
		UserID:      1,
		DisplayName: "Old Name",
		Gender:      entity.GenderMale,
		Biography:   "Old bio",
	}

	req := &usecase.UpdateProfileRequest{
		DisplayName: "New Name",
		Biography:   "New bio",
		Location:    "New Location",
	}

	// Mock expectations
	userRepo.On("GetByID", ctx, uint(1)).Return(user, nil)
	userProfileRepo.On("GetByUserID", ctx, uint(1)).Return(profile, nil)
	userProfileRepo.On("Update", ctx, mock.AnythingOfType("*entity.UserProfile")).Return(nil)

	// Execute
	result, err := service.UpdateProfile(ctx, 1, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "New Name", result.DisplayName)

	// Verify expectations
	userRepo.AssertExpectations(t)
	userProfileRepo.AssertExpectations(t)
}

func TestUserService_AddPhoto_Success(t *testing.T) {
	service, _, _, photoRepo, _, _ := setupUserService()
	ctx := context.Background()

	// Mock expectations
	photoRepo.On("GetByUserID", ctx, uint(1)).Return([]*entity.Photo{}, nil)
	photoRepo.On("Create", ctx, mock.AnythingOfType("*entity.Photo")).Return(nil).Run(func(args mock.Arguments) {
		photo := args.Get(1).(*entity.Photo)
		photo.ID = 1 // Simulate database ID assignment
	})

	// Execute
	result, err := service.AddPhoto(ctx, 1, "http://example.com/photo.jpg", "Test photo")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.UserID)
	assert.Equal(t, "http://example.com/photo.jpg", result.ImageURL)
	assert.Equal(t, "Test photo", result.Description)
	assert.True(t, result.IsPrimary) // First photo should be primary

	// Verify expectations
	photoRepo.AssertExpectations(t)
}

func TestUserService_AddPhoto_ExceedsLimit(t *testing.T) {
	service, _, _, photoRepo, _, _ := setupUserService()
	ctx := context.Background()

	// Create 10 existing photos (max limit)
	existingPhotos := make([]*entity.Photo, 10)
	for i := 0; i < 10; i++ {
		existingPhotos[i] = &entity.Photo{
			ID:     uint(i + 1),
			UserID: 1,
		}
	}

	// Mock expectations
	photoRepo.On("GetByUserID", ctx, uint(1)).Return(existingPhotos, nil)

	// Execute
	result, err := service.AddPhoto(ctx, 1, "http://example.com/photo.jpg", "Test photo")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "最多只能上傳 10 張照片")

	// Verify expectations
	photoRepo.AssertExpectations(t)
}

func TestUserService_GetAvailableInterests_Success(t *testing.T) {
	service, _, _, _, interestRepo, _ := setupUserService()
	ctx := context.Background()

	interests := []*entity.Interest{
		{ID: 1, Name: "音樂", Category: "娛樂"},
		{ID: 2, Name: "運動", Category: "健康"},
	}

	// Mock expectations
	interestRepo.On("GetAll", ctx).Return(interests, nil)

	// Execute
	result, err := service.GetAvailableInterests(ctx)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "音樂", result[0].Name)
	assert.Equal(t, "運動", result[1].Name)

	// Verify expectations
	interestRepo.AssertExpectations(t)
}

func TestUserService_SubmitAgeVerification_Success(t *testing.T) {
	service, _, _, _, _, ageVerificationRepo := setupUserService()
	ctx := context.Background()

	// Mock expectations
	ageVerificationRepo.On("GetByUserID", ctx, uint(1)).Return((*entity.AgeVerification)(nil), errors.New("not found"))
	ageVerificationRepo.On("Create", ctx, mock.AnythingOfType("*entity.AgeVerification")).Return(nil)

	// Execute
	err := service.SubmitAgeVerification(ctx, 1, entity.VerificationID, "A123456789", "/path/to/document.jpg")

	// Assert
	assert.NoError(t, err)

	// Verify expectations
	ageVerificationRepo.AssertExpectations(t)
}
