package repositories

import (
	"anime-score-backend/internal/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AnimeRepository struct {
	db *sqlx.DB
}

// NewAnimeRepository はDB接続を受け取ってリポジトリを生成する
func NewAnimeRepository(db *sqlx.DB) *AnimeRepository {
	return &AnimeRepository{db: db}
}

// Create はアニメ情報をDBに保存し、生成されたIDを返する
// ON CONFLICT (annict_id) DO NOTHING を使うことで、
// 万が一同時に同じアニメが保存されようとしてもエラーにならず、既存のIDを返すようにすると堅牢
func (r *AnimeRepository) Create(anime *models.Anime) error {
	// 重複チェックも兼ねたINSERT（Upsert的な挙動）
	// annict_idが既に存在する場合はtitleを更新するだけにする
	// IDを返すために RETURNING id を使用
	query := `
		INSERT INTO animes (annict_id, title, year, image_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (annict_id) DO UPDATE 
		SET title = EXCLUDED.title -- 既に存在する場合は情報を更新（念のため）
		RETURNING id
	`

	// QueryRowを使って、INSERTされた（または更新された）行のIDを取得
	err := r.db.QueryRow(
		query,
		anime.AnnictID,
		anime.Title,
		anime.Year,
		anime.ImageURL,
	).Scan(&anime.ID)

	if err != nil {
		return fmt.Errorf("failed to insert anime: %w", err)
	}

	return nil
}

// FindByAnnictID はAnnictID（外部ID）を使ってDBからアニメを探す
// レビュー投稿時に「このアニメは既にDBにあるか？」を調べるのに使う
func (r *AnimeRepository) FindByAnnictID(annictID int) (*models.Anime, error) {
	query := `SELECT id, annict_id, title, year, image_url, created_at 
	          FROM animes WHERE annict_id = $1`

	var anime models.Anime
	err := r.db.QueryRow(query, annictID).Scan(
		&anime.ID,
		&anime.AnnictID,
		&anime.Title,
		&anime.Year,
		&anime.ImageURL,
		&anime.CreatedAt,
	)

	// sql.ErrNoRowsは検索結果が0件の場合の特別なエラー変数
	// errors.Isでエラーの種類を判定している
	// データがないことをエラーとせず、nilを返すようにする
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // 見つからない場合はnilを返す
		}
		return nil, fmt.Errorf("failed to find anime: %w", err)
	}

	return &anime, nil
}

// FindByIDWithStats はIDを使ってアニメ情報とその統計情報を取得する
// FindByIDWithStats はアニメID(内部ID)を使ってDBからアニメと統計情報を一度に取得
func (r *AnimeRepository) FindByIDWithStats(id int64) (*models.Anime, *models.AnimeStats, error) {
	// LEFT JOINを使うことで、レビューがないアニメでも取得可能
	// COALESCEはリストの中から、最初に『NULLではない』値を返す関数
	// COALESCEを使って、統計情報がNULLの場合は0を返すようにしている
	query := `
        SELECT 
            a.id, a.annict_id, a.title, a.year, a.image_url, a.created_at,
            COALESCE(s.review_count, 0) as review_count,
            COALESCE(s.avg_score, 0) as avg_score
        FROM animes a
        LEFT JOIN anime_stats s ON a.id = s.anime_id
        WHERE a.id = $1
    `

	var a models.AnimeWithStats

	err := r.db.QueryRow(query, id).Scan(
		&a.ID,
		&a.AnnictID,
		&a.Title,
		&a.Year,
		&a.ImageURL,
		&a.CreatedAt,
		&a.ReviewCount,
		&a.AvgScore,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil // アニメが見つからない場合
		}
		return nil, nil, fmt.Errorf("failed to find anime with stats: %w", err)
	}

	return &a.Anime, &models.AnimeStats{
		AnimeID:     a.ID,
		ReviewCount: a.ReviewCount,
		AvgScore:    a.AvgScore,
	}, nil
}

// FindAllWithStats はアニメ一覧を統計情報付きで取得する
// 平均点順（降順）でソートし、ページネーションに対応
func (r *AnimeRepository) FindAllWithStats(limit, offset int) ([]models.AnimeWithStats, int, error) {
	// 総件数を取得
	var total int
	countQuery := `SELECT COUNT(*) FROM animes`
	if err := r.db.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count animes: %w", err)
	}

	// アニメ一覧を平均点順（降順）で取得
	// レビューがないアニメは avg_score = 0 として扱う
	// limitは何件取得するか、offsetは何件飛ばすか
	query := `
		SELECT 
			a.id, a.annict_id, a.title, a.year, a.image_url, a.created_at,
			COALESCE(s.review_count, 0) as review_count,
			COALESCE(s.avg_score, 0) as avg_score
		FROM animes a
		LEFT JOIN anime_stats s ON a.id = s.anime_id
		ORDER BY avg_score DESC, review_count DESC, a.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch animes: %w", err)
	}
	defer rows.Close()

	var animes []models.AnimeWithStats
	// 次の行がある限りループ
	for rows.Next() {
		var a models.AnimeWithStats
		err := rows.Scan(
			&a.ID,
			&a.AnnictID,
			&a.Title,
			&a.Year,
			&a.ImageURL,
			&a.CreatedAt,
			&a.ReviewCount,
			&a.AvgScore,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan anime: %w", err)
		}
		animes = append(animes, a)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows iteration error: %w", err)
	}

	return animes, total, nil
}
