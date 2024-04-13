package models

import (
	"time"

	"github.com/google/uuid"
)

type Like struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null" json:"likedByID"`
	EventID        *uuid.UUID `gorm:"type:uuid" json:"eventID"`
	AnnouncementID *uuid.UUID `gorm:"type:uuid" json:"announcementID"`
	PostID         *uuid.UUID `gorm:"type:uuid" json:"postID"`
	CommentID      *uuid.UUID `gorm:"type:uuid" json:"commentID"`
	ReviewID       *uuid.UUID `gorm:"type:uuid" json:"reviewID"`
	Status         int8       `gorm:"not null;default:0" json:"-"` //* 0 for like and -1 for dislike
	CreatedAt      time.Time  `gorm:"default:current_timestamp" json:"likedAt"`
}

func (likeModel *Like) SetItemID(likeType string, itemID uuid.UUID) {
	switch likeType {
	case "comment":
		likeModel.CommentID = &itemID
	case "event":
		likeModel.EventID = &itemID
	case "announcement":
		likeModel.AnnouncementID = &itemID
	}
}
