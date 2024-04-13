package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Event struct {
	ID            uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GroupID       uuid.UUID      `gorm:"type:uuid;not null" json:"groupID"`
	Group         Group          `gorm:"" json:"group"`
	DyteID        string         `gorm:"" json:"dyteID"`
	Title         string         `gorm:"type:text;not null" json:"title"`
	Tagline       string         `gorm:"type:text" json:"tagline"`
	CoverPic      string         `gorm:"type:text; default:default.jpg" json:"coverPic"`
	BlurHash      string         `gorm:"type:text; default:no-hash" json:"blurHash"`
	Description   string         `gorm:"type:text;not null" json:"description"`
	Links         pq.StringArray `gorm:"type:text[]" json:"links"`
	Tags          pq.StringArray `gorm:"type:text[]" json:"tags"`
	NoLikes       int            `gorm:"default:0" json:"noLikes"`
	NoComments    int            `gorm:"default:0" json:"noComments"`
	StartTime     time.Time      `gorm:"not null" json:"startTime"`
	EndTime       time.Time      `gorm:"not null" json:"endTime"`
	Location      string         `gorm:"not null" json:"location"`
	Category      string         `gorm:"type:text;not null" json:"category"`
	CreatedAt     time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
	Comments      []Comment      `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"comments"`
	Likes         []Like         `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"-"`
	Reports       []Report       `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"-"`
	Notifications []Notification `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE" json:"-"`
}
