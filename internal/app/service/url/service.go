package link

import (
	"context"

	"github.com/viettrung2103/bookmark-management/internal/app/repository/url"
	"github.com/viettrung2103/bookmark-management/pkg/stringutils"
)

// ShortenUrl represents the shorten url service
//
//go:generate mockery --name=URLService --filename=shortenurl.go
type URLService interface {
	ShortenUrlWithExpiringTime(ctx context.Context, url string, expireTime int) (string, error)
	GetLinkFromCode(ctx context.Context, urlCode string) (string, error)
}

type shortenUrlService struct {
	repo   url.URLRepository
	keygen stringutils.KeyGenerator
}

// NewShortenUrl returns a new ShortenUrl
func NewService(repo url.URLRepository, keygen stringutils.KeyGenerator) URLService {
	return &shortenUrlService{
		repo:   repo,
		keygen: keygen,
	}
}
