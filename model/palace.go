package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Palace struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	OriginalID    uuid.UUID    `json:"originalID" db:"originalID"`
	Name          string       `json:"name" db:"name"`
	CreatedBy     uuid.UUID    `json:"createdBy" db:"createdBy"`
	Image         string       `json:"image" db:"image"`
	EmbededPins   []EmbededPin `json:"embededPins"`
	Share         bool         `json:"share" db:"share"`
	SharedAt      time.Time    `db:"shared_at"`
	FirstSharedAt time.Time    `db:"firstshared_at"`
	SavedCount    int          `json:"savedCount"`
	CreaterName   string       `json:"createrName"`
}

func GetSharedPalaces(ctx context.Context) ([]*Palace, error) {
	var palaces []*Palace
	err := db.SelectContext(ctx, &palaces, "SELECT id, originalID, name, createdBy, image, share, shared_at, firstshared_at FROM palaces WHERE share=true")
	if err != nil {
		return nil, err
	}

	for _, palace := range palaces {
		savedCount, err := GetPalaceSavedCount(ctx, palace.ID)
		if err != nil {
			return nil, err
		}
		palace.SavedCount = *savedCount

		createrName, err := GetMe(ctx, palace.CreatedBy.String())
		if err != nil {
			return nil, err
		}
		palace.CreaterName = createrName
	}

	return palaces, nil
}

func GetMyPalaces(ctx context.Context, userID uuid.UUID) ([]*Palace, error) {
	var palaces []*Palace
	err := db.SelectContext(ctx, &palaces, "SELECT id, originalID,  name, createdBy, image, share FROM palaces WHERE heldBy=? ", userID)
	if err != nil {
		return nil, err
	}

	for _, palace := range palaces {
		savedCount, err := GetPalaceSavedCount(ctx, palace.ID)
		if err != nil {
			return nil, err
		}
		palace.SavedCount = *savedCount

		createrName, err := GetMe(ctx, palace.CreatedBy.String())
		if err != nil {
			return nil, err
		}
		palace.CreaterName = createrName
	}

	return palaces, nil
}

func GetPalace(ctx context.Context, palaceID uuid.UUID) (*Palace, error) {
	var palace Palace
	err := db.GetContext(ctx, &palace, "SELECT id, originalID, name, createdBy, image, share FROM palaces WHERE id=? ", palaceID)
	if err != nil {
		return nil, err
	}

	savedCount, err := GetPalaceSavedCount(ctx, palace.ID)
	if err != nil {
		return nil, err
	}
	palace.SavedCount = *savedCount

	createrName, err := GetMe(ctx, palace.CreatedBy.String())
	if err != nil {
		return nil, err
	}
	palace.CreaterName = createrName

	return &palace, nil
}

func CreatePalace(ctx context.Context, originalID *uuid.UUID, userID uuid.UUID, createdBy *uuid.UUID, name *string, path string) (*uuid.UUID, error) {
	palaceID := uuid.New()
	if originalID == nil {
		originalID = &palaceID
	}
	date := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO palaces (id, originalID, name, createdBy, heldBy, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ", palaceID, originalID, name, createdBy, userID, path, date, date)
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
	date := time.Now()
	_, err = db.ExecContext(ctx, "UPDATE palaces SET name=?, image=?, updated_at=? WHERE id=? ", name, image, date, palaceID)
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

	if heldBy.HeldBy != userID {
		return ErrUnauthorized
	}

	return nil
}

func GetPalaceSavedCount(ctx context.Context, palaceID uuid.UUID) (*int, error) {
	var savedCount int
	err := db.GetContext(ctx, &savedCount, "SELECT COUNT(*) FROM palace_user WHERE palaceID=? ", palaceID)
	if err != nil {
		return nil, err
	}

	return &savedCount, nil
}

func RecordPalaceSavingUser(ctx context.Context, palaceID, userID uuid.UUID) error {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM palace_user WHERE palaceID=? AND userID=? ", palaceID, userID)
	if err != nil {
		return nil
	}
	if count > 0 {
		return nil
	}

	_, err = db.ExecContext(ctx, "INSERT INTO palace_user (palaceID, userID) VALUES (?, ?) ", palaceID, userID)
	if err != nil {
		return err
	}

	return nil
}
