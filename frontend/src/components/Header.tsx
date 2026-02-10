"use client";

import Link from "next/link";
import { useAuth } from "@/contexts/AuthContext";
import { Button } from "@/components/ui/button";

export function Header() {
  const { user, logout, isLoading } = useAuth();

  const handleLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.error("ログアウトエラー:", error);
    }
  };

  return (
    <header className="border-b bg-white">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        {/* ロゴ */}
        <Link href="/" className="text-xl font-bold text-primary">
          AnimeScore
        </Link>

        {/* ナビゲーション */}
        <nav className="flex items-center gap-4">
          <Link href="/search" className="text-sm hover:text-primary">
            検索
          </Link>
          <Link href="/animes" className="text-sm hover:text-primary">
            アニメ(平均点順)
          </Link>

          {isLoading ? (
            <span className="text-sm text-gray-400">読み込み中...</span>
          ) : user ? (
            <>
              <Link href="/me" className="text-sm hover:text-primary">
                マイページ
              </Link>
              <span className="text-sm text-gray-600">{user.username}</span>
              <Button variant="outline" size="sm" onClick={handleLogout}>
                ログアウト
              </Button>
            </>
          ) : (
            <>
              <Link href="/login">
                <Button variant="outline" size="sm">
                  ログイン
                </Button>
              </Link>
              <Link href="/signup">
                <Button size="sm">会員登録</Button>
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
