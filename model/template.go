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

func CreateTemplate(ctx context.Context, name string, image string, pins []Pin, createdby uuid.UUID) (*uuid.UUID, error) {
	templateID := uuid.New()
	query := "INSERT INTO template (id, name, image, createdBy) VALUES (?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, query, templateID, name, image, createdby)
	if err != nil {
		return nil, err
	}
	err = CreateTemplatePins(ctx, templateID, pins)
	if err != nil {
		return nil, err
	}
	return &templateID, nil
}

// func GetTemplateFromId(ctx context.Context, id uuid.UUID) (*router.Template, error) {
// 	var template Template
// 	query := "SELECT * FROM template WHERE id=?"
// 	err := db.Get(&template, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var pins []Pin
// 	query = "SELECT * FROM pin WHERE template IN (SELECT id FROM template_pins WHERE id=?)"
// 	err = db.Get(&pins, query, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	res := router.Template{
// 		Id:        id,
// 		Name:      template.Name,
// 		Image:     string(template.Image),
// 		CreatedBy: template.CreatedBy,
// 	}
// 	res.Pins = pins
// }
