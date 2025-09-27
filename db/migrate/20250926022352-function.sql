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



-- +migrate Down
DROP FUNCTION IF EXISTS update_updated_at_column();