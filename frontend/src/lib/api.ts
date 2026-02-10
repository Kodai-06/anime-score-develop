import type {
  SignUpInput,
  SignUpResponse,
  LoginInput,
  LoginResponse,
  GetMeResponse,
  AnimeListResponse,
  AnimeSearchResponse,
  AnimeDetailResponse,
  ReviewListResponse,
  MyReviewListResponse,
  ReviewInput,
  ReviewCreateResponse,
} from "@/types";

// ========== 設定 ==========
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// ========== Fetch ラッパー ==========
class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {} // 通信の細かい設定,何も渡されなかったら空オブジェクト
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;

    const headers: HeadersInit = {
      "Content-Type": "application/json",
      ...options.headers, //スプレッド構文でオプションのヘッダーを展開(オプションがあればそれを優先するため最後に配置)
    };

    const response = await fetch(url, {
      ...options,
      credentials: 'include', //ブラウザはデフォルトでは別ドメインへのリクエストにクッキーを送信しないため、これを有効にする
      headers,
    });

    // エラーレスポンスの処理
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({})); // JSONパースに失敗した場合は空オブジェクトを返す
      const message = errorData.error || `HTTP Error: ${response.status}`;
      throw new ApiError(message, response.status);
    }

    return response.json();
  }

  // GET リクエスト
  get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: "GET" });
  }

  // POST リクエスト
  post<T>(endpoint: string, body?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: "POST",
      body: body ? JSON.stringify(body) : undefined,
    });
  }
}

// カスタムエラークラス
// Errorクラスにstatusプロパティを追加
export class ApiError extends Error {
  status: number;

  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

// APIクライアントのインスタンス
const api = new ApiClient(API_BASE_URL);

// ========== Auth API ==========
export async function signup(input: SignUpInput): Promise<SignUpResponse> {
  return api.post<SignUpResponse>("/api/signup", input);
}

export async function login(input: LoginInput): Promise<LoginResponse> {
  const response = await api.post<LoginResponse>("/api/login", input);
  
  return response;
}

export async function logout(): Promise<void> {
  await api.post<void>("/api/logout");
}

export async function getCurrentUser(): Promise<GetMeResponse> {
  return api.get<GetMeResponse>("/api/me");
}

// ========== Anime API ==========
export async function getAnimeList(
  page: number = 1,
  pageSize: number = 10
): Promise<AnimeListResponse> {
  return api.get<AnimeListResponse>(
    `/api/animes?page=${page}&pageSize=${pageSize}`
  );
}

// URLには特定の記号は使えないのでエンコードする
export async function searchAnimes(
  keyword: string,
  limit: number = 15,
  cursor?: string
): Promise<AnimeSearchResponse> {
  let url = `/api/animes/search?q=${encodeURIComponent(keyword)}&limit=${limit}`;
  if (cursor) {
    url += `&cursor=${encodeURIComponent(cursor)}`;
  }
  return api.get<AnimeSearchResponse>(url);
}

export async function getAnimeDetail(
  annictId: number
): Promise<AnimeDetailResponse> {
  return api.get<AnimeDetailResponse>(`/api/animes/${annictId}`);
}

// ========== Review API ==========
export async function getReviewsByAnime(
  animeId: number
): Promise<ReviewListResponse> {
  return api.get<ReviewListResponse>(`/api/reviews?anime_id=${animeId}`);
}

export async function createReview(
  input: ReviewInput
): Promise<ReviewCreateResponse> {
  return api.post<ReviewCreateResponse>("/api/reviews", input);
}

export async function getMyReviews(): Promise<MyReviewListResponse> {
  return api.get<MyReviewListResponse>("/api/me/reviews");
}

export async function getRecentReviews(): Promise<MyReviewListResponse> {
  return api.get<MyReviewListResponse>("/api/reviews/recent");
}