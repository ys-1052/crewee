# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリで作業する際のガイダンスを提供します。

## 言語ガイドライン

- **基本言語**: 基本的に日本語で応答してください
- **コード説明**: コードの説明などそのままがわかりやすい場合は英語でも構いません
- **コード実装**: コード内は全て英語で書いてください（変数名、関数名、コメントなど）
- **技術文書**: 技術仕様書やREADMEなどは日本語で記述してください

## プロジェクト概要

Creweeは、スポーツチームの代表者と選手をつなぐスポーツ特化のマッチングプラットフォームのMVPです。代表者が練習や試合の参加者を素早く集め、固定チーム未所属の選手が予定に合わせて参加できる場を提供します。

## アーキテクチャ

このプロジェクトは以下の技術で構築されるフルスタックWebアプリケーションです：

- **フロントエンド**: Next.js 15.4 (App Router), TypeScript, Tailwind CSS
- **バックエンド**: Go 1.24+, Echo Framework, sqlc
- **データベース**: PostgreSQL (Amazon RDS + RDS Proxy)
- **認証**: AWS Cognito
- **インフラ**: AWS CDK (TypeScript)
- **実行環境**: ECS Fargate
- **CI/CD**: GitHub Actions
- **ファイルストレージ**: Amazon S3
- **マイグレーション**: golang-migrate/migrate
- **監視**: CloudWatch + OpenTelemetry

## プロジェクト構造

現在、プロジェクトには `context/` ディレクトリ内の文書のみが含まれています：
- `context/design.md` - 技術設計書
- `context/requirements.md` - 要件定義書
- `context/tasks.md` - 実装計画とタスク分解

## 主要な技術的決定事項

### データベース設計
- ステータスフィールドにPostgreSQLのENUM型を使用
- 定員管理に楽観的ロックを実装
- JISコード準拠の地域データ構造
- 階層的地域データ（都道府県 → 市区町村）

### 認証・認可
- AWS Cognitoによるユーザー認証とメール認証
- JWT検証ミドルウェア
- 年齢制限（18歳以上のみ）
- ロールベースアクセス（チームオーナー、イベント作成者権限）

### API設計
- `/api/v1` ベースパスのRESTful API
- data/error/metaフィールドを持つ構造化JSONレスポンス
- 具体的なエラーコードによる包括的エラーハンドリング
- イベントと応募の状態遷移検証

### 主要エンティティ
- **Users**: ホーム地域を含むプロフィール管理
- **Teams**: ユーザーが所有し、スポーツに関連付けられる
- **Events**: チームが作成し、定員管理機能付き
- **Applications**: 承認ワークフロー付きのイベント応募
- **Sports**: 分類用マスタデータ
- **Regions**: 階層的な位置データ（JISコード準拠）

### ビジネスルール
- イベントの状態遷移: OPEN → CLOSED → CANCELED（ロールバック不可）
- 定員オーバーフローを防ぐ楽観的ロックによる応募承認
- 階層的地域データを使用した地理的検索
- イベントコメントによる自動通知システム

## 開発ワークフロー

現在は文書のみの計画段階のプロジェクトのため、実際の開発コマンドは実装に依存します。設計書に基づく想定コマンド：

### バックエンド (Go)
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

### フロントエンド (Next.js)
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

### インフラ (AWS CDK)
```bash
# 依存関係インストール
npm install

# インフラデプロイ
npx cdk deploy

# CloudFormation生成
npx cdk synth
```

## 主要機能の実装優先度

1. **認証システム** - AWS Cognito統合とJWT検証
2. **ユーザープロフィール管理** - 地域選択付き基本CRUD
3. **チーム管理** - チーム作成と所有者モデル
4. **イベント管理** - 状態遷移付きイベントCRUD
5. **検索・絞り込み** - 地理的・スポーツベースのイベント発見
6. **応募フロー** - 承認ワークフロー付きイベント応募
7. **コミュニケーション** - 基本的なイベントコメントと通知
8. **モデレーション** - ユーザー通報とブロック機能

## セキュリティ考慮事項

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

## 開発における注意事項

- 全ての文書は日本語で記述されています
- プロジェクトは `context/tasks.md` の詳細な実装計画に従います
- データベーススキーマ設計は参照整合性とパフォーマンスを重視
- システムは定員管理において防御的プログラミング手法を実装
- 将来的にはスポーツコミュニティ構築のためのSNS的機能への拡張を予定

## コードスタイル・品質

### Go (バックエンド)
- **Linter**: golangci-lint (.golangci.yml設定ファイル使用)
- **Formatter**: gofmt + goimports (自動実行)
- **Pre-commit**: git hooksでフォーマット・リント自動実行

### TypeScript/JavaScript (フロントエンド)
- **Linter**: ESLint (@typescript-eslint, @next/eslint-config-next)
- **Formatter**: Prettier (エディタ統合 + 自動実行)
- **Pre-commit**: Husky + lint-staged でコミット前チェック

### 自動フォーマット
Claude Codeの設定により、ファイル保存時にGoファイルとTypeScript/JavaScriptファイルが自動フォーマットされます。

## 開発ワークフローの重要なルール

### タスク管理
- **必須**: タスク完了時は必ず `context/tasks.md` ファイルの該当項目にチェック `[x]` を入れること
- **必須**: TodoWriteツールでの進捗管理と併せて、tasks.mdファイルも同時に更新すること
- 大きなタスクが完了したら、親タスクにもチェックを入れること

### コード品質・リントルール
- **必須**: 作業完了前にlintを実行し、すべてのlintエラーを解決すること
- **バックエンド**: `make lint` コマンドを実行してgolangci-lintが passすること
- **フロントエンド**: `npm run lint` コマンドを実行してESLintが passすること
- lintエラーが残ったままコミットしてはいけない

### コミットルール
- 論理的にまとまった単位でコミットを行うこと
