package db

import (
    "context"
    "time"

    "ruuviLogger/ruuviSensorProtocol"
)

func (db *sensorDb) AddMeasurement(data ruuviSensorProtocol.SensorData, ruuvitagID int) error {
    query := `INSERT INTO measurements
                (time, temperature, humidity, pressure, accel_x, accel_y, accel_z, battery_voltage, tx_power, movements, ruuvitag_id)
              VALUES
                ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

    movements := ^uint8(0)
    if data.MovementCounter != nil {
        movements = *data.MovementCounter
    }

    _, err := db.Pool.Exec(context.Background(), query,
        time.Now().UTC(),
        ruuviSensorProtocol.ValueOrNan(data.Temperature),
        ruuviSensorProtocol.ValueOrNan(data.Humidity),
        ruuviSensorProtocol.ValueOrNan(data.Pressure),
        ruuviSensorProtocol.ValueOrNan(data.Acceleration.X),
        ruuviSensorProtocol.ValueOrNan(data.Acceleration.Y),
        ruuviSensorProtocol.ValueOrNan(data.Acceleration.Z),
        ruuviSensorProtocol.ValueOrNan(data.BatteryVoltage),
        ruuviSensorProtocol.ValueOrNan(data.TXPower),
        movements,
        ruuvitagID)

    return err
}
