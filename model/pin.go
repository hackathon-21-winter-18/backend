package model

import (
	"context"

	"github.com/google/uuid"
)

func CreateEmbededPin(ctx context.Context, number int, palaceID uuid.UUID, x, y float32, word string, memo string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO embededpins (number, x, y, word, memo, palaceID) VALUES (?, ?, ?, ?, ?, ?) ", number, x, y, word, memo, palaceID)
	if err != nil {
		return err
	}

	return nil
}