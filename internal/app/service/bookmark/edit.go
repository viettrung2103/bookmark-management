package bookmark

import "context"

// EditBookmarkByID edits a bookmark by ID
func (b *bookmarkService) EditBookmarkByID(ctx context.Context, userID string, bookmarkID string, newDescription string, newURL string) error {
	return b.bookmarkRepo.EditBookmarkByID(ctx, userID, bookmarkID, newDescription, newURL)
}
