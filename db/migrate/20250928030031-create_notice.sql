-- +migrate Up
CREATE TABLE notice(
    id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE, -- イベントIDへの外部キー
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,  -- 読んだかどうかのフラグ
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN notice.id IS '通知ID';
COMMENT ON COLUMN notice.event_id IS 'イベントID';
COMMENT ON COLUMN notice.title IS '通知タイトル';
COMMENT ON COLUMN notice.content IS '通知内容';
COMMENT ON COLUMN notice.is_read IS '読んだかどうかのフラグ';
COMMENT ON COLUMN notice.created_at IS '作成日時';


-- +migrate Down
DROP TABLE IF EXISTS notice; -- テーブルが存在する場合のみ削除する
