package bookmark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
	"github.com/viettrung2103/bookmark-management/internal/app/repository/cache"
	"github.com/viettrung2103/bookmark-management/pkg/common"
)

const (
	cacheGroupKeyFormatGetBookmarks = "get_bookmarks_%s" // %s is userID
	cacheKeyFormatGetBookmarks      = "page_%d_limit_%d" // %d is page and limit
	CacheExpireDuration             = 24 * time.Hour
)

type bookmarkCacheService struct {
	s     Service  // bookmark service
	cache cache.DB //cache repo
}

func NewBookmarkCacheService(s Service, cache cache.DB) Service {
	return &bookmarkCacheService{
		s:     s,
		cache: cache,
	}
}

func (s *bookmarkCacheService) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarkResult, error) {
	// tao cache key
	cacheGroupKey := fmt.Sprintf(cacheGroupKeyFormatGetBookmarks, userID)
	cacheKey := fmt.Sprintf(cacheKeyFormatGetBookmarks, page, limit)

	// check cache data
	cacheDataRaw, err := s.cache.GetCacheData(ctx, cacheGroupKey, cacheKey)
	if err != nil {
		log.Err(err).Msg("Failed to get cache data")
	}

	// if yes >> return cached data
	if len(cacheDataRaw) > 0 && err == nil {
		result := &GetBookmarkResult{}
		err := json.Unmarshal(cacheDataRaw, result)
		if err != nil {
			log.Err(err).Msg("Failed to unmarshal cache data")
			err := s.cache.DeleteCacheKey(ctx, cacheGroupKey, cacheKey)
			if err != nil {
				log.Err(err).Msg("Failed to delete cache key")
			}
		}
		return result, nil
	}

	// if no >> call service to get data from db
	result, err := s.s.GetBookmarks(ctx, userID, page, limit)
	if err != nil {
		log.Err(err).Msg("Failed to get bookmarks after failed to get cache data")
		return nil, err
	}

	// set cache data
	resultInBytes, err := json.Marshal(result)
	if err != nil {
		log.Err(err).Msg("Failed to marshal result data")
	}
	if len(resultInBytes) > 0 && err == nil {
		err := s.cache.SetCacheData(ctx, cacheGroupKey, cacheKey, resultInBytes, CacheExpireDuration)
		if err != nil {
			log.Err(err).Msg("Failed to set cache data")
		}
	}

	// return
	return result, nil

}

func (s *bookmarkCacheService) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	cacheGroupKey := fmt.Sprintf(cacheGroupKeyFormatGetBookmarks, userID)
	err := s.cache.DeleteCacheGroupKey(ctx, cacheGroupKey)
	if err != nil {
		log.Err(err).Msg(common.DeleteCacheError)
		return nil, err
	}
	return s.s.AddBookmark(ctx, description, url, userID)
}
func (s *bookmarkCacheService) EditBookmarkByID(ctx context.Context, userID string, bookmarkID string, newDescription string, newURL string) error {
	cacheGroupKey := fmt.Sprintf(cacheGroupKeyFormatGetBookmarks, userID)
	err := s.cache.DeleteCacheGroupKey(ctx, cacheGroupKey)
	if err != nil {
		log.Err(err).Msg(common.DeleteCacheError)
		return err
	}

	return s.s.EditBookmarkByID(ctx, userID, bookmarkID, newDescription, newURL)
}
func (s *bookmarkCacheService) DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error {
	cacheGroupKey := fmt.Sprintf(cacheGroupKeyFormatGetBookmarks, userID)
	err := s.cache.DeleteCacheGroupKey(ctx, cacheGroupKey)
	if err != nil {
		log.Err(err).Msg(common.DeleteCacheError)
		return err
	}

	return s.s.DeleteBookmarkByID(ctx, userID, bookmarkID)
}
