package mysql

import (
	"database/sql"
	"encoding/json"
	"golang_dev_docker/domain/entity"
	"golang_dev_docker/domain/repository"
)

type userProfileRepository struct {
	baseRepository
}

type userPhotoRepository struct {
	baseRepository
}

type userPreferenceRepository struct {
	baseRepository
}

// NewUserProfileRepository 創建用戶資料儲存庫
func NewUserProfileRepository(db *sql.DB) repository.UserProfileRepository {
	base := baseRepository{db: db}
	return &userProfileRepository{baseRepository: base}
}

// NewUserPhotoRepository 創建用戶照片儲存庫
func NewUserPhotoRepository(db *sql.DB) repository.UserPhotoRepository {
	base := baseRepository{db: db}
	return &userPhotoRepository{baseRepository: base}
}

// NewUserPreferenceRepository 創建用戶偏好儲存庫
func NewUserPreferenceRepository(db *sql.DB) repository.UserPreferenceRepository {
	base := baseRepository{db: db}
	return &userPreferenceRepository{baseRepository: base}
}

// UserProfile 相關實現
func (r *userProfileRepository) CreateProfile(profile *entity.UserProfile) error {
	lookingForJSON, _ := json.Marshal(profile.LookingFor)
	languagesJSON, _ := json.Marshal(profile.Languages)
	hobbiesJSON, _ := json.Marshal(profile.Hobbies)
	lifestyleJSON, _ := json.Marshal(profile.Lifestyle)

	query := `
		INSERT INTO user_profiles (
			user_id, height, weight, education, occupation, company, relationship,
			looking_for, languages, hobbies, lifestyle, pet_preference, 
			drinking_habit, smoking_habit, exercise_habit, social_media_link,
			personality_type, zodiac, religion, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		profile.UserID, profile.Height, profile.Weight, profile.Education,
		profile.Occupation, profile.Company, profile.Relationship,
		lookingForJSON, languagesJSON, hobbiesJSON, lifestyleJSON,
		profile.PetPreference, profile.DrinkingHabit, profile.SmokingHabit,
		profile.ExerciseHabit, profile.SocialMediaLink, profile.PersonalityType,
		profile.Zodiac, profile.Religion,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	profile.ID = int(id)
	return nil
}

func (r *userProfileRepository) GetProfileByUserID(userID int) (*entity.UserProfile, error) {
	query := `
		SELECT id, user_id, height, weight, education, occupation, company, relationship,
			   looking_for, languages, hobbies, lifestyle, pet_preference,
			   drinking_habit, smoking_habit, exercise_habit, social_media_link,
			   personality_type, zodiac, religion, created_at, updated_at
		FROM user_profiles
		WHERE user_id = ?
	`

	profile := &entity.UserProfile{}
	var lookingForJSON, languagesJSON, hobbiesJSON, lifestyleJSON []byte

	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Height, &profile.Weight,
		&profile.Education, &profile.Occupation, &profile.Company, &profile.Relationship,
		&lookingForJSON, &languagesJSON, &hobbiesJSON, &lifestyleJSON,
		&profile.PetPreference, &profile.DrinkingHabit, &profile.SmokingHabit,
		&profile.ExerciseHabit, &profile.SocialMediaLink, &profile.PersonalityType,
		&profile.Zodiac, &profile.Religion, &profile.CreatedAt, &profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 解析 JSON 欄位
	json.Unmarshal(lookingForJSON, &profile.LookingFor)
	json.Unmarshal(languagesJSON, &profile.Languages)
	json.Unmarshal(hobbiesJSON, &profile.Hobbies)
	json.Unmarshal(lifestyleJSON, &profile.Lifestyle)

	return profile, nil
}

func (r *userProfileRepository) UpdateProfile(profile *entity.UserProfile) error {
	lookingForJSON, _ := json.Marshal(profile.LookingFor)
	languagesJSON, _ := json.Marshal(profile.Languages)
	hobbiesJSON, _ := json.Marshal(profile.Hobbies)
	lifestyleJSON, _ := json.Marshal(profile.Lifestyle)

	query := `
		UPDATE user_profiles SET
			height = ?, weight = ?, education = ?, occupation = ?, company = ?,
			relationship = ?, looking_for = ?, languages = ?, hobbies = ?,
			lifestyle = ?, pet_preference = ?, drinking_habit = ?, smoking_habit = ?,
			exercise_habit = ?, social_media_link = ?, personality_type = ?,
			zodiac = ?, religion = ?, updated_at = NOW()
		WHERE user_id = ?
	`

	_, err := r.db.Exec(query,
		profile.Height, profile.Weight, profile.Education, profile.Occupation,
		profile.Company, profile.Relationship, lookingForJSON, languagesJSON,
		hobbiesJSON, lifestyleJSON, profile.PetPreference, profile.DrinkingHabit,
		profile.SmokingHabit, profile.ExerciseHabit, profile.SocialMediaLink,
		profile.PersonalityType, profile.Zodiac, profile.Religion, profile.UserID,
	)

	return err
}

func (r *userProfileRepository) DeleteProfile(userID int) error {
	query := `DELETE FROM user_profiles WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}

