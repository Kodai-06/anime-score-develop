--  Usersテーブル 
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--  Animesテーブル (Annict APIデータのキャッシュ)
CREATE TABLE animes (
    id SERIAL PRIMARY KEY,
    annictId INTEGER UNIQUE NOT NULL,   -- Annict API の作品ID
    title VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL,              -- 放送年 (例: 2024)
    image_url VARCHAR(500),             -- 作品画像URL (Annict APIから取得)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

--  Reviewsテーブル (0-100点, コメント任意)
CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    anime_id INTEGER NOT NULL REFERENCES animes(id) ON DELETE CASCADE,
    score INTEGER NOT NULL CHECK (score >= 0 AND score <= 100), -- 0~100点
    comment TEXT, -- NOT NULLを付けないので、NULL(未入力)が許可されます
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- 1ユーザー1アニメにつき1レビューのみの制約
    UNIQUE(user_id, anime_id)
);

--  インデックス (クエリパフォーマンス向上)
-- インデックスはinsertやupdateが遅くなる
CREATE INDEX idx_reviews_user_id ON reviews(user_id);
CREATE INDEX idx_reviews_anime_id ON reviews(anime_id);
CREATE INDEX idx_animes_title ON animes(title);
CREATE INDEX idx_animes_annictId ON animes(annictId);

--  アニメごとの統計情報を表示するビュー
-- ビューは簡単に言えばよく使う長いクエリをショートカット化するもの
-- ビューに含まれるORDER BY は必ずしも保証されないのでここで書かない
CREATE VIEW anime_stats AS
SELECT 
    anime_id,
    COUNT(id) AS review_count,          -- レビュー数
    ROUND(AVG(score), 1) AS avg_score   -- 平均点 (小数第1位まで)
FROM 
    reviews
GROUP BY 
    anime_id;