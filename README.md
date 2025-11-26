# goNexttemp

Go + Next.js によるフルスタックWebアプリケーションのテンプレートリポジトリ。
PWA対応、独自認証システム搭載。

## 技術スタック

### フロントエンド (`frontend/`)

| カテゴリ | 技術 |
|---------|------|
| パッケージ管理 | pnpm |
| フレームワーク | Next.js (App Router) |
| 言語 | TypeScript |
| Linter/Formatter | BiomeJS |
| スタイリング | TailwindCSS |
| UIコンポーネント | ShadCN/UI |
| PWA | next-pwa |
| 単体・統合テスト | Vitest |
| E2Eテスト | Playwright |

### バックエンド (`backend/`)

| カテゴリ | 技術 |
|---------|------|
| 言語 | Go 1.22+ |
| フレームワーク | Gin |
| ORM | GORM |
| DB | PostgreSQL |
| 認証 | 独自実装 (JWT + bcrypt) |
| バリデーション | go-playground/validator |
| ログ | slog (標準ライブラリ) |
| 設定管理 | envconfig |
| マイグレーション | golang-migrate |
| ホットリロード | Air |

### インフラ

| カテゴリ | 技術 |
|---------|------|
| タスクランナー | Task (go-task) |
| コンテナ | Podman (podman-compose) |
| バックエンドホスティング | Firebase (Cloud Run) |
| フロントエンドホスティング | Vercel |
| CI/CD | GitHub Actions |

## ディレクトリ構成

```
goNexttemp/
├── frontend/                    # Next.js プロジェクト（独立）
│   ├── src/
│   │   ├── app/                # App Router
│   │   ├── components/
│   │   │   ├── ui/            # ShadCN/UI
│   │   │   └── features/      # 機能別コンポーネント
│   │   ├── lib/               # ユーティリティ
│   │   ├── hooks/             # カスタムフック
│   │   └── types/             # 型定義
│   ├── tests/
│   │   ├── unit/              # Vitest
│   │   └── e2e/               # Playwright
│   ├── public/
│   │   ├── manifest.json      # PWAマニフェスト
│   │   └── icons/             # PWAアイコン
│   ├── package.json
│   ├── biome.json
│   └── next.config.ts
│
├── backend/                     # Go プロジェクト（独立）
│   ├── cmd/
│   │   └── server/            # main.go
│   ├── internal/
│   │   ├── handler/           # HTTPハンドラ
│   │   ├── service/           # ビジネスロジック
│   │   ├── repository/        # データアクセス
│   │   ├── model/             # データモデル
│   │   ├── middleware/        # 認証ミドルウェア等
│   │   └── auth/              # 認証ロジック (JWT, bcrypt)
│   ├── pkg/
│   │   └── response/          # 共通レスポンス
│   ├── migrations/            # DBマイグレーション
│   ├── go.mod
│   └── go.sum
│
├── infra/
│   ├── podman/
│   │   ├── compose.yaml       # frontend + backend + postgres
│   │   ├── Containerfile.frontend
│   │   └── Containerfile.backend
│   └── firebase/
│
├── .github/
│   └── workflows/
│       ├── frontend.yml
│       └── backend.yml
│
├── .env.example
├── .gitignore
├── Taskfile.yml               # タスクランナー設定
├── CLAUDE.md
└── README.md
```

## クイックスタート

### 前提条件

- Node.js 20+
- pnpm 9+
- Go 1.22+
- Podman + podman-compose
- Task (go-task)

**Task のインストール:**

```bash
# macOS
brew install go-task

# その他: https://taskfile.dev/installation/
```

### セットアップ（ワンコマンド）

```bash
# リポジトリをクローン
git clone https://github.com/your-username/goNexttemp.git
cd goNexttemp

# 一括セットアップ（依存チェック → 環境構築 → 起動 → マイグレーション）
task setup
```

### 開発コマンド

```bash
# コマンド一覧を表示
task

# === 一括実行 ===
task setup              # 初回セットアップ（全自動）
task dev                # 開発環境起動 + ログフォロー
task ci                 # CI用チェック（lint + test + build）
task check              # 全体lint
task test               # 全体テスト
task build              # 全体ビルド

# === Podman ===
task up                 # 全サービス起動
task down               # 全サービス停止
task restart            # 全サービス再起動
task logs               # ログ確認（フォロー）
task ps                 # コンテナ状態確認
task clean              # コンテナ・ボリューム・キャッシュ削除

# === フロントエンド ===
task frontend:dev       # 開発サーバー (localhost:3000)
task frontend:build     # ビルド
task frontend:lint      # BiomeJS lint
task frontend:lint:fix  # BiomeJS lint + 自動修正
task frontend:test      # Vitest テスト
task frontend:test:e2e  # Playwright E2E

# === バックエンド ===
task backend:dev        # Air ホットリロード (localhost:8080)
task backend:build      # ビルド
task backend:lint       # go vet + staticcheck
task backend:test       # Go テスト

# === データベース ===
task migrate:up         # マイグレーション適用
task migrate:down       # ロールバック（1つ）
task migrate:status     # 状態確認
task migrate:create NAME=xxx  # 新規マイグレーション作成
task db:shell           # PostgreSQL に接続
task db:reset           # データベースリセット

# === ユーティリティ ===
task check:deps         # 依存ツールの確認
task env:init           # .env ファイル初期化
```

## 環境変数

`.env.example` をコピーして `.env` を作成：

```env
# === 共通 ===
NODE_ENV=development

# === ポート ===
FRONTEND_PORT=3000
BACKEND_PORT=8080
POSTGRES_PORT=5432

# === PostgreSQL ===
POSTGRES_HOST=localhost
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=goNexttemp

# === JWT認証 ===
JWT_SECRET=your-super-secret-key-change-in-production
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h

# === Firebase (本番用) ===
FIREBASE_PROJECT_ID=
FIREBASE_API_KEY=

# === CORS ===
CORS_ORIGINS=http://localhost:3000
```

## 認証システム

独自JWT認証を実装：

| トークン | 有効期限 | 保存場所 |
|---------|---------|---------|
| アクセストークン | 15分 | メモリ / localStorage |
| リフレッシュトークン | 7日 | httpOnly Cookie |

### エンドポイント

```
POST /api/v1/auth/register  # ユーザー登録
POST /api/v1/auth/login     # ログイン
POST /api/v1/auth/refresh   # トークン更新
POST /api/v1/auth/logout    # ログアウト
GET  /api/v1/auth/me        # 現在のユーザー情報
```

## PWA対応

- オフラインキャッシュ対応
- ホーム画面へのインストール
- プッシュ通知対応（要設定）

## デプロイ

### フロントエンド → Vercel

1. Vercelでプロジェクト作成
2. Root Directory: `frontend`
3. Framework Preset: Next.js
4. 環境変数を設定

### バックエンド → Firebase Cloud Run

1. GitHub Actions で自動デプロイ
2. `backend.yml` ワークフローが実行
3. Firebase プロジェクト設定が必要

## ライセンス

MIT
