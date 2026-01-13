-- +migrate Up
CREATE EXTENSION IF NOT EXISTS pg_search;

CREATE INDEX IF NOT EXISTS idx_events_fulltext_search
ON events
USING bm25 (title, memo)
WITH (key_field = 'id');

-- +migrate Down
DROP INDEX IF EXISTS idx_events_fulltext_search;
DROP EXTENSION IF EXISTS pg_search;
