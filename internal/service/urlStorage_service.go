package service

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/repository"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

const (
	urlCodeLength = 7
)

//go:generate mockery --name=ShortenUrl --filename=shortenurl.go

// ShortenUrl represents the shorten url service
type ShortenUrl interface {
	ShortenUrl(ctx context.Context, url string) (string, error)
	ShortenUrlWithExpiringTime(ctx context.Context, url string, expireTime int) (string, error)
	CheckHealth(ctx context.Context) error
}

type shortenUrlService struct {
	repo repository.UrlStorage
}

// NewShortenUrl returns a new ShortenUrl
func NewShortenUrl(repo repository.UrlStorage) ShortenUrl {
	return &shortenUrlService{repo: repo}
}

// ShortenUrl shortens a url
func (s *shortenUrlService) ShortenUrl(ctx context.Context, url string) (string, error) {

	// tao key
	urlCode, err := stringutils.GenerateCode(urlCodeLength)
	if err != nil {
		return "", err
	}
	// add vao repo
	err = s.repo.StoreURL(ctx, urlCode, url)
	if err != nil {
		return "", nil
	}

	// tra ve key
	return urlCode, nil
}

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

func (s *shortenUrlService) CheckHealth(ctx context.Context) error {
	return s.repo.CheckHealth(ctx)
}
