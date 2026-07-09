package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

//go:generate mockery --name=UserRepository --filename=user.go

// user Repository interface
type UserRepository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByUserId(ctx context.Context, userId string) (*model.User, error)
	GetUserByUniqueField(ctx context.Context, uniqueCol, value string) (*model.User, error)
	EditUserByID(ctx context.Context, userId string, inputDisplayName string, inputEmail string) error
}

type userRepository struct {
	db *gorm.DB
}

// NewRepository returns a new repository
func NewRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
