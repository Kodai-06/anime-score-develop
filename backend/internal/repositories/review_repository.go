package repositories

import (
	"anime-score-backend/internal/models"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ReviewRepository struct {
	db *sqlx.DB
}

// NewReviewRepository はDB接続を受け取ってリポジトリを生成する
func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

// Create はレビューをDBに保存する
func (r *ReviewRepository) Create(review *models.Review) error {
	query := `
		INSERT INTO reviews (user_id, anime_id, score, comment)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		query,
		review.UserID,
		review.AnimeID,
		review.Score,
		review.Comment,
	).Scan(&review.ID, &review.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

// FindByUserAndAnime は特定のユーザーが特定のアニメに対して既にレビューしているか確認する
// 1ユーザー1作品1レビューの制約チェックに使用
func (r *ReviewRepository) FindByUserAndAnime(userID, animeID int64) (*models.Review, error) {
	query := `
		SELECT id, user_id, anime_id, score, comment, created_at
		FROM reviews
		WHERE user_id = $1 AND anime_id = $2
	`

	var review models.Review
	err := r.db.QueryRow(query, userID, animeID).Scan(
		&review.ID,
		&review.UserID,
		&review.AnimeID,
		&review.Score,
		&review.Comment,
		&review.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // レビューがない場合はnilを返す
		}
		return nil, fmt.Errorf("failed to find review: %w", err)
	}

	return &review, nil
}

// FindByAnimeID は特定のアニメのレビュー一覧を取得する（新着順）
func (r *ReviewRepository) FindByAnimeID(animeID int64) ([]models.Review, error) {
	query := `
		SELECT id, user_id, anime_id, score, comment, created_at
		FROM reviews
		WHERE anime_id = $1
		ORDER BY created_at DESC
	`

	var reviews []models.Review
	err := r.db.Select(&reviews, query, animeID)
	if err != nil {
		return nil, fmt.Errorf("failed to find reviews: %w", err)
	}

	return reviews, nil
}

// FindByUserID は特定のユーザーのレビュー一覧を取得する（新着順）
func (r *ReviewRepository) FindByUserID(userID int64) ([]models.Review, error) {
	query := `
		SELECT id, user_id, anime_id, score, comment, created_at
		FROM reviews
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var reviews []models.Review
	err := r.db.Select(&reviews, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find reviews: %w", err)
	}

	return reviews, nil
}
