
-- +migrate Up

-- for GitHub SSO
CREATE TABLE user_githubs(
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    github_id VARCHAR(255) NOT NULL,
    github_name VARCHAR(255) NOT NULL,
    github_email VARCHAR(255) NOT NULL,
    sso_login BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- for GitHub Apps
CREATE TABLE calendar_githubs(
    calendar_id UUID PRIMARY KEY REFERENCES calendars(id) ON DELETE CASCADE,
    github_id VARCHAR(255) NOT NULL,
    github_name VARCHAR(255) NOT NULL,
    installation_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE user_githubs;
DROP TABLE calendar_githubs;
