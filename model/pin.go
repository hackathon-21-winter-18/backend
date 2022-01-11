package model

import (
	"context"

	"github.com/google/uuid"
)

type EmbededPin struct {
	Number      *int     `json:"number,omitempty" db:"number"`
	X           *float32 `json:"x,omitempty" db:"x"`
	Y           *float32 `json:"y,omitempty" db:"y"`
	Word        string   `json:"word" db:"word"`
	Place       string   `json:"place" db:"place"`
	Situation   string   `json:"situation" db:"situation"`
	GroupName   string   `json:"groupName" db:"groupName"`
	GroupNumber int      `json:"groupNumber" db:"groupNumber"`
}

type Pin struct {
	Number *int     `json:"number,omitempty" db:"number"`
	X      *float32 `json:"x" db:"x"`
	Y      *float32 `json:"y" db:"y"`
}

func GetEmbededPins(ctx context.Context, PalaceID uuid.UUID) ([]EmbededPin, error) {
	var embededPins []EmbededPin
	err := db.SelectContext(ctx, &embededPins, "SELECT * FROM embededpins WHERE palaceID=? ORDER BY number ASC ", PalaceID)
	if err != nil {
		return nil, err
	}

	return embededPins, nil
}

func CreateEmbededPin(ctx context.Context, number *int, palaceID uuid.UUID, x, y *float32, word, place, condition string, groupName string, groupNumber int) error {
	_, err := db.ExecContext(ctx, "INSERT INTO embededpins (number, x, y, word, place, situation, palaceID, groupName, groupNumber) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ", number, x, y, word, place, condition, palaceID, groupName, groupNumber)
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

func GetPins(ctx context.Context, templateID uuid.UUID) ([]Pin, error) {
	var templatePins []Pin
	err := db.SelectContext(ctx, &templatePins, "SELECT number, x, y FROM pins WHERE templateID=? ORDER BY number ASC ", templateID)
	if err != nil {
		return nil, err
	}

	return templatePins, nil
}

func CreatePin(ctx context.Context, number *int, templateID uuid.UUID, x, y *float32) error {
	_, err := db.ExecContext(ctx, "INSERT INTO pins (number, x, y, templateID) VALUES (?, ?, ?, ?) ", number, x, y, templateID)
	if err != nil {
		return err
	}
	return nil
}

func DeletePins(ctx context.Context, templateID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM pins WHERE templateID=? ", templateID)
	if err != nil {
		return err
	}
	return nil
}
