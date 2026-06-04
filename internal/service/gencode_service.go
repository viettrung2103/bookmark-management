package service

import "github.com/viettrung2103/bookmark-management/pkg/stringutils"

const (
	passLength = 10
)

// GenPass represents the genpass service
//
//go:generate mockery --name=Code --filename=gencode.go
type Code interface {
	GenerateCode() (string, error)
}

type codeService struct{}

// NewGenPass return a GenPassService
func NewCode() Code {
	return &codeService{}
}

// GeneratePassword generates a random password
func (s *codeService) GenerateCode() (string, error) {

	return stringutils.GenerateCode(passLength)

}
