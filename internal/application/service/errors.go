package service

import "errors"

// 應用層錯誤定義
var (
	ErrUserNotFound     = errors.New("user not found")
	ErrRoomNotFound     = errors.New("room not found")
	ErrUserNotInRoom    = errors.New("user is not a member of this room")
	ErrInvalidMessage   = errors.New("invalid message content")
	ErrPermissionDenied = errors.New("permission denied")
)
