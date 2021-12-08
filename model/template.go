package model

import (
	"context"

	"github.com/google/uuid"
)

type Template struct {
	ID           uuid.UUID     `json:"id" db:"id"`
	Name         string        `json:"name" db:"name"`
	Image        string        `json:"image" db:"image"`
	TemplatePins []TemplatePin `json:"templatePins"`
}

type TemplatePin struct {
	Number int     `json:"number" db:"number"`
	X      float32 `json:"x" db:"x"`
	Y      float32 `json:"y" db:"y"`
}

func CreateTemplate(ctx context.Context, userID, createdBy uuid.UUID, name, path string) (*uuid.UUID, error) {
	templateID := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO templates (id, name, createdBy, image) VALUES (?, ?, ?, ?) ", templateID, name, createdBy, path)
	if err != nil {
		return nil, err
	}
	return &templateID, nil
}
