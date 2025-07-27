# Crewee

Creweeは、スポーツチームの代表者と選手をつなぐスポーツ特化のマッチングプラットフォームのMVPです。代表者が練習や試合の参加者を素早く集め、固定チーム未所属の選手が予定に合わせて参加できる場を提供します。

## 技術スタック

### フロントエンド
- Next.js 15.4 (App Router)
- TypeScript
- Tailwind CSS
- React Query

### バックエンド
- Go 1.24+
- Echo Framework
- sqlc
- golang-migrate

### データベース
- PostgreSQL
- Amazon RDS + RDS Proxy

### 認証
- AWS Cognito

### インフラ
- AWS CDK (TypeScript)
- ECS Fargate
- Amazon S3
- CloudWatch + OpenTelemetry

### CI/CD
- GitHub Actions

## プロジェクト構造

```
crewee/
├── backend/           # Go backend application
├── frontend/          # Next.js frontend application
├── infrastructure/    # AWS CDK infrastructure code
├── docker/           # Docker configuration files
├── docs/             # Additional documentation
└── context/          # Project documentation
    ├── design.md     # 技術設計書
    ├── requirements.md # 要件定義書
    └── tasks.md      # 実装計画とタスク分解
```

## 開発環境のセットアップ

### 前提条件
- Docker & Docker Compose
- Go 1.24+
- Node.js 18+
- PostgreSQL (Dockerで提供)

### 環境構築手順

1. リポジトリのクローン
```bash
git clone <repository-url>
cd crewee
```

2. 開発用データベースの起動
```bash
docker-compose up -d
```

3. バックエンドの起動
```bash
cd backend
go mod tidy
make dev
```

4. フロントエンドの起動
```bash
cd frontend
npm install
npm run dev
```

### 利用可能なコマンド

#### バックエンド (Go)
```bash
# 依存関係インストール
go mod tidy

# マイグレーション実行
migrate -path ./migrations -database "postgres://..." up

# sqlcコード生成
sqlc generate

# テスト実行
go test ./...

# リント実行
golangci-lint run

# フォーマット実行
gofmt -w . && goimports -w .

# 開発サーバー起動
go run cmd/server/main.go
```

#### フロントエンド (Next.js)
```bash
# 依存関係インストール
npm install

# 開発サーバー起動
npm run dev

# プロダクションビルド
npm run build

# テスト実行
npm test

# コードリント・フォーマット
npm run lint
npm run lint:fix
npx prettier --write .
```

#### インフラ (AWS CDK)
```bash
# 依存関係インストール
npm install

# インフラデプロイ
npx cdk deploy

# CloudFormation生成
npx cdk synth
```

## 主要機能

1. **認証システム** - AWS Cognito統合とJWT検証
2. **ユーザープロフィール管理** - 地域選択付き基本CRUD
3. **チーム管理** - チーム作成と所有者モデル
4. **イベント管理** - 状態遷移付きイベントCRUD
5. **検索・絞り込み** - 地理的・スポーツベースのイベント発見
6. **応募フロー** - 承認ワークフロー付きイベント応募
7. **コミュニケーション** - 基本的なイベントコメントと通知
8. **モデレーション** - ユーザー通報とブロック機能

## セキュリティ

- HTTPS強制
- AWS WAFレート制限
- CORS設定
- 入力検証とサニタイゼーション
- XSSとCSRF保護
- JWT署名検証
- データベースクエリのパラメータ化

## パフォーマンス要件

- ページ読み込み時間（LCP）< 2.5秒
- 重要な操作のAPI応答時間 < 200ms
- 適切なインデックスによるデータベースクエリ最適化
- RDS Proxyによるコネクションプーリング

## 監視・可観測性

- AWS X-RayとのOpenTelemetry統合
- CloudWatchメトリクスとアラーム
- トレース相関付き構造化ログ
- エラー追跡とパフォーマンス監視

## ライセンス

MIT License