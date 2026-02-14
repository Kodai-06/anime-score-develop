import { cookies } from "next/headers";
import { NextRequest, NextResponse } from "next/server";

// バックエンドのURL（サーバーサイドのみで使用するため NEXT_PUBLIC_ は不要）
const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

// Cookie の設定値
const COOKIE_NAME = "auth_token";
const COOKIE_MAX_AGE = 60 * 60 * 24; // 24時間（秒）

// ========== 認証用 Cookie を設定するヘルパー ==========
function setAuthCookie(response: NextResponse, token: string) {
  response.cookies.set(COOKIE_NAME, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: COOKIE_MAX_AGE,
  });
}

function clearAuthCookie(response: NextResponse) {
  response.cookies.set(COOKIE_NAME, "", {
    httpOnly: true,
    secure: process.env.NODE_ENV === "production",
    sameSite: "lax",
    path: "/",
    maxAge: 0,
  });
}

// ========== BFF プロキシ本体 ==========
async function proxyToBackend(
  req: NextRequest,
  params: { path: string[] }
) {
    //Next.jsは、APIルートのパスパラメータを配列として提供する
  const path = params.path.join("/"); // [...path] の配列を / で結合して文字列にする（例: ["users", "profile"] → "users/profile"）
  const queryString = new URL(req.url).search; // ?page=1&pageSize=10 等
  const backendUrl = `${BACKEND_URL}/api/${path}${queryString}`;

  // ── Cookie からトークンを読み取り ──
  const cookieStore = await cookies();
  const token = cookieStore.get(COOKIE_NAME)?.value;

  // ── バックエンドへ送るヘッダーを構築 ──
  const headers: HeadersInit = {
    "Content-Type": "application/json",
  };

  if (token) {
    // Cookie のトークンを Authorization: Bearer ヘッダーに乗せ換え
    headers["Authorization"] = `Bearer ${token}`;
  }

  // ── リクエストボディの転送（GET/HEAD 以外） ──
  let body: string | undefined;
  if (req.method !== "GET" && req.method !== "HEAD") {
    body = await req.text();
  }

  // ── バックエンドへリクエスト ──
  const backendRes = await fetch(backendUrl, {
    method: req.method,
    headers,
    body,
  });

  // ── レスポンスの処理 ──
  // JSON パースを試みる（失敗した場合は空オブジェクト）
  let data: Record<string, unknown>;
  try {
    data = await backendRes.json();
  } catch {
    data = {};
  }

  // ── login / signup: レスポンスからトークンを取り出して Cookie にセット ──
  if (
    (path === "login" || path === "signup") &&
    backendRes.ok &&
    typeof data.token === "string"
  ) {
    // トークンをクライアントに返さないようコピーから除去
    // data オブジェクトから token プロパティを取り出し、残りを rest に格納する
    const { token: extractedToken, ...rest } = data;
    const response = NextResponse.json(rest, { status: backendRes.status });
    setAuthCookie(response, extractedToken);
    return response;
  }

  // ── logout: Cookie を削除 ──
  if (path === "logout" && backendRes.ok) {
    const response = NextResponse.json(data, { status: backendRes.status });
    clearAuthCookie(response);
    return response;
  }

  // ── その他: そのまま返す ──
  return NextResponse.json(data, { status: backendRes.status });
}

// ========== HTTP メソッドごとのハンドラー ==========
type RouteContext = { params: Promise<{ path: string[] }> };

// GETやPOSTの第２引数にcontextオブジェクトが渡され、その中にparamsが含まれる。
export async function GET(req: NextRequest, context: RouteContext) {
  return proxyToBackend(req, await context.params);
}

export async function POST(req: NextRequest, context: RouteContext) {
  return proxyToBackend(req, await context.params);
}

export async function PUT(req: NextRequest, context: RouteContext) {
  return proxyToBackend(req, await context.params);
}

export async function DELETE(req: NextRequest, context: RouteContext) {
  return proxyToBackend(req, await context.params);
}

export async function PATCH(req: NextRequest, context: RouteContext) {
  return proxyToBackend(req, await context.params);
}
