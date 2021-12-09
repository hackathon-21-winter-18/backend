package model

import (
	"errors"
)

var (
	ErrNotFound = errors.New("Not Found")
	ErrNoChange = errors.New("No Change")
	ErrForbidden = errors.New("Forbidden")
)