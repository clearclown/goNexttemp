# CLAUDE.md

このファイルはAI（Claude）がこのリポジトリで開発を行う際のガイドラインです。

## プロジェクト概要

Go (Gin) + Next.js によるフルスタックWebアプリケーションのテンプレート。
`frontend/` と `backend/` は独立したプロジェクトとして管理。

## アーキテクチャ

```
┌─────────────────┐         ┌─────────────────┐
│    Vercel       │         │  Fly.io      │
│  (Frontend)     │◀───────▶│  Go (Gin)      │
│   Next.js       │  REST   │                 │
│   PWA対応       │  API    │                 │
└─────────────────┘         └────────┬────────┘
                                     │
                            ┌────────▼────────┐
                            │   PostgreSQL    │
                            │  (Fly.io)   │
                            └─────────────────┘
```

## 技術スタック詳細

### フロントエンド
- Next.js (App Router) + TypeScript
- BiomeJS (Prettier/ESLint不使用)
- TailwindCSS + ShadCN/UI
- PWA (next-pwa)
- Vitest + Playwright

### バックエンド
- Go 1.22+ / Gin
- GORM (PostgreSQL)
- 独自JWT認証 (bcrypt + リフレッシュトークン)
- golang-migrate
- Air (ホットリロード)

---

## 開発ルール

### フロントエンド (`frontend/`)

#### ディレクトリ規約
```
src/
├── app/              # App Router（ページ・レイアウト）
├── components/
│   ├── ui/           # ShadCN/UI ベース（汎用）
│   └── features/     # 機能単位（認証、ダッシュボード等）
├── lib/              # ユーティリティ（api client, cn等）
├── hooks/            # カスタムフック
└── types/            # 型定義
```

#### コーディング規約
```typescript
// コンポーネント: 関数コンポーネント + named export
export function UserCard({ user }: UserCardProps) { ... }

// hooks: use プレフィックス
export function useAuth() { ... }

// 型: interface 推奨
interface UserCardProps {
  user: User;
}

// API呼び出し: lib/api.ts に集約
// 環境変数: NEXT_PUBLIC_ プレフィックス（クライアント用）
```

#### PWA関連
- `public/manifest.json` - アプリ情報
- `public/icons/` - 各サイズアイコン (192x192, 512x512等)
- Service Worker は next-pwa が自動生成

#### テスト
```bash
pnpm test          # Vitest 単体テスト
pnpm test:e2e      # Playwright E2E
```

### バックエンド (`backend/`)

#### ディレクトリ規約
```
backend/
├── cmd/server/           # main.go（エントリポイントのみ）
├── internal/
│   ├── config/           # 設定読み込み (envconfig)
│   ├── handler/          # HTTPハンドラ（入力検証・レスポンス）
│   ├── service/          # ビジネスロジック
│   ├── repository/       # DB操作（GORM）
│   ├── model/            # ドメインモデル・エンティティ
│   ├── middleware/       # Gin ミドルウェア
│   └── auth/             # 認証ロジック (JWT, bcrypt)
├── pkg/
│   └── response/         # 共通JSONレスポンス
├── migrations/           # golang-migrate SQL
├── go.mod
└── .air.toml             # Air設定
```

#### 依存関係（go.mod）
```go
require (
    github.com/gin-gonic/gin
    github.com/gin-contrib/cors
    gorm.io/gorm
    gorm.io/driver/postgres
    github.com/golang-jwt/jwt/v5
    golang.org/x/crypto              // bcrypt
    github.com/go-playground/validator/v10
    github.com/kelseyhightower/envconfig
    github.com/stretchr/testify
)
```

