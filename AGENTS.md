# 🚀 快速啟動

## 1. 確保 Docker 服務已啟動

```bash
# 檢查 Docker 狀態
docker info

# 如果 Docker 未啟動，請開啟 Docker Desktop
open /Applications/Docker.app
```

## 2. 編譯 Docker 映像檔

```bash
# 編譯所有服務的映像檔
docker compose -f build/docker-compose.yaml build

# 強制重新編譯（不使用快取）
docker compose -f build/docker-compose.yaml build --no-cache
```

## 3. 啟動開發環境

```bash
# 在背景啟動所有服務
docker compose -f build/docker-compose.yaml up -d

# 或者在前台啟動（可以看到即時日誌）
docker compose -f build/docker-compose.yaml up
```

## 📋 服務說明

啟動後將包含以下服務：

| 服務 | 端口 | 說明 |
|------|------|------|
| **Go 應用程式** | 8080 | 主要的 Web 應用程式 |
| **MySQL 資料庫** | 3306 | 資料庫服務 |
| **phpMyAdmin** | 8081 | 資料庫管理介面 |

## 🔧 常用指令

## 檢查服務狀態

```bash
# 查看運行中的容器
docker compose -f build/docker-compose.yaml ps

# 查看應用程式日誌
docker compose -f build/docker-compose.yaml logs app

# 查看所有服務日誌
docker compose -f build/docker-compose.yaml logs
```

## 停止和重啟服務

```bash
# 停止所有服務
docker compose -f build/docker-compose.yaml down

# 停止並移除所有資料（包括資料庫）
docker compose -f build/docker-compose.yaml down -v

# 重啟特定服務
docker compose -f build/docker-compose.yaml restart app

# 關閉本機開發環境後，清除已無使用的 image 和 volume 釋放本機腦容量
docker system prune -a --volumes
docker image prune -a
```

## 🌐 存取應用程式

啟動成功後，您可以透過以下網址存取：

- **主應用程式**: <http://localhost:8080>
- **聊天室**: <http://localhost:8080/chat>
- **API 健康檢查**: <http://localhost:8080/api/status>
- **phpMyAdmin**: <http://localhost:8081>

## 🐛 故障排除

## 常見問題

1. **Docker daemon 連線錯誤**

   ```bash
   # 確保 Docker Desktop 正在運行
   docker info
   ```

2. **端口已被占用**

   ```bash
   # 檢查端口使用情況
   lsof -i :8080
   lsof -i :3306
   lsof -i :8081
   # 停止占用端口的程序或修改 docker-compose.yaml 中的端口配置
   ```

3. **資料庫連線失敗**

   ```bash
   # 檢查 MySQL 容器是否健康
   docker compose -f build/docker-compose.yaml ps
   
   # 查看 MySQL 日誌
   docker compose -f build/docker-compose.yaml logs mysql
   ```

4. **程式碼變更未生效**

   ```bash
   # 重新編譯並啟動
   docker compose -f build/docker-compose.yaml down
   docker compose -f build/docker-compose.yaml build --no-cache
   docker compose -f build/docker-compose.yaml up -d
   ```

## ⚙️ 環境配置

應用程式使用 YAML 配置檔案管理不同環境：

本機開發環境會對應以下配置檔案：

- **開發環境**: `config/development.docker.yaml`

### 雲端管理平台會對應以下兩個環境配置檔案

- **生產環境**: `config/production.yaml`
- **測試環境**: `config/test.yaml`

Docker 環境預設使用 `production.yaml` 配置，可以通過修改 `docker-compose.yaml` 中的 `APP_ENV` 環境變數來變更。

## 📝 開發流程

1. **修改程式碼** → 在本機編輯器中修改
2. **重新編譯** → `docker compose -f build/docker-compose.yaml build`
3. **重啟服務** → `docker compose -f build/docker-compose.yaml up -d`
4. **測試應用** → 訪問 <http://localhost:8080>

## 🔗 相關資源

- [Docker Compose 文件](https://docs.docker.com/compose/)
- [Go 官方文件](https://golang.org/doc/)
- [Gin 框架文件](https://gin-gonic.com/docs/)
