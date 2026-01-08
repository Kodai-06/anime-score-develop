package repositories

import (
	"anime-score-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

// コンストラクタ関数
// UserRepository型の物体を作っている(構造体の初期化)
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create: ユーザーをDBに保存
// *UserRepository型にメソッド(Create())を追加している
func (r *UserRepository) Create(user *models.User) error {
	query := `
        INSERT INTO users (username, email, password_hash) 
        VALUES ($1, $2, $3) 
        RETURNING id, created_at`

	// Scanを使って、生成されたIDと作成日時をuser構造体に書き戻す
	// SQLの実行結果はScan()を実行するときまで持ち越される
	// QueryRow()にはクエリとプレースホルダに入る変数を渡している
	err := r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt)
	return err
}

// GetByEmail: メールアドレスからユーザーを取得
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`

	// sqlxのGetを使うと構造体にマッピングしてくれる
	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
