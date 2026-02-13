package model

import (
	gorm "gorm.io/gorm"
	"time"
)

/*
posts
- id (PK)
- user_id (FK -> users.id)
- parent_id (nullable, FK -> posts.id)
- content (varchar 200 or text with constraint)
- created_at
- deleted_at (nullable)
*/

type Post struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"not null;index"`
	User      User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ParentID  *uint  `gorm:"index"`
	Parent    *Post  `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL"`
	Content   string `gorm:"size:200;not null"`
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
