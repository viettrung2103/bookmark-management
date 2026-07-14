package test_data

import (
	"github.com/google/uuid"
	"github.com/viettrung2103/bookmark-management/internal/app/model"
)

// Global test variables/constants accessible to all tests in this package
const (
	TestUUIDText       = "12345678-1234-1234-1234-123456789012"
	TestUsername       = "testuser"
	TestPassword       = "plainpassword"
	TestHashedPassword = "hashedpassword"
	TestExpectedToken  = "fake-jwt-token"
)

var (
	TestUUID = uuid.MustParse(TestUUIDText)

	// Default mock user helper
	TestUser = &model.User{
		Base: model.Base{
			ID: TestUUID,
		},
		Username: TestUsername,
		Password: TestHashedPassword,
	}
)
