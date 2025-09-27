-- 共通使用する関数を記述

-- +migrate Up
-- トリガー関数を作成
-- google_idsとusersで共通
-- updated_atカラムを自動更新するための関数
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';

-- テーブルごとにトリガーを作成する関数
CREATE OR REPLACE FUNCTION create_updated_at_trigger(table_name TEXT)
RETURNS VOID AS $$
BEGIN
    EXECUTE format('
        CREATE TRIGGER update_%I_updated_at
            BEFORE UPDATE ON %I
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    ', table_name, table_name);
END;
$$ LANGUAGE 'plpgsql';

-- +migrate Down
DROP FUNCTION IF EXISTS create_updated_at_trigger(TEXT);
DROP FUNCTION IF EXISTS update_updated_at_column();