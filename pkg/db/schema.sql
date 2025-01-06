CREATE TABLE IF NOT EXISTS sessions (
    name TEXT PRIMARY KEY,
    command TEXT NOT NULL,
    args TEXT,
    status TEXT
);

