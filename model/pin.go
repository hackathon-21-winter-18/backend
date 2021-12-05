package model

import (
	"context"

	"github.com/google/uuid"
)

func CreateEmbededPin(ctx context.Context, userID uuid.UUID, x, y float32, word string, memo string) error {
	embededPinID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO pins ()")
}