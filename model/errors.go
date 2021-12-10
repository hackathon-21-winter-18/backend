package model

import (
	"errors"
)

var (
	ErrNotFound = errors.New("Not Found")
	ErrForbidden = errors.New("Forbidden")
	ErrUnauthorized = errors.New("Unauthorized")
)