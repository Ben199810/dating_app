# 手動測試場景 - Dating App 綜合測試指南

# Manual Testing Scenarios - Comprehensive Testing Guide

## 測試環境設定

### Environment Setup

**前置條件 Prerequisites:**

- Docker 與 Docker Compose 已安裝
- 測試資料庫已重置
- 應用程式運行於 localhost:8080
- 至少2個測試用戶帳號可用

**測試資料準備 Test Data Preparation:**

```bash
# 啟動測試環境
docker compose -f build/docker-compose.yaml up -d

# 重置測試資料庫
docker compose exec mysql mysql -u root -p -e "DROP DATABASE IF EXISTS dating_app; CREATE DATABASE dating_app;"

# 啟動應用程式
go run .
```

---

## 測試場景 1: 用戶註冊與驗證

### Test Scenario 1: User Registration and Verification

### 🎯 測試目標 Test Objectives

驗證用戶註冊流程、年齡驗證、及帳號安全性

### 📋 測試步驟 Test Steps

#### TC001: 成功註冊測試

**步驟 Steps:**

1. 開啟瀏覽器訪問 `http://localhost:8080`
2. 點擊「註冊」按鈕
3. 填寫有效的註冊表單：
   - 電子郵件: `testuser1@example.com`
   - 密碼: `SecurePass123!`
   - 確認密碼: `SecurePass123!`
   - 姓名: `測試用戶1`
   - 生日: `1995-06-15` (確保年滿18歲)
   - 性別: `male`
4. 提交表單

**預期結果 Expected Results:**

- ✅ 註冊成功訊息顯示
- ✅ 重定向至個人資料設定頁面
- ✅ 資料庫中建立新用戶記錄
- ✅ 密碼已加密儲存

#### TC002: 年齡驗證失敗測試

**步驟 Steps:**

1. 嘗試註冊未滿18歲用戶
2. 生日設定為 `2010-01-01`
3. 填寫其他有效資訊並提交

**預期結果 Expected Results:**

- ❌ 顯示年齡不符合要求錯誤
- ❌ 註冊被拒絕
- ✅ 錯誤訊息清楚說明年齡限制

#### TC003: 重複電子郵件測試

**步驟 Steps:**

1. 使用已註冊的電子郵件嘗試再次註冊
2. 電子郵件: `testuser1@example.com`

**預期結果 Expected Results:**

- ❌ 顯示「電子郵件已存在」錯誤
- ❌ 註冊被拒絕

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] 用戶資料正確儲存至資料庫
- [ ] 密碼哈希化處理正確
- [ ] 年齡驗證邏輯運作正常
- [ ] 錯誤處理訊息用戶友善

---

## 測試場景 2: 用戶認證與登入

### Test Scenario 2: User Authentication and Login

### 🎯 測試目標 Test Objectives

驗證登入機制、JWT token 生成、會話管理

### 📋 測試步驟 Test Steps

#### TC004: 成功登入測試

**步驟 Steps:**

1. 訪問登入頁面 `/auth/login`
2. 輸入有效憑證：
   - 電子郵件: `testuser1@example.com`
   - 密碼: `SecurePass123!`
3. 點擊「登入」按鈕

**預期結果 Expected Results:**

- ✅ 成功登入訊息
- ✅ JWT token 在 cookie 中設定
- ✅ 重定向至主頁面
- ✅ 用戶狀態顯示為已登入

#### TC005: 登入失敗測試

**步驟 Steps:**

1. 使用錯誤密碼嘗試登入
   - 電子郵件: `testuser1@example.com`
   - 密碼: `WrongPassword`
2. 提交登入表單

**預期結果 Expected Results:**

- ❌ 顯示「憑證無效」錯誤
- ❌ 未設定 JWT token
- ❌ 保持在登入頁面

#### TC006: JWT Token 驗證測試

**步驟 Steps:**

1. 成功登入後查看瀏覽器開發者工具
2. 檢查 Cookies 中的 JWT token
3. 訪問需要認證的頁面 `/user/profile`

**預期結果 Expected Results:**

