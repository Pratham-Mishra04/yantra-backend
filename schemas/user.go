package schemas

import (
	"github.com/lib/pq"
)

type UserCreateSchema struct {
	Name        string `json:"name" validate:"required,max=25"`
	Username    string `json:"username" validate:"required,max=16"` //alphanum+_
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	IsModerator bool   `json:"isModerator" validate:"required"`
	IsDoctor    bool   `json:"isDoctor" validate:"required"`
	IsStudent   bool   `json:"isStudent" validate:"required"`
}

type UserUpdateSchema struct {
	Name       *string         `json:"name" validate:"max=25"`
	ProfilePic *string         `json:"profilePic" validate:"image"`
	CoverPic   *string         `json:"coverPic" validate:"image"`
	Bio        *string         `json:"bio" validate:"max=500"`
	Location   *string         `json:"location" validate:"max=50"`
	Tags       *pq.StringArray `json:"tags"`
}
