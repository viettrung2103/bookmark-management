package urlstorage

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/repository/urlstorage"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

// ShortenUrl represents the shorten url service
//
//go:generate mockery --name=ShortenUrl --filename=shortenurl.go
type Service interface {
	ShortenUrlWithExpiringTime(ctx context.Context, url string, expireTime int) (string, error)
	GetLinkFromCode(ctx context.Context, urlCode string) (string, error)
}

type shortenUrlService struct {
	repo   urlstorage.Repository
	keygen stringutils.KeyGenerator
}

// NewShortenUrl returns a new ShortenUrl
func NewService(repo urlstorage.Repository, keygen stringutils.KeyGenerator) Service {
	return &shortenUrlService{
		repo:   repo,
		keygen: keygen,
	}
}
