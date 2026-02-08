"use client";

import { useState } from "react";
import Link from "next/link";
import { Header } from "@/components/Header";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { searchAnimes } from "@/lib/api";
import type { AnnictWork } from "@/types";

export default function SearchPage() {
  const [keyword, setKeyword] = useState("");
  const [results, setResults] = useState<AnnictWork[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasSearched, setHasSearched] = useState(false);

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!keyword.trim()) return; // 空のキーワードは無視,trim()で前後の空白を削除

    setIsLoading(true);
    setError(null);
    setHasSearched(true);

    try {
      const response = await searchAnimes(keyword);
      setResults(response.data || []);
    } catch (err) {
      setError("検索に失敗しました");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <h1 className="mb-6 text-2xl font-bold">アニメ検索</h1>

        {/* 検索フォーム */}
        <form onSubmit={handleSearch} className="mb-8 flex gap-2">
          <Input
            type="text"
            value={keyword}
            onChange={(e) => setKeyword(e.target.value)}
            placeholder="アニメのタイトルを入力"
            className="max-w-md"
          />
          <Button type="submit" disabled={isLoading || !keyword.trim()}>
            {isLoading ? "検索中..." : "検索"}
          </Button>
        </form>

        {/* エラー */}
        {error && (
          <div className="rounded bg-red-50 p-4 text-red-600">{error}</div>
        )}

        {/* ローディング */}
        {isLoading && (
          <div className="text-center py-12">
            <p className="text-gray-500">検索中...</p>
          </div>
        )}

        {/* 検索結果 */}
        {!isLoading && hasSearched && (
          <>
            {results.length === 0 ? (
              <div className="text-center py-12">
                <p className="text-gray-500">検索結果が見つかりませんでした</p>
              </div>
            ) : (
              <div className="space-y-3">
                <p className="text-sm text-gray-600">
                  {results.length}件の結果
                </p>
                {results.map((anime) => (
                  <Link key={anime.annictId} href={`/animes/${anime.annictId}`}>
                    <Card className="hover:shadow-md transition-shadow cursor-pointer">
                      <CardContent className="flex items-center gap-4 py-4">
                        <div className="flex-1">
                          <h3 className="font-medium">{anime.title}</h3>
                          <p className="text-sm text-gray-500">
                            {anime.seasonYear ? `${anime.seasonYear}年` : "放送年不明"}
                          </p>
                        </div>
                        <Button variant="outline" size="sm">
                          詳細を見る
                        </Button>
                      </CardContent>
                    </Card>
                  </Link>
                ))}
              </div>
            )}
          </>
        )}

        {/* 未検索時のヒント */}
        {!hasSearched && (
          <div className="text-center py-12">
            <p className="text-gray-500">
              アニメのタイトルを入力して検索してください
            </p>
          </div>
        )}
      </main>
    </div>
  );
}
