# backend/Dockerfile

# マルチステージビルド: ビルドステージと実行ステージを分離してイメージサイズを削減

# ========== ビルドステージ ==========
FROM golang:1.25-alpine AS builder

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係のダウンロード（キャッシュ利用のため先にコピー）
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
# CGO_ENABLED=0: 静的リンクでバイナリを作成（依存ライブラリを含める）
# -o main: 出力ファイル名を "main" に指定
# GOOS=linux: Linux向けにビルド(一応指定)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# ========== 実行ステージ ==========
# 実行するときは軽量なAlpineを使用(Goはバイナリ化されているためランタイム不要)
FROM alpine:latest

# セキュリティアップデートと証明書のインストール
# ca-certificates: HTTPS通信に必要
# これはroot権限が必要なのでユーザー切り替え前に実行
RUN apk --no-cache add ca-certificates

# 作業ディレクトリの設定
WORKDIR /app

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# 非rootユーザーの作成と切り替え（rootユーザーでアプリを実行するとセキュリティリスクがあるため）
# -D: パスワードなしのユーザーを作成
# -u 1000: ユーザーIDを1000に設定(rootユーザーは0)
# chown: ディレクトリ全体の所有者を新しいユーザーに変更
RUN adduser -D -u 1000 appuser && \
    chown -R appuser:appuser /app

USER appuser

# ポートの公開
EXPOSE 8080

# 実行
CMD ["./main"]