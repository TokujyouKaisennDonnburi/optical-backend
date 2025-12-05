
-- +migrate Up
CREATE TABLE github_ids(
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    github_id VARCHAR(255) NOT NULL,
    github_name VARCHAR(255) NOT NULL,
    installation_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE github_ids;
