-- 建立使用者資料表
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    avatar_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    INDEX idx_username (username),
    INDEX idx_email (email)
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
    users (
        username,
        email,
        password_hash,
        display_name
    )
VALUES (
        'admin',
        'admin@example.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7J.YY8vC8.7dQ7dvAv7L5M5H5sHQyNW',
        '管理員'
    ),
    (
        'demo_user',
        'demo@example.com',
        '$2a$10$N9qo8uLOickgx2ZMRZoMye7J.YY8vC8.7dQ7dvAv7L5M5H5sHQyNW',
        '演示使用者'
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