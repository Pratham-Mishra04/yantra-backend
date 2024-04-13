package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID                               uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	CreatingUserID                   uuid.UUID  `gorm:"type:uuid;not null" json:"createdByID"`
	CreatingUser                     User       `gorm:"" json:"createdBy"`
	AcceptingUserID                  uuid.UUID  `gorm:"type:uuid;not null" json:"acceptedByID"`
	AcceptingUser                    User       `gorm:"" json:"acceptedBy"`
	CreatedAt                        time.Time  `gorm:"default:current_timestamp" json:"createdAt"`
	LastResetByCreatingUser          time.Time  `gorm:"default:current_timestamp" json:"-"`
	LastResetByAcceptingUser         time.Time  `gorm:"default:current_timestamp" json:"-"`
	BlockedByCreatingUser            bool       `gorm:"default:false" json:"blockedByCreatingUser"`
	BlockedByAcceptingUser           bool       `gorm:"default:false" json:"blockedByAcceptingUser"`
	Messages                         []Message  `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID                  *uuid.UUID `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage                    *Message   `gorm:"foreignKey:LatestMessageID;constraint:OnDelete:CASCADE" json:"latestMessage"`
	Accepted                         bool       `gorm:"default:false" json:"accepted"`
	LastReadMessageByCreatingUserID  *uuid.UUID `gorm:"type:uuid" json:"lastReadMessageByCreatingUserID"`
	LastReadMessageByAcceptingUserID *uuid.UUID `gorm:"type:uuid" json:"lastReadMessageByAcceptingUserID"`
	LastReadMessageByCreatingUser    *Message   `gorm:"foreignKey:LastReadMessageByCreatingUserID;constraint:OnDelete:CASCADE" json:"lastReadMessageByCreatingUser"`
	LastReadMessageByAcceptingUser   *Message   `gorm:"foreignKey:LastReadMessageByAcceptingUserID;constraint:OnDelete:CASCADE" json:"lastReadMessageByAcceptingUser"`
}

type GroupChat struct { //TODO24 store number of members in model to show in invitation
	ID              uuid.UUID          `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Title           string             `gorm:"type:varchar(50);" json:"title"`
	Description     string             `gorm:"type:text" json:"description"`
	AdminOnly       bool               `gorm:"default:false" json:"adminOnly"`
	CoverPic        string             `gorm:"type:text; default:default.jpg" json:"coverPic"`
	UserID          uuid.UUID          `gorm:"type:uuid;not null" json:"userID"`
	User            User               `gorm:"" json:"user"`
	CreatedAt       time.Time          `gorm:"default:current_timestamp" json:"createdAt"`
	Messages        []GroupChatMessage `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE" json:"messages"`
	LatestMessageID *uuid.UUID         `gorm:"type:uuid" json:"latestMessageID"`
	LatestMessage   *GroupChatMessage  `gorm:"foreignKey:LatestMessageID;constraint:OnDelete:CASCADE" json:"latestMessage"`
}
