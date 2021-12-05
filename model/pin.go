package model

import (
	"context"

	"github.com/google/uuid"
)

func CreateEmbededPin(ctx context.Context, palaceID uuid.UUID, x, y float32, word string, memo string) error {
	embededPinID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO embededpins (id, x, y, word, memo, palaceID) VALUES(?, ?, ?, ?, ?, ?) ", embededPinID, x, y, word, memo, palaceID)
	if err != nil {
		return err
	}

	return nil
}