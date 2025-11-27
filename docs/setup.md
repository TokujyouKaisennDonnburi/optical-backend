# Setup

## sql-migrateのインストール

```bash
go install github.com/rubenv/sql-migrate/...@latest
```

## MinioMCのインストール

```bash
# macOS
brew install minio/stable/mc

# linux
wget https://dl.min.io/client/mc/release/linux-amd64/mc
chmod +x mc
sudo mv mc /usr/local/bin
```

## 環境変数の設定

```bash
cp .env.example .env
```

## データベース立ち上げ
```bash
docker compose up -d
```

## マイグレーション

※  Windowsの場合は`migrate.bat`を実行する

```bash
./migrate.sh
```

