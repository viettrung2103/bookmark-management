package link

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	//repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/mocks"
	repoMocks "github.com/viettrung2103/bookmark-management/internal/app/repository/url/mocks"
)

var redisTestErr = errors.New("test error")

// TestService_GetLinkFromKey tests the GetLinkFromKey method
func TestService_GetLinkFromKey(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRepo   func(ctx context.Context) *repoMocks.URLRepository
		expectedUrl string
		expectedErr error
	}{
		{
			name: "normal case",
			setupRepo: func(ctx context.Context) *repoMocks.URLRepository {
				mock := repoMocks.NewURLRepository(t)
				mock.On("GetURL", ctx, "test").Return("https://test.com", nil)
				return mock
			},

			expectedUrl: "https://test.com",
			expectedErr: nil,
		},
		{
			name: "empty case",
			setupRepo: func(ctx context.Context) *repoMocks.URLRepository {
				mock := repoMocks.NewURLRepository(t)
				mock.On("GetURL", ctx, "test").Return("", redisTestErr)
				return mock
			},
			expectedUrl: "",
			expectedErr: redisTestErr,
		},
		{
			name: "err case ",
			setupRepo: func(ctx context.Context) *repoMocks.URLRepository {
				mock := repoMocks.NewURLRepository(t)
				mock.On("GetURL", ctx, "test").Return("", redisTestErr)
				return mock
			},

			expectedUrl: "",
			expectedErr: redisTestErr,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockRepo := tc.setupRepo(ctx)
			testService := NewService(mockRepo, nil)
			result, err := testService.GetLinkFromCode(ctx, "test")
			assert.Equal(t, result, tc.expectedUrl)
			assert.ErrorIs(t, err, tc.expectedErr)

		})
	}
}
