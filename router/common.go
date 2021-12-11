package router

import "github.com/google/uuid"

type Share struct {
	Share bool `json:"share"`
}

type ID struct {
	ID uuid.UUID `json:"id"`
}
