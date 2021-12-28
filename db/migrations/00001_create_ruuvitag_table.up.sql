CREATE TABLE IF NOT EXISTS ruuvitag
(
    id      serial PRIMARY KEY,
    name TEXT NOT NULL,
    address TEXT NOT NULL,

    UNIQUE (address)
);
