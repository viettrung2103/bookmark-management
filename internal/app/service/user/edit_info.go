package user

import (
	"context"
)

// EditInfoByID edits a user by ID
func (s *userService) EditInfoByID(ctx context.Context, userId string, inputDisplayName string, inputEmail string) error {
	return s.userRepo.EditUserByID(ctx, userId, inputDisplayName, inputEmail)
}
