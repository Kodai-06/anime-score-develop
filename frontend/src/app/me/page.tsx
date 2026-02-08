"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Header } from "@/components/Header";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/contexts/AuthContext";
import { getMyReviews } from "@/lib/api";
import type { ReviewWithAnime } from "@/types";

export default function MyPage() {
  const router = useRouter();
  const { user, isLoading: authLoading } = useAuth(); // isLoadingをauthLoadingに名前変更

  const [reviews, setReviews] = useState<ReviewWithAnime[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // 未ログインならログインページへリダイレクト
  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [user, authLoading, router]); // ESLintのルールでuseEffect ではeffect 内で使用される外部の値をすべて依存配列に含めることが推奨のためrouterも依存配列に追加

  // レビュー履歴を取得
  useEffect(() => {
    const fetchReviews = async () => {
      if (!user) return;

      setIsLoading(true);
      setError(null);
      try {
        const response = await getMyReviews();
        setReviews(response.data || []);
      } catch (err) {
        setError("レビューの取得に失敗しました");
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchReviews();
  }, [user]);

  // 認証確認中
  if (authLoading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Header />
        <main className="container mx-auto px-4 py-8">
          <div className="text-center py-12">
            <p className="text-gray-500">読み込み中...</p>
          </div>
        </main>
      </div>
    );
  }

  // 未ログイン（リダイレクト前の表示）
  if (!user) {
    return null;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-2xl font-bold">マイページ</h1>

        {/* ユーザー情報 */}
        <Card className="mb-8">
          <CardHeader>
            <CardTitle className="text-lg">プロフィール</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p>
                <span className="text-gray-500">ユーザー名：</span>
                {user.username}
              </p>
              {/* <p>
                <span className="text-gray-500">メール：</span>
                {user.email}
              </p> */}
              <p>
                <span className="text-gray-500">投稿数：</span>
                {reviews.length}件
              </p>
            </div>
          </CardContent>
        </Card>

        {/* レビュー履歴 */}
        <h2 className="text-lg font-bold mb-4">レビュー履歴</h2>

        {/* ローディング */}
        {isLoading && (
          <div className="text-center py-8">
            <p className="text-gray-500">読み込み中...</p>
          </div>
        )}

        {/* エラー */}
        {error && (
          <div className="rounded bg-red-50 p-4 text-red-600">{error}</div>
        )}

        {/* レビュー一覧 */}
        {!isLoading && !error && (
          <>
            {reviews.length === 0 ? (
              <div className="text-center py-8">
                <p className="text-gray-500 mb-4">まだレビューがありません</p>
                <Link
                  href="/search"
                  className="text-primary hover:underline"
                >
                  アニメを検索してレビューを投稿しましょう
                </Link>
              </div>
            ) : (
              <div className="space-y-4">
                {reviews.map((review) => (
                  <Card key={review.id}>
                    <CardContent className="py-4">
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <Link
                            href={`/animes/${review.animeId}`}
                            className="font-medium hover:text-primary"
                          >
                            {review.animeTitle}
                          </Link>
                          <p className="text-sm text-gray-500">
                            {review.animeYear}年
                          </p>
                          <div className="mt-2">
                            <span className="text-lg font-bold text-primary">
                              {review.score}点
                            </span>
                          </div>
                          {review.comment && (
                            <p className="mt-2 text-gray-700">
                              {review.comment}
                            </p>
                          )}
                        </div>
                        <span className="text-sm text-gray-500">
                          {new Date(review.createdAt).toLocaleDateString("ja-JP")}
                        </span>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
