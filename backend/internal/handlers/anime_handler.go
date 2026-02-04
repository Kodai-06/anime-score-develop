package handlers

import (
	"anime-score-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnimeHandler struct {
	service *services.AnimeService
}

// NewAnimeHandler はハンドラのインスタンスを生成
func NewAnimeHandler(service *services.AnimeService) *AnimeHandler {
	return &AnimeHandler{
		service: service,
	}
}

// Search は /api/animes/search へのリクエストを処理
func (h *AnimeHandler) Search(c *gin.Context) {
	// 1. クエリパラメータの取得
	// URL: /api/animes/search?q=ガンダム&limit=20&cursor=xxx
	keyword := c.Query("q")
	limitStr := c.DefaultQuery("limit", "15") // 指定がなければ "15"
	cursor := c.Query("cursor")

	// 2. バリデーション（簡易）
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search keyword 'q' is required"})
		return
	}

	// クエリパラメータは文字列なので
	// limitを数値に変換（失敗したらデフォルト値を使うなどの安全策）
	// AtoiはASCII to Integer
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 15
	}

	// 3. Service層の呼び出し
	works, nextCursor, err := h.service.SearchAnimes(keyword, limit, cursor)
	if err != nil {
		// 外部APIのエラーなどはここでログに出し、ユーザーには500を返す
		// 本番では詳細なエラーメッセージを隠蔽するのが一般的
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search animes"})
		return
	}

	// 4. レスポンスの返却
	c.JSON(http.StatusOK, gin.H{
		"data":       works,
		"nextCursor": nextCursor,
	})
}

// GetDetail は /api/animes/:id へのリクエストを処理
func (h *AnimeHandler) GetDetail(c *gin.Context) {
	// ID をパスパラメータから取得
	idStr := c.Param("id")
	annictID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid anime ID"})
		return
	}

	// Service を呼び出してレスポンスを返す
	anime, stats, err := h.service.GetAnimeDetail(int64(annictID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get anime details"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"anime": anime,
		"stats": stats,
	})
}

// GetList は /api/animes へのリクエストを処理（アニメ一覧取得）
func (h *AnimeHandler) GetList(c *gin.Context) {
	// クエリパラメータの取得
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	// 数値に変換
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	// Service呼び出し
	result, err := h.service.GetAnimeList(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get anime list"})
		return
	}

	c.JSON(http.StatusOK, result)
}
