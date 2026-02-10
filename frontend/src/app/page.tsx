"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Header } from "@/components/Header";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { getRecentReviews } from "@/lib/api";
import type { ReviewWithAnime } from "@/types";

export default function Home() {
  const [reviews, setReviews] = useState<ReviewWithAnime[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // レビュー新着順を取得
  useEffect(() => {
    const fetchReviews = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await getRecentReviews();
        setReviews(response.data || []);
      } catch (err) {
        setError("レビューの取得に失敗しました");
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchReviews();
  }, []);

  // 日付フォーマット
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString("ja-JP", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-2xl font-bold">新着レビュー</h1>

        {/* ローディング */}
        {isLoading && (
          <div className="text-center py-12">
            <p className="text-gray-500">読み込み中...</p>
          </div>
        )}

        {/* エラー */}
        {error && (
          <div className="text-center py-12">
            <p className="text-red-500">{error}</p>
          </div>
        )}

        {/* レビュー一覧 */}
        {!isLoading && !error && (
          <>
            {reviews.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500">まだレビューがありません</p>
                <Link href="/search" className="text-primary hover:underline mt-2 inline-block">
                  アニメを検索してレビューをしよう
                </Link>
              </div>
            ) : (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                {reviews.map((review) => (
                  <Link key={review.id} href={`/animes/${review.animeAnnictId}`}>
                    <Card className="h-full hover:shadow-md transition-shadow cursor-pointer">
                      <CardHeader className="pb-2">
                        <CardTitle className="text-base line-clamp-2">
                          {review.animeTitle}
                        </CardTitle>
                        <p className="text-sm text-gray-500">{review.animeYear}年</p>
                      </CardHeader>
                      <CardContent>
                        <div className="flex items-center gap-2 mb-2">
                          <span className="text-lg font-bold text-primary">
                            {review.score}点
                          </span>
                        </div>
                        {review.comment && (
                          <p className="text-sm text-gray-600 line-clamp-3">
                            {review.comment}
                          </p>
                        )}
                        <p className="text-xs text-gray-400 mt-2">
                          {formatDate(review.createdAt)}
                        </p>
                      </CardContent>
                    </Card>
                  </Link>
                ))}
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
