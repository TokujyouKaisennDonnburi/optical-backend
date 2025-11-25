
-- +migrate Up
CREATE TABLE event_locations(
    event_id UUID NOT NULL UNIQUE REFERENCES events(id) ON DELETE CASCADE,
    location VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX idx_event_locations_event_id_and_location ON event_locations(event_id, location);
CREATE INDEX idx_event_locations_event_id ON event_locations(event_id);

-- +migrate Down

DROP INDEX IF EXISTS idx_event_locations_event_id_and_location 
DROP INDEX IF EXISTS idx_event_locations_event_id;
DROP TABLE IF EXISTS event_locaitons;
