package models

import "github.com/google/uuid"

type Event struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GroupID uuid.UUID `gorm:"type:uuid;not null" json:"groupID"`
	Group   Group     `gorm:"" json:"group"`
}
