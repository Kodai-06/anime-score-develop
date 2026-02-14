package handlers

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/services"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// クッキーセット用のヘルパー関数
func (h *AuthHandler) setAuthCookie(c *gin.Context, token string) {
	// 1. SameSite属性の設定
	// CSRF対策のため、SameSite属性を設定
	// c.SetCookieを呼ぶ「前」に設定する必要がある
	// laxモードは、外部サイトからのGETリクエスト以外にはクッキーを送信しない
	c.SetSameSite(http.SameSiteLaxMode) // または http.SameSiteStrictMode

	// 環境変数でSecureフラグを判定（本番ではHTTPS必須）
	isProduction := os.Getenv("ENV") == "production"

	// 2. クッキーにトークンをセット
	// c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	c.SetCookie(
		"auth_token", // cookie名
		token,        // 値（JWTなど）
		3600*24,      // 有効期限（秒）: ここでは24時間
		"/",          // パス: サイト全体で有効
		"",           // ドメイン: 空文字だと現在のドメインのみ
		isProduction, // Secure: ENV=productionのときtrue
		true,         // HttpOnly: JavaScriptからのアクセス禁止（必須）
	)
}

// Signup ハンドラー
// c.ShouldBindJSON(&input) はリクエストボディのJSONを構造体にマッピングしている
func (h *AuthHandler) Signup(c *gin.Context) {
	var input models.SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.service.Signup(input)
	if err != nil {
		// ユーザー名重複・バリデーションエラーはクライアントに内容を返す
		errMsg := err.Error()
		if errMsg == "このユーザー名は既に使用されています" || errMsg == "ユーザー名は1〜50文字で入力してください" {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// クッキーにトークンをセット
	h.setAuthCookie(c, token)

	// レスポンスボディにトークンを含める（BFF がトークンを受け取り Cookie に変換する）
	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user, "token": token})
}

// Login ハンドラー
func (h *AuthHandler) Login(c *gin.Context) {
	var input models.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := h.service.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 2. クッキーにトークンをセット
	h.setAuthCookie(c, token)

	// 3. レスポンスボディにトークンを含める（BFF がトークンを受け取り Cookie に変換する）
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "token": token})
}

// Logout ハンドラー
func (h *AuthHandler) Logout(c *gin.Context) {
	// 環境変数でSecureフラグを判定
	isProduction := os.Getenv("ENV") == "production"

	// クッキーの有効期限を過去に設定して削除
	c.SetCookie(
		"auth_token", // cookie名
		"",           // 値を空に
		-1,           // 有効期限を過去に設定
		"/",          // パス
		"",           // ドメイン
		isProduction, // Secure: ENV=productionのときtrue
		true,         // HttpOnly
	)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
