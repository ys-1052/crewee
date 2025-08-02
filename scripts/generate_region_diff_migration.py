#!/usr/bin/env python3
"""
CSVファイルの差分からマイグレーション用SQLファイルを生成するスクリプト
変更された地域データのみをマイグレーション対象とする
"""

import hashlib
import re
import sys
from datetime import datetime
from pathlib import Path
from typing import Dict, List

import pandas as pd


def get_next_migration_number():
    """既存のマイグレーションファイルから次の番号を取得"""
    migration_dir = Path("backend/migrations")
    if not migration_dir.exists():
        return "001"

    max_number = 0
    pattern = re.compile(r"^\d+_(\d{3})_.*\.sql$")

    for file_path in migration_dir.glob("*.sql"):
        match = pattern.match(file_path.name)
        if match:
            number = int(match.group(1))
            max_number = max(max_number, number)

    return f"{max_number + 1:03d}"


def generate_diff_migration_sql():
    """CSVの差分からマイグレーション用SQLを生成"""

    # CSVファイルパス
    main_csv = Path("backend/sql/local_government_codes.csv")
    backup_csv = Path("backend/sql/local_government_codes.csv.backup")

    if not main_csv.exists():
        print(f"エラー: {main_csv} が見つかりません")
        sys.exit(1)

    try:
        # 新しいCSVファイルを読み込み
        df_new = pd.read_csv(main_csv)
        df_new = df_new.fillna("")  # NaNを空文字に統一

        # バックアップファイルが存在する場合は差分を検出
        if backup_csv.exists():
            df_old = pd.read_csv(backup_csv)
            df_old = df_old.fillna("")  # NaNを空文字に統一
            changes = detect_changes(df_old, df_new)
        else:
            # 初回実行時は全データを新規として扱う
            print("初回実行：全データを新規追加として処理します。")
            changes = {
                "added": df_new.to_dict("records"),
                "updated": [],
                "deleted": [],
            }

        # 変更がない場合は処理を終了
        if not any(changes.values()):
            print("変更がありません。マイグレーションファイルは生成されませんでした。")
            return

        # マイグレーションファイルのタイムスタンプ（UTC）
        timestamp = datetime.utcnow().strftime("%Y%m%d%H%M%S")

        # マイグレーションファイルのパス
        migration_dir = Path("backend/migrations")
        migration_number = get_next_migration_number()
        migration_up = migration_dir / (
            f"{timestamp}_{migration_number}_update_regions_data.up.sql"
        )
        migration_down = migration_dir / (
            f"{timestamp}_{migration_number}_update_regions_data.down.sql"
        )

        # UP migration SQL生成
        up_sql = generate_diff_up_sql(changes)

        # DOWN migration SQL生成
        down_sql = generate_diff_down_sql(
            changes, df_old if backup_csv.exists() else pd.DataFrame()
        )

        # ファイルに出力
        with open(migration_up, "w", encoding="utf-8") as f:
            f.write(up_sql)

        with open(migration_down, "w", encoding="utf-8") as f:
            f.write(down_sql)

        # 統計情報を表示
        print("✓ 差分マイグレーションファイルを生成しました:")
        print(f"  UP:   {migration_up}")
        print(f"  DOWN: {migration_down}")
        print("✓ 変更統計:")
        print(f"  新規追加: {len(changes['added'])} 件")
        print(f"  更新: {len(changes['updated'])} 件")
        print(f"  削除: {len(changes['deleted'])} 件")

        # バックアップを作成
        create_backup(main_csv, backup_csv)

    except Exception as e:
        print(f"エラー: {e}")
        sys.exit(1)


def detect_changes(df_old: pd.DataFrame, df_new: pd.DataFrame) -> Dict[str, List]:
    """新旧CSVの差分を検出"""

    # データフレームをレコード辞書のリストに変換
    old_records = {row["団体コード"]: row for row in df_old.to_dict("records")}
    new_records = {row["団体コード"]: row for row in df_new.to_dict("records")}

    old_codes = set(old_records.keys())
    new_codes = set(new_records.keys())

    # 新規追加されたレコード
    added_codes = new_codes - old_codes
    added = [new_records[code] for code in added_codes]

    # 削除されたレコード
    deleted_codes = old_codes - new_codes
    deleted = [old_records[code] for code in deleted_codes]

    # 更新されたレコード（内容が変わったもの）
    common_codes = old_codes & new_codes
    updated = []
    for code in common_codes:
        if record_hash(old_records[code]) != record_hash(new_records[code]):
            updated.append(new_records[code])

    return {"added": added, "updated": updated, "deleted": deleted}


