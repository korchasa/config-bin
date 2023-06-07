CREATE TABLE bins
(
    uuid              TEXT PRIMARY KEY,
    current_version INTEGER
);

CREATE TABLE configurations
(
    uuid         TEXT,
    data       TEXT,
    version    INTEGER,
    created_at TIMESTAMP,
    format     TEXT,
    FOREIGN KEY (uuid) REFERENCES bins (uuid)
);
