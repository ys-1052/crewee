---
allowed-tools: Bash, Read, Glob, Grep
description: 'Git commit workflow following Conventional Commits specification'
---

このコマンドはConventional Commitsの仕様に従ってGitコミットを作成します。

## 実行内容

1. **現在の状態を確認**
   - `git status` でステージング状況を確認
   - `git diff --cached` でステージされた変更内容を確認
   - 最近のコミット履歴を確認してコミットメッセージのスタイルを把握

2. **変更内容を分析してコミットメッセージを作成**
   - ステージされた変更を分析
   - Conventional Commitsの仕様に従ってメッセージを作成
   - 変更の種類（feat, fix, docs, etc.）を自動判定

3. **コミットを実行**
   - 適切なコミットメッセージでコミットを作成
   - Claude Codeの署名を追加

4. **リモートにプッシュ**
   - 現在のブランチをリモートリポジトリにプッシュ
   - 新しいブランチの場合は上流ブランチを設定

## 前提条件

- **重要**: コミットしたいファイルを事前に`git add`でステージングエリアに追加しておく必要があります
- **注意**: Claudeは自動的に`git add`を実行しません。ユーザーが明示的にファイルをステージングすることが必要です
- ステージされたファイルがない場合は何もしません
- Claudeはステージングエリアにある変更のみを対象としてコミットを作成します
- 例: `git add .` または `git add <ファイル名>` でファイルをステージング

## 参考: Conventional Commits 仕様

このコマンドでは以下の仕様に従ってコミットメッセージを自動生成します：

### 基本形式
```
<type>(<scope>): <description>
```

### Type（変更の種類）
| Type | 説明 | 例 |
|------|------|-----|
| `feat` | 新機能の追加 | `feat(auth): add user login`
| `fix` | バグの修正 | `fix(api): handle null response`
| `docs` | ドキュメントの変更 | `docs: update installation guide`
| `style` | コード整形（動作に影響なし） | `style: format code with prettier`
| `refactor` | リファクタリング | `refactor(utils): simplify validation`
| `test` | テストの追加・修正 | `test: add unit tests for auth`
| `chore` | ビルド・ツール関連 | `chore: update dependencies`

### Scope（影響範囲）
- 変更が影響するコンポーネントやファイル
- 例: `auth`, `api`, `utils`, `frontend`, `backend`

### Description（説明）
- 英語で簡潔に記述
- 命令形を使用（"add", "fix", "update"など）
- 小文字で始まり、末尾にピリオドなし
