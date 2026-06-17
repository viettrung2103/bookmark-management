package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string `gorm:"type:uuid;primarykey;column:id" json:"id"`
	DisplayName string `gorm:"column:display_name" json:"display_name"`
	Username    string `gorm:"unique;column:username" json:"username"`
	Password    string `gorm:"column:password" json:"_"`
	Email       string `gorm:"unique;column:email" json:"email"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
