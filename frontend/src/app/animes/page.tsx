"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Header } from "@/components/Header";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { getAnimeList } from "@/lib/api";
import type { AnimeWithStats } from "@/types";

export default function AnimesPage() {
  const [animes, setAnimes] = useState<AnimeWithStats[]>([]);
  const [page, setPage] = useState(1);
  const [totalPage, setTotalPage] = useState(1);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // アニメ一覧を取得
  useEffect(() => {
    const fetchAnimes = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const response = await getAnimeList(page, 12);
        setAnimes(response.data || []);
        setTotalPage(response.pagination.totalPage);
      } catch (err) {
        setError("アニメの取得に失敗しました");
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchAnimes();
  }, [page]);

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-2xl font-bold">アニメ一覧(平均点順)</h1>

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

        {/* アニメ一覧 */}
        {!isLoading && !error && (
          <>
            {animes.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500">まだアニメがありません</p>
                <Link href="/search" className="text-primary hover:underline mt-2 inline-block">
                  アニメを検索してレビューをしよう
                </Link>
              </div>
            ) : (
              <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                {animes.map((anime) => (
                  <Link key={anime.id} href={`/animes/${anime.annictId}`}>
                    <Card className="h-full hover:shadow-md transition-shadow cursor-pointer">
                      <CardHeader className="pb-2">
                        <CardTitle className="text-base line-clamp-2">
                          {anime.title}
                        </CardTitle>
                      </CardHeader>
                      <CardContent>
                        <p className="text-sm text-gray-500">{anime.year}年</p>
                        <div className="mt-2 flex items-center gap-4">
                          <span className="text-lg font-bold text-primary">
                            {anime.avgScore ? `${anime.avgScore}点` : "- 点"}
                          </span>
                          <span className="text-sm text-gray-500">
                            {anime.reviewCount}件のレビュー
                          </span>
                        </div>
                      </CardContent>
                    </Card>
                  </Link>
                ))}
              </div>
            )}

            {/* ページネーション */}
            {totalPage > 1 && (
              <div className="mt-8 flex justify-center gap-2">
                <Button
                  variant="outline"
                  disabled={page === 1}
                  onClick={() => setPage(page - 1)}
                >
                  前へ
                </Button>
                <span className="flex items-center px-4">
                  {page} / {totalPage}
                </span>
                <Button
                  variant="outline"
                  disabled={page === totalPage}
                  onClick={() => setPage(page + 1)}
                >
                  次へ
                </Button>
              </div>
            )}
          </>
        )}
      </main>
    </div>
  );
}
