package models

import (
	"time"

	"github.com/google/uuid"
)

/*
notification type:
*-1 - Welcome to the platform
*/

type Notification struct {
	ID               uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	NotificationType int16        `json:"notificationType"`
	UserID           uuid.UUID    `gorm:"type:uuid;not null" json:"userID"`
	User             User         `json:"user"`
	SenderID         uuid.UUID    `gorm:"type:uuid;not null" json:"senderID"`
	Sender           User         `json:"sender"`
	EventID          *uuid.UUID   `gorm:"type:uuid" json:"eventID"`
	Event            Event        `json:"event"`
	AnnouncementID   *uuid.UUID   `gorm:"type:uuid" json:"announcementID"`
	Announcement     Announcement `json:"announcement"`
	Read             bool         `gorm:"default:false" json:"isRead"`
	CreatedAt        time.Time    `gorm:"default:current_timestamp" json:"createdAt"`
}
