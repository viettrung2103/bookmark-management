package user

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/model"
)

func (s *userService) SelfInfo(ctx context.Context, userId string) (*model.User, error) {
	return s.userRepo.GetUserByUserId(ctx, userId)
}
