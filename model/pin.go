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

func CreateEmbededPin(ctx context.Context, number *int, palaceID uuid.UUID, x, y *float32, word, place, do string) error {
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

func GetTemplatePins(ctx context.Context, TemplateID uuid.UUID) ([]TemplatePin, error) {
	var templatePins []TemplatePin
	err := db.SelectContext(ctx, &templatePins, "SELECT number, x, y FROM templatepins WHERE templateID=? ORDER BY number ASC ", TemplateID)
	if err != nil {
		return nil, err
	}

	return templatePins, nil
}

func CreateTemplatePin(ctx context.Context, number int, templateID uuid.UUID, x, y float32) error {
	_, err := db.ExecContext(ctx, "INSERT INTO templatepins (number, x, y, templateID) VALUES (?, ?, ?, ?) ", number, x, y, templateID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTemplatePins(ctx context.Context, templateID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM templatepins WHERE templateID=? ", templateID)
	if err != nil {
		return err
	}
	return nil
}
