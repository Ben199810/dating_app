-- 建立使用者資料表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    age INT,
    gender ENUM('male', 'female', 'other'),
    is_verified BOOLEAN DEFAULT FALSE,
    status ENUM(
        'active',
        'inactive',
        'banned'
    ) DEFAULT 'active',
    last_active_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_gender (gender),
    INDEX idx_age (age),
    INDEX idx_status (status),
    INDEX idx_last_active_at (last_active_at)
);

-- 建立用戶詳細資料表
CREATE TABLE IF NOT EXISTS user_profiles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    bio TEXT,
    interests JSON,
    location_lat DECIMAL(10, 8),
    location_lng DECIMAL(11, 8),
    city VARCHAR(100),
    country VARCHAR(100),
    height INT,
    weight INT,
    education VARCHAR(255),
    occupation VARCHAR(255),
    company VARCHAR(255),
    relationship VARCHAR(100),
    looking_for JSON,
    languages JSON,
    hobbies JSON,
    lifestyle JSON,
    pet_preference VARCHAR(100),
    drinking_habit VARCHAR(100),
    smoking_habit VARCHAR(100),
    exercise_habit VARCHAR(100),
    social_media_link VARCHAR(500),
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