package models

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

// ReviewInput はレビュー投稿時の入力データ
type ReviewInput struct {
	AnnictID int     `json:"annictId" binding:"required"` // Annict APIのアニメID
	Score    int     `json:"score" binding:"required,min=0,max=100"`
	Comment  *string `json:"comment"`
}

// ReviewWithAnime はレビュー情報とアニメ情報を組み合わせた構造体
type ReviewWithAnime struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	AnimeID   int64     `db:"anime_id" json:"animeId"`
	Score     int       `db:"score" json:"score"`
	Comment   *string   `db:"comment" json:"comment"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	// アニメ情報
	Animetitle    string  `db:"anime_title" json:"animeTitle"`
	AnimeYear     int     `db:"anime_year" json:"animeYear"`
	AnimeImageURL *string `db:"anime_image_url" json:"animeImageUrl"`
}
