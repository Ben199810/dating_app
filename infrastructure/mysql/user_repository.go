package mysql

import (
	"database/sql"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
	"strings"
)

// baseRepository 提供基礎的資料庫操作能力
type baseRepository struct {
	db *sql.DB
}

// userRepository 專注於用戶相關的 CRUD 操作
type userRepository struct {
	baseRepository
	entityMapper *UserEntityMapper
}

// authRepository 專注於認證相關操作
type authRepository struct {
	baseRepository
}

// repository.UserRepository 定義了這個 struct 需要實現的方法
func NewUserRepository(db *sql.DB) (repository.UserRepository, repository.AuthRepository) {
	base := baseRepository{db: db}
	userRepo := &userRepository{
		baseRepository: base,
		entityMapper:   NewUserEntityMapper(),
	}
	authRepo := &authRepository{baseRepository: base}
	return userRepo, authRepo
}

func (r *userRepository) CreateUser(user *entity.UserInformation) error {
	query := `
        INSERT INTO users (
            username, email, password, age, gender, is_verified,
            status, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
    `

	// 使用 EntityMapper 準備插入資料
	values, err := r.entityMapper.PrepareForInsert(user)
	if err != nil {
		return err
	}

	result, err := r.db.Exec(query, values...)
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

func (r *userRepository) GetUserByEmail(email string) (*entity.UserInformation, error) {
	query := `
        SELECT id, username, email, password, age, gender, is_verified,
               status, last_active_at, created_at, updated_at
        FROM users 
        WHERE email = ?
    `

	dbModel := &UserDBModel{}

	err := r.db.QueryRow(query, email).Scan(
		&dbModel.ID, &dbModel.Username, &dbModel.Email, &dbModel.Password,
		&dbModel.Age, &dbModel.Gender,
		&dbModel.IsVerified, &dbModel.Status, &dbModel.LastActiveAt,
		&dbModel.CreatedAt, &dbModel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // 找不到用戶，返回 nil 而不是錯誤
		}
		return nil, err
	}

	// 使用 EntityMapper 轉換為領域實體
	user, err := r.entityMapper.ToEntity(dbModel)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetUserByID(id int) (*entity.UserInformation, error) {
	query := `
        SELECT id, username, email, password, age, gender, is_verified,
               status, last_active_at, created_at, updated_at
        FROM users 
        WHERE id = ?
    `

	dbModel := &UserDBModel{}

	err := r.db.QueryRow(query, id).Scan(
		&dbModel.ID, &dbModel.Username, &dbModel.Email, &dbModel.Password,
		&dbModel.Age, &dbModel.Gender,
		&dbModel.IsVerified, &dbModel.Status, &dbModel.LastActiveAt,
		&dbModel.CreatedAt, &dbModel.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 使用 EntityMapper 轉換為領域實體
	user, err := r.entityMapper.ToEntity(dbModel)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user *entity.UserInformation) error {
	// 使用 EntityMapper 轉換為資料庫模型
	dbModel, err := r.entityMapper.ToDBModel(user)
	if err != nil {
		return err
	}

	query := `
        UPDATE users 
        SET username = ?, email = ?, age = ?, gender = ?,
            is_verified = ?, status = ?, updated_at = NOW()
        WHERE id = ?
    `

	_, err = r.db.Exec(query,
		dbModel.Username, dbModel.Email, dbModel.Age, dbModel.Gender,
		dbModel.IsVerified, dbModel.Status, dbModel.ID,
	)
	return err
}

func (r *userRepository) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// AuthRepository 介面實作 - 移到 authRepository 結構體
func (r *authRepository) GetUserByEmail(email string) (*entity.UserInformation, error) {
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

func (r *authRepository) GetUserByUsername(username string) (*entity.UserInformation, error) {
	query := `
        SELECT id, username, email, password, created_at, updated_at 
        FROM users 
        WHERE username = ?
    `

	user := &entity.UserInformation{}
	err := r.db.QueryRow(query, username).Scan(
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

func (r *authRepository) UpdateLastLoginTime(userID int) error {
	query := `
        UPDATE users 
        SET updated_at = NOW() 
        WHERE id = ?
    `

	_, err := r.db.Exec(query, userID)
	return err
}

func (r *authRepository) UserExists(email, username string) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM users 
        WHERE email = ? OR username = ?
    `

	var count int
	err := r.db.QueryRow(query, email, username).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 新增的交友軟體功能實現
func (r *userRepository) GetUsersByLocation(lat, lng float64, radiusKm int, limit int) ([]*entity.UserInformation, error) {
	// 使用 Haversine 公式計算距離
	query := `
		SELECT id, username, email, age, gender, is_verified,
			   status, last_active_at, created_at, updated_at
		FROM users 
		WHERE id IN (
			SELECT DISTINCT up.user_id 
			FROM user_profiles up 
			WHERE up.location_lat IS NOT NULL 
			  AND up.location_lng IS NOT NULL
			  AND (6371 * acos(cos(radians(?)) * cos(radians(up.location_lat)) * 
			      cos(radians(up.location_lng) - radians(?)) + sin(radians(?)) * 
			      sin(radians(up.location_lat)))) <= ?
		)
		AND status = 'active'
		ORDER BY id
		LIMIT ?
	`

	rows, err := r.db.Query(query, lat, lng, lat, radiusKm, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.UserInformation
	for rows.Next() {
		dbModel := &UserDBModel{}

		err := rows.Scan(
			&dbModel.ID, &dbModel.Username, &dbModel.Email, &dbModel.Age, &dbModel.Gender,
			&dbModel.IsVerified, &dbModel.Status,
			&dbModel.LastActiveAt, &dbModel.CreatedAt, &dbModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 使用 EntityMapper 轉換為領域實體
		user, err := r.entityMapper.ToEntity(dbModel)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetUsersByAgeRange(minAge, maxAge int, limit int) ([]*entity.UserInformation, error) {
	query := `
		SELECT id, username, email, age, gender, is_verified,
			   status, last_active_at, created_at, updated_at
		FROM users 
		WHERE age BETWEEN ? AND ?
		  AND status = 'active'
		ORDER BY last_active_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, minAge, maxAge, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.UserInformation
	for rows.Next() {
		dbModel := &UserDBModel{}

		err := rows.Scan(
			&dbModel.ID, &dbModel.Username, &dbModel.Email, &dbModel.Age, &dbModel.Gender,
			&dbModel.IsVerified, &dbModel.Status,
			&dbModel.LastActiveAt, &dbModel.CreatedAt, &dbModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 使用 EntityMapper 轉換為領域實體
		user, err := r.entityMapper.ToEntity(dbModel)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) GetUsersByGender(gender entity.Gender, limit int) ([]*entity.UserInformation, error) {
	query := `
		SELECT id, username, email, age, gender, is_verified,
			   status, last_active_at, created_at, updated_at
		FROM users 
		WHERE gender = ?
		  AND status = 'active'
		ORDER BY last_active_at DESC
		LIMIT ?
	`

	rows, err := r.db.Query(query, gender, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.UserInformation
	for rows.Next() {
		dbModel := &UserDBModel{}

		err := rows.Scan(
			&dbModel.ID, &dbModel.Username, &dbModel.Email, &dbModel.Age, &dbModel.Gender,
			&dbModel.IsVerified, &dbModel.Status,
			&dbModel.LastActiveAt, &dbModel.CreatedAt, &dbModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 使用 EntityMapper 轉換為領域實體
		user, err := r.entityMapper.ToEntity(dbModel)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) UpdateLastActiveTime(userID int) error {
	query := `
		UPDATE users 
		SET last_active_at = NOW(), updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query, userID)
	return err
}

func (r *userRepository) SearchUsers(filters map[string]interface{}, limit, offset int) ([]*entity.UserInformation, error) {
	baseQuery := `
		SELECT id, username, email, age, gender, is_verified,
			   status, last_active_at, created_at, updated_at
		FROM users 
		WHERE status = 'active'
	`

	var conditions []string
	var args []interface{}

	// 動態建構 WHERE 條件
	if ageMin, ok := filters["age_min"]; ok {
		conditions = append(conditions, "age >= ?")
		args = append(args, ageMin)
	}
	if ageMax, ok := filters["age_max"]; ok {
		conditions = append(conditions, "age <= ?")
		args = append(args, ageMax)
	}
	if gender, ok := filters["gender"]; ok {
		conditions = append(conditions, "gender = ?")
		args = append(args, gender)
	}
	if isVerified, ok := filters["is_verified"]; ok {
		conditions = append(conditions, "is_verified = ?")
		args = append(args, isVerified)
	}

	// 移除 city 過濾器，因為 city 現在在 UserProfile 表中

	// 組合完整查詢
	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += " ORDER BY last_active_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.UserInformation
	for rows.Next() {
		dbModel := &UserDBModel{}

		err := rows.Scan(
			&dbModel.ID, &dbModel.Username, &dbModel.Email, &dbModel.Age, &dbModel.Gender,
			&dbModel.IsVerified, &dbModel.Status,
			&dbModel.LastActiveAt, &dbModel.CreatedAt, &dbModel.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// 使用 EntityMapper 轉換為領域實體
		user, err := r.entityMapper.ToEntity(dbModel)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
