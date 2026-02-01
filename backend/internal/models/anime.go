package models

import "time"

// Anime は Annict から取得した作品データをローカルにキャッシュするモデル。
type Anime struct {
	ID        int64     `db:"id" json:"id"`
	AnnictID  int64     `db:"annictid" json:"annictId"`
	Title     string    `db:"title" json:"title"`
	Year      int       `db:"year" json:"year"`
	ImageURL  *string   `db:"image_url" json:"imageUrl"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// AnimeStats はビュー anime_stats の集計結果を表すモデル。
type AnimeStats struct {
	AnimeID     int64   `db:"anime_id" json:"animeId"`
	ReviewCount int     `db:"review_count" json:"reviewCount"`
	AvgScore    float64 `db:"avg_score" json:"avgScore"`
}

// intは環境依存で最大値が異なるので数が大きくなる可能性のあるIDはint64を使う
