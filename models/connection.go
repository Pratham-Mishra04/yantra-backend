package models

import (
	"time"

	"github.com/google/uuid"
)

type Connection struct {
	SenderID   uuid.UUID `json:"senderID"`
	Sender     User      `gorm:"foreignKey:SenderID" json:"sender"`
	ReceiverID uuid.UUID `json:"receiverID"`
	Receiver   User      `gorm:"foreignKey:ReceiverID" json:"receiver"`
	Status     int8      `gorm:"not null;default:0" json:"-"` //* -1 rejected, 0 waiting, 1 accepted
	CreatedAt  time.Time `gorm:"default:current_timestamp" json:"-"`
}
