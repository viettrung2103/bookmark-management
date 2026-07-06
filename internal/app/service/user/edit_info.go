package user

import (
	"context"
)

//err := h.service.EditInfoByID(c,uid, input.DisplayName, input.Email)

//func (s *userService) SelfInfo(ctx context.Context, userId string) (*model.User, error) {
//	return s.userRepo.GetUserByUserId(ctx, userId)
//}

func (s *userService) EditInfoByID(ctx context.Context, userId string, inputDisplayName string, inputEmail string) error {
	return s.userRepo.EditUserByID(ctx, userId, inputDisplayName, inputEmail)
}
