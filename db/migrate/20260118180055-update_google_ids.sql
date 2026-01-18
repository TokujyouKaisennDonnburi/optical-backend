
-- +migrate Up
DROP INDEX IF EXISTS idx_google_ids_user_id;
DROP INDEX IF EXISTS idx_google_ids_google_id;
DROP TABLE IF EXISTS google_ids;

CREATE TABLE google_ids(
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    google_id VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN google_ids.user_id IS 'ユーザーID';
COMMENT ON COLUMN google_ids.google_id IS 'GoogleアカウントID';
COMMENT ON COLUMN google_ids.created_at IS '作成日時';
COMMENT ON COLUMN google_ids.updated_at IS '更新日時';

CREATE INDEX idx_google_ids_google_id ON google_ids(google_id);

-- +migrate Down
DROP TABLE IF EXISTS google_ids;

CREATE TABLE google_ids(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    google_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN google_ids.id IS 'グーグルID（主キー）';
COMMENT ON COLUMN google_ids.user_id IS 'ユーザーID（外部キー）';
COMMENT ON COLUMN google_ids.google_id IS 'GoogleアカウントのユニークID';
COMMENT ON COLUMN google_ids.created_at IS '作成日時';
COMMENT ON COLUMN google_ids.updated_at IS '更新日時';

CREATE UNIQUE INDEX idx_google_ids_user_id ON google_ids(user_id);
CREATE INDEX idx_google_ids_google_id ON google_ids(google_id);
