-- user系 プロフィール

-- +migrate Up

CREATE TABLE user_profiles(
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,    -- usersテーブルのidを参照と削除動作
    image_url TEXT
);

-- テーブルとカラムにコメント追加
COMMENT ON TABLE user_profiles IS 'ユーザープロフィールテーブル';
COMMENT ON COLUMN user_profiles.user_id IS 'ユーザーID（主キー、外部キー）';
COMMENT ON COLUMN user_profiles.image_url IS 'プロフィール画像のURL';

-- +migrate Down
DROP TABLE IF EXISTS user_profiles; -- テーブルが存在する場合のみ削除する