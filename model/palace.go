package model

import (
	"context"

	"github.com/google/uuid"
)

func CreatePalace(ctx context.Context, userID uuid.UUID, name string, image string) (*uuid.UUID, error) {
	palaceID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO palaces (id, name, createdBy, image) VALUES (?, ?, ?, ?) ", palaceID, name, userID, image)
	if err != nil {
		return nil, err
	}

	return &palaceID, nil
}
