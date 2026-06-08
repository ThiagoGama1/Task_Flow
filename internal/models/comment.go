package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	TaskID  uint   `gorm:"not null"`
	Task    Task
	UserID  uint `gorm:"not null"`
	User    User
}
