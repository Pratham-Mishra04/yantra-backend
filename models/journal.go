package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Journal struct {
	ID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	Pages  []Page    `gorm:"foreignKey:JournalID;constraint:OnDelete:CASCADE" json:"pages"`
}

type Page struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	JournalID uuid.UUID       `gorm:"type:uuid;not null" json:"journalID"`
	Emotions  pq.StringArray  `gorm:"type:text[]" json:"-"`
	Scores    pq.Float64Array `gorm:"type:text[]" json:"-"`
	NER       pq.StringArray  `gorm:"type:text[]" json:"-"`
	Title     string          `gorm:"type:text" json:"title"`
	Content   string          `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time       `gorm:"default:current_timestamp;index:idx_created_at,sort:desc" json:"createdAt"`
	UpdatedAt time.Time       `gorm:"default:current_timestamp;index:idx_created_at,sort:desc" json:"updatedAt"`
}
