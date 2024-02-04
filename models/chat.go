package model

import (
	"time"

	"gorm.io/gorm"
)

type ChatAttachment struct {
	gorm.Model
	MessageID string
	Name      string
	Size      int
	Type      string
	Path      string
}

type ChatMessage struct {
	gorm.Model
	ConversationID string
	Body           string
	ContentType    string
	Attachments    []ChatAttachment `gorm:"foreignKey:MessageID"`
	SenderID       uint
	CreatedAt      string
	UpdatedAt      string
}

type ChatConversation struct {
	gorm.Model
	Participants []User `gorm:"many2many:conversation_users;"`
	Type         string
	UnreadCount  int
	Messages     []ChatMessage `gorm:"foreignKey:ConversationID"`
}

type ConversationUser struct {
	ChatConversationID uint64 `gorm:"primaryKey"`
	UserID             uint64 `gorm:"primaryKey"`
}

type ChatSendMessage struct {
	ConversationID string           `json:"conversationId"`
	MessageID      string           `json:"messageId"`
	Message        string           `json:"message"`
	ContentType    string           `json:"contentType"`
	Attachments    []ChatAttachment `json:"attachments"`
	CreatedAt      time.Time        `json:"createdAt"`
	SenderID       uint             `json:"senderId"`
}
