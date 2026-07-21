# デプロイ用ドキュメント

本番環境の構成と、新規にデプロイし直す際の手順をまとめる。

## 構成

| 要素 | 使用サービス | 備考 |
| --- | --- | --- |
| APIサーバー | Render (Web Service) | Dockerfile からビルド |
| PostgreSQL | Neon (ap-southeast-1) | `pg_search` 拡張が必要 |
| Redis | Render (Key Value) | キャッシュ用途のみ |
| 画像ストレージ | Cloudflare R2 | S3互換。環境変数名は `MINIO_*` を流用 |
| メール送信 | Resend | HTTP API |
| DNS | Cloudflare | `api.opti-cal.org` |

構成は [`render.yaml`](../render.yaml) で管理する。ダッシュボードで直接設定せず、可能な限りこのファイルに寄せること。

## 構成上の制約

デプロイ先を変更する際は、以下の制約を満たす必要がある。

### PostgreSQL は `pg_search` が使えること

`db/migrate/20260114001101-add_pg_search.sql` が `CREATE EXTENSION pg_search` を実行し、
`events` と `event_locations` に BM25 インデックスを作成する。
イベント検索 (`internal/calendar/gateway/event_search.go`) が `paradedb.parse()` を直接呼ぶため、
素の PostgreSQL では検索エンドポイントが動作しない。

`RUNTIME_MIGRATION=1` の場合、拡張が無いDBに接続すると起動時のマイグレーションが失敗し、
`cmd/api/main.go` の `panic` でプロセスごと停止する。

Render の Managed PostgreSQL は拡張を任意に追加できないため使用できない。

### メール送信に SMTP を使わないこと

Render の無料プランは outbound SMTP (25/465/587) を遮断している。
このため SMTP ではなく Resend の HTTP API を使用する (#341)。
他の PaaS でも同様の制限があることが多い。

### 画像ストレージを変更しないこと

`avatars.is_relative_path` が `true` のレコードは、`IMAGE_STORAGE_BASE_URL` を基準とした
相対パスで画像を参照している。このためベースURLを変更すると既存の画像が表示されなくなる。

## 新規デプロイ手順

### 1. Resend の準備

1. アカウントを作成する
2. Domains から `opti-cal.org` を追加する
3. 表示された DKIM / SPF レコードを Cloudflare の DNS に登録する
4. Verify が完了することを確認する
5. API キーを発行する (権限は Sending access で足りる)

APIキーの発行自体はドメイン認証を待たずに行える。
ただし認証が完了するまで、実際の宛先にメールを送ることはできない。

### 2. Render のサービス作成

1. Blueprint から `render.yaml` を読み込む
2. `sync: false` の環境変数を入力する (後述)
3. Secret File を登録する
   - ファイル名: `optical-github.private-key.pem`
   - 内容: GitHub App の秘密鍵
   - 実行時に `/etc/secrets/` 配下へ配置される

### 3. 環境変数の設定

`render.yaml` に値を書いていない変数は、Blueprint 作成時に入力を求められる。

PostgreSQL は接続文字列ではなく個別の変数で受け取る実装のため、
Neon の接続情報を分解して設定する。

```
POSTGRES_HOST=<プロジェクト固有のホスト名>.ap-southeast-1.aws.neon.tech
POSTGRES_PORT=5432
POSTGRES_DB=neondb
POSTGRES_USER=neondb_owner
POSTGRES_PASSWORD=<Neonのダッシュボードから取得>
POSTGRES_SSLMODE=require
```

Neon が提示する接続文字列には `channel_binding=require` が含まれることがあるが、
`lib/pq` はこのパラメータを解釈できずサーバーへそのまま送るため、接続が拒否される。
接続文字列を直接使う場合は取り除くこと。

Redis も同様に、Key Value の接続文字列から `REDIS_ADDRESS` (host:port) と
`REDIS_PASSWORD` を分解して設定する。

### 4. 切り替え

1. デプロイ完了後、`/options` が 401 を返すことを確認する (500 や 502 でないこと)
2. GitHub App の Callback URL を更新する
3. Google OAuth のリダイレクトURIを更新する
4. `cmd/api/main.go` の CORS 許可オリジンを確認する
5. Cloudflare で `api.opti-cal.org` を新しいサービスに向ける
6. 旧サービスは切り戻し先としてしばらく残す

## 暗号鍵の扱い

| 鍵 | 暗号化対象の保存先 | 新規発行 |
| --- | --- | --- |
| `INSTALLATION_ID_ENCRYPTION_KEY` | PostgreSQL (永続) | 不可 |
| `JWT_SECRET_KEY` | 保存されない (署名のみ) | 可 |
| `REDIS_ENCRYPTION_KEY` | Redis (TTL付き) | 可 |

`INSTALLATION_ID_ENCRYPTION_KEY` は `calendar_githubs.installation_id` の暗号化に使う
(`internal/github/gateway/install_calendar.go`)。
既存のDBを引き継ぐ場合、鍵を変更すると復号できなくなる。

鍵を引き継げない場合は、新しい鍵を発行した上で対象カレンダーの GitHub App を再連携する。
`INSERT` は `installation_id` を上書きする実装のため、再連携すれば新しい鍵の暗号文に置き換わる。

`JWT_SECRET_KEY` を変更すると全ユーザーが再ログインを求められるが、データは壊れない。
`REDIS_ENCRYPTION_KEY` はキャッシュミス時に GitHub API から取り直すフォールバックがあるため、
変更しても実質的な影響はない。

## 無料プランの制約

- **スピンダウン** — 15分間アクセスがないとインスタンスが停止し、次回リクエストでコールドスタートする。
  メール送信は goroutine による fire-and-forget のため (`internal/calendar/service/command/invite_members.go`)、
  停止のタイミング次第で送信されないことがある
- **永続ディスクなし** — Web Service のファイルシステムは揮発する
- **メンバー招待不可** — Hobby プランのワークスペースには1人しか所属できない。
  複数人で運用する場合は Pro 以上が必要
- **Resend の送信上限** — 無料プランは月3,000通、日100通
