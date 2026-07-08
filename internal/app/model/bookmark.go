package model

type Bookmark struct {
	Base
	Description string `json:"description"`
	URL         string `json:"url"`
	Code        string `json:"code"`
	UserID      string `json:"-" gorm:"type:uuid;column:user_id"`
	User        User   `gorm:"references:ID" json:"-"`
}