- ✅ JWT token 存在於 cookie
- ✅ Token 格式正確 (3段式結構)
- ✅ 可以訪問受保護頁面
- ✅ Token 包含正確的用戶資訊

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] JWT token 正確生成並設定
- [ ] 認證中間件正常運作
- [ ] 會話超時處理正確
- [ ] 安全標頭設定適當

---

## 測試場景 3: 個人資料管理

### Test Scenario 3: Profile Management

### 🎯 測試目標 Test Objectives

驗證個人資料設定、照片上傳、興趣管理功能

### 📋 測試步驟 Test Steps

#### TC007: 完整個人資料設定

**步驟 Steps:**

1. 登入後訪問 `/user/profile`
2. 填寫完整個人資料：
   - 職業: `軟體工程師`
   - 教育: `大學`
   - 身高: `175`
   - 個人簡介: `喜歡旅行和程式設計的工程師`
   - 興趣: 選擇3-5個興趣標籤
3. 保存設定

**預期結果 Expected Results:**

- ✅ 資料成功保存
- ✅ 顯示成功更新訊息
- ✅ 資料庫正確更新
- ✅ 個人資料完整度提升

#### TC008: 照片上傳測試

**步驟 Steps:**

1. 在個人資料頁面點擊「上傳照片」
2. 選擇有效的圖片檔案 (JPG, PNG, WebP)
3. 確保檔案小於設定的大小限制
4. 上傳照片

**預期結果 Expected Results:**

- ✅ 照片成功上傳
- ✅ 圖片正確顯示在預覽區
- ✅ 檔案儲存至指定目錄
- ✅ 資料庫記錄照片資訊

#### TC009: 照片上傳限制測試

**步驟 Steps:**

1. 嘗試上傳超大檔案 (>10MB)
2. 嘗試上傳不支援的檔案格式 (如 PDF)
3. 嘗試上傳超過6張照片

**預期結果 Expected Results:**

- ❌ 大檔案被拒絕並顯示錯誤訊息
- ❌ 不支援格式被拒絕
- ❌ 超過數量限制的照片被拒絕

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] 個人資料驗證規則正確
- [ ] 照片上傳安全性檢查
- [ ] 檔案儲存路徑正確
- [ ] 資料完整性約束有效

---

## 測試場景 4: 配對系統測試

### Test Scenario 4: Matching System Testing

### 🎯 測試目標 Test Objectives

驗證配對演算法、用戶篩選、配對結果準確性

### 📋 測試步驟 Test Steps

#### TC010: 基本配對功能測試

**步驟 Steps:**

1. 確保有至少2個完整個人資料的用戶
2. 用戶1登入並訪問配對頁面 `/matching`
3. 查看推薦的配對用戶
4. 對推薦用戶進行「喜歡」或「不喜歡」操作

**預期結果 Expected Results:**

- ✅ 顯示符合篩選條件的用戶
- ✅ 用戶資訊完整顯示（照片、基本資料）
- ✅ 滑動操作正常運作
- ✅ 配對狀態正確記錄

#### TC011: 雙向配對測試

**步驟 Steps:**

1. 用戶1對用戶2按「喜歡」
2. 切換至用戶2帳號
3. 用戶2對用戶1也按「喜歡」
4. 檢查配對結果

**預期結果 Expected Results:**

- ✅ 產生雙向配對 (Match)
- ✅ 兩位用戶都收到配對通知
- ✅ 配對記錄儲存至資料庫
- ✅ 啟用聊天功能

#### TC012: 配對篩選條件測試

**步驟 Steps:**

1. 設定配對篩選條件：
   - 年齡範圍: 25-35
   - 距離: 50公里內
   - 有照片的用戶
2. 執行配對搜尋
3. 檢查結果是否符合篩選條件

**預期結果 Expected Results:**

- ✅ 只顯示符合年齡範圍的用戶
- ✅ 距離計算正確
- ✅ 只顯示有照片的用戶
- ✅ 篩選邏輯運作正常

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] 配對演算法邏輯正確
- [ ] 地理位置計算準確
- [ ] 配對狀態管理完整
- [ ] 用戶偏好設定有效

---

## 測試場景 5: 即時聊天系統

### Test Scenario 5: Real-time Chat System

### 🎯 測試目標 Test Objectives

