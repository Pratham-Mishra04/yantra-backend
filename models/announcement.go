package models

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GroupID       uuid.UUID      `gorm:"type:uuid;not null" json:"groupID"`
	Group         Group          `gorm:"" json:"group"`
	Title         string         `gorm:"" json:"title"`
	Content       string         `gorm:"not null" json:"content"`
	IsEdited      bool           `gorm:"default:false" json:"isEdited"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
	NoLikes       int            `gorm:"default:0" json:"noLikes"`
	NoComments    int            `gorm:"default:0" json:"noComments"`
	Comments      []Comment      `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	Notifications []Notification `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	Likes         []Like         `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
	Reports       []Report       `gorm:"foreignKey:AnnouncementID;constraint:OnDelete:CASCADE" json:"-"`
}
