package model

import (
	"context"
	"fmt"
	"time"

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
	Place  string  `json:"place" db:"place"`
	Do     string  `json:"do" db:"do"`
}

type palaceImagePath struct {
	path string
}

type firstShared struct {
	FirstShared bool `db:"firstshared"`
}

func GetPalaces(ctx context.Context, userID uuid.UUID) ([]*Palace, error) {
	var palaces []*Palace
	err := db.SelectContext(ctx, &palaces, "SELECT id, name, image FROM palaces WHERE createdBy=? ", userID)
	if err != nil {
		return nil, err
	}

	return palaces, nil
}

func CreatePalace(ctx context.Context, userID, createdBy uuid.UUID, name, path string) (*uuid.UUID, error) {
	palaceID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO palaces (id, name, createdBy, heldBy, image) VALUES (?, ?, ?, ?, ?) ", palaceID, name, createdBy, userID, path)
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

func SharePalace(ctx context.Context, palaceID uuid.UUID, share bool) error {
	var firstShared firstShared
	if share {
		err := db.GetContext(ctx, &firstShared, "SELECT firstshared FROM palaces WHERE id=? ", palaceID)
		if err != nil {
			return err
		}
		if firstShared.FirstShared {
			date := time.Now()
			_, err := db.ExecContext(ctx, "UPDATE palaces SET share=?, shared_at=? WHERE id=? ", share, date, palaceID)
			if err != nil {
				return err
			}
		} else {
			date := time.Now()
			_, err := db.ExecContext(ctx, "UPDATE palaces SET share=true, firstshared=true, firstshared_at=?, shared_at=? WHERE id=? ", date, date, palaceID)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := db.ExecContext(ctx, "UPDATE palaces SET share=false WHERE id=? ", palaceID)
		if err != nil {
			return err
		}
	}

	return nil
}

func Location() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
