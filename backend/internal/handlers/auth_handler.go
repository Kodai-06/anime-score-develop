package handlers

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/services"
	"net/http"

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
	// 2. クッキーにトークンをセット
	// c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	c.SetCookie(
		"auth_token", // cookie名
		token,        // 値（JWTなど）
		3600*24,      // 有効期限（秒）: ここでは24時間
		"/",          // パス: サイト全体で有効
		"",           // ドメイン: 空文字だと現在のドメインのみ
		false,        // Secure: localhost開発時はfalse, 本番(HTTPS)はtrueにする
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// クッキーにトークンをセット
	h.setAuthCookie(c, token)

	// レスポンスボディにはトークンを含めず、成功メッセージとユーザー情報のみ返す
	c.JSON(http.StatusCreated, gin.H{"message": "User created", "user": user})
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

	// 3. レスポンスボディにはトークンを含めず、成功メッセージのみ返す
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}

// Logout ハンドラー
func (h *AuthHandler) Logout(c *gin.Context) {
	// クッキーの有効期限を過去に設定して削除
	c.SetCookie(
		"auth_token", // cookie名
		"",           // 値を空に
		-1,           // 有効期限を過去に設定
		"/",          // パス
		"",           // ドメイン
		false,        // Secure
		true,         // HttpOnly
	)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
