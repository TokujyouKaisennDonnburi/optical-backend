-- スケジュールに追加するユーザーを管理するテーブル

-- +migrate Up
CREATE TABLE schedule_members(
    schedule_id UUID REFERENCES schedules(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (schedule_id, user_id)  -- 複合主キー
);

-- コメント
COMMENT ON COLUMN schedule_members.schedule_id IS 'スケジュールID';
COMMENT ON COLUMN schedule_members.user_id IS 'ユーザーID';

-- +migrate Down
DROP TABLE IF EXISTS schedule_members; -- テーブルが存在する場合のみ削除する

