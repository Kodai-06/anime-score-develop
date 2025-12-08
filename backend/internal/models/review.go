package model

import "time"

// Review はユーザーがアニメに付けたスコアと任意コメントを保持するモデル。
type Review struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	AnimeID   int64     `db:"anime_id" json:"animeId"`
	Score     int       `db:"score" json:"score"`
	Comment   *string   `db:"comment" json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
