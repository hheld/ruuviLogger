CREATE TABLE IF NOT EXISTS measurement_kind
(
    id   serial PRIMARY KEY,
    name TEXT NOT NULL,
    unit TEXT NOT NULL,

    UNIQUE (name)
);

INSERT INTO measurement_kind (name, unit)
VALUES ('temperature', '°C'),
       ('air_humidity', '%'),
       ('air_pressure', 'hPa');