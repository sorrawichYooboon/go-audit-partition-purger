package domain

import (
	"time"
)

type AuditLog struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Action    string    `json:"action" gorm:"not null"`
	Details   string    `json:"details" gorm:"type:jsonb"`
	CreatedAt time.Time `json:"created_at" gorm:"primaryKey;not null;default:CURRENT_TIMESTAMP"`
}
