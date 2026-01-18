
-- +migrate Up
INSERT INTO options(id, name, deprecated)
VALUES (7, 'todo', FALSE)
ON CONFLICT (id) DO NOTHING;

-- +migrate Down
DELETE FROM options WHERE id = 7;
