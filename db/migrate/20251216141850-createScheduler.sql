
-- +migrate Up
CREATE TABLE scheduler(
    id UUID PRIMARY KEY,
    calendar_id UUID,
    title VARCHAR(255) NOT NULL,
    memo VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NULL,
    end_title TIMESTAMP NULL,
    is_allday BOOLEAN NOT NULL,
);

CREATE TABLE scheduler_attendance(
    id UUID PRIMARY KEY,
    scheduler_id UUID PRIMARY KEY,
    user_id UUID,
    comment VARCHAR(255) NOT NULL,
);

CREATE TABLE scheduler_status(
    attendance_id UUID,
    date DATE NULL,
    status INTEGER,
);

-- TODO 外部で使うIDは外部キー指定
-- TODO status attendanceIDとdateで重複しないようにする


COMMENT ON COLUMN scheduler.id IS 'スケジューラーID';
COMMENT ON COLUMN scheduler.calendar_id IS 'カレンダーID';
COMMENT ON COLUMN scheduler.title IS 'スケジューラータイトル';
COMMENT ON COLUMN scheduler.memo IS 'メモ';
COMMENT ON COLUMN scheduler.start_time IS '開始時間';
COMMENT ON COLUMN scheduler.end_time IS '終了時間';
COMMENT ON COLUMN scheduler.is_allday IS '終日チェック';

-- +migrate Down
DROP TABLE IF EXISTS scheduler;
DROP TABLE IF EXISTS scheduler_attendance;
