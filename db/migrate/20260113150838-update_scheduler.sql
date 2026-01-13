
-- +migrate Up
ALTER TABLE scheduler ADD COLUMN is_done BOOLEAN NOT NULL DEFAULT FALSE;
COMMENT ON COLUMN scheduler.is_done IS '完了チェック';


-- +migrate Down
ALTER TABLE scheduler DROP COLUMN is_done;