// UserPhoto 相關實現
func (r *userPhotoRepository) CreatePhoto(photo *entity.UserPhoto) error {
	query := `
		INSERT INTO user_photos (
			user_id, photo_url, thumbnail_url, is_primary, ` + "`order`" + `,
			status, caption, is_verified, uploaded_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		photo.UserID, photo.PhotoURL, photo.ThumbnailURL, photo.IsPrimary,
		photo.Order, photo.Status, photo.Caption, photo.IsVerified,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	photo.ID = int(id)
	return nil
}

func (r *userPhotoRepository) GetPhotosByUserID(userID int) ([]*entity.UserPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, is_primary, ` + "`order`" + `,
			   status, caption, is_verified, uploaded_at, created_at, updated_at
		FROM user_photos
		WHERE user_id = ?
		ORDER BY ` + "`order`" + ` ASC, created_at ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var photos []*entity.UserPhoto
	for rows.Next() {
		photo := &entity.UserPhoto{}
		err := rows.Scan(
			&photo.ID, &photo.UserID, &photo.PhotoURL, &photo.ThumbnailURL,
			&photo.IsPrimary, &photo.Order, &photo.Status, &photo.Caption,
			&photo.IsVerified, &photo.UploadedAt, &photo.CreatedAt, &photo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		photos = append(photos, photo)
	}

	return photos, nil
}

func (r *userPhotoRepository) GetPhotoByID(id int) (*entity.UserPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, is_primary, ` + "`order`" + `,
			   status, caption, is_verified, uploaded_at, created_at, updated_at
		FROM user_photos
		WHERE id = ?
	`

	photo := &entity.UserPhoto{}
	err := r.db.QueryRow(query, id).Scan(
		&photo.ID, &photo.UserID, &photo.PhotoURL, &photo.ThumbnailURL,
		&photo.IsPrimary, &photo.Order, &photo.Status, &photo.Caption,
		&photo.IsVerified, &photo.UploadedAt, &photo.CreatedAt, &photo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return photo, nil
}

func (r *userPhotoRepository) UpdatePhoto(photo *entity.UserPhoto) error {
	query := `
		UPDATE user_photos SET
			photo_url = ?, thumbnail_url = ?, is_primary = ?, ` + "`order`" + ` = ?,
			status = ?, caption = ?, is_verified = ?, updated_at = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		photo.PhotoURL, photo.ThumbnailURL, photo.IsPrimary, photo.Order,
		photo.Status, photo.Caption, photo.IsVerified, photo.ID,
	)

	return err
}

func (r *userPhotoRepository) DeletePhoto(id int) error {
	query := `DELETE FROM user_photos WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userPhotoRepository) SetPrimaryPhoto(userID, photoID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 將該用戶的所有照片設為非主要照片
	_, err = tx.Exec(`UPDATE user_photos SET is_primary = FALSE WHERE user_id = ?`, userID)
	if err != nil {
		return err
	}

	// 設定指定照片為主要照片
	_, err = tx.Exec(`UPDATE user_photos SET is_primary = TRUE WHERE id = ? AND user_id = ?`, photoID, userID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *userPhotoRepository) GetPrimaryPhoto(userID int) (*entity.UserPhoto, error) {
	query := `
		SELECT id, user_id, photo_url, thumbnail_url, is_primary, ` + "`order`" + `,
			   status, caption, is_verified, uploaded_at, created_at, updated_at
		FROM user_photos
		WHERE user_id = ? AND is_primary = TRUE
		LIMIT 1
	`

	photo := &entity.UserPhoto{}
	err := r.db.QueryRow(query, userID).Scan(
		&photo.ID, &photo.UserID, &photo.PhotoURL, &photo.ThumbnailURL,
		&photo.IsPrimary, &photo.Order, &photo.Status, &photo.Caption,
		&photo.IsVerified, &photo.UploadedAt, &photo.CreatedAt, &photo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return photo, nil
}

// UserPreference 相關實現
func (r *userPreferenceRepository) CreatePreference(preference *entity.UserPreference) error {
	educationJSON, _ := json.Marshal(preference.Education)
	interestsJSON, _ := json.Marshal(preference.Interests)
	lifestyleJSON, _ := json.Marshal(preference.Lifestyle)

	query := `
		INSERT INTO user_preferences (
			user_id, preferred_gender, age_min, age_max, distance_max,
			height_min, height_max, education, interests, lifestyle,
			show_me, show_distance, show_age, show_last_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query,
		preference.UserID, preference.PreferredGender, preference.AgeMin, preference.AgeMax,
		preference.DistanceMax, preference.HeightMin, preference.HeightMax,
		educationJSON, interestsJSON, lifestyleJSON, preference.ShowMe,
		preference.ShowDistance, preference.ShowAge, preference.ShowLastActive,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	preference.ID = int(id)
	return nil
}

func (r *userPreferenceRepository) GetPreferenceByUserID(userID int) (*entity.UserPreference, error) {
	query := `
		SELECT id, user_id, preferred_gender, age_min, age_max, distance_max,
			   height_min, height_max, education, interests, lifestyle,
			   show_me, show_distance, show_age, show_last_active, created_at, updated_at
		FROM user_preferences
		WHERE user_id = ?
	`

	preference := &entity.UserPreference{}
	var educationJSON, interestsJSON, lifestyleJSON []byte

	err := r.db.QueryRow(query, userID).Scan(
		&preference.ID, &preference.UserID, &preference.PreferredGender,
		&preference.AgeMin, &preference.AgeMax, &preference.DistanceMax,
		&preference.HeightMin, &preference.HeightMax, &educationJSON,
		&interestsJSON, &lifestyleJSON, &preference.ShowMe, &preference.ShowDistance,
		&preference.ShowAge, &preference.ShowLastActive, &preference.CreatedAt, &preference.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 解析 JSON 欄位
	json.Unmarshal(educationJSON, &preference.Education)
	json.Unmarshal(interestsJSON, &preference.Interests)
	json.Unmarshal(lifestyleJSON, &preference.Lifestyle)

	return preference, nil
}

func (r *userPreferenceRepository) UpdatePreference(preference *entity.UserPreference) error {
	educationJSON, _ := json.Marshal(preference.Education)
	interestsJSON, _ := json.Marshal(preference.Interests)
	lifestyleJSON, _ := json.Marshal(preference.Lifestyle)

	query := `
		UPDATE user_preferences SET
			preferred_gender = ?, age_min = ?, age_max = ?, distance_max = ?,
			height_min = ?, height_max = ?, education = ?, interests = ?,
			lifestyle = ?, show_me = ?, show_distance = ?, show_age = ?,
			show_last_active = ?, updated_at = NOW()
		WHERE user_id = ?
	`

	_, err := r.db.Exec(query,
		preference.PreferredGender, preference.AgeMin, preference.AgeMax, preference.DistanceMax,
		preference.HeightMin, preference.HeightMax, educationJSON, interestsJSON,
		lifestyleJSON, preference.ShowMe, preference.ShowDistance, preference.ShowAge,
		preference.ShowLastActive, preference.UserID,
	)

	return err
}

func (r *userPreferenceRepository) DeletePreference(userID int) error {
	query := `DELETE FROM user_preferences WHERE user_id = ?`
	_, err := r.db.Exec(query, userID)
	return err
}
