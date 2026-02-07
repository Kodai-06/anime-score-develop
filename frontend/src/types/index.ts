// ========== User ==========
export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
}

// ========== Auth ==========
export interface SignUpInput {
  username: string;
  email: string;
  password: string;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface SignUpResponse {
  message: string;
  user: User;
}

export interface LoginResponse {
  message: string;
  user: User;
}

export interface GetMeResponse {
  user: User;
}

// ========== Anime ==========
export interface Anime {
  id: number;
  annictId: number;
  title: string;
  year: number;
  imageUrl: string | null;
  createdAt: string;
}

export interface AnimeStats {
  animeId: number;
  reviewCount: number;
  avgScore: number;
}

export interface AnimeWithStats extends Anime {
  reviewCount: number;
  avgScore: number;
}

export interface Pagination {
  page: number;
  pageSize: number;
  total: number;
  totalPage: number;
}

export interface AnimeListResponse {
  data: AnimeWithStats[];
  pagination: Pagination;
}

// Annict API からの検索結果
export interface AnnictWork {
  annictId: number;
  title: string;
  seasonYear: number | null;
  image: {
    recommendedImageUrl: string;
  };
}

export interface AnimeSearchResponse {
  data: AnnictWork[];
  nextCursor: string | null;
}

export interface AnimeDetailResponse {
  anime: Anime;
  stats: AnimeStats | null;
}

// ========== Review ==========
export interface Review {
  id: number;
  userId: number;
  animeId: number;
  score: number;
  comment: string | null;
  createdAt: string;
}

export interface ReviewInput {
  annictId: number;
  score: number;
  comment?: string;
}

export interface ReviewWithAnime extends Review {
  animeTitle: string;
  animeYear: number;
  animeImageUrl: string | null;
}

export interface ReviewListResponse {
  data: Review[];
}

export interface MyReviewListResponse {
  data: ReviewWithAnime[];
}

export interface ReviewCreateResponse {
  message: string;
  review: Review;
}

// ========== API Error ==========
export interface ApiError {
  error: string;
}
