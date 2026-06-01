package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string  `gorm:"not null"`
	Description string
	Status      string  `gorm:"not null;default:'todo'"`
	ProjectID   uint    `gorm:"not null"`
	Project     Project
	AssigneeID  *uint
	Assignee    *User
}
