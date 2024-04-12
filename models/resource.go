package models

import (
	"time"

	"github.com/google/uuid"
)

type ResourceBucket struct {
	ID                  uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	GroupID             uuid.UUID      `gorm:"type:uuid;not null" json:"groupID"`
	Title               string         `gorm:"type:text;not null" json:"title"`
	Description         string         `gorm:"type:text" json:"description"`
	OnlyAdminViewAccess bool           `gorm:"" json:"onlyAdminViewAccess"`
	OnlyAdminEditAccess bool           `gorm:"" json:"onlyAdminEditAccess"`
	CreatedAt           time.Time      `gorm:"default:current_timestamp" json:"createdAt"`
	NumberOfFiles       int16          `gorm:"default:0" json:"noFiles"`
	ResourceFiles       []ResourceFile `gorm:"foreignKey:ResourceBucketID;constraint:OnDelete:CASCADE" json:"resourceFiles"`
}

type ResourceFile struct {
	ID               uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User             User      `json:"user"`
	ResourceBucketID uuid.UUID `gorm:"type:uuid;not null" json:"resourceBucketID"`
	Title            string    `gorm:"type:text;not null" json:"title"`
	Description      string    `gorm:"type:text" json:"description"`
	Type             string    `gorm:"type:text" json:"type"`
	FileUploaded     bool      `gorm:"default:false" json:"isFileUploaded"`
	Path             string    `gorm:"type:text;not null" json:"path"`
	CreatedAt        time.Time `gorm:"default:current_timestamp" json:"createdAt"`
}
