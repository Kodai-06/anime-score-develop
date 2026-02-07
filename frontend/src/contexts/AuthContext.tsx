"use client";

import { createContext, useContext, useState, useEffect } from "react";
import type { User, LoginInput, SignUpInput } from "@/types";
import {
  login as loginApi,
  signup as signupApi,
  logout as logoutApi,
  getCurrentUser,
  ApiError,
} from "@/lib/api";

// ========== Contextで共有する値の型 ==========
interface AuthContextType {
  user: User | null; // ログイン中のユーザー情報（未ログインならnull）
  isLoading: boolean; // 認証状態を確認中かどうか
  login: (input: LoginInput) => Promise<void>;
  signup: (input: SignUpInput) => Promise<void>;
  logout: () => Promise<void>;
}

// ========== Contextの作成 ==========
// データを共有するためのコンテキスト(箱のようなもの)を作る
// 初期値はundefined（Providerの外で使われた場合のエラー検出用）
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// ========== Provider コンポーネント ==========
// アプリ全体をこれで囲むと、どこからでも認証情報にアクセスできる
// childrenはreactが自動で渡してくれる開始タグと終了タグので挟んだ要素が入る特別な変数
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // ページ読み込み時に認証状態を確認
  useEffect(() => {
    getCurrentUser()
      .then((response) => {
        setUser(response.user);
      })
      .catch((error) => {
        // 401エラーは未ログイン状態なので正常
        if (error instanceof ApiError && error.status === 401) {
          setUser(null);
        } else {
          console.error("認証確認エラー:", error);
        }
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, []);

  // ログイン処理
  const login = async (input: LoginInput) => {
    const response = await loginApi(input);
    setUser(response.user);
  };

  // サインアップ処理
  const signup = async (input: SignUpInput) => {
    const response = await signupApi(input);
    setUser(response.user);
  };

  // ログアウト処理
  const logout = async () => {
    await logoutApi();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, login, signup, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// ========== カスタムフック ==========
// コンポーネントからContextの値を取得するためのフック
export function useAuth() {
  const context = useContext(AuthContext);

  // Providerの外で使われた場合はエラー
  if (context === undefined) {
    throw new Error("useAuthはAuthProviderの中で使ってください");
  }

  return context;
}
