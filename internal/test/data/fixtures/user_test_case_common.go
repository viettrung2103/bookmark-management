package fixtures

import (
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

// UserCommonTestDB struct for user test data

type UserCommonTestDB struct {
	base
}

// Migrate migrates the user table
func (u *UserCommonTestDB) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

// GenerateData generates test data for users
func (u *UserCommonTestDB) GenerateData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []*model.User{
		{
			ID:          "133b3b42-70b9-456c-82e7-bf1b570e6c51",
			DisplayName: "John Doe",
			Username:    "johndoe99",
			Password:    "$2a$12$R9h/lIPbuRRvvV1GG9Z.6uCWMUa.6H7H0S.aK92F7t6GgS2T7M3Ki", // "hashed_password_1"
			Email:       "john.doe@example.com",
		},
		{
			ID:          "87a3cb94-d2e8-422d-bb91-fc5215949eb8",
			DisplayName: "Jane Smith",
			Username:    "janesmith_dev",
			Password:    "$2a$12$K3vX5YwOmP2bZ9rQ3nU7Xe8YvMw9T6uC9iK2oP1lRmSzTxVuWxYz.", // "hashed_password_2"
			Email:       "jane.smith@example.com",
		},
		{
			ID:          "999a9a99-9999-9999-9999-999999999999",
			DisplayName: "test",
			Username:    "test",
			Password:    "$2a$12$R9h/lIPbuRRvvV1GG9Z.6uCWMUa.6H7H0S.aK92F7t6GgS2T7M3Ki",
			Email:       "test@mail.com",
		},
	}

	return db.CreateInBatches(users, 10).Error
}
