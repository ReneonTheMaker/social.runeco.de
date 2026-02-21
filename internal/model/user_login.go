package model

import "time"

type UserLogin struct {
	ID        uint   `gorm:"primaryKey"`
	TokenHash string `gorm:"uniqueIndex;size:64;not null"`
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"constraint:OnDelete:CASCADE"`

	CreatedAt  time.Time
	LastSeenAt time.Time
	ExpiresAt  time.Time
}
