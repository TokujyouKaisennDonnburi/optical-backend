-- +migrate Up
DROP TABLE IF EXISTS calendar_options;
DROP TABLE IF EXISTS options;

CREATE TABLE options(
    id INT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    deprecated BOOLEAN DEFAULT FALSE
);
COMMENT ON COLUMN options.id IS 'オプションID';
COMMENT ON COLUMN options.name IS 'オプション名';
COMMENT ON COLUMN options.deprecated IS 'オプション推奨フラグ';

CREATE TABLE calendar_options(
    calendar_id UUID REFERENCES calendars(id) ON DELETE CASCADE,
    option_id INT REFERENCES options(id) ON DELETE CASCADE,
    PRIMARY KEY (calendar_id, option_id)
);
COMMENT ON COLUMN calendar_options.calendar_id IS 'カレンダーID';
COMMENT ON COLUMN calendar_options.option_id IS 'オプションID';

INSERT INTO options(id,name,deprecated)VALUES
(1, 'pr_review_pending_count', FALSE),
(2, 'review_load_status', FALSE);

-- +migrate Down

  DROP TABLE IF EXISTS calendar_options;
  DROP TABLE IF EXISTS options;

  CREATE TABLE options(
      id UUID PRIMARY KEY,
      name VARCHAR(255) NOT NULL
  );

  CREATE TABLE calendar_options(
      calendar_id UUID REFERENCES calendars(id) ON
      DELETE CASCADE,
      option_id UUID REFERENCES options(id) ON
      DELETE CASCADE,
      PRIMARY KEY (calendar_id, option_id)
  );

