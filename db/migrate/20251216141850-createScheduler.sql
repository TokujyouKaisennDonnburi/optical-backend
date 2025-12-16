-- +migrate Up
CREATE TABLE scheduler(
    id UUID PRIMARY KEY,
    calendar_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    memo VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NULL,
    end_time TIMESTAMP NULL,
    is_allday BOOLEAN NOT NULL,
    FOREIGN KEY (calendar_id) REFERENCES calendars(id) ON DELETE CASCADE
);

CREATE TABLE scheduler_attendance(
    id UUID PRIMARY KEY,
    scheduler_id UUID NOT NULL,
    user_id UUID NOT NULL,
    comment VARCHAR(255) NOT NULL,
    FOREIGN KEY (scheduler_id) REFERENCES scheduler(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE scheduler_status(
    attendance_id UUID NOT NULL,
    date DATE NOT NULL,
    status INTEGER,
    PRIMARY KEY (attendance_id, date),
    FOREIGN KEY (attendance_id) REFERENCES scheduler_attendance(id) ON DELETE CASCADE
);

-- scheduler
COMMENT ON COLUMN scheduler.id IS 'スケジューラーID';
COMMENT ON COLUMN scheduler.calendar_id IS 'カレンダーID';
COMMENT ON COLUMN scheduler.title IS 'スケジューラータイトル';
COMMENT ON COLUMN scheduler.memo IS 'メモ';
COMMENT ON COLUMN scheduler.start_time IS '開始時間';
COMMENT ON COLUMN scheduler.end_time IS '終了時間';
COMMENT ON COLUMN scheduler.is_allday IS '終日チェック';
-- attendance
COMMENT ON COLUMN scheduler_attendance.id IS '調整ID';
COMMENT ON COLUMN scheduler_attendance.scheduler_id IS 'スケジューラーID';
COMMENT ON COLUMN scheduler_attendance.user_id IS 'ユーザーID';
COMMENT ON COLUMN scheduler_attendance.comment IS 'コメント';
-- status
COMMENT ON COLUMN scheduler_status.attendance_id IS '調整ID';
COMMENT ON COLUMN scheduler_status.date IS '日付';
COMMENT ON COLUMN scheduler_status.status IS '参加確認';

-- +migrate Down

DROP TABLE IF EXISTS scheduler_status;
DROP TABLE IF EXISTS scheduler_attendance;
DROP TABLE IF EXISTS scheduler;
