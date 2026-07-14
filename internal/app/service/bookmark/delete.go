package bookmark

import "context"

func (s *bookmarkService) DeleteBookmarkByID(ctx context.Context, userID string, bookmarkID string) error {
	println("is problem in service")

	return s.bookmarkRepo.DeleteBookmarkByID(ctx, userID, bookmarkID)
}
