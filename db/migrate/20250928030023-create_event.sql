-- +migrate Up
CREATE TABLE event(
    id UUID PRIMARY KEY ,
    schedule_id UUID NOT NULL REFERENCES schedules(id) ON DELETE CASCADE,
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

COMMENT ON COLUMN event.id IS 'イベントID';
COMMENT ON COLUMN event.schedule_id IS 'スケジュールID';
COMMENT ON COLUMN event.title IS 'タイトル';
COMMENT ON COLUMN event.memo IS 'メモ';
COMMENT ON COLUMN event.color IS 'カラー';
COMMENT ON COLUMN event.all_day IS '終日フラグ';
COMMENT ON COLUMN event.start_at IS '開始時刻';
COMMENT ON COLUMN event.end_at IS '終了時刻';
COMMENT ON COLUMN event.created_at IS '作成日時';
COMMENT ON COLUMN event.updated_at IS '更新日時';
COMMENT ON COLUMN event.deleted_at IS '削除日時';



-- +migrate Down
DROP TRIGGER IF EXISTS update_event_updated_at ON event; -- トリガーも削除
DROP TABLE IF EXISTS event; -- テーブルが存在する場合のみ削除する
