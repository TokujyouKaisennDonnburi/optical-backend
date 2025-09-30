-- +migrate Up
CREATE TABLE options(
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

COMMENT ON COLUMN options.id IS 'オプションID';
COMMENT ON COLUMN options.name IS 'オプション名';

-- +migrate Down
DROP TABLE IF EXISTS options; -- テーブルが存在する場合のみ削除する