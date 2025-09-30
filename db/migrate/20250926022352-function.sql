-- 共通使用する関数を記述

-- +migrate Up
-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $BODY$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$BODY$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION create_updated_at_trigger(table_name TEXT)
RETURNS VOID AS $BODY$
BEGIN
    EXECUTE format('
        CREATE TRIGGER update_%I_updated_at
            BEFORE UPDATE ON %I
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
    ', table_name, table_name);
END;
$BODY$ LANGUAGE 'plpgsql';
-- +migrate StatementEnd

-- +migrate Down
DROP FUNCTION IF EXISTS create_updated_at_trigger(TEXT);
DROP FUNCTION IF EXISTS update_updated_at_column();