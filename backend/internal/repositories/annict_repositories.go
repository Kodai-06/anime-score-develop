package repositories

import (
	"anime-score-backend/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Annict GraphQL API のエンドポイント
const annictEndpoint = "https://api.annict.com/graphql"

// AnnictRepository は Annict API と通信するためのリポジトリ
type AnnictRepository struct {
	token  string       // Annict API の認証トークン
	client *http.Client // HTTP リクエスト送信用クライアント
}

// NewAnnictRepository はリポジトリのインスタンスを作成
func NewAnnictRepository(token string) *AnnictRepository {
	return &AnnictRepository{
		token: token,
		client: &http.Client{
			Timeout: 10 * time.Second, // タイムアウト設定
		},
	}
}

// SearchWorks はタイトルでアニメを検索し、結果のリストを返す
// keyword: 検索キーワード
// limit: 取得したい件数
// afterCursor: "ここから後ろを取得したい"という場所のID（初回は空文字 "" でOK）
func (r *AnnictRepository) SearchWorks(keyword string, limit int, afterCursor string) ([]models.AnnictWork, string, error) {
	// Annict API に送信する GraphQL クエリを定義
	// 検索条件: 指定されたタイトル、件数制限、ページネーション
	// 結果をシーズンの降順でソート
	query := `
		query SearchWorks($title: String!, $limit: Int!, $after: String) {
			searchWorks(
				titles: [$title], 
				first: $limit, 
				after: $after, 
				orderBy: { field: SEASON, direction: DESC }
			) {
				nodes {
					annictId
					title
					seasonYear
					image {
						recommendedImageUrl
					}
				}
				pageInfo {
					hasNextPage
					endCursor
				}
			}
		}
	`

	// GraphQL クエリに渡す変数をセット
	variables := map[string]interface{}{
		"title": keyword, // 検索キーワード
		"limit": limit,   // 取得件数
	}
	// ページネーションが必要な場合、カーソルを追加
	if afterCursor != "" {
		variables["after"] = afterCursor
	}

	// GraphQLのリクエストはクエリと変数をJSON形式で送る必要がある
	// GraphQL リクエストを JSON 形式(バイト列)にマーシャル
	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// HTTP POST リクエストを作成
	// bytes.NewBuffer(requestBody)はただのバイト列であるreauestBodyをio.Readerインターフェースに変換する
	// io.ReadrerインターフェースはReadメソッドを持つインターフェースのこと
	// http.NewRequestの第3引数はio.Readerインターフェースを受け取るので、bytes.NewBufferで変換する必要がある
	req, err := http.NewRequest("POST", annictEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	// リクエストヘッダーを設定（認証とコンテンツタイプ）
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	// Annict API にリクエストを送信
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// ステータスコードをチェック
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("annict api returned non-200 status: %d", resp.StatusCode)
	}

	// レスポンスボディをパース
	// graphQLRespにレスポンスをマッピングしている
	var graphQLResp models.AnnictGraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphQLResp); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	// GraphQL エラーをチェック
	if len(graphQLResp.Errors) > 0 {
		return nil, "", fmt.Errorf("graphql error: %s", graphQLResp.Errors[0].Message)
	}

	// 次ページ用のカーソルを取得
	// ページネーションが存在すれば EndCursor を、ないなら空文字を返す
	nextCursor := ""
	if graphQLResp.Data.SearchWorks.PageInfo.HasNextPage {
		nextCursor = graphQLResp.Data.SearchWorks.PageInfo.EndCursor
	}

	// 検索結果とカーソルを返す
	return graphQLResp.Data.SearchWorks.Nodes, nextCursor, nil
}
