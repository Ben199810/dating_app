# 快速開始指南：18+ 交友聊天應用程式

**Phase 1 輸出** - 開發者快速上手指南

## 概述

本指南將幫助開發者快速建立和執行 18+ 交友聊天應用程式的開發環境，並驗證核心功能。

### 系統需求

- Go 1.23+
- Docker & Docker Compose
- Git
- 現代網頁瀏覽器（支援 WebSocket）

### 核心功能

1. **用戶註冊登入** - 18+ 年齡驗證
2. **個人檔案管理** - 照片、興趣、基本資料
3. **配對系統** - 雙向配對機制
4. **即時聊天** - WebSocket 即時通訊
5. **安全功能** - 檢舉、封鎖機制

---

## 快速安裝

### 1. 克隆專案

```bash
git clone <repository-url>
cd golang_dev_docker
```

### 2. 環境設定

```bash
# 複製環境配置（如果需要）
cp config/development.yaml.example config/development.yaml

# 確保 Docker 服務啟動
docker --version
docker-compose --version
```

### 3. 啟動開發環境

```bash
# 使用 Docker Compose 啟動所有服務
docker compose -f build/docker-compose.yaml up -d

# 檢查服務狀態
docker compose -f build/docker-compose.yaml ps
```

### 4. 驗證安裝

訪問以下 URL 確認服務正常：

- **主應用程式**: <http://localhost:8080>
- **API 健康檢查**: <http://localhost:8080/api/health>
- **phpMyAdmin**: <http://localhost:8081>

---

## 核心功能測試

### 測試場景 1：用戶註冊與登入

#### 1.1 註冊新用戶（成功案例）

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePassword123",
    "birth_date": "1995-06-15",
    "display_name": "John",
    "gender": "male"
  }'
```

**預期結果**: HTTP 201，返回用戶 ID

#### 1.2 年齡驗證測試（失敗案例）

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "young@example.com",
    "password": "SecurePassword123",
    "birth_date": "2010-01-01",
    "display_name": "TooYoung",
    "gender": "male"
  }'
```

**預期結果**: HTTP 400，錯誤訊息 "必須年滿18歲才能註冊"

#### 1.3 用戶登入

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "SecurePassword123"
  }'
```

**預期結果**: HTTP 200，返回 JWT token

### 測試場景 2：個人檔案管理

#### 2.1 查看個人資料

```bash
# 使用從登入獲得的 token
TOKEN="your-jwt-token-here"

curl -X GET http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer $TOKEN"
```

#### 2.2 更新個人資料

```bash
curl -X PUT http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "display_name": "John Updated",
    "bio": "Hello, I am looking for meaningful connections",
    "show_age": true,
    "max_distance": 25,
    "age_range_min": 25,
    "age_range_max": 35,
    "interests": [1, 2, 3]
  }'
```

### 測試場景 3：配對系統

#### 3.1 註冊第二個用戶

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "jane@example.com",
    "password": "SecurePassword456",
    "birth_date": "1992-08-20",
    "display_name": "Jane",
    "gender": "female"
  }'
```

#### 3.2 探索配對對象

```bash
# 使用 Jane 的 token
JANE_TOKEN="jane-jwt-token-here"

curl -X GET "http://localhost:8080/api/matches/discover?limit=10" \
  -H "Authorization: Bearer $JANE_TOKEN"
```

#### 3.3 表示興趣（Jane 對 John）

```bash
curl -X POST http://localhost:8080/api/matches/like \
  -H "Authorization: Bearer $JANE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "target_user_id": 1
  }'
```

**預期結果**: 創建 pending 配對記錄

#### 3.4 查看待處理配對（John）

```bash
curl -X GET http://localhost:8080/api/matches/pending \
  -H "Authorization: Bearer $TOKEN"
```

#### 3.5 回應配對（John 接受 Jane）

```bash
curl -X POST http://localhost:8080/api/matches/1/respond \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "action": "accept"
  }'
```

**預期結果**: 配對成功，狀態變為 "matched"

### 測試場景 4：即時聊天

#### 4.1 取得聊天列表

```bash
curl -X GET http://localhost:8080/api/chats \
  -H "Authorization: Bearer $TOKEN"
```

#### 4.2 發送聊天訊息

```bash
curl -X POST http://localhost:8080/api/chats/1/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hi Jane, nice to meet you!",
    "message_type": "text"
  }'
```

#### 4.3 取得聊天記錄

