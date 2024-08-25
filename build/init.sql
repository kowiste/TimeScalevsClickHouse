CREATE TABLE IF NOT EXISTS measures (
    id TEXT,
    asset TEXT,
    value INT NOT NULL,
    time TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (id, time)
);
-- Create an index on the asset column
CREATE INDEX idx_asset ON measures (asset);