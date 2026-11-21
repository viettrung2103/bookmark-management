package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/viettrung2103/bookmark-lib/pkg/requestutils"
	"github.com/viettrung2103/bookmark-management/internal/app/ratelimit"
)

type RateLimit interface {
	RateLimitPublic() gin.HandlerFunc  // byRequest
	RateLimitPrivate() gin.HandlerFunc // byUser
}

type rateLimit struct {
	repo ratelimit.Repository
}

func NewRateLimit(repo ratelimit.Repository) RateLimit {
	return &rateLimit{repo: repo}
}

const (
	rateLimitKeyFormat = "rate_limit:%s"
	rateLimitMaxCount  = 10
	rateLimitDuration  = 1 * time.Minute
)

func (r *rateLimit) RateLimitPublic() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create rate limit key
		// c.Request.RemoteAddr is user ip

		key := fmt.Sprintf(rateLimitKeyFormat, c.Request.RemoteAddr)

		// get current rate limit
		curRate, err := r.repo.GetCurrentRateLimit(c, key)
		if err != nil {
			log.Err(err).Str("key", key).Msg("get current rate limit")
		}

		// if rate lmit is exceeded, return 429
		if curRate > rateLimitMaxCount {
			c.AbortWithStatusJSON(429, gin.H{"Message": "Rate limit Excedede"})
			return
		}

		// if rate limit is not exceeded. continue
		r.repo.IncreaseRateLimit(c, key, rateLimitDuration)

		// increase rate limit
		c.Next()

	}
}

func (r *rateLimit) RateLimitPrivate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create rate limit key
		// c.Request.RemoteAddr is user ip
		// next is ban based on user id, not user ip
		uid, err := requestutils.GetUserIDFromRequest(c)
		if err != nil {
			// neu ko lay dc user, bi loi khi lay jwt, la loi nghiem trong
			return
		}

		key := fmt.Sprintf(rateLimitKeyFormat, uid)

		// get current rate limit
		curRate, err := r.repo.GetCurrentRateLimit(c, key)
		if err != nil {
			log.Err(err).Str("key", key).Msg("get current rate limit")
		}

		// if rate lmit is exceeded, return 429
		if curRate > rateLimitMaxCount {
			c.AbortWithStatusJSON(429, gin.H{"Message": "Rate limit Excedede"})
			return
		}

		// if rate limit is not exceeded. continue
		r.repo.IncreaseRateLimit(c, key, rateLimitDuration)

		// increase rate limit
		c.Next()

	}
}
