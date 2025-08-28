package chat

// ChatRoom 代表聊天室的領域實體
type ChatRoom struct {
	ID          RoomID   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []UserID `json:"members"`
	IsActive    bool     `json:"is_active"`
}

// NewChatRoom 創建新聊天室
func NewChatRoom(name, description string) *ChatRoom {
	return &ChatRoom{
		ID:          RoomID("room_" + name),
		Name:        name,
		Description: description,
		Members:     make([]UserID, 0),
		IsActive:    true,
	}
}

// AddMember 添加成員到聊天室
func (r *ChatRoom) AddMember(userID UserID) {
	for _, member := range r.Members {
		if member == userID {
			return // 用戶已存在
		}
	}
	r.Members = append(r.Members, userID)
}

// RemoveMember 從聊天室移除成員
func (r *ChatRoom) RemoveMember(userID UserID) {
	for i, member := range r.Members {
		if member == userID {
			r.Members = append(r.Members[:i], r.Members[i+1:]...)
			return
		}
	}
}

// IsMember 檢查用戶是否為聊天室成員
func (r *ChatRoom) IsMember(userID UserID) bool {
	for _, member := range r.Members {
		if member == userID {
			return true
		}
	}
	return false
}

// ChatRoomRepository 聊天室倉儲介面
type ChatRoomRepository interface {
	Save(room *ChatRoom) error
	FindByID(id RoomID) (*ChatRoom, error)
	FindAll() ([]*ChatRoom, error)
}
