-- 0001_init.sql
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS items (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name     TEXT NOT NULL UNIQUE,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    price    REAL NOT NULL CHECK (price > 0),
    tags     TEXT NOT NULL, -- JSON-массив строк в TEXT
    status   TEXT NOT NULL CHECK (status IN ('active','archived'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_items_name ON items(name);