驗證 WebSocket 連線、即時訊息傳送、聊天室管理

### 📋 測試步驟 Test Steps

#### TC013: WebSocket 連線測試

**步驟 Steps:**

1. 開啟兩個瀏覽器視窗
2. 分別登入已配對的兩個用戶
3. 進入聊天頁面 `/chat`
4. 檢查 WebSocket 連線狀態

**預期結果 Expected Results:**

- ✅ WebSocket 連線成功建立
- ✅ 連線狀態顯示為「已連線」
- ✅ 沒有連線錯誤訊息
- ✅ 用戶上線狀態正確顯示

#### TC014: 即時訊息傳送測試

**步驟 Steps:**

1. 用戶1傳送訊息: "Hello, nice to meet you!"
2. 檢查用戶2是否即時收到訊息
3. 用戶2回覆: "Hi there! How are you?"
4. 檢查用戶1是否即時收到回覆

**預期結果 Expected Results:**

- ✅ 訊息即時顯示（延遲 < 1秒）
- ✅ 訊息格式正確（時間戳、發送者）
- ✅ 訊息儲存至資料庫
- ✅ 雙向通訊正常

#### TC015: 聊天歷史記錄測試

**步驟 Steps:**

1. 傳送多條測試訊息
2. 關閉瀏覽器視窗
3. 重新登入並進入聊天室
4. 檢查訊息歷史是否保存

**預期結果 Expected Results:**

- ✅ 歷史訊息完整保存
- ✅ 訊息順序正確
- ✅ 時間戳正確顯示
- ✅ 載入速度合理

#### TC016: 多人聊天壓力測試

**步驟 Steps:**

1. 開啟多個瀏覽器標籤頁
2. 模擬5-10個用戶同時在線
3. 同時傳送訊息
4. 檢查系統穩定性

**預期結果 Expected Results:**

- ✅ 所有連線保持穩定
- ✅ 訊息傳送無遺失
- ✅ 系統反應時間正常
- ✅ 無記憶體洩漏或錯誤

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] WebSocket 連線穩定性
- [ ] 訊息傳送可靠性
- [ ] 資料持久化正確
- [ ] 併發處理能力

---

## 測試場景 6: 安全性與權限測試

### Test Scenario 6: Security and Authorization Testing

### 🎯 測試目標 Test Objectives

驗證安全防護機制、權限控制、輸入驗證

### 📋 測試步驟 Test Steps

#### TC017: 未授權訪問測試

**步驟 Steps:**

1. 未登入狀態下直接訪問受保護頁面
2. 嘗試訪問：
   - `/user/profile`
   - `/matching`
   - `/chat`
3. 檢查重定向和錯誤處理

**預期結果 Expected Results:**

- ❌ 拒絕未授權訪問
- ✅ 重定向至登入頁面
- ✅ 顯示適當的錯誤訊息
- ✅ 會話安全處理正確

#### TC018: SQL 注入防護測試

**步驟 Steps:**

1. 在登入表單中輸入 SQL 注入程式碼
2. 電子郵件欄位: `admin'; DROP TABLE users; --`
3. 密碼欄位: `' OR '1'='1`
4. 提交表單

**預期結果 Expected Results:**

- ✅ 輸入被正確轉義或拒絕
- ✅ 沒有執行危險的 SQL 指令
- ✅ 資料庫結構完整
- ✅ 記錄安全事件

#### TC019: XSS 防護測試

**步驟 Steps:**

1. 在個人簡介中輸入 JavaScript 程式碼
2. 輸入: `<script>alert('XSS Attack')</script>`
3. 保存並檢視個人資料
4. 檢查腳本是否被執行

**預期結果 Expected Results:**

- ✅ 腳本被轉義或移除
- ❌ 沒有執行 JavaScript 警告
- ✅ 內容安全顯示
- ✅ HTML 標籤被適當處理

#### TC020: 檔案上傳安全測試

**步驟 Steps:**

1. 嘗試上傳可執行檔案 (.exe, .php)
2. 嘗試上傳包含惡意程式碼的檔案
3. 檢查檔案類型驗證

**預期結果 Expected Results:**

