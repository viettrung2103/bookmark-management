package service

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/repository"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

const (
	urlCodeLength = 7
)

// ShortenUrl represents the shorten url service
//
//go:generate mockery --name=ShortenUrl --filename=shortenurl.go
type ShortenUrl interface {
	ShortenUrlWithExpiringTime(ctx context.Context, url string, expireTime int) (string, error)
}

type shortenUrlService struct {
	repo repository.UrlStorage
}

// NewShortenUrl returns a new ShortenUrl
func NewShortenUrl(repo repository.UrlStorage) ShortenUrl {
	return &shortenUrlService{repo: repo}
}

// ShortenUrlWithExpiringTime shortens a url with expiring time
func (s *shortenUrlService) ShortenUrlWithExpiringTime(ctx context.Context, url string, expireTime int) (string, error) {

	for {
		// tao key
		urlCode, err := stringutils.GenerateCode(urlCodeLength)
		if err != nil {
			return "", err
		}
		// add vao repo
		isUnique, err := s.repo.StoreUrlIfUniqueCode(ctx, urlCode, url, expireTime)
		if err != nil {
			return "", nil
		}

		// tra ve key
		if isUnique == true {
			return urlCode, nil
		}

	}
}
