
-- +migrate Up
CREATE TABLE calendar_images(
    id UUID PRIMARY KEY,
    url TEXT NOT NULL
);

ALTER TABLE calendars ADD COLUMN image_id UUID REFERENCES calendar_images(id) ON DELETE SET NULL;
COMMENT ON COLUMN calendars.image_id IS '画像ID';

COMMENT ON TABLE calendar_images IS '画像';
COMMENT ON COLUMN calendar_images.id IS '画像ID';
COMMENT ON COLUMN calendar_images.url IS '画像URL';

-- +migrate Down
ALTER TABLE DROP COLUMN image_id;
DROP TABLE IF EXISTS image_url;
