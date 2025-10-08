# Setup

## sql-migrateのインストール

```bash
go install github.com/rubenv/sql-migrate/...@latest
```

## 環境変数の設定

```bash
cp .env.example .env
```

## データベース立ち上げ
```bash
docker compose up -d
```
