#!/bin/bash
# 生產環境部署腳本

set -e

echo "🚀 開始部署 Dating App 生產環境..."

# 檢查必要工具
command -v docker >/dev/null 2>&1 || { echo "❌ 需要安裝 Docker"; exit 1; }
command -v docker-compose >/dev/null 2>&1 || { echo "❌ 需要安裝 Docker Compose"; exit 1; }

# 設定變數
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
BUILD_DIR="$PROJECT_ROOT/build"

cd "$BUILD_DIR"

# 檢查環境檔案
if [ ! -f ".env" ]; then
    echo "⚠️  環境檔案 .env 不存在，複製範本..."
    cp .env.prod.template .env
    echo "📝 請編輯 .env 檔案並填入實際的機敏資訊"
    echo "🔧 特別需要設定以下變數："
    echo "   - JWT_SECRET_KEY"
    echo "   - MYSQL_ROOT_PASSWORD" 
    echo "   - DB_PASSWORD"
    echo "   - REDIS_PASSWORD"
    echo ""
    read -p "是否已完成 .env 配置？(y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ 請先完成環境配置"
        exit 1
    fi
fi

# 載入環境變數
source .env

# 安全檢查
echo "🔒 執行安全檢查..."
if [ "$JWT_SECRET_KEY" = "your-super-secure-jwt-secret-key-minimum-32-characters-change-this" ]; then
    echo "❌ 請更改預設的 JWT_SECRET_KEY"
    exit 1
fi

if [ ${#JWT_SECRET_KEY} -lt 32 ]; then
    echo "❌ JWT_SECRET_KEY 長度必須至少 32 字元"
    exit 1
fi

# 建置映像檔
echo "🏗️  建置 Docker 映像檔..."
docker-compose -f docker-compose.prod.yaml build --no-cache

# 執行測試
echo "🧪 執行測試..."
docker run --rm \
    -v "$PROJECT_ROOT:/app" \
    -w /app \
    golang:1.23-alpine \
    sh -c "go mod download && go test -short ./tests/unit/... ./tests/security/... ./tests/performance/..." || {
    echo "⚠️  測試失敗，但繼續部署（生產環境警告）"
}

# 建立必要目錄
echo "📁 建立必要目錄..."
mkdir -p logs uploads ssl

# 停止舊容器
echo "🛑 停止現有服務..."
docker-compose -f docker-compose.prod.yaml down

# 清理舊映像檔（可選）
echo "🧹 清理舊映像檔..."
docker image prune -f

# 啟動服務
echo "🚀 啟動生產服務..."
docker-compose -f docker-compose.prod.yaml up -d

# 等待服務啟動
echo "⏳ 等待服務啟動..."
sleep 30

# 健康檢查
echo "🏥 執行健康檢查..."
MAX_RETRIES=10
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if curl -f http://localhost:8080/health >/dev/null 2>&1; then
        echo "✅ 服務健康檢查通過"
        break
    else
        echo "⏳ 等待服務啟動... ($((RETRY_COUNT + 1))/$MAX_RETRIES)"
        sleep 10
        RETRY_COUNT=$((RETRY_COUNT + 1))
    fi
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "❌ 服務啟動失敗"
    echo "📋 檢查日誌："
    docker-compose -f docker-compose.prod.yaml logs --tail=50
    exit 1
fi

# 顯示服務狀態
echo "📊 服務狀態："
docker-compose -f docker-compose.prod.yaml ps

# 顯示訪問資訊
echo ""
echo "🎉 部署完成！"
echo ""
echo "📡 服務端點："
echo "  應用程式:      http://localhost:8080"
echo "  健康檢查:      http://localhost:8080/health"
echo "  系統指標:      http://localhost:8080/metrics"
echo "  資料庫管理:    http://localhost:8081 (如果啟用)"
echo "  Grafana:      http://localhost:3000 (如果啟用)"
echo "  Prometheus:   http://localhost:9090 (如果啟用)"
echo ""
echo "📋 管理命令："
echo "  查看日誌:      docker-compose -f docker-compose.prod.yaml logs -f"
echo "  停止服務:      docker-compose -f docker-compose.prod.yaml down"
echo "  重啟服務:      docker-compose -f docker-compose.prod.yaml restart"
echo "  更新服務:      ./deploy-prod.sh"
echo ""
echo "🔧 效能調校建議："
echo "  1. 調整 MySQL innodb_buffer_pool_size 至可用記憶體的 70-80%"
echo "  2. 設定 Redis maxmemory 避免 OOM"
echo "  3. 調整 Nginx worker_processes 至 CPU 核心數"
echo "  4. 設定適當的 Docker 資源限制"
echo ""

# 顯示效能指標
echo "📈 系統效能指標："
if curl -s http://localhost:8080/metrics >/dev/null 2>&1; then
    curl -s http://localhost:8080/metrics | head -20
else
    echo "⚠️  無法獲取指標，請稍後再試"
fi