-- user系　googleids

-- +migrate Up

CREATE TABLE google_ids(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,   -- 外部キー制約とその削除動作
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
)

-- テーブルとカラムにコメント追加
COMMENT ON TABLE google_ids IS 'Google認証連携テーブル';
COMMENT ON COLUMN google_ids.id IS 'グーグルID（主キー）';
COMMENT ON COLUMN google_ids.user_id IS 'ユーザーID（外部キー）';
COMMENT ON COLUMN google_ids.google_id IS 'GoogleアカウントのユニークID';
COMMENT ON COLUMN google_ids.created_at IS '作成日時';
COMMENT ON COLUMN google_ids.updated_at IS '更新日時';

-- updated_atを自動更新するトリガーを作成
CREATE TRIGGER update_google_ids_updated_at
    BEFORE UPDATE ON google_ids
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


-- インデックス作成
CREATE UNIQUE INDEX idx_google_ids_user_id ON google_ids(user_id);
CREATE INDEX idx_google_ids_google_id ON google_ids(google_id); -- 1ユーザーが複数のGoogleIDを持つ場合に備えてインデックスを作成

-- +migrate Down
DROP TABLE IF EXISTS google_ids; -- テーブルが存在する場合のみ削除する
DROP TRIGGER IF EXISTS update_google_ids_updated_at ON google_ids; -- トリガーも削除
DROP INDEX IF EXISTS idx_google_ids_user_id; -- インデックス削除
DROP INDEX IF EXISTS idx_google_ids_google_id; -- インデックス削除