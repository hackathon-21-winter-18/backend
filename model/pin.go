package model

import (
	"context"

	"github.com/google/uuid"
)

type Pin struct {
	Id uuid.UUID `db:"id" json:"id"`
	X  float32   `db:"x" json:"x"`
	Y  float32   `db:"y" json:"y"`
}

func CreateTemplatePins(ctx context.Context, templateid uuid.UUID, pins []Pin) error {
	for _, v := range pins {
		pinid := uuid.New()
		_, err := db.ExecContext(ctx, "INSERT INTO pin (id, x, y) VALUES (?, ?, ?)", pinid, v.X, v.Y)
		if err != nil {
			return err
		}
		_, err = db.ExecContext(ctx, "INSERT INTO template_pins (template, pin) VALUES (?, ?)", templateid, pinid)
		if err != nil {
			return err
		}
	}
	return nil
}
