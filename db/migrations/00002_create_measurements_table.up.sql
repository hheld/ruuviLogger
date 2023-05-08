CREATE TABLE IF NOT EXISTS measurements
(
    time          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    temperature   float,
    humidity      float,
    pressure      float,
    accel_x       float,
    accel_y       float,
    accel_z       float,
    battery_voltage float,
    tx_power      float,
    movements     int,
    ruuvitag_id   int,

    CONSTRAINT fk_ruuvitag FOREIGN KEY (ruuvitag_id) REFERENCES ruuvitag (id)
);

-- SELECT create_hypertable('measurements', 'time');
