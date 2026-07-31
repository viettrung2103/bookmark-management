package bookmark

import (
	"context"

	"gorm.io/gorm"
)

func (r *bookmarkRepository) Transaction(ctx context.Context, fn func(txRepo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &bookmarkRepository{
			db: tx,
		}
		return fn(txRepo)
	})
}
