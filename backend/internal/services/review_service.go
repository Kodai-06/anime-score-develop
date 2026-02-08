package services

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/repositories"
	"errors"
)

type ReviewService struct {
	reviewRepo   *repositories.ReviewRepository
	animeService *AnimeService
}

// NewReviewService はReviewServiceのインスタンスを生成
func NewReviewService(
	reviewRepo *repositories.ReviewRepository,
	animeService *AnimeService,
) *ReviewService {
	return &ReviewService{
		reviewRepo:   reviewRepo,
		animeService: animeService,
	}
}

// CreateReview はレビューを投稿する
// 1. アニメがDBに存在するか確認（詳細ページ表示時に保存済みのはず）
// 2. 既に同じユーザーが同じアニメにレビューしていないかチェック
// 3. レビューを保存
func (s *ReviewService) CreateReview(userID int64, input models.ReviewInput) (*models.Review, error) {
	// 1. スコアのバリデーション
	if input.Score < 0 || input.Score > 100 {
		return nil, errors.New("スコアは0〜100の範囲で入力してください")
	}

	// 2. アニメをDBから探す（詳細ページ表示時に保存済みのはず）
	// 見つからない場合はAnnict APIから取得して保存（フォールバック）
	anime, err := s.animeService.FindOrCreateAnime(input.AnnictID)
	if err != nil {
		return nil, err
	}

	// 3. 既にこのユーザーがこのアニメにレビューしていないかチェック
	existingReview, err := s.reviewRepo.FindByUserAndAnime(userID, anime.ID)
	if err != nil {
		return nil, err
	}
	if existingReview != nil {
		return nil, errors.New("既にこのアニメにはレビューを投稿済みです")
	}

	// 4. レビューを作成
	review := &models.Review{
		UserID:  userID,
		AnimeID: anime.ID,
		Score:   input.Score,
		Comment: input.Comment,
	}

	if err := s.reviewRepo.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}

// GetReviewsByAnimeID は特定アニメのレビュー一覧を取得
// ※すべての操作をServiceを通して行うことで、コードの一貫性が保たれる
func (s *ReviewService) GetReviewsByAnimeID(animeID int64) ([]models.Review, error) {
	return s.reviewRepo.FindByAnimeID(animeID)
}

// GetReviewsByUserID は特定ユーザーのレビュー一覧を取得
func (s *ReviewService) GetReviewsByUserID(userID int64) ([]models.Review, error) {
	return s.reviewRepo.FindByUserID(userID)
}

// GetReviewsByUserIDWithAnime は特定ユーザーのレビュー一覧をアニメ情報と共に取得
func (s *ReviewService) GetReviewsByUserIDWithAnime(userID int64) ([]models.ReviewWithAnime, error) {
	return s.reviewRepo.FindByUserIDWithAnime(userID)
}
