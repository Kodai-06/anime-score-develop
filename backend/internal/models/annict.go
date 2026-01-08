package models

// AnnictGraphQLResponse はAnnict APIからのレスポンス全体を受け取る構造体
// アニメ情報は複数返ってくるからスライスにする
type AnnictGraphQLResponse struct {
	Data struct {
		SearchWorks struct {
			Nodes    []AnnictWork `json:"nodes"`
			PageInfo PageInfo     `json:"pageInfo"`
		} `json:"searchWorks"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// AnnictWork は単一のアニメ作品情報を表す構造体
// 要件定義書の「取得・利用する情報」に対応
type AnnictWork struct {
	AnnictID   int    `json:"annictId"`
	Title      string `json:"title"`
	SeasonYear *int   `json:"seasonYear"` // nullの場合があるためポインタ(int型にnullはない)
	Image      struct {
		RecommendedImageUrl string `json:"recommendedImageUrl"`
	} `json:"image"`
}

// PageInfo はページネーション情報を表す構造体
type PageInfo struct {
	HasNextPage bool   `json:"hasNextPage"` // 次のページがあるか
	EndCursor   string `json:"endCursor"`   // 次のページの開始位置
}
