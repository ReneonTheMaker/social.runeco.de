package model

import (
	"time"
)

/*
user_details
- user_id (PK, FK -> users.id)
- display_name
- bio
- profile_picture_url
- updated_at
*/

type UserInfo struct {
	UserID            uint    `gorm:"primaryKey;not null"`
	User              *User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	DisplayName       string  `gorm:"size:64;not null"`
	Bio               string  `gorm:"size:256"`
	ProfilePictureURL *string `gorm:"size:256"`
	UpdatedAt         time.Time
}
