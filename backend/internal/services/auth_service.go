package services

import (
	"anime-score-backend/internal/models"
	"anime-score-backend/internal/repositories"
	"errors"
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
	// 1. パスワードをハッシュ化
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 2. ユーザーモデル作成
	user := &models.User{
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(hashedPass),
	}

	// 3. DBに保存
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login: ログインロジック（JWTトークンを返す）
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

	// 秘密鍵で署名（本来は環境変数から読み込むべきですが、今は仮の文字列で）
	tokenString, err := token.SignedString([]byte("YOUR_SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
