CREATE TABLE IF NOT EXISTS posts
(
    uuid     CHAR(38) PRIMARY KEY,
    content  TEXT NOT NULL,
    metadata JSON
);