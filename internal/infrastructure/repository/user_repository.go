package repository

import (
	"golang_dev_docker/internal/domain/user"
	"sync"
)

// InMemoryUserRepository 記憶體用戶倉儲實現
type InMemoryUserRepository struct {
	users map[user.UserID]*user.User
	mu    sync.RWMutex
}

// NewInMemoryUserRepository 創建記憶體用戶倉儲
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[user.UserID]*user.User),
	}
}

// Save 儲存用戶
func (r *InMemoryUserRepository) Save(u *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[u.ID] = u
	return nil
}

// FindByID 根據ID查找用戶
func (r *InMemoryUserRepository) FindByID(id user.UserID) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, exists := r.users[id]; exists {
		return u, nil
	}
	return nil, ErrUserNotFound
}

// FindByEmail 根據郵箱查找用戶
func (r *InMemoryUserRepository) FindByEmail(email string) (*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}

// FindAll 查找所有用戶
func (r *InMemoryUserRepository) FindAll() ([]*user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var users []*user.User
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}
