package services

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/repositories"
)

type AnimeService struct {
	annictRepo *repositories.AnnictRepository
	animeRepo  *repositories.AnimeRepository
}

// NewAnimeService はAnimeServiceのインスタンスを生成
// 依存関係（ここではAnnictRepository）を注入
func NewAnimeService(annictRepo *repositories.AnnictRepository, animeRepo *repositories.AnimeRepository) *AnimeService {
	return &AnimeService{
		annictRepo: annictRepo,
		animeRepo:  animeRepo,
	}
}

// SearchAnimes はキーワードに基づいてアニメを検索
// 将来的にはここで「DB検索」と「外部API検索」を併用するロジックに拡張予定
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

// FindOrCreateAnime は指定されたAnnict IDのアニメを取得します。
// DBに存在しない場合はAnnict APIから取得してDBに保存します。
func (s *AnimeService) FindOrCreateAnime(annictID int) (*models.Anime, error) {
	// 1. まず自分のDBを探す (キャッシュチェック)
	localAnime, err := s.animeRepo.FindByAnnictID(annictID)
	if err != nil {
		return nil, err
	}
	// もし見つかったらそれを返す
	if localAnime != nil {
		return localAnime, nil
	}

	// 2. DBになければ、Annict APIから情報を取得
	annictWork, err := s.annictRepo.GetWorkByID(annictID)
	if err != nil {
		return nil, err
	}

	// 3. 取得した情報をDB保存用のモデルに変換
	// 画像URLが空なら nil を入れる
	var imageURL *string
	if annictWork.Image.RecommendedImageUrl != "" {
		url := annictWork.Image.RecommendedImageUrl
		imageURL = &url
	}

	// SeasonYearはポインタなのでnilチェックを行う（nilなら0を入れる）
	year := 0
	if annictWork.SeasonYear != nil {
		year = *annictWork.SeasonYear
	}

	newAnime := &models.Anime{
		AnnictID: int64(annictWork.AnnictID),
		Title:    annictWork.Title,
		Year:     year,
		ImageURL: imageURL,
	}

	// 4. DBに保存
	// Createメソッド内でIDが採番され、newAnime.IDにセットされます
	if err := s.animeRepo.Create(newAnime); err != nil {
		return nil, err
	}

	return newAnime, nil
}
