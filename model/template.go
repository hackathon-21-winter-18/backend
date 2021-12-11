package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	Name          string        `json:"name" db:"name"`
	Image         string        `json:"image" db:"image"`
	TemplatePins  []TemplatePin `json:"pins"`
	Share         bool          `json:"share" db:"share"`
	SharedAt      time.Time     `db:"shared_at"`
	FirstSharedAt time.Time     `db:"firstshared_at"`
}

func GetTemplates(ctx context.Context) ([]*Template, error) {
	var templates []*Template
	// if sort == "" || sort == "" {

	// }
	err := db.SelectContext(ctx, &templates, "SELECT id, name, image, shared_at, firstshared_at FROM templates WHERE share=true")
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func GetMyTemplates(ctx context.Context, userID uuid.UUID) ([]*Template, error) {
	var templates []*Template
	err := db.SelectContext(ctx, &templates, "SELECT id, name, image, share FROM templates WHERE heldBy=? ", userID)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func CreateTemplate(ctx context.Context, userID uuid.UUID, createdBy *uuid.UUID, name *string, path string) (*uuid.UUID, error) {
	templateID := uuid.New()
	date := time.Now()
	_, err := db.ExecContext(ctx, "INSERT INTO templates (id, name, createdBy, heldBy, image, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?) ", templateID, name, createdBy, userID, path, date, date)
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

func GetTemplateImagePath(ctx context.Context, templateID uuid.UUID) (string, error) {
	var path string
	err := db.GetContext(ctx, &path, "SELECT image FROM templates WHERE id=? ", templateID)
	if err != nil {
		return "", err
	}
	return path, nil
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
