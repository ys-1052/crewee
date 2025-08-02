#!/usr/bin/env python3
"""
地方公共団体コード.xlsxファイルをUTF-8 CSVに変換するスクリプト
都道府県・市区町村データのみを対象とする
"""

import os
import sys
from pathlib import Path

import pandas as pd


def convert_excel_to_csv():
    """ExcelファイルをCSVに変換"""
    excel_file = "local_government_codes.xlsx"
    output_dir = Path("backend/sql")

    # 出力ディレクトリを作成
    output_dir.mkdir(parents=True, exist_ok=True)

    if not os.path.exists(excel_file):
        print(f"エラー: {excel_file} が見つかりません")
        sys.exit(1)

    try:
        # Excelファイルを読み込み
        df_dict = pd.read_excel(excel_file, sheet_name=None)

        # 通常の地方公共団体データを処理
        if "R6.1.1現在の団体" in df_dict:
            df_main = df_dict["R6.1.1現在の団体"]

            # 列名をクリーンアップ
            df_main.columns = df_main.columns.str.replace("\n", "")

            # NaNを空文字に置換
            df_main = df_main.fillna("")

            # 団体コードを6桁に正規化（必要に応じて）
            df_main["団体コード"] = df_main["団体コード"].astype(str).str.zfill(6)

            # UTF-8でCSV出力
            output_path = output_dir / "local_government_codes.csv"
            df_main.to_csv(output_path, index=False, encoding="utf-8")
            print(
                f"✓ 通常の地方公共団体データを {output_path} に出力しました ({len(df_main)} 行)"
            )

        print("変換完了!")

    except Exception as e:
        print(f"エラー: {e}")
        sys.exit(1)


if __name__ == "__main__":
    convert_excel_to_csv()
