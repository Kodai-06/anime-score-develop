# frontend/Dockerfile

# Next.jsは依存パッケージのインストールが必要

# ベースイメージ
FROM node:25-alpine

WORKDIR /app

# パッケージ情報のコピーとインストール
COPY package*.json ./
RUN npm install

# ソースコードのコピー
COPY . .

# Next.jsのデフォルトポート
EXPOSE 3000

# 開発モードで起動
CMD ["npm", "run", "dev"]