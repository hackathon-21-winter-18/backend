package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Palace struct {
	ID          uuid.UUID    `json:"id" db:"id"`
	Name        string       `json:"name" db:"name"`
	Image       string       `json:"image" db:"image"`
	EmbededPins []EmbededPin `json:"embededPins"`
	Share       bool         `json:"share" db:"share"`
}

type firstShared struct {
	FirstShared bool `db:"firstshared"`
}

type heldBy struct {
	heldBy uuid.UUID `db:"heldBy"`
}

func GetPalaces(ctx context.Context, userID uuid.UUID) ([]*Palace, error) {
	var palaces []*Palace
	err := db.SelectContext(ctx, &palaces, "SELECT id, name, image, share FROM palaces WHERE heldBy=? ", userID)
	if err != nil {
		return nil, err
	}

	return palaces, nil
}

func CreatePalace(ctx context.Context, userID uuid.UUID, createdBy *uuid.UUID, name *string, path string) (*uuid.UUID, error) {
	palaceID := uuid.New()
	date := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO palaces (id, name, createdBy, heldBy, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) ", palaceID, name, createdBy, userID, path, date, date)
	if err != nil {
		return nil, err
	}
	return &palaceID, nil
}

func UpdatePalace(ctx context.Context, palaceID uuid.UUID, name *string, image string) error {
	var count int
	// TODO なくてもよさそう
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM palaces WHERE id=?", palaceID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
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

func CheckPalaceHeldBy(ctx context.Context, userID, palaceID uuid.UUID) error {
	var heldBy heldBy
	err := db.GetContext(ctx, &heldBy, "SELECT heldBy FROM palaces WHERE id=? ", palaceID)
	if err != nil {
		return err
	}

	if heldBy.heldBy != userID {
		return ErrUnauthorized
	}

	return nil
}
