-- +migrate Up
CREATE TABLE schedule_options(
    schedule_id UUID REFERENCES schedules(id) ON DELETE CASCADE,    -- 外部キー制約とカスケード削除
    option_id UUID REFERENCES options(id) ON DELETE CASCADE,
    PRIMARY KEY (schedule_id, option_name)  -- 複合主キー
);

COMMENT ON COLUMN schedule_options.schedule_id IS 'スケジュールID';
COMMENT ON COLUMN schedule_options.option_id IS 'オプション名';

-- +migrate Down
DROP TABLE IF EXISTS schedule_options; -- テーブルが存在する場合のみ削除する
