package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib" // pgxドライバー
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

	"anime-score-backend/internal/handlers"
	"anime-score-backend/internal/repositories"
	"anime-score-backend/internal/services"
)

func main() {
	// .envファイルの読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// データベース接続
	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN is not set in .env")
	}

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	defer db.Close()

	// 接続確認
	if err := db.Ping(); err != nil {
		log.Fatalln("Database ping failed:", err)
	}
	log.Println("Successfully connected to database!")

	// Ginルーターのセットアップ
	r := gin.Default()

	// 依存関係の注入 (DI)

	// 認証関連
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	// アニメ検索関連
	annictRepo := repositories.NewAnnictRepository(os.Getenv("ANNICT_ACCESS_TOKEN"))
	animeService := services.NewAnimeService(annictRepo)
	animeHandler := handlers.NewAnimeHandler(animeService)

	// ルーティング
	api := r.Group("/api")
	{
		api.POST("/signup", authHandler.Signup)
		api.POST("/login", authHandler.Login)

		// アニメ検索エンドポイント (GET /api/animes/search)
		api.GET("/animes/search", animeHandler.Search)
	}

	// ヘルスチェック用エンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     "connected",
		})
	})

	// サーバー起動
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalln("Failed to start server:", err)
	}
}
