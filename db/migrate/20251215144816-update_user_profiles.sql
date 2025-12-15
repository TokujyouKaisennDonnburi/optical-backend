
-- +migrate Up
CREATE TABLE avatars(
    id UUID PRIMARY KEY,
    url TEXT NOT NULL
);

ALTER TABLE user_profiles DROP COLUMN image_url;
ALTER TABLE user_profiles ADD COLUMN avatar_id UUID REFERENCES avatars(id);

-- +migrate Down

ALTER TABLE user_profiles DROP COLUMN avatar_id;
ALTER TABLE user_profiles ADD COLUMN image_url TEXT NOT NULL;

DROP TABLE avatars;
