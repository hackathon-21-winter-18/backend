package model

import (
	"context"

	"github.com/google/uuid"
)

func GetEmbededPins(ctx context.Context, PalaceID uuid.UUID) ([]EmbededPin, error) {
	var embededPins []EmbededPin
	err := db.SelectContext(ctx, &embededPins, "SELECT number, x, y, word, place, do FROM embededpins WHERE palaceID=? ORDER BY number ASC ", PalaceID)
	if err != nil {
		return nil, err
	}

	return embededPins, nil
}

func CreateEmbededPin(ctx context.Context, number int, palaceID uuid.UUID, x, y float32, word, place, do string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO embededpins (number, x, y, word, place, do, palaceID) VALUES (?, ?, ?, ?, ?, ?, ?) ", number, x, y, word, place, do, palaceID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEmbededPins(ctx context.Context, palaceID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM embededpins WHERE palaceID=? ", palaceID)
	if err != nil {
		return err
	}
	return nil
}