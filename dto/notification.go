package dto

import (
	"time"
)

type Notificationdto struct {
	Message   string    `json:"message"`
	UserID    uint16    `json:"userid"`
	CreatedAt time.Time `json:"created_at"`
}
