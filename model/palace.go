package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Palace struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Image       string       `json:"image" db:"image"`
	EmbededPins []EmbededPin `json:"embededPins"`
}

type EmbededPin struct {
	Number int     `json:"number" db:"number"`
	X      float32 `json:"x" db:"x"`
	Y      float32 `json:"y" db:"y"`
	Word   string  `json:"word" db:"word"`
	Memo   string  `json:"memo" db:"memo"`
}

type palaceImagePath struct {
	path string
}

func GetPalaces(ctx context.Context, userID uuid.UUID) ([]*Palace, error) {
	var palaces []*Palace
	err := db.SelectContext(ctx, &palaces, "SELECT id, name, image FROM palaces WHERE createdBy=? ", userID)
	if err != nil {
		return nil, err
	}

	return palaces, nil
}

func CreatePalace(ctx context.Context, userID uuid.UUID, name, path string) (*uuid.UUID, error) {
	palaceID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO palaces (id, name, createdBy, image) VALUES (?, ?, ?, ?) ", palaceID, name, userID, path)
	if err != nil {
		return nil, err
	}
	return &palaceID, nil
}

func UpdatePalace(ctx context.Context, palaceID uuid.UUID, name, image string) error {
	var count int

	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM palaces WHERE id=?", palaceID)
	if err != nil {
		return err
	}
	if count == 0 {
		// TODO badrequestは返せてるけどメッセージはいってない
		return fmt.Errorf("存在しない宮殿です")
	}
	_, err = db.ExecContext(ctx, "UPDATE palaces SET name=?, image=? WHERE id=? ", name, image, palaceID)
	if err != nil {
		return err
	}
	return nil
}

func DeletePalace(ctx context.Context, palaceID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM palaces WHERE id=? ", palaceID)
	if err != nil {
		return err
	}
	return nil
}

func GetPalaceImagePath(ctx context.Context, palaceID uuid.UUID) (string, error) {
	var path string
	err := db.GetContext(ctx, &path, "SELECT image FROM palaces WHERE id=? ", palaceID)
	if err != nil {
		return "", err
	}
	return path, nil
}
