package models

import (
	"time"

	"github.com/google/uuid"
)

type Connection struct {
	SenderID   uuid.UUID `json:"followerID"`
	Sender     User      `gorm:"foreignKey:FollowerID" json:"follower"`
	ReceiverID uuid.UUID `json:"followedID"`
	Receiver   User      `gorm:"foreignKey:FollowedID" json:"followed"`
	Status     int8      `gorm:"not null;default:0" json:"-"` //* -1 rejected, 0 waiting, 1 accepted
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"-"`
}
