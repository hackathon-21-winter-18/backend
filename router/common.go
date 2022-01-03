package router

import "github.com/google/uuid"

type Share struct {
	Share     bool      `json:"share"`
	CreatedBy uuid.UUID `json:"createdBy"`
}

type ID struct {
	ID uuid.UUID `json:"id"`
}
