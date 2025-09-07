# Calendar Backend

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

- **camelCase** で統一

## コミット規則

- `feat:` 機能追加
- `fix:` バグ修正
- `docs:` ドキュメント修正
- `refactor:` リファクタリング
- `test:` テスト
