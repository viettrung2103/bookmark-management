package dbutils

import (
	"errors"
	"strings"
)

var errorFilters = []func(err error) (bool, error){
	filterDuplicationError,
	filterForeignKeyError,
	filterRecordNotFoundErr,
}

// CatchDBError catches database errors and returns filtered errors
func CatchDBError(err error) error {
	if err == nil {
		return nil
	}
	for _, currentFilter := range errorFilters {
		match, filteredError := currentFilter(err)
		if match {
			return filteredError
		}
	}
	return err
}

var (
	ErrDuplication    = errors.New("duplicated column")
	ErrForeignKey     = errors.New("foreign key constraint")
	ErrRecordNotFound = errors.New("record not found")
)

func filterDuplicationError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "unique constraint"), ErrDuplication
}

func filterForeignKeyError(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "foreign key constraint"), ErrForeignKey
}

func filterRecordNotFoundErr(err error) (bool, error) {
	return strings.Contains(strings.ToLower(err.Error()), "record not found"), ErrRecordNotFound
}
