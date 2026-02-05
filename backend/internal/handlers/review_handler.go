package handlers

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/services"

	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service *services.ReviewService
}

// NewReviewHandler はハンドラのインスタンスを生成
func NewReviewHandler(service *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// Create は POST /api/reviews へのリクエストを処理する
// レビュー投稿（認証必須）
func (h *ReviewHandler) Create(c *gin.Context) {

	// 1. 認証ミドルウェアでセットされたユーザーIDを取得
	// (ミドルウェアが正常に動作していないときの場合のために一応)
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}
	userID := int64(userIDValue.(int))

	// 2. リクエストボディをパース
	var input models.ReviewInput
	// ShouldBindJSONはJSON形式のリクエストボディを構造体にバインドする
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です: " + err.Error()})
		return
	}

	// 3. サービス層でレビュー作成
	review, err := h.service.CreateReview(userID, input)
	if err != nil {
		// エラーメッセージに応じてステータスコードを変える
		// この書き方だとエラーメッセージの文言に依存してしまうので
		// 本当はカスタムエラー型を定義して判別したい
		if err.Error() == "既にこのアニメにはレビューを投稿済みです" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 成功レスポンス
	c.JSON(http.StatusCreated, gin.H{
		"message": "レビューを投稿しました",
		"review":  review,
	})
}

// 特定のアニメのレビュー一覧を取得するハンドラー
func (h *ReviewHandler) ListByAnime(c *gin.Context) {

	// 1. クエリパラメータからアニメIDを取得
	animeIDStr := c.Query("anime_id")

	// クエリパラメータは文字列なので
	// animeIDを数値に変換（失敗したらエラーを返すなどの安全策）
	animeID, err := strconv.ParseInt(animeIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime_id"})
		return
	}

	// 2. サービス層でレビュー一覧を取得
	reviews, err := h.service.GetReviewsByAnimeID(animeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reviews"})
		return
	}

	// 3. 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
	})

}

// 自分のレビュー一覧を取得するハンドラー
func (h *ReviewHandler) ListByMe(c *gin.Context) {

	// 1. 認証ミドルウェアでセットされたユーザーIDを取得
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証が必要です"})
		return
	}
	userID := int64(userIDValue.(int)) //interface{} を直接別の型に変換することはできないのでintであることを確認してからint64に変換

	// 2. サービス層でレビュー一覧をアニメ情報と共に取得
	reviews, err := h.service.GetReviewsByUserIDWithAnime(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reviews"})
		return
	}

	// 3. 成功レスポンス
	c.JSON(http.StatusOK, gin.H{
		"data": reviews,
	})
}
