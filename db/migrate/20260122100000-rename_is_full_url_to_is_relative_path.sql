-- +migrate Up
ALTER TABLE avatars RENAME COLUMN is_full_url TO is_relative_path;

-- +migrate Down
ALTER TABLE avatars RENAME COLUMN is_relative_path TO is_full_url;
