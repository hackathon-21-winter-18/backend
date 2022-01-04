package model

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID            uuid.UUID `json:"id" db:"id"`
	OriginalID    uuid.UUID `json:"originalID" db:"originalID"`
	Name          string    `json:"name" db:"name"`
	CreatedBy     uuid.UUID `json:"createdBy" db:"createdBy"`
	Image         string    `json:"image" db:"image"`
	HeldBy        uuid.UUID `json:"heldBy" db:"heldBy"`
	Pins          []Pin     `json:"pins"`
	Share         bool      `json:"share" db:"share"`
	SharedAt      time.Time `db:"shared_at"`
	FirstSharedAt time.Time `db:"firstshared_at"`
	SavedCount    int       `json:"savedCount" db:"savedCount"`
	CreatorName   string    `json:"creatorName"`
	EditorName    string    `json:"editorName"`
}

func GetSharedTemplates(ctx context.Context, requestQuery RequestQuery) ([]*Template, error) {
	var queryCondition string
	if requestQuery.MaxEmbededPins > 0 {
		queryCondition += " AND number_of_pins <= " + strconv.Itoa(requestQuery.MaxEmbededPins)
	}
	if requestQuery.MinEmbededPins > 0 {
		queryCondition += " AND number_of_pins >= " + strconv.Itoa(requestQuery.MinEmbededPins)
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

	var templates []*Template
	err := db.SelectContext(ctx, &templates, "SELECT id, originalID, name, createdBy, image, heldBy, share, savedCount, shared_at, firstshared_at FROM templates WHERE share=true"+queryCondition)
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		creatorName, err := GetMe(ctx, template.CreatedBy.String())
		if err != nil {
			return nil, err
		}
		template.CreatorName = creatorName

		editorName, err := GetMe(ctx, template.HeldBy.String())
		if err != nil {
			return nil, err
		}
		template.EditorName = editorName
	}

	return templates, nil
}

func GetMyTemplates(ctx context.Context, userID uuid.UUID, requestQuery RequestQuery) ([]*Template, error) {
	var queryCondition string
	if requestQuery.MaxEmbededPins > 0 {
		queryCondition += " AND number_of_pins <= " + strconv.Itoa(requestQuery.MaxEmbededPins)
	}
	if requestQuery.MinEmbededPins > 0 {
		queryCondition += " AND number_of_pins >= " + strconv.Itoa(requestQuery.MinEmbededPins)
	}
	if requestQuery.Sort == "updated_at" || requestQuery.Sort == "" {
		queryCondition += " ORDER BY updated_at DESC"
	} else if requestQuery.Sort == "-updated_at" {
		queryCondition += " ORDER BY updated_at ASC"
	} else {
		return nil, errors.New("invalid sort query")
	}

	var templates []*Template
	err := db.SelectContext(ctx, &templates, "SELECT id, originalID, name, createdBy, image, share, savedCount FROM templates WHERE heldBy=? "+queryCondition, userID)
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		savedCount, err := GetTemplateSavedCount(ctx, template.ID)
		if err != nil {
			return nil, err
		}
		template.SavedCount = *savedCount

		creatorName, err := GetMe(ctx, template.CreatedBy.String())
		if err != nil {
			return nil, err
		}
		template.CreatorName = creatorName
	}

	return templates, nil
}

func GetTemplate(ctx context.Context, templateID uuid.UUID) (*Template, error) {
	var template Template
	err := db.GetContext(ctx, &template, "SELECT id, originalID, name, createdBy, image, share FROM templates WHERE id=? ", templateID)
	if err != nil {
		return nil, err
	}

	savedCount, err := GetTemplateSavedCount(ctx, template.ID)
	if err != nil {
		return nil, err
	}
	template.SavedCount = *savedCount

	creatorName, err := GetMe(ctx, template.CreatedBy.String())
	if err != nil {
		return nil, err
	}
	template.CreatorName = creatorName

	return &template, nil
}

func CreateTemplate(ctx context.Context, originalID *uuid.UUID, userID uuid.UUID, createdBy *uuid.UUID, name *string, number_of_pins int, path string) (*uuid.UUID, error) {
	templateID := uuid.New()
	if originalID == nil {
		originalID = &templateID
	}
	date := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO templates (id, originalID, name, createdBy, heldBy, number_of_pins, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ", templateID, originalID, name, createdBy, userID, number_of_pins, path, date, date)
	if err != nil {
		return nil, err
	}
	return &templateID, nil
}

func UpdateTemplate(ctx context.Context, templateID uuid.UUID, name *string, number_of_pins int, image string) error {
	var count int

	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM templates WHERE id=?", templateID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	date := time.Now()
	_, err = db.ExecContext(ctx, "UPDATE templates SET name=?, number_of_Pins=?, image=?, updated_at=? WHERE id=? ", name, number_of_pins, image, date, templateID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM templates WHERE id=? ", templateID)
	if err != nil {
		return err
	}
	return nil
}

func ShareTemplate(ctx context.Context, templateID uuid.UUID, share bool) error {
	var firstShared firstShared
	if share {
		err := db.GetContext(ctx, &firstShared, "SELECT firstshared FROM templates WHERE id=? ", templateID)
		if err != nil {
			return err
		}
		if firstShared.FirstShared {
			date := time.Now()
			_, err := db.ExecContext(ctx, "UPDATE templates SET share=?, shared_at=? WHERE id=? ", share, date, templateID)
			if err != nil {
				return err
			}
		} else {
			date := time.Now()
			_, err := db.ExecContext(ctx, "UPDATE templates SET share=true, firstshared=true, firstshared_at=?, shared_at=? WHERE id=? ", date, date, templateID)
			if err != nil {
				return err
			}
		}
	} else {
		_, err := db.ExecContext(ctx, "UPDATE templates SET share=false WHERE id=? ", templateID)
		if err != nil {
			return err
		}
	}

	return nil
}

func CheckTemplateHeldBy(ctx context.Context, userID, templateID uuid.UUID) error {
	var heldBy heldBy
	err := db.GetContext(ctx, &heldBy, "SELECT heldBy FROM templates WHERE id=? ", templateID)
	if err != nil {
		return err
	}

	if heldBy.HeldBy != userID {
		return ErrUnauthorized
	}

	return nil
}

func GetTemplateSavedCount(ctx context.Context, templateID uuid.UUID) (*int, error) {
	var savedCount int
	err := db.GetContext(ctx, &savedCount, "SELECT COUNT(*) FROM template_user WHERE templateID=? ", templateID)
	if err != nil {
		return nil, err
	}

	return &savedCount, nil
}

func RecordTemplateSavingUser(ctx context.Context, templateID, userID uuid.UUID) error {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM template_user WHERE templateID=? AND userID=? ", templateID, userID)
	if err != nil {
		return nil
	}
	if count > 0 {
		return nil
	}

	_, err = db.ExecContext(ctx, "INSERT INTO template_user (templateID, userID) VALUES (?, ?) ", templateID, userID)
	if err != nil {
		return err
	}

	var savedCount int
	err = db.GetContext(ctx, &savedCount, "SELECT COUNT(*) FROM template_user WHERE templateID=? ", templateID)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, "UPDATE templates SET savedCount=? WHERE id=? ", savedCount, templateID)
	if err != nil {
		return err
	}

	return nil
}
