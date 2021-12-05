package model

import (
	"context"

	"github.com/google/uuid"
)

type Template struct {
	Id        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Image     []byte    `db:"image"`
	CreatedBy uuid.UUID `db:"createdBy"`
}

type TemplatePins struct {
	Template uuid.UUID `db:"template"`
	Pin      uuid.UUID `db:"pin"`
}

type Pin struct {
	Id uuid.UUID `db:"id"`
	X  float32   `db:"x"`
	Y  float32   `db:"y"`
}

func CreateTemplate(ctx context.Context, name string, image string, pins []Pin, createdby uuid.UUID) (*uuid.UUID, error) {
	templateID := uuid.New()
	query := "INSERT INTO template (id, name, image, createdBy) VALUES (?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, query, templateID, name, image, createdby)
	if err != nil {
		return nil, err
	}
	for _, v := range pins {
		pinid := uuid.New()
		queryPin := "INSERT INTO pin (id, x, y) VALUES (?, ?, ?)"
		_, err := db.ExecContext(ctx, queryPin, pinid, v.X, v.Y)
		if err != nil {
			return nil, err
		}
		queryTeamplatePin := "INSERT INTO template_pins (template, pin) VALUES (?, ?)"
		_, err = db.ExecContext(ctx, queryTeamplatePin, templateID, pinid)
		if err != nil {
			return nil, err
		}
	}
	return &templateID, nil
}
