
-- +migrate Up
ALTER TABLE events DROP COLUMN color;

-- +migrate Down
ALTER TABLE events ADD COLUMN color VARCHAR(50) NOT NULL;
COMMENT ON COLUMN events.color IS 'カラー';