```bash
curl -X GET "http://localhost:8080/api/chats/1/messages?limit=50" \
  -H "Authorization: Bearer $JANE_TOKEN"
```

#### 4.4 WebSocket 測試

使用瀏覽器開發者工具或 WebSocket 測試工具：

```javascript
// 在瀏覽器 Console 中執行
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = function(event) {
    console.log('WebSocket 連接成功');
    // 發送認證訊息
    ws.send(JSON.stringify({
        type: 'auth',
        data: { token: 'your-jwt-token' }
    }));
};

ws.onmessage = function(event) {
    console.log('收到訊息:', JSON.parse(event.data));
};

// 發送測試訊息
ws.send(JSON.stringify({
    type: 'chat_message',
    data: {
        match_id: 1,
        content: 'WebSocket test message'
    }
}));
```

### 測試場景 5：安全功能

#### 5.1 檢舉用戶

```bash
curl -X POST http://localhost:8080/api/reports \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "reported_user_id": 2,
    "report_type": "inappropriate",
    "reason": "不當內容測試"
  }'
```

#### 5.2 封鎖用戶

```bash
curl -X POST http://localhost:8080/api/blocks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "blocked_user_id": 2
  }'
```

#### 5.3 驗證封鎖效果

```bash
# 被封鎖的用戶不應出現在配對推薦中
curl -X GET "http://localhost:8080/api/matches/discover" \
  -H "Authorization: Bearer $TOKEN"
```

---

## 前端測試

### 1. 訪問網頁介面

打開瀏覽器訪問：<http://localhost:8080>

### 2. 註冊新帳戶

1. 點擊「註冊」按鈕
2. 填寫表單（確保年齡 >= 18）
3. 提交並驗證是否成功

### 3. 登入測試

1. 使用註冊的帳戶登入
2. 驗證是否跳轉到個人檔案頁面

### 4. 個人檔案設定

1. 上傳個人照片
2. 填寫基本資料和興趣
3. 設定配對偏好

### 5. 配對功能測試

1. 進入配對頁面
2. 測試滑動功能（喜歡/跳過）
3. 驗證配對通知

### 6. 聊天功能測試

1. 找到配對成功的用戶
2. 進入聊天室
3. 發送訊息測試即時性
4. 測試多標籤頁同步

---

## 故障排除

### 常見問題

#### 1. 資料庫連接失敗

```bash
# 檢查 MySQL 容器狀態
docker compose -f build/docker-compose.yaml logs mysql

# 重啟資料庫服務
docker compose -f build/docker-compose.yaml restart mysql
```

#### 2. WebSocket 連接失敗

```bash
# 檢查應用程式日誌
docker compose -f build/docker-compose.yaml logs app

# 檢查防火牆設定
netstat -tlnp | grep :8080
```

#### 3. JWT Token 無效

- 確認 token 格式正確
- 檢查 token 是否過期
- 驗證 Authorization header 格式

#### 4. 年齡驗證失敗

- 確認生日格式為 "YYYY-MM-DD"
- 計算年齡是否 >= 18
- 檢查時區設定

### 日誌查看

```bash
# 查看所有服務日誌
docker compose -f build/docker-compose.yaml logs

# 查看特定服務日誌
docker compose -f build/docker-compose.yaml logs app
docker compose -f build/docker-compose.yaml logs mysql

# 即時查看日誌
docker compose -f build/docker-compose.yaml logs -f app
```

### 重置環境

```bash
# 停止並移除所有容器和資料
docker compose -f build/docker-compose.yaml down -v

# 重新編譯和啟動
docker compose -f build/docker-compose.yaml build --no-cache
docker compose -f build/docker-compose.yaml up -d
```

---

## 下一步

### 開發建議

1. **設定開發環境**：配置 IDE 和調試工具
2. **閱讀程式碼**：理解 DDD 架構和專案結構
3. **執行測試**：`go test ./...` 確保所有測試通過
4. **查看 API 文件**：`contracts/api-spec.yaml` 詳細規格
5. **參考憲章**：`.specify/memory/constitution.md` 開發原則

### 擴展功能

- 照片上傳和處理
- 推播通知
- 地理位置服務
- 進階配對演算法
- 管理後台

### 效能優化

- Redis 快取策略
- 資料庫索引優化
- CDN 靜態資源
- 負載平衡

這個快速開始指南涵蓋了所有核心功能的測試場景，確保開發者可以快速驗證系統的完整性和正確性。
