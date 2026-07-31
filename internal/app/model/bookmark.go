package model

import "github.com/google/uuid"

type Bookmark struct {
	Base
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Code        string    `json:"code"`
	CodeInt     uint64    `json:"-" gorm:"autoIncrement;uniqueIndex;column:code_int"`
	UserID      uuid.UUID `json:"-" gorm:"type:uuid;column:user_id"`
	User        User      `gorm:"references:ID" json:"-"`
}
