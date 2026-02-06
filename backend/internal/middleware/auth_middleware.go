package middleware

import (
	"fmt"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 認証ミドルウェア
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. クッキーからトークンを取得
		tokenString, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is required"})
			c.Abort() // 処理をここで止める
			return
		}

		// 3. トークンの検証
		// ※ Login時と同じシークレットキーを使うこと！
		secret_key := os.Getenv("JWT_SECRET_KEY")
		secretKey := []byte(secret_key)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// アルゴリズムがHMACかどうか確認（セキュリティ対策）
			// 型アサーション
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		// 4. トークンが無効、または期限切れの場合
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 5. トークンからユーザーIDを取り出し、コンテキストにセットする
		// これにより、後のハンドラーで c.Get("userID") としてIDを使えるようになる
		// JWTはヘッダー、ペイロード、署名の3部分から構成される
		// JWTのクレームとは、ペイロード部分に含まれる情報(JSON形式)のこと
		// claims, ok := token.Claims.(jwt.MapClaims)は、トークンのクレームを
		// jwt.MapClaims型に変換し、okがtrueなら成功、falseなら失敗を示す
		// // token.Claims は interface{} 型であり、キーを指定できないからmap型に変換する
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// float64型にしないとint()を使えない
			// claims["user_id"]のuser_idはJWT生成時にペイロードに設定したキー
			if userID, ok := claims["user_id"].(float64); ok {
				// c.set(key string, value any) でコンテキストに値を保存する
				// 後のハンドラーで c.Get("userID") として取得可能
				c.Set("userID", int(userID))
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}
		}

		// 6. 次の処理へ進む
		c.Next()
	}
}
