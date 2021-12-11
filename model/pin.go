package model

import (
	"context"
	"sort"
	"strconv"

	"github.com/google/uuid"
)

type EmbededPin struct {
	Number *int     `json:"number,omitempty" db:"number"`
	X      *float32 `json:"x,omitempty" db:"x"`
	Y      *float32 `json:"y,omitempty" db:"y"`
	Word   string   `json:"word" db:"word"`
	Place  string   `json:"place" db:"place"`
	Do     string   `json:"do" db:"do"`
}

type TemplatePin struct {
	Number *int     `json:"number,omitempty" db:"number"`
	X      *float32 `json:"x" db:"x"`
	Y      *float32 `json:"y" db:"y"`
}

func ExtractFromPalacesBasedOnEmbededPins(palaces []*Palace, max, min string) []*Palace {
	sort.Slice(palaces, func(i, j int) bool {
		pini := len(palaces[i].EmbededPins)
		pinj := len(palaces[j].EmbededPins)
		return pini < pinj
	})
	minptr := 0
	maxptr := len(palaces)
	for i, v := range palaces {
		pin := len(v.EmbededPins)
		if minpin, err := strconv.Atoi(min); err == nil && pin < minpin {
			minptr = i + 1
		}
		if maxpin, err := strconv.Atoi(max); err == nil && pin > maxpin {
			maxptr = i
			break
		}
	}
	palaces = palaces[minptr:maxptr]
	return palaces
}

func ExtractFromTemplatesBasedOnTemplatePins(templates []*Template, max, min string) []*Template {
	sort.Slice(templates, func(i, j int) bool {
		pini := len(templates[i].TemplatePins)
		pinj := len(templates[j].TemplatePins)
		return pini < pinj
	})
	minptr := 0
	maxptr := len(templates)
	for i, v := range templates {
		pin := len(v.TemplatePins)
		if minpin, err := strconv.Atoi(min); err == nil && pin < minpin {
			minptr = i + 1
		}
		if maxpin, err := strconv.Atoi(max); err == nil && pin > maxpin {
			maxptr = i
			break
		}
	}
	templates = templates[minptr:maxptr]
	return templates
}

func GetEmbededPins(ctx context.Context, PalaceID uuid.UUID) ([]EmbededPin, error) {
	var embededPins []EmbededPin
	err := db.SelectContext(ctx, &embededPins, "SELECT number, x, y, word, place, do FROM embededpins WHERE palaceID=? ORDER BY number ASC ", PalaceID)
	if err != nil {
		return nil, err
	}

	return embededPins, nil
}

func CreateEmbededPin(ctx context.Context, number *int, palaceID uuid.UUID, x, y *float32, word, place, do string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO embededpins (number, x, y, word, place, do, palaceID) VALUES (?, ?, ?, ?, ?, ?, ?) ", number, x, y, word, place, do, palaceID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEmbededPins(ctx context.Context, palaceID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM embededpins WHERE palaceID=? ", palaceID)
	if err != nil {
		return err
	}
	return nil
}

func GetTemplatePins(ctx context.Context, TemplateID uuid.UUID) ([]TemplatePin, error) {
	var templatePins []TemplatePin
	err := db.SelectContext(ctx, &templatePins, "SELECT number, x, y FROM templatepins WHERE templateID=? ORDER BY number ASC ", TemplateID)
	if err != nil {
		return nil, err
	}

	return templatePins, nil
}

func CreateTemplatePin(ctx context.Context, number *int, templateID uuid.UUID, x, y *float32) error {
	_, err := db.ExecContext(ctx, "INSERT INTO templatepins (number, x, y, templateID) VALUES (?, ?, ?, ?) ", number, x, y, templateID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTemplatePins(ctx context.Context, templateID uuid.UUID) error {
	_, err := db.ExecContext(ctx, "DELETE FROM templatepins WHERE templateID=? ", templateID)
	if err != nil {
		return err
	}
	return nil
}
