-- +migrate Up
CREATE TABLE events(
    id UUID PRIMARY KEY ,
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    memo TEXT NOT NULL,
    color VARCHAR(50) NOT NULL,
    all_day BOOLEAN NOT NULL,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL   -- NULL許容
);

COMMENT ON COLUMN events.id IS 'イベントID';
COMMENT ON COLUMN events.calendar_id IS 'カレンダーID';
COMMENT ON COLUMN events.title IS 'タイトル';
COMMENT ON COLUMN events.memo IS 'メモ';
COMMENT ON COLUMN events.color IS 'カラー';
COMMENT ON COLUMN events.all_day IS '終日フラグ';
COMMENT ON COLUMN events.start_at IS '開始時刻';
COMMENT ON COLUMN events.end_at IS '終了時刻';
COMMENT ON COLUMN events.created_at IS '作成日時';
COMMENT ON COLUMN events.updated_at IS '更新日時';
COMMENT ON COLUMN events.deleted_at IS '削除日時';



-- +migrate Down
DROP TABLE IF EXISTS events; -- テーブルが存在する場合のみ削除する
