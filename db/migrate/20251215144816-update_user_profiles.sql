
-- +migrate Up
CREATE TABLE avatars(
    id UUID PRIMARY KEY,
    url TEXT NOT NULL
);

ALTER TABLE user_profiles ALTER COLUMN image_url TYPE UUID DEFAULT NULL REFERENCES avatars(id);
ALTER TABLE user_profiles RENAME COLUMN image_url TO avatar_id; 

-- +migrate Down

ALTER TABLE user_profiles RENAME COLUMN image_id TO image_url 
ALTER TABLE user_profiles ALTER COLUMN image_url TYPE TEXT NOT NULL;

DROP TABLE user_avatars;

