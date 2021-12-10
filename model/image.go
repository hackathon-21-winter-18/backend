package model

import (
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"

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

func EncodeToBase64(ctx context.Context, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	size := fi.Size()

	data := make([]byte, size)
	file.Read(data)

	return base64.StdEncoding.EncodeToString(data), nil
}

func DecodeToImageAndSave(ctx context.Context, encoded, path string) error {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("Failed to decode image")
	}

	file, err := os.Create(path) //TODO 既にあってもエラー返さない
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)
	return nil
}

func GetPalaceImagePath(ctx context.Context, palaceID uuid.UUID) (string, error) {
	var path string
	err := db.GetContext(ctx, &path, "SELECT image FROM palaces WHERE id=? ", palaceID)
	if err != nil {
		return "", err
	}
	return path, nil
}

func CreatePathName(ctx context.Context, base64 string) (string, error) { // go run main.goをやりなおしても値は変わらない
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, 25)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}

	var extension string
	head5 := string([]rune(base64)[:5])
	switch head5 {
	case "iVBOR":
		extension = ".png"
	case "R0lGO":
		extension = ".gif"
	case "/9j/4":
		extension = ".jpeg"
	default:
		err = fmt.Errorf("invalid image")
	}
	if err != nil {
		return "", err
	}

	return "./assets/" + result + extension, nil
}

func RemoveImage(ctx context.Context, path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}
