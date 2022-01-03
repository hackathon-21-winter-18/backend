package model

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type notice struct {
	ID uuid.UUID
	User uuid.UUID 
	Content string
	Read bool
}

func CreateNotice(ctx context.Context, editer, creater uuid.UUID, objectID uuid.UUID, palaceIssued bool) error {
	noticeID := uuid.New()
	content := "公開したものを元に他のユーザーが新たな"
	if palaceIssued {
		content += fmt.Sprintf("[宮殿]`(https://frontend-opal-delta-19.vercel.app/~/%s)", objectID)
	} else {
		content += fmt.Sprintf("[テンプレート]`(https://frontend-opal-delta-19.vercel.app/~/%s)", objectID)
	}
	content += "を作成しました。" + "\n" + "\n"
	_, err := db.ExecContext(ctx, "INSERT INTO notices (id, userID, content) VALUES (?, ?, ?) ", noticeID, creater, content)
	if err != nil {
		return err
	}

	return nil
}