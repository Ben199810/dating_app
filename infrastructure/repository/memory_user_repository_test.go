package repository

import (
	"golang_dev_docker/domain/entity"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMemoryUserRepository(t *testing.T) {
	repo := NewMemoryUserRepository()

	// 測試用戶
	user1 := &entity.User{
		ID:          uuid.New(),
		Username:    "testuser1",
		Email:       "test1@example.com",
		Password:    "hashedpassword1",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastLoginAt: nil,
	}

	user2 := &entity.User{
		ID:          uuid.New(),
		Username:    "testuser2",
		Email:       "test2@example.com",
		Password:    "hashedpassword2",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		LastLoginAt: nil,
	}

	t.Run("Save and FindByID", func(t *testing.T) {
		// 保存用戶
		err := repo.Save(user1)
		if err != nil {
			t.Fatalf("保存用戶失敗: %v", err)
		}

		// 根據ID查找
		found, err := repo.FindByID(user1.ID.String())
		if err != nil {
			t.Fatalf("根據ID查找用戶失敗: %v", err)
		}

		if found.ID != user1.ID {
			t.Errorf("期望用戶ID %v, 得到 %v", user1.ID, found.ID)
		}
		if found.Username != user1.Username {
			t.Errorf("期望用戶名 %s, 得到 %s", user1.Username, found.Username)
		}
	})

	t.Run("ExistsByUsername", func(t *testing.T) {
		// 檢查存在的用戶名
		exists, err := repo.ExistsByUsername(user1.Username)
		if err != nil {
			t.Fatalf("檢查用戶名存在性失敗: %v", err)
		}
		if !exists {
			t.Error("應該檢測到用戶名已存在")
		}

		// 檢查不存在的用戶名
		exists, err = repo.ExistsByUsername("nonexistent")
		if err != nil {
			t.Fatalf("檢查用戶名存在性失敗: %v", err)
		}
		if exists {
			t.Error("不應該檢測到不存在的用戶名")
		}

		// 測試大小寫不敏感
		exists, err = repo.ExistsByUsername(strings.ToUpper(user1.Username))
		if err != nil {
			t.Fatalf("檢查用戶名存在性失敗: %v", err)
		}
		if !exists {
			t.Error("用戶名檢查應該是大小寫不敏感的")
		}
	})

	t.Run("ExistsByEmail", func(t *testing.T) {
		// 檢查存在的電子郵件
		exists, err := repo.ExistsByEmail(user1.Email)
		if err != nil {
			t.Fatalf("檢查電子郵件存在性失敗: %v", err)
		}
		if !exists {
			t.Error("應該檢測到電子郵件已存在")
		}

		// 檢查不存在的電子郵件
		exists, err = repo.ExistsByEmail("nonexistent@example.com")
		if err != nil {
			t.Fatalf("檢查電子郵件存在性失敗: %v", err)
		}
		if exists {
			t.Error("不應該檢測到不存在的電子郵件")
		}

		// 測試大小寫不敏感
		exists, err = repo.ExistsByEmail(strings.ToUpper(user1.Email))
		if err != nil {
			t.Fatalf("檢查電子郵件存在性失敗: %v", err)
		}
		if !exists {
			t.Error("電子郵件檢查應該是大小寫不敏感的")
		}
	})

	t.Run("FindByUsername", func(t *testing.T) {
		found, err := repo.FindByUsername(user1.Username)
		if err != nil {
			t.Fatalf("根據用戶名查找用戶失敗: %v", err)
		}

		if found.Username != user1.Username {
			t.Errorf("期望用戶名 %s, 得到 %s", user1.Username, found.Username)
		}

		// 測試不存在的用戶名
		_, err = repo.FindByUsername("nonexistent")
		if err == nil {
			t.Error("應該返回錯誤，當用戶名不存在時")
		}
	})

	t.Run("FindByEmail", func(t *testing.T) {
		found, err := repo.FindByEmail(user1.Email)
		if err != nil {
			t.Fatalf("根據電子郵件查找用戶失敗: %v", err)
		}

		if found.Email != user1.Email {
			t.Errorf("期望電子郵件 %s, 得到 %s", user1.Email, found.Email)
		}

		// 測試不存在的電子郵件
		_, err = repo.FindByEmail("nonexistent@example.com")
		if err == nil {
			t.Error("應該返回錯誤，當電子郵件不存在時")
		}
	})

	t.Run("Save duplicate user", func(t *testing.T) {
		// 嘗試保存重複的用戶
		err := repo.Save(user1)
		if err == nil {
			t.Error("應該返回錯誤，當嘗試保存重複用戶時")
		}
	})

	t.Run("Save multiple users", func(t *testing.T) {
		// 保存第二個用戶
		err := repo.Save(user2)
		if err != nil {
			t.Fatalf("保存第二個用戶失敗: %v", err)
		}

		// 驗證兩個用戶都存在
		exists1, _ := repo.ExistsByUsername(user1.Username)
		exists2, _ := repo.ExistsByUsername(user2.Username)

		if !exists1 {
			t.Error("第一個用戶應該存在")
		}
		if !exists2 {
			t.Error("第二個用戶應該存在")
		}
	})
}
