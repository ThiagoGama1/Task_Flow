package models

import "gorm.io/gorm"

type ActivityLog struct {
	gorm.Model
	ProjectID uint   `gorm:"not null;index"`
	UserID    uint   `gorm:"not null"`
	User      User
	Action    string `gorm:"not null"`
}
