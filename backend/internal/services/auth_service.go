package services

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/repositories"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(repo *repositories.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Signup: ユーザー登録ロジック
func (s *AuthService) Signup(input models.SignUpInput) (*models.User, error) {
	// 1. ユーザー名のバリデーション（英数字と記号のみ）
	if !models.ValidateUsername(input.Username) {
		return nil, errors.New("ユーザー名は3〜50文字までで英数字とアンダースコア(_)、ハイフン(-)のみ使用できます")
	}

	// 2. ユーザー名の重複チェック
	exists, err := s.repo.ExistsByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("このユーザー名は既に使用されています")
	}

	// 3. パスワードをハッシュ化
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. ユーザーモデル作成
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPass),
	}

	// 5. DBに保存
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login: ログインロジック（JWTトークンを返す）
// JWTは
// Header（ヘッダー）: 暗号化方式などの情報
// Payload（ペイロード）: ユーザーIDや有効期限などのデータ(暗号化されていないので機密情報は入れないこと)
// Signature（署名）: シークレットキーを使って生成された暗号データ
// で構成される
func (s *AuthService) Login(input models.LoginInput) (string, error) {
	// 1. Emailでユーザー検索
	user, err := s.repo.GetByEmail(input.Email)
	if err != nil {
		return "", errors.New("ユーザーが見つかりません")
	}

	// 2. パスワード照合 (ハッシュ同士を比較)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		return "", errors.New("パスワードが間違っています")
	}

	// 3. JWTトークンの生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 72時間有効
	})

	// 秘密鍵で署名（環境変数から読み込む）
	secret_key := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