def record_hash(record: Dict) -> str:
    """レコードのハッシュ値を計算（変更検出用）"""
    # 団体コード以外のフィールドでハッシュを計算
    content = (
        f"{record.get('都道府県名（漢字）', '')}"
        f"{record.get('市区町村名（漢字）', '')}"
        f"{record.get('都道府県名（カナ）', '')}"
        f"{record.get('市区町村名（カナ）', '')}"
    )
    return hashlib.md5(content.encode()).hexdigest()


def generate_diff_up_sql(changes: Dict[str, List]) -> str:
    """差分更新用のUP SQLを生成"""
    sql_lines = [
        "-- Update regions data (differential migration)",
        "-- Generated automatically from local government codes diff",
        "",
        "BEGIN;",
        "",
    ]

    # 新規追加
    if changes["added"]:
        sql_lines.extend(
            [
                f"-- Insert new regions ({len(changes['added'])} records)",
                "INSERT INTO regions (jis_code, name, name_kana, region_type, "
                "parent_jis_code, created_at, updated_at) VALUES",
            ]
        )

        insert_values = []
        for record in changes["added"]:
            jis_code = str(record["団体コード"]).zfill(6)
            name = str(record.get("都道府県名（漢字）", "")).replace("'", "''")
            name_kana = str(record.get("都道府県名（カナ）", "")).replace("'", "''")

            # 市区町村名がある場合は市区町村、ない場合は都道府県
            if record.get("市区町村名（漢字）", ""):
                name = str(record["市区町村名（漢字）"]).replace("'", "''")
                name_kana = str(record["市区町村名（カナ）"]).replace("'", "''")
                region_type = "municipality"
                parent_jis_code = find_prefecture_jis_code_from_record(
                    jis_code, changes["added"]
                )
                parent_clause = (
                    f"'{parent_jis_code}'" if parent_jis_code != "NULL" else "NULL"
                )
            else:
                region_type = "prefecture"
                parent_clause = "NULL"

            insert_values.append(
                f"('{jis_code}', '{name}', '{name_kana}', '{region_type}', "
                f"{parent_clause}, NOW(), NOW())"
            )

        sql_lines.append(",\n".join(insert_values) + ";")
        sql_lines.append("")

    # 更新
    if changes["updated"]:
        sql_lines.append(
            f"-- Update existing regions ({len(changes['updated'])} records)"
        )

        for record in changes["updated"]:
            jis_code = str(record["団体コード"]).zfill(6)

            if record.get("市区町村名（漢字）", ""):
                name = str(record["市区町村名（漢字）"]).replace("'", "''")
                name_kana = str(record["市区町村名（カナ）"]).replace("'", "''")
            else:
                name = str(record.get("都道府県名（漢字）", "")).replace("'", "''")
                name_kana = str(record.get("都道府県名（カナ）", "")).replace("'", "''")

            sql_lines.append(
                f"UPDATE regions SET "
                f"name = '{name}', "
                f"name_kana = '{name_kana}', "
                f"updated_at = NOW() "
                f"WHERE jis_code = '{jis_code}';"
            )

        sql_lines.append("")

    # 削除（物理削除）
    if changes["deleted"]:
        sql_lines.extend(
            [
                f"-- Delete removed regions " f"({len(changes['deleted'])} records)",
                "-- Note: Physical deletion for master data",
            ]
        )

        deleted_codes = [
            f"'{str(record['団体コード']).zfill(6)}'" for record in changes["deleted"]
        ]

        sql_lines.append(
            f"DELETE FROM regions " f"WHERE jis_code IN ({', '.join(deleted_codes)});"
        )
        sql_lines.append("")

    sql_lines.extend(
        [
            "COMMIT;",
            "",
            "-- Update statistics",
            "ANALYZE regions;",
        ]
    )

    return "\n".join(sql_lines)


