-- usersのスキーマを作成するマイグレーションファイル(postgres用)

-- +migrate Up

-- トリガー関数を作成
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

CREATE TABLE users(
    id UUID PRIMARY KEY,    -- UUID 参考はBINARY
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

COMMENT ON COLUMN users.id IS 'ユーザーID';
COMMENT ON COLUMN users.name IS 'ユーザー名';
COMMENT ON COLUMN users.email IS 'メールアドレス';
COMMENT ON COLUMN users.password_hash IS 'パスワードハッシュ';
COMMENT ON COLUMN users.created_at IS '作成日時';
COMMENT ON COLUMN users.updated_at IS '更新日時';

-- updated_atを自動更新するトリガーを作成
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();


-- インデックスの作成(高速化のため)
CREATE INDEX idx_users_email ON users(email);


-- +migrate Down
-- DROP TABLE users;
DROP TABLE IF EXISTS users; -- テーブルが存在する場合のみ削除する
DROP FUNCTION IF EXISTS update_updated_at_column(); -- トリガー関数も削除
DROP TRIGGER IF EXISTS update_users_updated_at ON users; -- トリガーも削除
