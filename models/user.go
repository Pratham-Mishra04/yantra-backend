package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type User struct {
	ID                        uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name                      string           `gorm:"type:text;not null" json:"name"`
	Username                  string           `gorm:"type:text;unique;not null" json:"username"`
	Email                     string           `gorm:"unique;not null" json:"-"`
	Password                  string           `json:"-"`
	ProfilePic                string           `gorm:"default:default.jpg" json:"profilePic"`
	CoverPic                  string           `gorm:"default:default.jpg" json:"coverPic"`
	PhoneNo                   string           `json:"-"`
	Bio                       string           `json:"bio"`
	Tags                      pq.StringArray   `gorm:"type:text[]" json:"tags"`
	Location                  string           `gorm:"type:text" json:"location"`
	PasswordResetToken        string           `json:"-"`
	PasswordResetTokenExpires time.Time        `json:"-"`
	PasswordChangedAt         time.Time        `gorm:"default:current_timestamp" json:"-"`
	DeactivatedAt             time.Time        `gorm:"" json:"-"`
	IsModerator               bool             `gorm:"default:false" json:"isModerator"`
	IsVerified                bool             `gorm:"default:false" json:"isVerified"`
	IsOnboardingCompleted     bool             `gorm:"default:false" json:"isOnboardingComplete"`
	LastLoggedIn              time.Time        `gorm:"default:current_timestamp" json:"-"`
	IsActive                  bool             `gorm:"default:true" json:"-"`
	CreatedAt                 time.Time        `gorm:"default:current_timestamp;index:idx_created_at,sort:desc" json:"-"`
	OAuth                     OAuth            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Verification              UserVerification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Journal                   Journal          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Followers                 []Connection     `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE" json:"-"`
	Following                 []Connection     `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"-"`
}

type Moderator struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	User       User      `gorm:"" json:"user"`
	IsDoctor   bool      `gorm:"default:false" json:"isDoctor"`
	IsStudent  bool      `gorm:"default:false" json:"isStudent"`
	University string    `json:"university"`
}

func (u *User) AfterFind(tx *gorm.DB) error {
	if !u.IsActive {
		u.Username = "deactived"
		u.Name = "User"
		u.CoverPic = "default.jpg"
		u.ProfilePic = "default.jpg"
		u.Bio = ""
		u.Tags = nil
	}
	return nil
}

type Provider string

const (
	Google Provider = "Google"
)

type OAuth struct {
	ID                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID              uuid.UUID `gorm:"type:uuid;not null" json:"userID"`
	Provider            Provider  `gorm:"type:text" json:"provider"`
	OnBoardingCompleted bool      `gorm:"default:false" json:"-"`
	CreatedAt           time.Time `gorm:"default:current_timestamp" json:"-"`
}

type UserVerification struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;unique" json:"userID"`
	Code           string    `json:"code"`
	ExpirationTime time.Time `json:"expirationTime"`
}
