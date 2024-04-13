package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"userID"`
	User      User       `gorm:"" json:"user"`
	EventID   *uuid.UUID `gorm:"type:uuid" json:"eventID"`
	Event     Event      `gorm:"" json:"event"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	IsPending bool       `gorm:"default:false" json:"isPending"`
	CreatedAt time.Time  `gorm:"default:current_timestamp" json:"createdAt"`
}
