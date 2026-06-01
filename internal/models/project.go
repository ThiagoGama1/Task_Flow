package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	OwnerID     uint
	Owner       User
	Members     []User `gorm:"many2many:project_members;"`
	Tasks       []Task
}
