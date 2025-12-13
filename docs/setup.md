# Setup

## sql-migrateのインストール

```bash
go install github.com/rubenv/sql-migrate/...@latest
```

## airのインストール

```bash
go install github.com/air-verse/air@latest
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

### Github

- `GITHUB_APP_ID`及び、`GITHUB_CLIENT_ID`には[アプリ設定ページ](https://github.com/organizations/TokujyouKaisennDonnburi/settings/apps/optical-github)から`Client ID`をコピーして設定する
- `GITHUB_CLIENT_SECRET`には`Client secrets`を生成して設定する
- 設定ページの`Private keys`タブからプライベートキーを作成し、`optical-github.private-key.pem`に名前を変更しリポジトリのルートに配置する

## データベース立ち上げ
```bash
docker compose up -d
```

## マイグレーション

```bash
./migrate.sh
```
