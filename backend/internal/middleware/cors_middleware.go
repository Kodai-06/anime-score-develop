package middleware

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORSミドルウェアの設定
func CORSMiddleware() gin.HandlerFunc {
	// 許可するオリジン（フロントエンドのURL）を定義
	// 環境変数 FRONTEND_URL があればそれを使い、なければデフォルトで localhost:3000 を許可
	allowOrigins := []string{"http://localhost:3000"}

	if url := os.Getenv("FRONTEND_URL"); url != "" {
		allowOrigins = append(allowOrigins, url)
	}

	config := cors.Config{
		// 許可するドメイン
		AllowOrigins: allowOrigins,

		// 許可するHTTPメソッド
		AllowMethods: []string{
			"POST", "GET", "OPTIONS", "PUT", "DELETE", "PATCH",
		},

		// 許可するヘッダー
		AllowHeaders: []string{
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
		},

		// クッキーの送受信を許可
		AllowCredentials: true,

		// プリフライトリクエストのキャッシュ時間
		MaxAge: 12 * time.Hour,
	}

	return cors.New(config)
}
