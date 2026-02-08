"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Header } from "@/components/Header";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { useAuth } from "@/contexts/AuthContext";
import { getAnimeDetail, getReviewsByAnime, createReview, ApiError } from "@/lib/api";
import type { Anime, AnimeStats, Review } from "@/types";

// レビュー投稿のバリデーションスキーマ
const reviewSchema = z.object({
  score: z
    .number({ error: "数値を入力してください" })
    .min(0, "0以上を入力してください")
    .max(100, "100以下を入力してください"),
  comment: z.string().optional(),
});

type ReviewFormData = z.infer<typeof reviewSchema>;

export default function AnimeDetailPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const annictId = parseInt(id, 10);
  const router = useRouter();
  const { user } = useAuth();

  // 状態
  const [anime, setAnime] = useState<Anime | null>(null);
  const [stats, setStats] = useState<AnimeStats | null>(null);
  const [reviews, setReviews] = useState<Review[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // レビュー投稿フォーム
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitSuccess, setSubmitSuccess] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ReviewFormData>({
    resolver: zodResolver(reviewSchema),
  });

  // アニメ詳細とレビューを取得
  useEffect(() => {
    const fetchData = async () => {
      setIsLoading(true);
      setError(null);
      try {
        const detailResponse = await getAnimeDetail(annictId);
        setAnime(detailResponse.anime);
        setStats(detailResponse.stats);

        // アニメがDBに存在する場合のみレビューを取得
        if (detailResponse.anime?.id) {
          const reviewsResponse = await getReviewsByAnime(detailResponse.anime.id);
          setReviews(reviewsResponse.data || []);
        }
      } catch (err) {
        if (err instanceof ApiError && err.status === 404) {
          setError("アニメが見つかりませんでした");
        } else {
          setError("データの取得に失敗しました");
        }
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchData();
  }, [annictId]);

  // レビュー投稿
  const onSubmitReview = async (data: ReviewFormData) => {
    setSubmitError(null);
    setSubmitSuccess(false);

    try {
      await createReview({
        annictId,
        score: data.score,
        comment: data.comment || undefined,
      });
      setSubmitSuccess(true);
      reset();

      // データを再取得
      const detailResponse = await getAnimeDetail(annictId);
      setAnime(detailResponse.anime);
      setStats(detailResponse.stats);
      if (detailResponse.anime?.id) {
        const reviewsResponse = await getReviewsByAnime(detailResponse.anime.id);
        setReviews(reviewsResponse.data || []);
      }
    } catch (err) {
      if (err instanceof ApiError) {
        setSubmitError(err.message);
      } else {
        setSubmitError("レビューの投稿に失敗しました");
      }
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        {/* ローディング */}
        {isLoading && (
          <div className="text-center py-12">
            <p className="text-gray-500">読み込み中...</p>
          </div>
        )}

        {/* エラー */}
        {error && (
          <div className="text-center py-12">
            <p className="text-red-500 mb-4">{error}</p>
            <Button variant="outline" onClick={() => router.back()}>
              戻る
            </Button>
          </div>
        )}

        {/* アニメ詳細 */}
        {!isLoading && !error && anime && (
          <div className="space-y-8">
            {/* アニメ情報 */}
            <Card>
              <CardHeader>
                <CardTitle className="text-xl">{anime.title}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-6">
                  <div>
                    <p className="text-sm text-gray-500">放送年</p>
                    <p className="font-medium">{anime.year}年</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">平均スコア</p>
                    <p className="text-2xl font-bold text-primary">
                      {stats?.avgScore ? `${stats.avgScore}点` : "- 点"}
                    </p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-500">レビュー数</p>
                    <p className="font-medium">{stats?.reviewCount || 0}件</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* レビュー投稿フォーム */}
            {user ? (
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">レビューを投稿</CardTitle>
                </CardHeader>
                <CardContent>
                  <form onSubmit={handleSubmit(onSubmitReview)} className="space-y-4">
                    {submitError && (
                      <div className="rounded bg-red-50 p-3 text-sm text-red-600">
                        {submitError}
                      </div>
                    )}
                    {submitSuccess && (
                      <div className="rounded bg-green-50 p-3 text-sm text-green-600">
                        レビューを投稿しました！
                      </div>
                    )}

                    <div className="space-y-2">
                      <Label htmlFor="score">スコア（0〜100）</Label>
                      <Input
                        id="score"
                        type="number"
                        placeholder="80"
                        className="max-w-32"
                        {...register("score", { valueAsNumber: true })}
                      />
                      {errors.score && (
                        <p className="text-sm text-red-500">{errors.score.message}</p>
                      )}
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="comment">コメント（任意）</Label>
                      <Textarea
                        id="comment"
                        placeholder="感想を書いてください"
                        rows={3}
                        {...register("comment")}
                      />
                    </div>

                    <Button type="submit" disabled={isSubmitting}>
                      {isSubmitting ? "投稿中..." : "投稿する"}
                    </Button>
                  </form>
                </CardContent>
              </Card>
            ) : (
              <Card>
                <CardContent className="py-6 text-center">
                  <p className="text-gray-600 mb-4">
                    レビューを投稿するにはログインが必要です
                  </p>
                  <Button onClick={() => router.push("/login")}>
                    ログインする
                  </Button>
                </CardContent>
              </Card>
            )}

            {/* レビュー一覧 */}
            <div>
              <h2 className="text-lg font-bold mb-4">レビュー一覧</h2>
              {reviews.length === 0 ? (
                <p className="text-gray-500">まだレビューがありません</p>
              ) : (
                <div className="space-y-4">
                  {reviews.map((review) => (
                    <Card key={review.id}>
                      <CardContent className="py-4">
                        <div className="flex items-start justify-between">
                          <div className="flex-1">
                            <span className="text-lg font-bold text-primary">
                              {review.score}点
                            </span>
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
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
