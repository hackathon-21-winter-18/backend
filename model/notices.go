package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Notice struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Checked   bool      `json:"checked" db:"checked"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func GetNotices(ctx context.Context, userID uuid.UUID) ([]Notice, error) {
	var notices []Notice
	err := db.SelectContext(ctx, &notices, "SELECT id, content, checked, created_at FROM notices WHERE userID=? ORDER BY created_at DESC ", userID)
	if err != nil {
		//TODO return userID not found
		return nil, err
	}

	for _, notice := range notices {
		if !notice.Checked {
			_, err = db.ExecContext(ctx, "UPDATE notices SET checked=?, updated_at=? WHERE id=? ", true, time.Now(), notice.ID)
			if err != nil {
				return nil, err
			}
		}
	}

	return notices, nil
}

func CreateNotice(ctx context.Context, createrID uuid.UUID, objectID uuid.UUID, palaceIssued bool) error {
	noticeID := uuid.New()
	var firstshared bool
	var err error
	content := "公開したものを元に他のユーザーが新たな"

	if palaceIssued {
		err = db.GetContext(ctx, &firstshared, "SELECT firstshared FROM palaces WHERE id=? ", objectID)
		if err != nil {
			//TODO return ID not found
			return err
		}
		if firstshared {
			return nil
		}
		content += "宮殿"
	} else {
		err = db.GetContext(ctx, &firstshared, "SELECT firstshared FROM templates WHERE id=? ", objectID)
		if err != nil {
			//TODO return ID not found
			return err
		}
		if firstshared {
			return nil
		}
		content += "テンプレート"
	}
	content += "を公開しました。"

	date := time.Now()
	_, err = db.ExecContext(ctx, "INSERT INTO notices (id, userID, content, created_at, updated_at) VALUES (?, ?, ?, ?, ?) ", noticeID, createrID, content, date, date)
	if err != nil {
		return err
	}

	return nil
}

func GetCountOfUnreadNotices(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(*) FROM notices WHERE userID=? AND checked=False ", userID)
	if err != nil {
		//TODO return userID not found
		return 0, err
	}

	return count, nil
}
