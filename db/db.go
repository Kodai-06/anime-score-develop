package db

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// DB はアプリケーション全体で使うデータベース接続
var DB *sqlx.DB

// Init はPostgreSQLへの接続を初期化する
func Init() error {
	// 環境変数から接続情報を取得（.envで設定）
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname := getEnv("DB_NAME", "anime_review_db")

	// PostgreSQL接続文字列 (DSN)
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// データベース接続
	var err error
	DB, err = sqlx.Connect("pgx", dsn)
	if err != nil {
		return fmt.Errorf("データベース接続エラー: %w", err)
	}

	// 接続テスト
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("データベースPingエラー: %w", err)
	}

	log.Println("✅ データベース接続成功")
	return nil
}

// Close はデータベース接続を閉じる
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// getEnv は環境変数を取得、なければデフォルト値を返す
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
