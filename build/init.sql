-- 18+ 交友聊天應用程式資料庫初始化腳本-- 18+ 交友聊天應用程式資料庫初始化腳本

-- 基於 specs/001-18/data-model.md 規格設計-- 基於 specs/001-18/data-model.md 規格設計



-- 建立資料庫 (如果不存在)-- 建立資料庫 (如果不存在)

CREATE DATABASE IF NOT EXISTS dating_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;CREATE DATABASE IF NOT EXISTS dating_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;



USE dating_app;USE dating_app;



-- 創建應用程式用戶（如果不存在）-- 創建應用程式用戶（如果不存在）

-- 由於 MySQL 容器已通過環境變數創建了用戶，這裡只是確保權限-- 由於 MySQL 容器已通過環境變數創建了用戶，這裡只是確保權限

GRANT ALL PRIVILEGES ON dating_app.* TO 'dating_user'@'%';GRANT ALL PRIVILEGES ON dating_app.* TO 'dating_user'@'%';

FLUSH PRIVILEGES;FLUSH PRIVILEGES;



-- 1. 用戶基本資料表-- 1. 用戶基本資料表-- 1. 用戶基本資料表

CREATE TABLE IF NOT EXISTS users (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    email VARCHAR(255) UNIQUE NOT NULL,CREATE TABLE IF NOT EXISTS users (CREATE TABLE IF NOT EXISTS users (

    password_hash VARCHAR(255) NOT NULL,

    birth_date DATE NOT NULL,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否通過年齡驗證',

    is_active BOOLEAN DEFAULT TRUE COMMENT '帳戶是否啟用',    email VARCHAR(255) UNIQUE NOT NULL,    email VARCHAR(255) UNIQUE NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    password_hash VARCHAR(255) NOT NULL,    password_hash VARCHAR(255) NOT NULL,

    INDEX idx_email (email),

    INDEX idx_is_active (is_active),    birth_date DATE NOT NULL,    birth_date DATE NOT NULL,

    INDEX idx_is_verified (is_verified)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否通過年齡驗證',    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否通過年齡驗證',



-- 2. 用戶詳細資料表 (一對一關係)    is_active BOOLEAN DEFAULT TRUE COMMENT '帳戶是否啟用',    is_active BOOLEAN DEFAULT TRUE COMMENT '帳戶是否啟用',

CREATE TABLE IF NOT EXISTS user_profiles (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    user_id INT UNSIGNED NOT NULL,

    name VARCHAR(100) NOT NULL,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    gender ENUM('male', 'female', 'other') NOT NULL,

    bio TEXT,        

    location_city VARCHAR(100),

    location_country VARCHAR(100) DEFAULT '台灣',    INDEX idx_email (email),    INDEX idx_email (email),

    height INT COMMENT '身高 (cm)',

    job_title VARCHAR(100),    INDEX idx_is_verified (is_verified),    INDEX idx_is_verified (is_verified),

    education_level ENUM('high_school', 'bachelor', 'master', 'phd', 'other'),

    interests JSON COMMENT '興趣愛好標籤陣列',    INDEX idx_is_active (is_active),    INDEX idx_is_active (is_active),

    looking_for ENUM('serious', 'casual', 'friends', 'unsure') DEFAULT 'serious',

    min_age_preference INT DEFAULT 18,    INDEX idx_created_at (created_at)    INDEX idx_created_at (created_at)

    max_age_preference INT DEFAULT 100,

    max_distance_km INT DEFAULT 50 COMMENT '最大距離偏好(公里)',) ENGINE=InnoDB COMMENT='用戶基本資料';

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,) ENGINE=InnoDB COMMENT='用戶基本資料';

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    UNIQUE KEY unique_user_profile (user_id),-- 2. 用戶檔案詳細資料表-- 2. 用戶檔案詳細資料表

    INDEX idx_gender (gender),

    INDEX idx_location (location_city, location_country),

    INDEX idx_age_prefs (min_age_preference, max_age_preference),CREATE TABLE IF NOT EXISTS user_profiles (CREATE TABLE IF NOT EXISTS user_profiles (

    INDEX idx_looking_for (looking_for)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,



-- 3. 照片管理表    user_id INT UNSIGNED NOT NULL UNIQUE,    user_id INT UNSIGNED NOT NULL UNIQUE,

CREATE TABLE IF NOT EXISTS photos (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    display_name VARCHAR(50) NOT NULL COMMENT '顯示名稱',    display_name VARCHAR(50) NOT NULL COMMENT '顯示名稱',

    user_id INT UNSIGNED NOT NULL,

    file_path VARCHAR(500) NOT NULL,    bio TEXT COMMENT '個人簡介，最多500字元',    bio TEXT COMMENT '個人簡介，最多500字元',

    file_size INT UNSIGNED,

    mime_type VARCHAR(100),    gender ENUM('male', 'female', 'other') NOT NULL COMMENT '性別',    gender ENUM('male', 'female', 'other') NOT NULL COMMENT '性別',

    is_primary BOOLEAN DEFAULT FALSE COMMENT '是否為主要照片',

    is_verified BOOLEAN DEFAULT FALSE COMMENT '是否通過審核',    show_age BOOLEAN DEFAULT TRUE COMMENT '是否顯示年齡',    show_age BOOLEAN DEFAULT TRUE COMMENT '是否顯示年齡',

    upload_ip VARCHAR(45),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    location_lat DECIMAL(10, 8) NULL COMMENT '緯度座標',    location_lat DECIMAL(10, 8) NULL COMMENT '緯度座標',

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    location_lng DECIMAL(11, 8) NULL COMMENT '經度座標',    location_lng DECIMAL(11, 8) NULL COMMENT '經度座標',

    INDEX idx_user_id (user_id),

    INDEX idx_is_primary (is_primary),    max_distance INT DEFAULT 50 COMMENT '配對最大距離(km)',    max_distance INT DEFAULT 50 COMMENT '配對最大距離(km)',

    INDEX idx_is_verified (is_verified)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;    age_range_min INT DEFAULT 18 COMMENT '期望年齡範圍最小值',    age_range_min INT DEFAULT 18 COMMENT '期望年齡範圍最小值',



-- 4. 年齡驗證表    age_range_max INT DEFAULT 99 COMMENT '期望年齡範圍最大值',    age_range_max INT DEFAULT 99 COMMENT '期望年齡範圍最大值',

CREATE TABLE IF NOT EXISTS age_verifications (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    user_id INT UNSIGNED NOT NULL,

    verification_method ENUM('birth_date', 'id_document', 'credit_card') NOT NULL,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    document_path VARCHAR(500),

    verification_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',        

    verified_at TIMESTAMP NULL,

    expires_at TIMESTAMP NULL,    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    admin_notes TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    INDEX idx_user_id (user_id),    INDEX idx_user_id (user_id),

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    INDEX idx_gender (gender),    INDEX idx_gender (gender),

    INDEX idx_user_id (user_id),

    INDEX idx_status (verification_status),    INDEX idx_location (location_lat, location_lng),    INDEX idx_location (location_lat, location_lng),

    INDEX idx_expires_at (expires_at)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;    INDEX idx_age_range (age_range_min, age_range_max)    INDEX idx_age_range (age_range_min, age_range_max)



-- 5. 配對表 (雙向關係)) ENGINE=InnoDB COMMENT='用戶詳細檔案';

CREATE TABLE IF NOT EXISTS matches (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,) ENGINE=InnoDB COMMENT='用戶詳細檔案';

    user1_id INT UNSIGNED NOT NULL COMMENT '發起配對的用戶',

    user2_id INT UNSIGNED NOT NULL COMMENT '被配對的用戶',-- 3. 興趣標籤表-- 3. 興趣標籤表

    user1_action ENUM('like', 'pass', 'super_like') NOT NULL,

    user2_action ENUM('like', 'pass', 'super_like', 'pending') DEFAULT 'pending',

    match_status ENUM('pending', 'matched', 'unmatched') DEFAULT 'pending',CREATE TABLE IF NOT EXISTS interests (CREATE TABLE IF NOT EXISTS interests (

    matched_at TIMESTAMP NULL COMMENT '配對成功時間',

    user1_actioned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    user2_actioned_at TIMESTAMP NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    name VARCHAR(50) NOT NULL UNIQUE COMMENT '興趣名稱',    name VARCHAR(50) NOT NULL UNIQUE COMMENT '興趣名稱',

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,    category VARCHAR(30) NOT NULL COMMENT '興趣分類',    category VARCHAR(30) NOT NULL COMMENT '興趣分類',

    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,

    UNIQUE KEY unique_match_pair (user1_id, user2_id),    is_active BOOLEAN DEFAULT TRUE,    is_active BOOLEAN DEFAULT TRUE,

    INDEX idx_user1_id (user1_id),

    INDEX idx_user2_id (user2_id),    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_match_status (match_status),

    INDEX idx_matched_at (matched_at),        

    INDEX idx_actions (user1_action, user2_action)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;    INDEX idx_name (name),    INDEX idx_name (name),



-- 6. 聊天訊息表    INDEX idx_category (category),    INDEX idx_category (category),

CREATE TABLE IF NOT EXISTS chat_messages (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    INDEX idx_is_active (is_active)    INDEX idx_is_active (is_active)

    match_id INT UNSIGNED NOT NULL,

    sender_id INT UNSIGNED NOT NULL,) ENGINE=InnoDB COMMENT='興趣標籤';

    message_type ENUM('text', 'image', 'emoji', 'system') DEFAULT 'text',

    content TEXT NOT NULL,) ENGINE=InnoDB COMMENT='興趣標籤';

    media_path VARCHAR(500),

    is_read BOOLEAN DEFAULT FALSE,-- 4. 用戶興趣關聯表-- 4. 用戶興趣關聯表

    read_at TIMESTAMP NULL,

    is_deleted BOOLEAN DEFAULT FALSE COMMENT '軟刪除標記',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,CREATE TABLE IF NOT EXISTS user_interests (CREATE TABLE IF NOT EXISTS user_interests (

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_match_id (match_id),    user_id INT UNSIGNED NOT NULL,    user_id INT UNSIGNED NOT NULL,

    INDEX idx_sender_id (sender_id),

    INDEX idx_created_at (created_at),    interest_id INT UNSIGNED NOT NULL,    interest_id INT UNSIGNED NOT NULL,

    INDEX idx_is_read (is_read),

    INDEX idx_message_type (message_type)    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

        

-- 7. 舉報表

CREATE TABLE IF NOT EXISTS reports (    UNIQUE KEY unique_user_interest (user_id, interest_id),    UNIQUE KEY unique_user_interest (user_id, interest_id),

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    reporter_id INT UNSIGNED NOT NULL COMMENT '舉報人',    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    reported_id INT UNSIGNED NOT NULL COMMENT '被舉報人',

    report_type ENUM('fake_profile', 'inappropriate_content', 'harassment', 'spam', 'underage', 'other') NOT NULL,    FOREIGN KEY (interest_id) REFERENCES interests(id) ON DELETE CASCADE,    FOREIGN KEY (interest_id) REFERENCES interests(id) ON DELETE CASCADE,

    description TEXT,

    evidence_paths JSON COMMENT '證據文件路徑陣列',    INDEX idx_user_id (user_id),    INDEX idx_user_id (user_id),

    status ENUM('pending', 'investigating', 'resolved', 'dismissed') DEFAULT 'pending',

    admin_notes TEXT,    INDEX idx_interest_id (interest_id)    INDEX idx_interest_id (interest_id)

    resolved_at TIMESTAMP NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,) ENGINE=InnoDB COMMENT='用戶興趣關聯';

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,) ENGINE=InnoDB COMMENT='用戶興趣關聯';

    FOREIGN KEY (reported_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_reporter_id (reporter_id),-- 5. 用戶照片表-- 5. 用戶照片表

    INDEX idx_reported_id (reported_id),

    INDEX idx_report_type (report_type),

    INDEX idx_status (status),CREATE TABLE IF NOT EXISTS photos (CREATE TABLE IF NOT EXISTS photos (

    INDEX idx_created_at (created_at)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,



-- 8. 封鎖表    user_id INT UNSIGNED NOT NULL,    user_id INT UNSIGNED NOT NULL,

CREATE TABLE IF NOT EXISTS blocks (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    file_path VARCHAR(500) NOT NULL COMMENT '照片檔案路徑',    file_path VARCHAR(500) NOT NULL COMMENT '照片檔案路徑',

    blocker_id INT UNSIGNED NOT NULL COMMENT '封鎖人',

    blocked_id INT UNSIGNED NOT NULL COMMENT '被封鎖人',    is_primary BOOLEAN DEFAULT FALSE COMMENT '是否為主要照片',    is_primary BOOLEAN DEFAULT FALSE COMMENT '是否為主要照片',

    reason ENUM('harassment', 'spam', 'inappropriate', 'not_interested', 'other'),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    display_order INT DEFAULT 0 COMMENT '顯示順序',    display_order INT DEFAULT 0 COMMENT '顯示順序',

    FOREIGN KEY (blocker_id) REFERENCES users(id) ON DELETE CASCADE,

    FOREIGN KEY (blocked_id) REFERENCES users(id) ON DELETE CASCADE,    upload_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',    upload_status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',

    UNIQUE KEY unique_block_pair (blocker_id, blocked_id),

    INDEX idx_blocker_id (blocker_id),    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX idx_blocked_id (blocked_id)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,



-- 建立一些基本的觸發器和約束        

-- 確保配對表的用戶不能自己配對自己

ALTER TABLE matches ADD CONSTRAINT chk_no_self_match CHECK (user1_id != user2_id);    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,



-- 確保封鎖表的用戶不能封鎖自己    INDEX idx_user_id (user_id),    INDEX idx_user_id (user_id),

ALTER TABLE blocks ADD CONSTRAINT chk_no_self_block CHECK (blocker_id != blocked_id);

    INDEX idx_is_primary (is_primary),    INDEX idx_is_primary (is_primary),

-- 確保舉報表的用戶不能舉報自己

ALTER TABLE reports ADD CONSTRAINT chk_no_self_report CHECK (reporter_id != reported_id);    INDEX idx_display_order (display_order),    INDEX idx_display_order (display_order),

    INDEX idx_upload_status (upload_status)    INDEX idx_upload_status (upload_status)

) ENGINE=InnoDB COMMENT='用戶照片';

) ENGINE=InnoDB COMMENT='用戶照片';

-- 6. 年齡驗證記錄表-- 6. 年齡驗證記錄表


CREATE TABLE IF NOT EXISTS age_verifications (CREATE TABLE IF NOT EXISTS age_verifications (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    user_id INT UNSIGNED NOT NULL UNIQUE,    user_id INT UNSIGNED NOT NULL UNIQUE,

    verification_method ENUM('birth_date', 'id_document', 'phone', 'email') NOT NULL,    verification_method ENUM('birth_date', 'id_document', 'phone', 'email') NOT NULL,

    verification_data JSON COMMENT '驗證資料 (加密儲存)',    verification_data JSON COMMENT '驗證資料 (加密儲存)',

    status ENUM('pending', 'approved', 'rejected', 'expired') DEFAULT 'pending',    status ENUM('pending', 'approved', 'rejected', 'expired') DEFAULT 'pending',

    verified_at TIMESTAMP NULL,    verified_at TIMESTAMP NULL,

    expires_at TIMESTAMP NULL,    expires_at TIMESTAMP NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

        

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_user_id (user_id),    INDEX idx_user_id (user_id),

    INDEX idx_status (status),    INDEX idx_status (status),

    INDEX idx_verified_at (verified_at),    INDEX idx_verified_at (verified_at),

    INDEX idx_expires_at (expires_at)    INDEX idx_expires_at (expires_at)

) ENGINE=InnoDB COMMENT='年齡驗證記錄';

) ENGINE=InnoDB COMMENT='年齡驗證記錄';

-- 7. 配對記錄表-- 7. 配對記錄表


CREATE TABLE IF NOT EXISTS matches (CREATE TABLE IF NOT EXISTS matches (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    user1_id INT UNSIGNED NOT NULL COMMENT '發起配對用戶',    user1_id INT UNSIGNED NOT NULL COMMENT '發起配對用戶',

    user2_id INT UNSIGNED NOT NULL COMMENT '被配對用戶',    user2_id INT UNSIGNED NOT NULL COMMENT '被配對用戶',

    user1_action ENUM('like', 'dislike', 'super_like', 'pending') DEFAULT 'pending',    user1_action ENUM('like', 'dislike', 'super_like', 'pending') DEFAULT 'pending',

    user2_action ENUM('like', 'dislike', 'super_like', 'pending') DEFAULT 'pending',    user2_action ENUM('like', 'dislike', 'super_like', 'pending') DEFAULT 'pending',

    status ENUM('pending', 'matched', 'unmatched', 'expired') DEFAULT 'pending',    status ENUM('pending', 'matched', 'unmatched', 'expired') DEFAULT 'pending',

    matched_at TIMESTAMP NULL COMMENT '配對成功時間',    matched_at TIMESTAMP NULL COMMENT '配對成功時間',

    expires_at TIMESTAMP NULL COMMENT '配對過期時間',    expires_at TIMESTAMP NULL COMMENT '配對過期時間',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

        

    UNIQUE KEY unique_match (user1_id, user2_id),    UNIQUE KEY unique_match (user1_id, user2_id),

    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,

    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_user1_id (user1_id),    INDEX idx_user1_id (user1_id),

    INDEX idx_user2_id (user2_id),    INDEX idx_user2_id (user2_id),

    INDEX idx_status (status),    INDEX idx_status (status),

    INDEX idx_matched_at (matched_at),    INDEX idx_matched_at (matched_at),

    INDEX idx_expires_at (expires_at),    INDEX idx_expires_at (expires_at),

        

    CONSTRAINT chk_different_users CHECK (user1_id != user2_id)    CONSTRAINT chk_different_users CHECK (user1_id != user2_id)

) ENGINE=InnoDB COMMENT='配對記錄';

) ENGINE=InnoDB COMMENT='配對記錄';

-- 8. 聊天訊息表-- 8. 聊天訊息表


CREATE TABLE IF NOT EXISTS chat_messages (CREATE TABLE IF NOT EXISTS chat_messages (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    match_id INT UNSIGNED NOT NULL COMMENT '配對ID',    match_id INT UNSIGNED NOT NULL COMMENT '配對ID',

    sender_id INT UNSIGNED NOT NULL COMMENT '發送者用戶ID',    sender_id INT UNSIGNED NOT NULL COMMENT '發送者用戶ID',

    content TEXT NOT NULL COMMENT '訊息內容',    content TEXT NOT NULL COMMENT '訊息內容',

    message_type ENUM('text', 'image', 'emoji', 'system') DEFAULT 'text',    message_type ENUM('text', 'image', 'emoji', 'system') DEFAULT 'text',

    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已讀',    is_read BOOLEAN DEFAULT FALSE COMMENT '是否已讀',

    is_deleted BOOLEAN DEFAULT FALSE COMMENT '是否已刪除',    is_deleted BOOLEAN DEFAULT FALSE COMMENT '是否已刪除',

    read_at TIMESTAMP NULL COMMENT '已讀時間',    read_at TIMESTAMP NULL COMMENT '已讀時間',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

        

    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,    FOREIGN KEY (match_id) REFERENCES matches(id) ON DELETE CASCADE,

    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_match_id (match_id),    INDEX idx_match_id (match_id),

    INDEX idx_sender_id (sender_id),    INDEX idx_sender_id (sender_id),

    INDEX idx_created_at (created_at),    INDEX idx_created_at (created_at),

    INDEX idx_is_read (is_read),    INDEX idx_is_read (is_read),

    INDEX idx_is_deleted (is_deleted)    INDEX idx_is_deleted (is_deleted)

) ENGINE=InnoDB COMMENT='聊天訊息';

) ENGINE=InnoDB COMMENT='聊天訊息';

-- 9. 檢舉記錄表-- 9. 檢舉記錄表


CREATE TABLE IF NOT EXISTS reports (CREATE TABLE IF NOT EXISTS reports (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    reporter_id INT UNSIGNED NOT NULL COMMENT '檢舉者用戶ID',    reporter_id INT UNSIGNED NOT NULL COMMENT '檢舉者用戶ID',

    reported_user_id INT UNSIGNED NULL COMMENT '被檢舉用戶ID',    reported_user_id INT UNSIGNED NULL COMMENT '被檢舉用戶ID',

    reported_message_id INT UNSIGNED NULL COMMENT '被檢舉訊息ID',    reported_message_id INT UNSIGNED NULL COMMENT '被檢舉訊息ID',

    report_type ENUM('inappropriate_content', 'harassment', 'spam', 'fake_profile', 'underage', 'other') NOT NULL,    report_type ENUM('inappropriate_content', 'harassment', 'spam', 'fake_profile', 'underage', 'other') NOT NULL,

    reason TEXT NOT NULL COMMENT '檢舉原因',    reason TEXT NOT NULL COMMENT '檢舉原因',

    status ENUM('pending', 'investigating', 'resolved', 'dismissed') DEFAULT 'pending',    status ENUM('pending', 'investigating', 'resolved', 'dismissed') DEFAULT 'pending',

    resolution_notes TEXT NULL COMMENT '處理說明',    resolution_notes TEXT NULL COMMENT '處理說明',

    resolved_by INT UNSIGNED NULL COMMENT '處理人員ID',    resolved_by INT UNSIGNED NULL COMMENT '處理人員ID',

    resolved_at TIMESTAMP NULL COMMENT '處理時間',    resolved_at TIMESTAMP NULL COMMENT '處理時間',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

        

    FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,

    FOREIGN KEY (reported_user_id) REFERENCES users(id) ON DELETE SET NULL,    FOREIGN KEY (reported_user_id) REFERENCES users(id) ON DELETE SET NULL,

    FOREIGN KEY (reported_message_id) REFERENCES chat_messages(id) ON DELETE SET NULL,    FOREIGN KEY (reported_message_id) REFERENCES chat_messages(id) ON DELETE SET NULL,

    INDEX idx_reporter_id (reporter_id),    INDEX idx_reporter_id (reporter_id),

    INDEX idx_reported_user_id (reported_user_id),    INDEX idx_reported_user_id (reported_user_id),

    INDEX idx_reported_message_id (reported_message_id),    INDEX idx_reported_message_id (reported_message_id),

    INDEX idx_report_type (report_type),    INDEX idx_report_type (report_type),

    INDEX idx_status (status),    INDEX idx_status (status),

    INDEX idx_created_at (created_at)    INDEX idx_created_at (created_at)

) ENGINE=InnoDB COMMENT='檢舉記錄';

) ENGINE=InnoDB COMMENT='檢舉記錄';

-- 10. 封鎖記錄表-- 10. 封鎖記錄表


CREATE TABLE IF NOT EXISTS blocks (CREATE TABLE IF NOT EXISTS blocks (

    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,

    blocker_id INT UNSIGNED NOT NULL COMMENT '封鎖者用戶ID',    blocker_id INT UNSIGNED NOT NULL COMMENT '封鎖者用戶ID',

    blocked_user_id INT UNSIGNED NOT NULL COMMENT '被封鎖用戶ID',    blocked_user_id INT UNSIGNED NOT NULL COMMENT '被封鎖用戶ID',

    reason TEXT NULL COMMENT '封鎖原因',    reason TEXT NULL COMMENT '封鎖原因',

    is_active BOOLEAN DEFAULT TRUE COMMENT '封鎖是否生效',    is_active BOOLEAN DEFAULT TRUE COMMENT '封鎖是否生效',

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

        

    UNIQUE KEY unique_block (blocker_id, blocked_user_id),    UNIQUE KEY unique_block (blocker_id, blocked_user_id),

    FOREIGN KEY (blocker_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (blocker_id) REFERENCES users(id) ON DELETE CASCADE,

    FOREIGN KEY (blocked_user_id) REFERENCES users(id) ON DELETE CASCADE,    FOREIGN KEY (blocked_user_id) REFERENCES users(id) ON DELETE CASCADE,

    INDEX idx_blocker_id (blocker_id),    INDEX idx_blocker_id (blocker_id),

    INDEX idx_blocked_user_id (blocked_user_id),    INDEX idx_blocked_user_id (blocked_user_id),

    INDEX idx_is_active (is_active),    INDEX idx_is_active (is_active),

        

    CONSTRAINT chk_different_block_users CHECK (blocker_id != blocked_user_id)    CONSTRAINT chk_different_block_users CHECK (blocker_id != blocked_user_id)

) ENGINE=InnoDB COMMENT='封鎖記錄';

) ENGINE=InnoDB COMMENT='封鎖記錄';

-- 預設興趣標籤數據-- 預設興趣標籤數據

INSERT IGNORE INTO
    interests (name, category) VALUESINSERT IGNORE INTO interests (name, category)
VALUES

-- 運動類-- 運動類

('健身', 'sports'),
('健身', 'sports'),
('跑步', 'sports'),
('跑步', 'sports'),
('游泳', 'sports'),
('游泳', 'sports'),
('瑜伽', 'sports'),
('瑜伽', 'sports'),
('登山', 'sports'),
('登山', 'sports'),
('籃球', 'sports'),
('籃球', 'sports'),
('足球', 'sports'),
('足球', 'sports'),
('網球', 'sports'),
('網球', 'sports'),

-- 藝術文化-- 藝術文化

('音樂', 'arts'),
('音樂', 'arts'),
('畫畫', 'arts'),
('畫畫', 'arts'),
('攝影', 'arts'),
('攝影', 'arts'),
('舞蹈', 'arts'),
('舞蹈', 'arts'),
('電影', 'arts'),
('電影', 'arts'),
('閱讀', 'arts'),
('閱讀', 'arts'),
('戲劇', 'arts'),
('戲劇', 'arts'),
('博物館', 'arts'),
('博物館', 'arts'),

-- 美食-- 美食

('烹飪', 'food'),
('烹飪', 'food'),
('品酒', 'food'),
('品酒', 'food'),
('咖啡', 'food'),
('咖啡', 'food'),
('甜點', 'food'),
('甜點', 'food'),
('素食', 'food'),
('素食', 'food'),
('異國料理', 'food'),
('異國料理', 'food'),

-- 旅遊-- 旅遊

('旅遊', 'travel'),
('旅遊', 'travel'),
('露營', 'travel'),
('露營', 'travel'),
('海邊', 'travel'),
('海邊', 'travel'),
('山區', 'travel'),
('山區', 'travel'),
('城市探索', 'travel'),
('城市探索', 'travel'),
('背包旅行', 'travel'),
('背包旅行', 'travel'),

-- 科技-- 科技

('程式設計', 'technology'),
('程式設計', 'technology'),
('遊戲', 'technology'),
('遊戲', 'technology'),
('科技產品', 'technology'),
('科技產品', 'technology'),
('人工智慧', 'technology'),
('人工智慧', 'technology'),

-- 生活-- 生活

('寵物', 'lifestyle'),
('寵物', 'lifestyle'),
('園藝', 'lifestyle'),
('園藝', 'lifestyle'),
('手作', 'lifestyle'),
('手作', 'lifestyle'),
('收集', 'lifestyle'),
('收集', 'lifestyle'),
('志工服務', 'lifestyle'),
('志工服務', 'lifestyle'),
('冥想', 'lifestyle');

('冥想', 'lifestyle');

-- 創建管理員用戶 (測試用)-- 創建管理員用戶 (測試用)

INSERT IGNORE INTO
    users (
        email,
        password_hash,
        birth_date,
        is_verified,
        is_active
    ) VALUESINSERT IGNORE INTO users (
        email,
        password_hash,
        birth_date,
        is_verified,
        is_active
    )
VALUES (
        'admin@dating.app',
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPFPMB7xOQPhm',
        '1990-01-01',
        TRUE,
        TRUE
    );

(
    'admin@dating.app',
    '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPFPMB7xOQPhm',
    '1990-01-01',
    TRUE,
    TRUE
);

SET @admin_user_id = LAST_INSERT_ID();

SET @admin_user_id = LAST_INSERT_ID();

INSERT IGNORE INTO
    user_profiles (
        user_id,
        display_name,
        bio,
        gender,
        show_age
    ) VALUESINSERT IGNORE INTO user_profiles (
        user_id,
        display_name,
        bio,
        gender,
        show_age
    )
VALUES (
        @admin_user_id,
        '系統管理員',
        '系統管理員帳戶',
        'other',
        FALSE
    );

( @admin_user_id, '系統管理員', '系統管理員帳戶', 'other', FALSE );

personality_type VARCHAR(50),
    zodiac VARCHAR(50),
    religion VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_profile (user_id),
    INDEX idx_location (location_lat, location_lng),
    INDEX idx_city (city),
    INDEX idx_country (country)
);

-- 建立用戶照片表
CREATE TABLE IF NOT EXISTS user_photos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    photo_url VARCHAR(500) NOT NULL,
    thumbnail_url VARCHAR(500),
    is_primary BOOLEAN DEFAULT FALSE,
    `order` INT DEFAULT 0,
    status ENUM(
        'pending',
        'approved',
        'rejected'
    ) DEFAULT 'pending',
    caption TEXT,
    is_verified BOOLEAN DEFAULT FALSE,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_is_primary (is_primary),
    INDEX idx_status (status),
    INDEX idx_order (`order`)
);

-- 建立用戶偏好設定表
CREATE TABLE IF NOT EXISTS user_preferences (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    preferred_gender ENUM('male', 'female', 'other'),
    age_min INT,
    age_max INT,
    distance_max INT,
    height_min INT,
    height_max INT,
    education JSON,
    interests JSON,
    lifestyle JSON,
    show_me BOOLEAN DEFAULT TRUE,
    show_distance BOOLEAN DEFAULT TRUE,
    show_age BOOLEAN DEFAULT TRUE,
    show_last_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_preference (user_id)
);

-- 建立聊天訊息資料表（可選，用於儲存聊天記錄）
CREATE TABLE IF NOT EXISTS chat_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT,
    username VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL,
    INDEX idx_created_at (created_at),
    INDEX idx_user_id (user_id)
);

-- 建立聊天室資料表（可選，用於多聊天室功能）
CREATE TABLE IF NOT EXISTS chat_rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_public BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE SET NULL,
    INDEX idx_name (name)
);

-- 建立使用者聊天室關聯表（可選）
CREATE TABLE IF NOT EXISTS user_chat_rooms (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    room_id INT NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role ENUM('member', 'admin', 'owner') DEFAULT 'member',
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (room_id) REFERENCES chat_rooms (id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_room (user_id, room_id)
);

-- 插入一些示例資料
INSERT INTO
    users (username, email, password)
VALUES (
        'admin',
        'admin@example.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7J.YY8vC8.7dQ7dvAv7L5M5H5sHQyNW'
    ),
    (
        'demo_user',
        'demo@example.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7J.YY8vC8.7dQ7dvAv7L5M5H5sHQyNW'
    )
ON DUPLICATE KEY UPDATE
    username = username;

-- 插入預設聊天室
INSERT INTO
    chat_rooms (name, description, created_by)
VALUES ('一般聊天室', '歡迎來到一般聊天室！', 1),
    ('技術討論', '討論技術相關話題', 1)
ON DUPLICATE KEY UPDATE
    name = name;