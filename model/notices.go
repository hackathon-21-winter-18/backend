package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Notice struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Content   string    `json:"content" db:"content"`
	Read      bool      `json:"read" db:"read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func GetNotices(ctx context.Context, userID uuid.UUID) ([]Notice, error) {
	var notices []Notice
	err := db.SelectContext(ctx, &notices, "SELECT id, content, read, created_at FROM notices WhERE userID=? ORDER BY createdBy DESC ", userID)
	if err != nil {
		//TODO return userID not found
		return nil, err
	}

	for _, notice := range notices {
		if !notice.Read {
			_, err = db.ExecContext(ctx, "UPDATE notices SET read=? ", true)
			if err != nil {
				return nil, err
			}
		}
	}

	return notices, nil
}

func CreateNotice(ctx context.Context, createrID uuid.UUID, objectID uuid.UUID, palaceIssued bool) error {
	noticeID := uuid.New()
	content := "公開したものを元に他のユーザーが新たな"
	if palaceIssued {
		content += "宮殿"
	} else {
		content += "テンプレート"
	}
	content += "を公開しました。"
	_, err := db.ExecContext(ctx, "INSERT INTO notices (id, userID, content) VALUES (?, ?, ?) ", noticeID, createrID, content)
	if err != nil {
		return err
	}

	return nil
}
