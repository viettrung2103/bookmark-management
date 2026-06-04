package service

import "github.com/google/uuid"

// GenId interface implements by GenIdService
//
//go:generate mockery --name=GenId --filename=genid.go
type GenId interface {
	GenerateId() string
}

// GenIdService struct implements GenId interface
type GenIdService struct {
}

// NewGenId function implements GenId interface
func NewGenId() GenId {
	return &GenIdService{}
}

// GenerateId function implements GenId interface
func (s *GenIdService) GenerateId() string {

	id := uuid.New()

	return id.String()

}