def generate_diff_down_sql(changes: Dict[str, List], df_old: pd.DataFrame) -> str:
    """差分更新用のDOWN SQLを生成"""
    sql_lines = [
        "-- Rollback regions data changes",
        "-- This will undo the differential migration",
        "",
        "BEGIN;",
        "",
    ]

    # 追加されたレコードを削除
    if changes["added"]:
        added_codes = [
            f"'{str(record['団体コード']).zfill(6)}'" for record in changes["added"]
        ]
        sql_lines.extend(
            [
                f"-- Remove newly added regions ({len(changes['added'])} records)",
                f"DELETE FROM regions WHERE jis_code IN ({', '.join(added_codes)});",
                "",
            ]
        )

    # 更新されたレコードを元に戻す
    if changes["updated"] and not df_old.empty:
        sql_lines.append(
            f"-- Restore updated regions ({len(changes['updated'])} records)"
        )

        old_records = {
            str(row["団体コード"]).zfill(6): row for row in df_old.to_dict("records")
        }

        for record in changes["updated"]:
            jis_code = str(record["団体コード"]).zfill(6)
            if jis_code in old_records:
                old_record = old_records[jis_code]

                if old_record.get("市区町村名（漢字）", ""):
                    name = str(old_record["市区町村名（漢字）"]).replace("'", "''")
                    name_kana = str(old_record["市区町村名（カナ）"]).replace("'", "''")
                else:
                    name = str(old_record.get("都道府県名（漢字）", "")).replace(
                        "'", "''"
                    )
                    name_kana = str(old_record.get("都道府県名（カナ）", "")).replace(
                        "'", "''"
                    )

                sql_lines.append(
                    f"UPDATE regions SET "
                    f"name = '{name}', "
                    f"name_kana = '{name_kana}', "
                    f"updated_at = NOW() "
                    f"WHERE jis_code = '{jis_code}';"
                )

        sql_lines.append("")

    # 削除されたレコードを復元（再挿入）
    if changes["deleted"]:
        sql_lines.append(
            f"-- Restore deleted regions ({len(changes['deleted'])} records)"
        )

        insert_values = []
        for record in changes["deleted"]:
            jis_code = str(record["団体コード"]).zfill(6)
            name = str(record.get("都道府県名（漢字）", "")).replace("'", "''")
            name_kana = str(record.get("都道府県名（カナ）", "")).replace("'", "''")

            # 市区町村名がある場合は市区町村、ない場合は都道府県
            if record.get("市区町村名（漢字）", ""):
                name = str(record["市区町村名（漢字）"]).replace("'", "''")
                name_kana = str(record["市区町村名（カナ）"]).replace("'", "''")
                region_type = "municipality"
                parent_jis_code = find_prefecture_jis_code_from_record(
                    jis_code, changes["deleted"]
                )
                parent_clause = (
                    f"'{parent_jis_code}'" if parent_jis_code != "NULL" else "NULL"
                )
            else:
                region_type = "prefecture"
                parent_clause = "NULL"

            insert_values.append(
                f"('{jis_code}', '{name}', '{name_kana}', '{region_type}', "
                f"{parent_clause}, NOW(), NOW())"
            )

        sql_lines.extend(
            [
                "INSERT INTO regions (jis_code, name, name_kana, region_type, "
                "parent_jis_code, created_at, updated_at) VALUES",
                ",\n".join(insert_values) + ";",
                "",
            ]
        )

    sql_lines.extend(
        [
            "COMMIT;",
            "",
            "-- Update statistics",
            "ANALYZE regions;",
        ]
    )

    return "\n".join(sql_lines)


def find_prefecture_jis_code_from_record(
    municipality_jis_code: str, records: List[Dict]
) -> str:
    """レコードリストから都道府県のJISコードを取得"""
    prefecture_code = municipality_jis_code[:2]

    for record in records:
        jis_code = str(record["団体コード"]).zfill(6)
        # 都道府県レコード（市区町村名が空）で、コードが一致するもの
        if not record.get("市区町村名（漢字）", "") and jis_code.startswith(
            prefecture_code
        ):
            return jis_code

    return "NULL"


def create_backup(source_path: Path, backup_path: Path):
    """CSVファイルのバックアップを作成"""
    try:
        import shutil

        shutil.copy2(source_path, backup_path)
        print(f"✓ バックアップを作成しました: {backup_path}")
    except Exception as e:
        print(f"警告: バックアップの作成に失敗しました: {e}")


if __name__ == "__main__":
    generate_diff_migration_sql()
