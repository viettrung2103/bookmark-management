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
type ShortenUrl interface {
	ShortenUrl(ctx context.Context, url string) (string, error)
}

type shortenUrl struct {
	repo repository.UrlStorage
}

// NewShortenUrl returns a new ShortenUrl
func NewShortenUrl(repo repository.UrlStorage) ShortenUrl {
	return &shortenUrl{repo: repo}
}

// ShortenUrl shortens a url
func (s *shortenUrl) ShortenUrl(ctx context.Context, url string) (string, error) {
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
