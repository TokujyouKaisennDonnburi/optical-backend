-- +migrate Up
CREATE TABLE calendar_members(
    calendar_id UUID REFERENCES calendars(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ,  -- NULL許容
    PRIMARY KEY (calendar_id, user_id)  -- 複合主キー
);

COMMENT ON COLUMN calendar_members.calendar_id IS 'カレンダーID';
COMMENT ON COLUMN calendar_members.user_id IS 'ユーザーID';
COMMENT ON COLUMN calendar_members.joined_at IS '参加日時';

-- +migrate Down
DROP TABLE IF EXISTS calendar_members; -- テーブルが存在する場合のみ削除する

