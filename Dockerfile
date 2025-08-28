# 構建階段
FROM golang:1.22-alpine AS builder

WORKDIR /app

# 複製 go mod 檔案
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製原始碼
COPY . .

# 構建應用程式
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# 運行階段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 複製二進制檔案
COPY --from=builder /app/main .

# 複製靜態檔案
COPY --from=builder /app/web ./web

# 暴露埠號
EXPOSE 8080

# 執行應用程式
CMD ["./main"]
