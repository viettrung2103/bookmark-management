package service

import "github.com/viettrung2103/bookmark-management/pkg/stringutils"

const (
	passLength = 10
)

// GenPass represents the genpass service
//
//go:generate mockery --name=GenCode --filename=gencode.go
type GenCode interface {
	GenerateCode() string
}

type codeService struct{}

// NewGenPass return a GenPassService
func NewCode() GenCode {
	return &codeService{}
}

// GeneratePassword generates a random password
func (s *codeService) GenerateCode() string {

	return stringutils.GenerateRandomString(passLength)

}
