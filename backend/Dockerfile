# backend/Dockerfile

# Goはコンパイル言語なので、ビルドと実行の2段階が必要

# ベースイメージ（Goのバージョンはgo.modに合わせる）
FROM golang:1.25-alpine

# 作業ディレクトリの設定
WORKDIR /app

# 依存関係のダウンロード（キャッシュ利用のため先にコピー）
COPY go.mod go.sum ./
RUN go mod download

# ソースコードのコピー
COPY . .

# アプリケーションのビルド
# "go build" コマンドで実行ファイルを作る
# -o main : 出力ファイル名を "main" に指定
# ./cmd/api/main.go : エントリーポイント（プログラムの開始地点）のパス
RUN go build -o main ./cmd/api/main.go

# ポートの公開
# コンテナの外（ホストPCや他のコンテナ）からアクセスできるようにポートを開ける
# Ginフレームワークのデフォルトは8080
EXPOSE 8080

# 実行
CMD ["./main"]