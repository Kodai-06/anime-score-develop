package models

import (
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

// ValidateUsername: ユーザー名が有効かチェック（文字数のみ）
func ValidateUsername(username string) bool {
	if len(username) < 1 || len(username) > 50 {
		return false
	}
	return true
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
