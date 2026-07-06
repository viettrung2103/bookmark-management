package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User struct for user model
type User struct {
	ID          string    `gorm:"type:uuid;primarykey;column:id" json:"id"`
	DisplayName string    `gorm:"column:display_name" json:"display_name"`
	Username    string    `gorm:"unique;column:username" json:"username"`
	Password    string    `gorm:"column:password" json:"-"`
	Email       string    `gorm:"unique;column:email" json:"email"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"` // GORM auto-manages this on creation
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"` // GORM auto-manages this on updates
}

// BeforeCreate is a callback function that is called before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