- ❌ 危險檔案類型被拒絕
- ✅ 顯示適當錯誤訊息
- ✅ 檔案掃描機制有效
- ✅ 上傳路徑安全

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] 認證機制安全性
- [ ] 輸入驗證完整性
- [ ] 檔案上傳安全性
- [ ] 錯誤處理安全性

---

## 測試場景 7: 效能與負載測試

### Test Scenario 7: Performance and Load Testing

### 🎯 測試目標 Test Objectives

驗證系統效能、併發處理能力、響應時間

### 📋 測試步驟 Test Steps

#### TC021: 頁面載入效能測試

**步驟 Steps:**

1. 使用瀏覽器開發者工具
2. 測量各頁面載入時間：
   - 首頁載入時間
   - 配對頁面載入時間
   - 聊天頁面載入時間
3. 記錄 Network 和 Performance 指標

**預期結果 Expected Results:**

- ✅ 首頁載入時間 < 2秒
- ✅ 配對頁面載入時間 < 3秒
- ✅ 聊天頁面載入時間 < 2秒
- ✅ 資源載入優化良好

#### TC022: 資料庫查詢效能測試

**步驟 Steps:**

1. 建立測試資料 (100+ 用戶)
2. 執行配對查詢
3. 檢查查詢執行時間
4. 監控資料庫連線數

**預期結果 Expected Results:**

- ✅ 配對查詢時間 < 500ms
- ✅ 資料庫連線池管理正常
- ✅ 索引使用效率良好
- ✅ 記憶體使用量合理

#### TC023: WebSocket 併發連線測試

**步驟 Steps:**

1. 建立多個 WebSocket 連線 (50+)
2. 同時傳送訊息
3. 監控系統資源使用
4. 檢查連線穩定性

**預期結果 Expected Results:**

- ✅ 支援50+併發連線
- ✅ 訊息傳送無遺失
- ✅ CPU 和記憶體使用正常
- ✅ 連線管理機制有效

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] 響應時間符合要求
- [ ] 併發處理能力足夠
- [ ] 資源使用量合理
- [ ] 系統穩定性良好

---

## 測試場景 8: 跨平台與相容性測試

### Test Scenario 8: Cross-Platform and Compatibility Testing

### 🎯 測試目標 Test Objectives

驗證不同瀏覽器、裝置、作業系統的相容性

### 📋 測試步驟 Test Steps

#### TC024: 多瀏覽器相容性測試

**測試瀏覽器 Browsers to Test:**

- Chrome (最新版本)
- Firefox (最新版本)
- Safari (最新版本)
- Edge (最新版本)

**步驟 Steps:**

1. 在每個瀏覽器中執行核心功能測試
2. 檢查 UI 顯示是否一致
3. 測試 JavaScript 功能是否正常

**預期結果 Expected Results:**

- ✅ 所有主要功能在各瀏覽器正常運作
- ✅ UI 外觀保持一致
- ✅ WebSocket 連線在各瀏覽器穩定

#### TC025: 響應式設計測試

**測試裝置 Device Sizes:**

- 手機 (360x640)
- 平板 (768x1024)
- 桌面 (1920x1080)

**步驟 Steps:**

1. 調整瀏覽器視窗大小
2. 測試各頁面在不同尺寸下的顯示
3. 檢查觸控操作和滑動功能

**預期結果 Expected Results:**

- ✅ 所有頁面響應式設計正常
- ✅ 觸控操作體驗良好
- ✅ 內容在各尺寸下完整顯示

### 🔍 驗證檢查點 Validation Checkpoints

- [ ] 跨瀏覽器相容性
- [ ] 響應式設計完整性
- [ ] 觸控體驗優化
- [ ] 視覺一致性

---

## 測試報告範本

### Test Report Template

### 測試執行摘要 Test Execution Summary

- **測試日期 Test Date:** [填入測試日期]
- **測試人員 Tester:** [填入測試人員姓名]
- **測試環境 Test Environment:** [Development/Staging/Production]
- **測試版本 Version Tested:** [填入應用程式版本]

### 測試結果統計 Test Results Statistics

- **總測試案例數 Total Test Cases:** 25
- **通過 Passed:** [X]
- **失敗 Failed:** [X]
- **跳過 Skipped:** [X]
- **通過率 Pass Rate:** [X%]

