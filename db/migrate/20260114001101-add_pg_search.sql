-- +migrate Up
-- pg_search拡張機能を有効化
CREATE EXTENSION IF NOT EXISTS pg_search;

-- eventsテーブルにBM25インデックスを作成
CREATE INDEX search_events_idx ON events
USING bm25 (id, title, memo)
WITH (key_field='id');

-- event_locationsテーブルにBM25インデックスを作成
CREATE INDEX search_event_locations_idx ON event_locations
USING bm25 (event_id, location)
WITH (key_field='event_id');

-- +migrate Down
DROP INDEX IF EXISTS search_event_locations_idx;
DROP INDEX IF EXISTS search_events_idx;
DROP EXTENSION IF EXISTS pg_search;
