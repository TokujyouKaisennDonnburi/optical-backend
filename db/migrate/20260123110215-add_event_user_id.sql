
-- +migrate Up
ALTER TABLE events ADD COLUMN user_id UUID NOT NULL;
COMMENT ON COLUMN events.user_id IS '作成者ID';

-- +migrate Down
ALTER TABLE events DROP COLUMN user_id;
