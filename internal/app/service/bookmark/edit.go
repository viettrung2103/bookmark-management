package bookmark

import "context"

// EditInfoByID edits a user by ID
//func (s *userService) EditInfoByID(ctx context.Context, userId string, inputDisplayName string, inputEmail string) error {
//	return s.userRepo.EditUserByID(ctx, userId, inputDisplayName, inputEmail)
//}

func (b *bookmarkService) EditBookmarkByID(ctx context.Context, userID string, bookmarkID string, newDescription string, newURL string) error {
	return b.bookmarkRepo.EditBookmarkByID(ctx, userID, bookmarkID, newDescription, newURL)
}
