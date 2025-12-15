
-- +migrate Up
CREATE TABLE avatars(
    id UUID PRIMARY KEY,
    url TEXT NOT NULL
);

ALTER TABLE user_profiles DROP COLUMN image_url;
ALTER TABLE user_profiles ADD COLUMN avatar_id UUID REFERENCES avatars(id);

-- +migrate Down

ALTER TABLE user_profiles RENAME COLUMN avatar_id TO image_url 
ALTER TABLE user_profiles ALTER COLUMN image_url TYPE TEXT NOT NULL;

DROP TABLE avatars;
