# 使用官方 Go 映像檔作為基礎映像檔
FROM golang:1.23-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製 go.mod 和 go.sum 檔案
COPY go.mod go.sum ./

# 下載依賴項
RUN go mod download

# 複製原始碼
COPY . .

# 建構應用程式
RUN go build -o main .

# 使用 alpine 作為最終映像檔
FROM alpine:latest

# 安裝 ca-certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從 builder 階段複製執行檔
COPY --from=builder /app/main .

# 複製靜態檔案資料夾
COPY --from=builder /app/static ./static

# 暴露端口
EXPOSE 8080

# 執行應用程式
CMD ["./main"]
