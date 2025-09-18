# golang_dev_docker

以golang為例建立docker的本機開發環境

## 配置管理

本專案使用 YAML 檔案管理不同環境的配置設定。

### 配置檔案

- `config/development.yaml` - 開發環境配置
- `config/production.yaml` - 生產環境配置  
- `config/test.yaml` - 測試環境配置

### 環境變數

設定 `APP_ENV` 環境變數來指定要載入的配置檔案：

```bash
export APP_ENV=development  # 載入 development.yaml
export APP_ENV=production   # 載入 production.yaml
export APP_ENV=test        # 載入 test.yaml
```

如果未設定 `APP_ENV`，預設會載入 `development.yaml`。

### 配置結構

```yaml
database:
  host: localhost
  port: 3306
  user: chat_user
  password: chat_password
  dbname: chat_app
  charset: utf8mb4
  parseTime: true
  loc: Local

server:
  port: 8080
  mode: debug  # gin 模式: debug, release, test

logging:
  level: info
  format: json
```

## 執行應用程式

```bash
# 使用預設環境 (development)
go run .

# 指定特定環境
APP_ENV=production go run .
```