### 關鍵問題報告 Critical Issues Found

1. **問題描述 Issue Description:**
   - **嚴重程度 Severity:** [High/Medium/Low]
   - **重現步驟 Reproduction Steps:** [詳細步驟]
   - **預期結果 Expected Result:** [預期行為]
   - **實際結果 Actual Result:** [實際行為]
   - **建議修復 Suggested Fix:** [修復建議]

### 效能測試結果 Performance Test Results

- **平均頁面載入時間 Average Page Load Time:** [X]ms
- **併發用戶支援數 Concurrent Users Supported:** [X]
- **資料庫查詢平均時間 Average DB Query Time:** [X]ms
- **記憶體使用量 Memory Usage:** [X]MB

### 安全性測試結果 Security Test Results

- **認證機制 Authentication:** [PASS/FAIL]
- **授權控制 Authorization:** [PASS/FAIL]
- **輸入驗證 Input Validation:** [PASS/FAIL]
- **檔案上傳安全 File Upload Security:** [PASS/FAIL]

### 建議與結論 Recommendations and Conclusions

[填入測試結論和改進建議]

---

## 自動化測試腳本範例

### Automation Test Script Examples

### 使用 Selenium WebDriver (Python)

```python
from selenium import webdriver
from selenium.webdriver.common.by import By
import time

def test_user_registration():
    driver = webdriver.Chrome()
    driver.get("http://localhost:8080")
    
    # 點擊註冊按鈕
    register_btn = driver.find_element(By.LINK_TEXT, "註冊")
    register_btn.click()
    
    # 填寫註冊表單
    driver.find_element(By.NAME, "email").send_keys("test@example.com")
    driver.find_element(By.NAME, "password").send_keys("SecurePass123!")
    driver.find_element(By.NAME, "name").send_keys("測試用戶")
    
    # 提交表單
    submit_btn = driver.find_element(By.TYPE, "submit")
    submit_btn.click()
    
    # 驗證結果
    assert "註冊成功" in driver.page_source
    
    driver.quit()

if __name__ == "__main__":
    test_user_registration()
```

### 使用 curl 進行 API 測試

```bash
#!/bin/bash

# 測試用戶註冊 API
curl -X POST http://localhost:8080/api/user/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!",
    "name": "測試用戶",
    "birth_date": "1995-06-15",
    "gender": "male"
  }'

# 測試登入 API
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'
```

---

## 檢查清單 Checklist

### 測試前準備 Pre-Test Preparation

- [ ] 測試環境已設定完成
- [ ] 測試資料已準備
- [ ] 應用程式正常運行
- [ ] 測試工具已安裝

### 功能測試 Functional Testing

- [ ] 用戶註冊功能
- [ ] 用戶登入功能
- [ ] 個人資料管理
- [ ] 照片上傳功能
- [ ] 配對系統功能
- [ ] 即時聊天功能

### 安全性測試 Security Testing

- [ ] 認證與授權
- [ ] 輸入驗證
- [ ] SQL 注入防護
- [ ] XSS 防護
- [ ] 檔案上傳安全

### 效能測試 Performance Testing

- [ ] 頁面載入速度
- [ ] 資料庫查詢效能
- [ ] WebSocket 效能
- [ ] 併發處理能力

### 相容性測試 Compatibility Testing

- [ ] 多瀏覽器測試
- [ ] 響應式設計測試
- [ ] 裝置相容性測試

### 測試完成 Test Completion

- [ ] 測試報告已完成
- [ ] 問題已記錄並分類
- [ ] 修復建議已提供
- [ ] 回歸測試已規劃

---

**📝 注意事項 Notes:**

1. 所有測試都應在隔離的測試環境中執行
2. 測試資料不應包含真實用戶資訊
3. 關鍵安全測試必須由專業安全人員執行
4. 效能測試應在符合生產環境規格的硬體上執行
5. 所有測試結果都應詳細記錄並歸檔

**🎯 測試完成標準:**

- 所有關鍵功能測試通過率 ≥ 95%
- 無高嚴重度安全漏洞
- 效能指標符合預期要求
- 跨平台相容性良好
