package mysql

import (
	"database/sql"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

type userRepository struct {
	db *sql.DB
}

// repository.UserRepository 定義了這個 struct 需要實現的方法
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *entity.UserInformation) error {
	query := `
        INSERT INTO users (username, email, password, created_at, updated_at) 
        VALUES (?, ?, ?, NOW(), NOW())
    `

	result, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return err
	}

	// 取得新建立的 ID
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *userRepository) GetByEmail(email string) (*entity.UserInformation, error) {
	query := `
        SELECT id, username, email, password, created_at, updated_at 
        FROM users 
        WHERE email = ?
    `

	user := &entity.UserInformation{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 找不到用戶，返回 nil 而不是錯誤
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetByID(id int) (*entity.UserInformation, error) {
	query := `
        SELECT id, username, email, password, created_at, updated_at 
        FROM users 
        WHERE id = ?
    `

	user := &entity.UserInformation{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(user *entity.UserInformation) error {
	query := `
        UPDATE users 
        SET username = ?, email = ?, password = ?, updated_at = NOW() 
        WHERE id = ?
    `

	_, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.ID)
	return err
}

func (r *userRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userRepository) GetUserProfile(id int) (*entity.UserInformation, error) {
	query := `
				SELECT id, username, email, created_at, updated_at 
				FROM users 
				WHERE id = ?
		`
	user := &entity.UserInformation{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 找不到用戶，返回 nil 而不是錯誤
		}
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUserProfile(user *entity.UserInformation) error {
	query := `
				UPDATE users 
				SET username = ?, email = ?, updated_at = NOW() 
				WHERE id = ?
		`

	_, err := r.db.Exec(query, user.Username, user.Email, user.ID)
	return err
}
