# AnimeScore

アニメ作品にスコアとコメントを投稿できるレビューサイト

URL: https://anime-score-develop-141202940381.asia-northeast1.run.app

## 使用技術
- **フロントエンド**: Next.js(TypeScript)
- **UI**: Tailwind CSS, shadcn/ui
- **バックエンド**: Gin(Go)
- **データベース**: supabase(PostgreSQL)
- **インフラ**: Docker,Vercel,GCP(Cloud Run)

## ディレクトリ構成

```
.
├── backend/           # Go API サーバー
│   ├── Dockerfile     # バックエンド用Dockerfile
│   ├── cmd/
│   │   └── api/       # エントリーポイント
│   └── internal/
│       ├── handlers/  # HTTPハンドラー
│       ├── models/    # 構造体定義
│       ├── services/  # ロジック
│       ├── repositories/ # データアクセス層
│       └── middleware/ # ミドルウェア(認証,CORS設定)
├── frontend/          # Next.js アプリケーション
│   ├── Dockerfile     # フロントエンド用Dockerfile
│   └── src/
│       ├── app/       # ページとレイアウト
│       ├── components/ # Reactコンポーネント
│       ├── contexts/  # React Context(認証)
│       ├── lib/       # ユーティリティ
│       └── types/     # TypeScript型定義
├── initdb/           # データベース初期化スクリプト
└── docker-compose.yml # Docker Compose設定
```

## 主な機能

- **認証**: JWT認証(HttpOnly属性のCookieに保存)
- **アニメ検索**: [Annict](https://annict.com/) のAPIを利用したアニメタイトル検索
- **レビュー**: 0〜100点のスコア＋任意コメントでレビューを投稿
- **アニメ詳細**: 平均スコア・レビュー数・レビュー一覧を確認
- **マイページ**: マイページで自分のレビュー履歴を確認

### Annict GraphQL API
アニメ情報の取得に [Annict](https://annict.com/) の GraphQL API を使用しています。

- **エンドポイント**: `https://api.annict.com/graphql`
- **用途**:
  - タイトルによるアニメ検索
  - Annict IDによるアニメ情報の取得
- **取得データ**: 作品ID、タイトル、放送年、画像URL
- **キャッシュ**: 取得したアニメ情報はDBにキャッシュし、2回目以降はDBから取得

