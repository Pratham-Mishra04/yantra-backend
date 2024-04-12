package models

import (
	"github.com/google/uuid"
)

type Group struct {
	ID              uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string           `gorm:"type:text;not null" json:"title"`
	Description     string           `gorm:"type:text;not null" json:"description"`
	Moderator       Moderator        `gorm:""`
	NumberOfMembers int16            `gorm:"default:0" json:"noMembers"`
	ResourceBucket  []ResourceBucket `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Polls           []Poll           `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Announcements   []Announcement   `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
	Posts           []Post           `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE" json:"-"`
}
