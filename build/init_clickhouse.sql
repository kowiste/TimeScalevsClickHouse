CREATE TABLE IF NOT EXISTS measures (
    id String,
    asset String,
    value Int32,
    time DateTime,
    PRIMARY KEY (id, time)
) ENGINE = MergeTree()
ORDER BY (id, time);