#### コーディング規約
```go
// パッケージ名: 小文字、単数形
package handler

// 構造体: PascalCase
type UserHandler struct {
    userService service.UserService
}

// コンストラクタ: New プレフィックス
func NewUserHandler(s service.UserService) *UserHandler {
    return &UserHandler{userService: s}
}

// Gin ハンドラ: *gin.Context
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.userService.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, response.Error("User not found"))
        return
    }
    c.JSON(http.StatusOK, response.Success(user))
}

// エラーハンドリング: 早期リターン、エラーラップ
if err != nil {
    return fmt.Errorf("failed to get user: %w", err)
}
```

#### レイヤー間の依存
```
handler → service → repository
              ↓
           model
```
※ 逆方向の依存禁止、interfaceで疎結合に

#### 認証実装

**トークン仕様:**
| 種類 | 有効期限 | 保存 | 用途 |
|-----|---------|-----|-----|
| Access Token | 15分 | Authorization Header | API認証 |
| Refresh Token | 7日 | httpOnly Cookie + DB | トークン更新 |

**JWT Claims:**
```go
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}
```

**パスワード:** bcrypt (cost=12)

#### マイグレーション
```bash
# 新規作成
migrate create -ext sql -dir migrations -seq create_users

# 適用
migrate -path migrations -database "postgres://..." up

# ロールバック
migrate -path migrations -database "postgres://..." down 1
```

### インフラ (`infra/`)

#### Podman Compose
```yaml
# compose.yaml の構成
services:
  postgres:
    image: postgres:16-alpine
    ports:
      - "${POSTGRES_PORT:-5432}:5432"

  backend:
    build:
      context: ../../backend
      dockerfile: ../infra/podman/Containerfile.backend
    ports:
      - "${BACKEND_PORT:-8080}:8080"
    depends_on:
      - postgres

  frontend:
    build:
      context: ../../frontend
      dockerfile: ../infra/podman/Containerfile.frontend
    ports:
      - "${FRONTEND_PORT:-3000}:3000"
```

#### 環境変数管理
- `.env.example` を常に最新に
- 機密情報は `.env`（Git管理外）
- Podman内では `env_file` で読み込み

---

## API設計

### 基本ルール
- ベースパス: `/api/v1`
- JSON形式のみ
- 認証必須エンドポイントは `Authorization: Bearer <token>`

### レスポンス形式
```json
// 成功
{
  "success": true,
  "data": { ... }
}

// エラー
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Email is required"
  }
}
```

### 認証エンドポイント
```
POST   /api/v1/auth/register   # 登録
POST   /api/v1/auth/login      # ログイン
POST   /api/v1/auth/refresh    # トークン更新
POST   /api/v1/auth/logout     # ログアウト
GET    /api/v1/auth/me         # 自分の情報
```

---

## ファイルテンプレート

### React コンポーネント
```typescript
interface ComponentNameProps {
  // props定義
}

export function ComponentName({ }: ComponentNameProps) {
  return (
    <div>
      {/* content */}
    </div>
  );
}
```

### Gin ハンドラ
```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type ResourceHandler struct {
    service service.ResourceService
}

func NewResourceHandler(s service.ResourceService) *ResourceHandler {
    return &ResourceHandler{service: s}
}

func (h *ResourceHandler) Get(c *gin.Context) {
    id := c.Param("id")

    resource, err := h.service.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, response.Error("Resource not found"))
        return
    }

    c.JSON(http.StatusOK, response.Success(resource))
}
```

### GORM モデル
```go
package model

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Email     string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

### マイグレーションSQL
```sql
-- migrations/000001_create_users.up.sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

```sql
-- migrations/000001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

---

## 禁止事項

- `frontend/` に Prettier/ESLint 設定を追加（BiomeJS統一）
- Pages Router の使用（App Router統一）
- `any` 型の多用
- `internal/` パッケージの外部公開
- 環境変数のハードコーディング
- パスワードの平文保存
- JWT秘密鍵のコミット

## 推奨事項

- Server Components を積極活用
- 型定義は厳密に
- エラーメッセージは具体的に
- コミットは Conventional Commits 形式
- テストは実装と同時に書く
