package models

import (
	"time"

	"github.com/google/uuid"
)

type Report struct { //TODO32 notifications will be sent after reporting
	ID             uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ReportType     int8         `json:"reportType"`
	ReporterID     uuid.UUID    `gorm:"type:uuid;not null" json:"reporterID"`
	Reporter       User         `json:"report"`
	UserID         *uuid.UUID   `gorm:"type:uuid;" json:"userID"`
	User           User         `json:"user"`
	EventID        *uuid.UUID   `gorm:"type:uuid" json:"eventID"`
	Event          Event        `gorm:"" json:"event"`
	AnnouncementID *uuid.UUID   `gorm:"type:uuid" json:"announcementID"`
	Announcement   Announcement `gorm:"" json:"announcement"`
	GroupID        *uuid.UUID   `gorm:"type:uuid" json:"groupID"`
	Group          Group        `json:"group"`
	PostID         *uuid.UUID   `gorm:"type:uuid" json:"postID"`
	Post           Post         `json:"post"`
	Content        string       `json:"content"`
	CreatedAt      time.Time    `gorm:"default:current_timestamp" json:"createdAt"`
}
