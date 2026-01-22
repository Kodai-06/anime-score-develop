package repositories

import (
	"anime-score-backend/internal/models"
	"database/sql"
	"errors"
	"fmt"
)

type AnimeRepository struct {
	db *sql.DB
}

// NewAnimeRepository はDB接続を受け取ってリポジトリを生成する
func NewAnimeRepository(db *sql.DB) *AnimeRepository {
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
