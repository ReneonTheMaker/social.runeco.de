package model

import "time"

/*
users
- id (PK)
- username (unique, indexed)
- password_hash
- created_at
*/

type User struct {
	ID           uint     `gorm:"primaryKey"`
	Username     string   `gorm:"uniqueIndex;size:32;not null"`
	PasswordHash string   `gorm:"not null"`
	UserInfo     UserInfo `gorm:"constraint:OnDelete:CASCADE"`
	Mod          bool     `gorm:"default:false"`
	CreatedAt    time.Time
}
