package models

import (
	"time"

	"github.com/google/uuid"
)

type Option struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	PollID    uuid.UUID `gorm:"type:uuid;not null" json:"pollID"`
	Content   string    `gorm:"not null" json:"content"`
	Votes     int       `gorm:"type:int;default:0" json:"noVotes"`
	VotedBy   []User    `gorm:"many2many:voted_by;constraint:OnDelete:CASCADE" json:"votedBy"`
	CreatedAt time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}

type Poll struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GroupID       uuid.UUID `gorm:"type:uuid;not null" json:"groupID"`
	Group         Group     `gorm:"" json:"group"`
	Title         string    `gorm:"" json:"title"`
	Content       string    `gorm:"not null" json:"content"`
	Options       []Option  `gorm:"foreignKey:PollID;constraint:OnDelete:CASCADE" json:"options"`
	IsMultiAnswer bool      `gorm:"default:false" json:"isMultiAnswer"`
	IsEdited      bool      `gorm:"default:false" json:"isEdited"`
	TotalVotes    int       `gorm:"default:0" json:"totalVotes"`
	CreatedAt     time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}
