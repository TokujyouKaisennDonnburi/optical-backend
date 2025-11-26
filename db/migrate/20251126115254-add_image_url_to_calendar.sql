
-- +migrate Up
CREATE TABLE calendar_images(
    id UUID PRIMARY KEY,
    url TEXT NOT NULL,
);

ALTER TABLE calendars ADD COLUMN image_id UUID REFERENCES calendar_images(id) ON DELETE SET NULL;
COMMENT ON COLUMN calendars.url IS "カレンダー画像URL";

-- +migrate Down
ALTER TABLE DROP COLUMN image_id;
DROP TABLE image_url;
