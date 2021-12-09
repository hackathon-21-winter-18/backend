package model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Template struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	Name         string        `json:"name" db:"name"`
	Image        string        `json:"image" db:"image"`
	TemplatePins []TemplatePin `json:"pins"`
}

type TemplatePin struct {
	Number int     `json:"number" db:"number"`
	X      float32 `json:"x" db:"x"`
	Y      float32 `json:"y" db:"y"`
}

func GetAllTemplates(ctx context.Context) ([]*Template, error) {
	var templates []*Template
	err := db.SelectContext(ctx, &templates, "SELECT id, name, image FROM templates")
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func GetTemplates(ctx context.Context, userID uuid.UUID) ([]*Template, error) {
	var templates []*Template
	err := db.SelectContext(ctx, &templates, "SELECT id, name, image FROM templates WHERE createdBy=? ", userID)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func CreateTemplate(ctx context.Context, userID, createdBy uuid.UUID, name, path string) (*uuid.UUID, error) {
	templateID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO templates (id, name, createdBy, image) VALUES (?, ?, ?, ?) ", templateID, name, createdBy, path)
	if err != nil {
		return nil, err
	}
	return &templateID, nil
}

func UpdateTemplate(ctx context.Context, templateID uuid.UUID, name, image string) error {
	var count int

	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM templates WHERE id=?", templateID)
	if err != nil {
		return err
	}
	if count == 0 {
		// TODO badrequestは返せてるけどメッセージはいってない
		return fmt.Errorf("存在しない宮殿です")
	}
	date := time.Now()
	_, err = db.ExecContext(ctx, "UPDATE templates SET name=?, image=?, update_at=? WHERE id=? ", name, image, date, templateID)
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
