-- +migrate Up
CREATE TABLE schedules(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP NULL   -- NULL許容
);

COMMENT ON COLUMN schedules.id IS 'スケジュールID';
COMMENT ON COLUMN schedules.name IS 'スケジュール名';
COMMENT ON COLUMN schedules.deleted_at IS '削除日時';

-- +migrate Down
DROP TABLE IF EXISTS schedules; -- テーブルが存在する場合のみ削除する
