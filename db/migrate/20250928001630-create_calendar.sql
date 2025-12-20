-- +migrate Up
CREATE TABLE calendars(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMPTZ NULL   -- NULL許容
);

COMMENT ON COLUMN calendars.id IS 'カレンダーID';
COMMENT ON COLUMN calendars.name IS 'カレンダー名';
COMMENT ON COLUMN calendars.color IS 'カレンダーカラー';
COMMENT ON COLUMN calendars.deleted_at IS '削除日時';

-- +migrate Down
DROP TABLE IF EXISTS calendars; -- テーブルが存在する場合のみ削除する
