package repository

import (
	"errors"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// MemoryUserRepository 內存用戶儲存庫實現，適用於開發和測試
type MemoryUserRepository struct {
	users map[uuid.UUID]*entity.User
	mutex sync.RWMutex
}

// NewMemoryUserRepository 創建新的內存用戶儲存庫
func NewMemoryUserRepository() repository.UserRepository {
	return &MemoryUserRepository{
		users: make(map[uuid.UUID]*entity.User),
		mutex: sync.RWMutex{},
	}
}

// ExistsByUsername 檢查用戶名是否已存在
func (r *MemoryUserRepository) ExistsByUsername(username string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	username = strings.TrimSpace(username)
	for _, user := range r.users {
		if strings.EqualFold(user.Username, username) {
			return true, nil
		}
	}
	return false, nil
}

// ExistsByEmail 檢查電子郵件是否已存在
func (r *MemoryUserRepository) ExistsByEmail(email string) (bool, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			return true, nil
		}
	}
	return false, nil
}

// Save 保存用戶
func (r *MemoryUserRepository) Save(user *entity.User) error {
	if user == nil {
		return errors.New("用戶不能為空")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 檢查是否已存在（避免重複保存）
	if _, exists := r.users[user.ID]; exists {
		return errors.New("用戶已存在")
	}

	// 複製用戶以避免外部修改
	userCopy := *user
	r.users[user.ID] = &userCopy

	return nil
}

// FindByID 根據ID查找用戶
func (r *MemoryUserRepository) FindByID(id string) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("無效的用戶ID格式")
	}

	user, exists := r.users[userID]
	if !exists {
		return nil, errors.New("找不到用戶")
	}

	// 返回用戶的副本
	userCopy := *user
	return &userCopy, nil
}

// FindByUsername 根據用戶名查找用戶
func (r *MemoryUserRepository) FindByUsername(username string) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	username = strings.TrimSpace(username)
	for _, user := range r.users {
		if strings.EqualFold(user.Username, username) {
			// 返回用戶的副本
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, errors.New("找不到用戶")
}

// FindByEmail 根據電子郵件查找用戶
func (r *MemoryUserRepository) FindByEmail(email string) (*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	email = strings.ToLower(strings.TrimSpace(email))
	for _, user := range r.users {
		if strings.EqualFold(user.Email, email) {
			// 返回用戶的副本
			userCopy := *user
			return &userCopy, nil
		}
	}

	return nil, errors.New("找不到用戶")
}

// GetAllUsers 獲取所有用戶（額外方法，用於調試）
func (r *MemoryUserRepository) GetAllUsers() []*entity.User {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		userCopy := *user
		users = append(users, &userCopy)
	}

	return users
}

// Clear 清空所有用戶（用於測試）
func (r *MemoryUserRepository) Clear() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.users = make(map[uuid.UUID]*entity.User)
}
