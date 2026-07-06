package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/app/repository/user"
	"github.com/viettrung2103/bookmark-management/pkg/jwtutils"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

// Service interface for user service
//
//go:generate mockery --name=UserService --filename=../../mocks/user.go
type UserService interface {
	CreateUser(ctx context.Context, displayName, username, password, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	SelfInfo(ctx context.Context, userId string) (*model.User, error)
	EditInfoByID(ctx context.Context, userId string, inputDisplayName string, inputEmail string) error
}

type userService struct {
	userRepo        user.Repository
	passwordHashing stringutils.PasswordHashing
	jwtGenerator    jwtutils.JWTGenerator
}

type UserServiceOpts struct {
	UserRepo        user.Repository
	PasswordHashing stringutils.PasswordHashing
	JwtGenerator    jwtutils.JWTGenerator
}

// NewService creates a new user service
//func NewService(repo user.Repository) Service {
//	return &userService{repo: repo}
//}

// func NewService(repo user.Repository, passwordHashing utils.PasswordHashing, jwtGen jwtutils.JWTGenerator) Service {
func NewService(opts *UserServiceOpts) UserService {

	return &userService{
		userRepo:        opts.UserRepo,
		passwordHashing: opts.PasswordHashing,
		jwtGenerator:    opts.JwtGenerator,
	}
}
