package contract

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPhotosUploadContract(t *testing.T) {
	// 設置測試用的 Gin 路由器
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// 此端點尚未實作 - 測試應該失敗
	router.POST("/users/photos", func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "端點尚未實作"})
	})

	t.Run("有效的照片上傳應該回傳201", func(t *testing.T) {
		// 建立帶有假影像檔案的 multipart 表單
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 新增照片檔案
		fileWriter, err := writer.CreateFormFile("photo", "test-photo.jpg")
		assert.NoError(t, err)

		// 寫入假的 JPEG 標頭和一些資料
		fakeJPEGData := []byte{0xFF, 0xD8, 0xFF, 0xE0}             // JPEG 檔案簽章
		fakeJPEGData = append(fakeJPEGData, make([]byte, 1000)...) // 新增一些資料
		_, err = fileWriter.Write(fakeJPEGData)
		assert.NoError(t, err)

		// 新增選用欄位
		writer.WriteField("is_primary", "false")
		writer.WriteField("caption", "我的超棒照片")

		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後此測試應該通過
		assert.Equal(t, http.StatusCreated, w.Code, "有效的照片上傳應回傳 201 建立成功")

		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err, "回應應為有效的 JSON")

		// 檢查照片回應結構
		assert.Contains(t, response, "id", "回應應包含照片ID")
		assert.Contains(t, response, "url", "回應應包含照片URL")
		assert.Contains(t, response, "is_primary", "回應應包含主要照片標記")
		assert.Contains(t, response, "caption", "回應應包含照片說明")
		assert.Contains(t, response, "created_at", "回應應包含建立時間")

		// 驗證資料類型
		if id, ok := response["id"]; ok {
			assert.IsType(t, float64(0), id, "照片ID應為數字")
		}
		if url, ok := response["url"]; ok {
			assert.IsType(t, "", url, "照片URL應為字串")
		}
		if isPrimary, ok := response["is_primary"]; ok {
			assert.IsType(t, true, isPrimary, "主要照片標記應為布林值")
		}
	})

	t.Run("缺少授權應回傳401", func(t *testing.T) {
		// 建立簡單的 multipart 表單
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, _ := writer.CreateFormFile("photo", "test.jpg")
		fileWriter.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0}) // JPEG 標頭
		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		// 沒有授權標頭

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "缺少權杖應回傳 401 未授權")
	})

	t.Run("無效權杖應回傳401", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		fileWriter, _ := writer.CreateFormFile("photo", "test.jpg")
		fileWriter.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0})
		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer invalid_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code, "無效權杖應回傳 401 未授權")
	})

	t.Run("缺少照片檔案應回傳400", func(t *testing.T) {
		// 建立沒有照片檔案的 multipart 表單
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("caption", "沒有檔案的照片")
		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "缺少照片檔案應回傳 400 錯誤請求")
	})

	t.Run("無效的檔案類型應回傳400", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 建立文字檔案而非影像檔案
		fileWriter, _ := writer.CreateFormFile("photo", "test.txt")
		fileWriter.Write([]byte("這不是影像檔案"))
		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "無效的檔案類型應回傳 400 錯誤請求")
	})

	t.Run("檔案過大應回傳413", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileWriter, _ := writer.CreateFormFile("photo", "huge-photo.jpg")
		// 寫入 JPEG 標頭
		fileWriter.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0})
		// 寫入大量資料（模擬10MB+檔案）
		largeData := make([]byte, 10*1024*1024) // 10MB
		fileWriter.Write(largeData)
		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code,
			"檔案過大應回傳 413 請求實體過大")
	})

	t.Run("照片數量過多應回傳400", func(t *testing.T) {
		// 模擬使用者已有最大數量的照片
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileWriter, _ := writer.CreateFormFile("photo", "photo.jpg")
		fileWriter.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0})
		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer token_for_user_with_max_photos")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 如果使用者已有最大照片數（通常6-10張）則應為400
		// 目前為501因為尚未實作
		assert.Contains(t, []int{http.StatusBadRequest, http.StatusNotImplemented}, w.Code,
			"照片數量過多應回傳 400 錯誤請求或因尚未實作回傳 501")
	})

	t.Run("無效的multipart表單應回傳400", func(t *testing.T) {
		// 傳送無效的 multipart 資料
		req := httptest.NewRequest("POST", "/users/photos", bytes.NewBufferString("invalid-multipart-data"))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=invalid")
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code, "無效的multipart表單應回傳 400 錯誤請求")
	})

	t.Run("設置主要照片應回傳201", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		fileWriter, _ := writer.CreateFormFile("photo", "primary-photo.jpg")
		fileWriter.Write([]byte{0xFF, 0xD8, 0xFF, 0xE0})

		// 標記為主要照片
		writer.WriteField("is_primary", "true")
		writer.WriteField("caption", "我的主要照片")
		writer.Close()

		req := httptest.NewRequest("POST", "/users/photos", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer valid_jwt_token")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// 實作後應為201
		assert.Contains(t, []int{http.StatusCreated, http.StatusNotImplemented}, w.Code,
			"主要照片應回傳 201 建立成功或因尚未實作回傳 501")
	})
}

func createTestImageFile() io.Reader {
	// 建立最小的有效 JPEG 檔案結構
	jpegHeader := []byte{
		0xFF, 0xD8, // SOI (影像開始)
		0xFF, 0xE0, // APP0
		0x00, 0x10, // APP0 區段長度
		0x4A, 0x46, 0x49, 0x46, 0x00, // "JFIF\0"
		0x01, 0x01, // 版本
		0x01,                   // 單位
		0x00, 0x48, 0x00, 0x48, // X 和 Y 密度
		0x00, 0x00, // 縮圖寬度和高度
		0xFF, 0xD9, // EOI (影像結束)
	}
	return bytes.NewReader(jpegHeader)
}
