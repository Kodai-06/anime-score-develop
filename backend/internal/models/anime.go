package models

import "time"

// Anime は Annict から取得した作品データをローカルにキャッシュするモデル。
type Anime struct {
	ID        int64     `db:"id" json:"id"`
	AnnictID  int64     `db:"annictid" json:"annictId"`
	Title     string    `db:"title" json:"title"`
	Year      int       `db:"year" json:"year"` // intは環境依存で最大値が異なるので数が大きくなる可能性のあるIDはint64を使う
	ImageURL  *string   `db:"image_url" json:"imageUrl"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// AnimeStats はビュー anime_stats の集計結果を表すモデル。
type AnimeStats struct {
	AnimeID     int64   `db:"anime_id" json:"animeId"`
	ReviewCount int     `db:"review_count" json:"reviewCount"`
	AvgScore    float64 `db:"avg_score" json:"avgScore"`
}

// AnimeWithStats はアニメ情報と統計情報を一緒に持つ構造体
type AnimeWithStats struct {
	Anime               // フィールド名を書かずに型名だけを書くと埋め込みとなり、子のフィールドにあたかも親のフィールドのようにアクセスできる
	ReviewCount int     `db:"review_count" json:"reviewCount"`
	AvgScore    float64 `db:"avg_score" json:"avgScore"`
}

// AnimeListResponse はアニメ一覧のレスポンス形式
type AnimeListResponse struct {
	Data       []AnimeWithStats `json:"data"`
	Pagination Pagination       `json:"pagination"`
}

// Pagination はページネーション情報
type Pagination struct {
	Page      int `json:"page"`
	PageSize  int `json:"pageSize"`
	Total     int `json:"total"`
	TotalPage int `json:"totalPage"`
}
