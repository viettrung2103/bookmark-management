package bookmark

import "context"

// DeleteBookmarkByID delete bookmark of id of current user
func (s *bookmarkService) DeleteBookmarkByID(ctx context.Context, userID string, bookmarkID string) error {

	return s.bookmarkRepo.DeleteBookmarkByID(ctx, userID, bookmarkID)
}
