package models

import (
	"regexp"
	"time"
)

// User 構造体: DBのusersテーブルに対応
type User struct {
	ID           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"` // JSONには出力しない設定
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// ユーザー名のバリデーション用正規表現（英数字と記号のみ許可）
// 許可する記号: _-
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// ValidateUsername: ユーザー名が有効かチェック
func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	return usernameRegex.MatchString(username)
}

// SignUpInput: フロントエンドから送られてくる登録用データ
type SignUpInput struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginInput: フロントエンドから送られてくるログイン用データ
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
