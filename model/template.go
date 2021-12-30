package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID            uuid.UUID `json:"id" db:"id"`
	OriginalID    uuid.UUID `json:"originalID" db:"originalID"`
	Name          string    `json:"name" db:"name"`
	CreatedBy     uuid.UUID `json:"createdBy" db:"createdBy"`
	Image         string    `json:"image" db:"image"`
	Pins          []Pin     `json:"pins"`
	Share         bool      `json:"share" db:"share"`
	SharedAt      time.Time `db:"shared_at"`
	FirstSharedAt time.Time `db:"firstshared_at"`
	SavedCount    int       `json:"savedCount"`
	CreaterName   string    `json:"createrName"`
}

func GetSharedTemplates(ctx context.Context) ([]*Template, error) {
	var templates []*Template
	// if sort == "" || sort == "" {

	// }
	err := db.SelectContext(ctx, &templates, "SELECT id, originalID, name, createdBy, image, share, shared_at, firstshared_at FROM templates WHERE share=true")
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		savedCount, err := GetPalaceSavedCount(ctx, template.ID)
		if err != nil {
			return nil, err
		}
		template.SavedCount = *savedCount
		
		createrName, err := GetMe(ctx, template.CreatedBy.String())
		if err != nil {
			return nil, err
		}
		template.CreaterName = createrName
	}

	return templates, nil
}

func GetMyTemplates(ctx context.Context, userID uuid.UUID) ([]*Template, error) {
	var templates []*Template
	err := db.SelectContext(ctx, &templates, "SELECT id, originalID, name, createdBy, image, share FROM templates WHERE heldBy=? ", userID)
	if err != nil {
		return nil, err
	}

	for _, template := range templates {
		savedCount, err := GetPalaceSavedCount(ctx, template.ID)
		if err != nil {
			return nil, err
		}
		template.SavedCount = *savedCount
		
		createrName, err := GetMe(ctx, template.CreatedBy.String())
		if err != nil {
			return nil, err
		}
		template.CreaterName = createrName
	}

	return templates, nil
}

func GetTemplate(ctx context.Context, templateID uuid.UUID) (*Template, error) {
	var template Template
	err := db.GetContext(ctx, &template, "SELECT id, originalID, name, createdBy, image, share FROM templates WHERE id=? ", templateID)
	if err != nil {
		return nil, err
	}

	savedCount, err := GetPalaceSavedCount(ctx, template.ID)
	if err != nil {
		return nil, err
	}
	template.SavedCount = *savedCount

	createrName, err := GetMe(ctx, template.CreatedBy.String())
	if err != nil {
		return nil, err
	}
	template.CreaterName = createrName

	return &template, nil
}

func CreateTemplate(ctx context.Context, originalID *uuid.UUID, userID uuid.UUID, createdBy *uuid.UUID, name *string, path string) (*uuid.UUID, error) {
	templateID := uuid.New()
	if originalID == nil {
		originalID = &templateID
	}
	date := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO templates (id, originalID, name, createdBy, heldBy, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?) ", templateID, originalID, name, createdBy, userID, path, date, date)
	if err != nil {
		return nil, err
	}
	return &templateID, nil
}

func UpdateTemplate(ctx context.Context, templateID uuid.UUID, name *string, image string) error {
	var count int

	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM templates WHERE id=?", templateID)
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNotFound
	}
	date := time.Now()
	_, err = db.ExecContext(ctx, "UPDATE templates SET name=?, image=?, updated_at=? WHERE id=? ", name, image, date, templateID)
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

	return nil
}