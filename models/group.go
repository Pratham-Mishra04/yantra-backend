package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string           `gorm:"type:text;not null" json:"title"`
	Description     string           `gorm:"type:text;not null" json:"description"`
	ModeratorID     uuid.UUID        `gorm:"type:uuid;not null" json:"journalID"`
	Moderator       Moderator        `gorm:"" json:"moderator"`
	NumberOfMembers int16            `gorm:"default:1" json:"noMembers"`
	ResourceBucket  []ResourceBucket `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Polls           []Poll           `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Announcements   []Announcement   `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Posts           []Post           `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Events          []Event          `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	CreatedAt       time.Time        `gorm:"default:current_timestamp" json:"createdAt"`
}

type GroupMembership struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User      User      `gorm:"" json:"user"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null" json:"groupID"`
	Group     Group     `gorm:"" json:"group"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}
