-- +migrate Up
ALTER TABLE avatars ADD COLUMN is_full_url BOOLEAN NOT NULL DEFAULT FALSE;

-- +migrate Down
ALTER TABLE avatars DROP COLUMN is_full_url;
