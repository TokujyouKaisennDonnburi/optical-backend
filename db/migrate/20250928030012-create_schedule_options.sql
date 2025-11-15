-- +migrate Up
CREATE TABLE calendar_options(
    calendar_id UUID REFERENCES calendars(id) ON DELETE CASCADE,    -- 外部キー制約とカスケード削除
    option_id UUID REFERENCES options(id) ON DELETE CASCADE,
    PRIMARY KEY (calendar_id, option_name)  -- 複合主キー
);

COMMENT ON COLUMN calendar_options.calendar_id IS 'カレンダーID';
COMMENT ON COLUMN calendar_options.option_id IS 'オプション名';

-- +migrate Down
DROP TABLE IF EXISTS calendar_options; -- テーブルが存在する場合のみ削除する
