package model

import (
	"context"
	"errors"
	"strconv"
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
	SavedCount    int          `json:"savedCount" db:"savedCount"`
	SharedAt      time.Time    `db:"shared_at"`
	FirstSharedAt time.Time    `db:"firstshared_at"`
	CreaterName   string       `json:"createrName"`
}

type RequestQuery struct {
	Sort           string
	MaxEmbededPins int
	MinEmbededPins int
}

func GetSharedPalaces(ctx context.Context, requestQuery RequestQuery) ([]*Palace, error) {
	var queryCondition string
	if requestQuery.MaxEmbededPins > 0 {
		queryCondition += " AND number_of_embededPins <= " + strconv.Itoa(requestQuery.MaxEmbededPins)
	}
	if requestQuery.MinEmbededPins > 0 {
		queryCondition += " AND number_of_embededPins >= " + strconv.Itoa(requestQuery.MinEmbededPins)
	}
	if requestQuery.Sort == "first_shared_at" || requestQuery.Sort == "" {
		queryCondition += " ORDER BY firstshared_at DESC"
	} else if requestQuery.Sort == "shared_at" {
		queryCondition += " ORDER BY shared_at DESC"
	} else if requestQuery.Sort == "savedCount" {
		queryCondition += " ORDER BY savedCount DESC"
	} else {
		return nil, errors.New("invalid sort query")
	}

	var palaces []*Palace
	err := db.SelectContext(ctx, &palaces, "SELECT id, originalID, name, createdBy, image, share, savedCount, shared_at, firstshared_at FROM palaces WHERE share=true" + queryCondition)
	if err != nil {
		return nil, err
	}

	for _, palace := range palaces {
		createrName, err := GetMe(ctx, palace.CreatedBy.String())
		if err != nil {
			return nil, err
		}
		palace.CreaterName = createrName
	}

	return palaces, nil
}

func GetMyPalaces(ctx context.Context, userID uuid.UUID, requestQuery RequestQuery) ([]*Palace, error) {
	var queryCondition string
	if requestQuery.MaxEmbededPins > 0 {
		queryCondition += " AND number_of_embededPins <= " + strconv.Itoa(requestQuery.MaxEmbededPins)
	}
	if requestQuery.MinEmbededPins > 0 {
		queryCondition += " AND number_of_embededPins >= " + strconv.Itoa(requestQuery.MinEmbededPins)
	}
	if requestQuery.Sort == "updated_at" || requestQuery.Sort == "" {
		queryCondition += " ORDER BY updated_at DESC"
	} else if requestQuery.Sort == "-updated_at" {
		queryCondition += " ORDER BY updated_at ASC"
	} else {
		return nil, errors.New("invalid sort query")
	}
	
	var palaces []*Palace
	err := db.SelectContext(ctx, &palaces, "SELECT id, originalID,  name, createdBy, image, share, savedCount FROM palaces WHERE heldBy=? " + queryCondition, userID)
	if err != nil {
		return nil, err
	}

	for _, palace := range palaces {
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

func CreatePalace(ctx context.Context, originalID *uuid.UUID, userID uuid.UUID, createdBy *uuid.UUID, name *string, number_of_embededPins int, path string) (*uuid.UUID, error) {
	palaceID := uuid.New()
	if originalID == nil {
		originalID = &palaceID
	}
	date := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO palaces (id, originalID, name, createdBy, heldBy, number_of_embededPins, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ", palaceID, originalID, name, createdBy, userID, number_of_embededPins, path, date, date)
	if err != nil {
		return nil, err
	}
	return &palaceID, nil
}

func UpdatePalace(ctx context.Context, palaceID uuid.UUID, name *string, number_of_embededPins int, image string) error {
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
	_, err = db.ExecContext(ctx, "UPDATE palaces SET name=?, number_of_embededPins=?, image=?, updated_at=? WHERE id=? ", name, number_of_embededPins, image, date, palaceID)
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

	var savedCount int
	err = db.GetContext(ctx, &savedCount, "SELECT COUNT(*) FROM palace_user WHERE palaceID=? ", palaceID)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "UPDATE palaces SET savedCount=? WHERE id=? ", savedCount + 1, palaceID)
	if err != nil {
		return err
	}

	return nil
}
