package services

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/repositories"
)

type AnimeService struct {
	annictRepo *repositories.AnnictRepository
}

// NewAnimeService はAnimeServiceのインスタンスを生成します
// 依存関係（ここではAnnictRepository）を注入します
func NewAnimeService(annictRepo *repositories.AnnictRepository) *AnimeService {
	return &AnimeService{
		annictRepo: annictRepo,
	}
}

// SearchAnimes はキーワードに基づいてアニメを検索します
// 将来的にはここで「DB検索」と「外部API検索」を併用するロジックに拡張します
func (s *AnimeService) SearchAnimes(keyword string, limit int, cursor string) ([]models.AnnictWork, string, error) {
	// 1. バリデーション（安全対策）
	// 極端に大きなリクエストが来ないように制限をかける
	if limit <= 0 {
		limit = 10 // デフォルト値
	}
	if limit > 50 {
		limit = 50 // 上限値
	}

	// 2. Annict APIへの問い合わせ
	// Repository層を呼び出してデータを取得
	works, nextCursor, err := s.annictRepo.SearchWorks(keyword, limit, cursor)
	if err != nil {
		return nil, "", err
	}

	// 3. 必要に応じたデータ加工（あればここに記述）
	// 現状は取得したデータをそのまま返しているが、
	// 将来的には「画像がない場合はデフォルト画像を入れる」などの処理をここに書ける

	return works, nextCursor, nil
}
