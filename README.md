# Optical Backend

Goプロジェクトのバックエンドリポジトリです

## セットアップ

### 依存関係のインストール
```bash
go mod download
```

### 環境変数の設定
```bash
cp .env.example .env
```

`.env` ファイルを編集して必要な環境変数を設定してください。

### 実行
```bash
go run cmd/api/main.go
```

## フレームワーク

このプロジェクトでは **chi** を使用しています

## 命名規則

- ファイル名・ディレクトリ名は**snake_case**で統一
- パッケージ外部に公開する変数・構造体・関数などは**PascalCase**で統一
- 上記に当てはまらないものは**camelCase** で統一

## コミット規則

- `feat:` 機能追加
- `update:` 機能変更
- `fix:` バグ修正
- `docs:` ドキュメント修正
- `refactor:` リファクタリング
- `test:` テスト

### Issue

コミットメッセージの最後に、対象の**Issue**番号を記載してください。
記載方法は`#番号`の形式です。

#### 例

```yaml
docs: githubのドキュメント作成 #3
```
