# 地域マスタデータの変換・投入手順

このドキュメントでは、local_government_codes.xlsxファイルを使用して、地域マスタデータをデータベースに投入する手順を説明します。

## 概要

local_government_codes.xlsxファイルには、JISコード準拠の都道府県・市区町村データが含まれています。
このシステムは**差分更新**に対応しており、変更されたデータのみを効率的に処理します。

### データソース

**地方公共団体コードファイル（Excel形式）**
- **提供元**: 総務省
- **ダウンロード先URL**: https://www.soumu.go.jp/denshijiti/code.html
- **ファイル名**: 地方公共団体コード住所.xlsx → `local_government_codes.xlsx`にリネーム
- **更新頻度**: 年1回程度（市町村合併等により更新）

### 処理方式

**自動差分検出方式**: 初回は全データ投入、以降は変更されたデータのみをマイグレーション

### データ処理フロー

1. Excel → UTF-8 CSV変換（一時ファイル）
2. 新旧CSV比較による差分検出（現在のCSV vs 前回のバックアップ）
3. 差分マイグレーションSQL生成（INSERT/UPDATE/DELETE）
4. 現在のCSVを次回用バックアップとして保存
5. **一時CSVファイルの自動削除**（処理完了後）
6. 手動でマイグレーション実行（`make migrate-up`）

## ファイル構成

```
crewee/
├── local_government_codes.xlsx      # 元データ（総務省からダウンロード・リネーム）
├── scripts/
│   ├── convert_lgcode_to_csv.py          # Excel→CSV変換スクリプト
│   └── generate_region_diff_migration.py # マイグレーションSQL生成スクリプト（自動差分検出）
├── backend/
│   ├── sql/
│   │   └── local_government_codes.csv.backup # 前回処理済みデータ（差分比較用）
│   │       # 注：.csvファイルは処理完了後に自動削除される
│   └── migrations/
│       └── *_XXX_update_regions_data.up.sql  # 地域データ更新（初回・差分共通）
└── venv/                            # Python仮想環境
```

## 使用方法

### 基本コマンド

**前提条件**: 総務省から`local_government_codes.xlsx`をダウンロードしてプロジェクトルートに配置

住所マスタの投入・更新（初回・継続共通）：

```bash
cd backend
make region-data-migrate
```

このコマンドは以下を自動実行します：
- Excel→CSV変換（一時ファイル）
- 差分検出（初回は全データ、2回目以降は変更分のみ）
- マイグレーションSQL生成（変更がある場合のみ）
- 一時CSVファイルの自動削除

**重要**: 
- 変更がない場合は「No changes detected. Migration files not generated.」と表示され、マイグレーションファイルは生成されません
- マイグレーションSQL生成後、手動で`make migrate-up`を実行してデータベースに反映してください

### 段階的実行

変換のみ実行する場合：

```bash
cd backend
make region-data-convert  # CSVファイル変換のみ
```

マイグレーションをデータベースに反映する場合：

```bash
cd backend
# DATABASE_URLを設定
export DATABASE_URL="postgres://user:password@localhost:5432/crewee_dev"
# マイグレーション実行
make migrate-up
```

### 手動実行

Python環境での個別実行：

```bash
# 仮想環境のセットアップ
python3 -m venv venv
source venv/bin/activate
pip install pandas openpyxl xlrd

# Excel→CSV変換（バックアップ付き）
python3 scripts/convert_lgcode_to_csv.py

# マイグレーションSQL生成（自動差分検出）
python3 scripts/generate_region_diff_migration.py

# マイグレーション実行
cd backend
migrate -path migrations -database "${DATABASE_URL}" up
```

## データ構造

### 元データ（Excel）

**シート: R6.1.1現在の団体**
- 団体コード（6桁JISコード）
- 都道府県名（漢字）
- 市区町村名（漢字）
- 都道府県名（カナ）
- 市区町村名（カナ）

### 変換後データベース構造

```sql
CREATE TABLE regions (
    jis_code VARCHAR(6) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    name_kana VARCHAR(100) NOT NULL,
    region_type region_type_enum NOT NULL,
    parent_jis_code VARCHAR(6) REFERENCES regions(jis_code),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

- **都道府県**: region_type='prefecture', parent_jis_code=NULL
- **市区町村**: region_type='municipality', parent_jis_code=都道府県のJISコード

## 処理内容

### 1. Excel→CSV変換
- Excel形式をUTF-8 CSVに変換（一時ファイル）
- 列名の正規化（改行文字削除）
- JISコードの6桁ゼロパディング
- 空文字の統一処理

### 2. 差分検出システム
- **ハッシュベース比較**: レコード内容をMD5ハッシュで比較
- **新規追加**: 新しいJISコードのレコード
- **更新**: 既存JISコードで内容が変更されたレコード
- **削除**: 新CSVに存在しないJISコードのレコード

### 3. 差分マイグレーションSQL生成
- **新規追加**: INSERT文で新規レコード投入
- **更新**: UPDATE文で既存レコード更新
- **削除**: 物理削除（DELETE文で完全削除）
- **ロールバック**: DOWN SQLで変更を完全に戻す

### 4. ファイル管理とクリーンアップ
- **一時CSVファイル**: 処理完了後に自動削除
- **バックアップファイル**: 次回の差分検出用に保持
- **マイグレーションファイル**: 永続的に保持（データベース履歴管理）

### 5. データ品質保証
- SQLエスケープ処理
- トランザクション処理（BEGIN/COMMIT）
- インデックス最適化
- ANALYZE実行（統計情報更新）

## 注意事項

- **DATABASE_URL環境変数が必要**: マイグレーション実行時
- **既存データの削除**: DOWN マイグレーションは全地域データを削除
- **依存関係**: regionsテーブルが事前に作成されている必要がある
- **文字エンコーディング**: 全てUTF-8で統一

## トラブルシューティング

### よくあるエラー

1. **pandas not found**
   ```bash
   # 仮想環境の再作成
   rm -rf venv
   python3 -m venv venv
   source venv/bin/activate
   pip install pandas openpyxl xlrd
   ```

2. **DATABASE_URLが未設定**
   ```bash
   # 例：PostgreSQLの場合
   export DATABASE_URL="postgres://user:password@localhost:5432/crewee_dev"
   
   # .envファイルでも設定可能
   echo 'DATABASE_URL="postgres://user:password@localhost:5432/crewee_dev"' >> .env
   ```

3. **migrationファイルが重複**
   ```bash
   # 古いマイグレーションファイルを削除してから再生成
   rm backend/migrations/*005_insert_regions_master_data.*
   ```

## データ更新

地方公共団体コードが更新された場合：

### 新しいデータの取得
1. **総務省サイトから最新ファイルをダウンロード**
   - URL: https://www.soumu.go.jp/denshijiti/code.html
   - 「地方公共団体コード住所.xlsx」をダウンロード
2. **ファイル名を変更**
   - `地方公共団体コード住所.xlsx` → `local_government_codes.xlsx`
3. **プロジェクトルートに配置**
   - 既存の`local_government_codes.xlsx`を上書き

### マイグレーション実行
1. `make region-data-migrate` で差分検出・マイグレーション生成
2. 新しいマイグレーションファイルが生成される
3. `make migrate-up` でデータベースに反映

## パフォーマンス

- **初回投入時間**: 約10-30秒（約1,800レコード）
- **インデックス作成**: 並行作成（CONCURRENTLY）でダウンタイム最小化
- **メモリ使用量**: Python処理時に約50MB
