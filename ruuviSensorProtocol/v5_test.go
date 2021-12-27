package ruuviSensorProtocol_test

import (
    "encoding/hex"
    "testing"

    "ruuviLogger/ruuviSensorProtocol"
)

// cf. https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md#test-vectors
const (
    validData     = "0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F"
    maxValues     = "057FFFFFFEFFFE7FFF7FFF7FFFFFDEFEFFFECBB8334C884F"
    minValues     = "058001000000008001800180010000000000CBB8334C884F"
    invalidValues = "058000FFFFFFFF800080008000FFFFFFFFFFFFFFFFFFFFFF"
)

func Test_valid_data_from_Ruuvi_project_can_be_interpreted(t *testing.T) {
    v, _ := hex.DecodeString(validData)

    sd, err := ruuviSensorProtocol.NewSensorData(v)
    if err != nil {
        t.Errorf("got an error initializing new sensor data from a valid byte string: %s", err)
    }

    if *sd.Temperature != 24.3 {
        t.Errorf("expected temperature 24.3, got %f", *sd.Temperature)
    }

    if *sd.Humidity != 53.49 {
        t.Errorf("expected humidity 53.49, got %f", *sd.Humidity)
    }

    if *sd.Pressure != 1000.44 {
        t.Errorf("expected air pressure 1000.44, got %f", *sd.Pressure)
    }

    if *sd.Acceleration.X != 0.004 {
        t.Errorf("expected x acceleration 0.004, got %f", *sd.Acceleration.X)
    }

    if *sd.Acceleration.Y != -0.004 {
        t.Errorf("expected y acceleration -0.004, got %f", *sd.Acceleration.Y)
    }

    if *sd.Acceleration.Z != 1.036 {
        t.Errorf("expected z acceleration 1.036, got %f", *sd.Acceleration.Z)
    }

    if *sd.BatteryVoltage != 2.977 {
        t.Errorf("expected battery voltage 2.977, got %f", *sd.BatteryVoltage)
    }

    if *sd.TXPower != 4.0 {
        t.Errorf("expected TX power 4.0, got %f", *sd.TXPower)
    }
}

func Test_max_data_from_Ruuvi_project_can_be_interpreted(t *testing.T) {
    v, _ := hex.DecodeString(maxValues)

    sd, err := ruuviSensorProtocol.NewSensorData(v)
    if err != nil {
        t.Errorf("got an error initializing new sensor data from a valid byte string: %s", err)
    }

    if *sd.Temperature != 163.835 {
        t.Errorf("expected temperature 163.835, got %f", *sd.Temperature)
    }

    if *sd.Humidity != 163.835 {
        t.Errorf("expected humidity 163.835, got %f", *sd.Humidity)
    }

    if *sd.Pressure != 1155.34 {
        t.Errorf("expected air pressure 1155.34, got %f", *sd.Pressure)
    }

    if *sd.Acceleration.X != 32.767 {
        t.Errorf("expected x acceleration 32.767, got %f", *sd.Acceleration.X)
    }

    if *sd.Acceleration.Y != 32.767 {
        t.Errorf("expected y acceleration 32.767, got %f", *sd.Acceleration.Y)
    }

    if *sd.Acceleration.Z != 32.767 {
        t.Errorf("expected z acceleration 32.767, got %f", *sd.Acceleration.Z)
    }

    if *sd.BatteryVoltage != 3.646 {
        t.Errorf("expected battery voltage 3.646, got %f", *sd.BatteryVoltage)
    }

    if *sd.TXPower != 20.0 {
        t.Errorf("expected TX power 20.0, got %f", *sd.TXPower)
    }
}

func Test_min_data_from_Ruuvi_project_can_be_interpreted(t *testing.T) {
    v, _ := hex.DecodeString(minValues)

    sd, err := ruuviSensorProtocol.NewSensorData(v)
    if err != nil {
        t.Errorf("got an error initializing new sensor data from a valid byte string: %s", err)
    }

    if *sd.Temperature != -163.835 {
        t.Errorf("expected temperature -163.835, got %f", *sd.Temperature)
    }

    if *sd.Humidity != 0.0 {
        t.Errorf("expected humidity 0.0, got %f", *sd.Humidity)
    }

    if *sd.Pressure != 500.0 {
        t.Errorf("expected air pressure 500.0, got %f", *sd.Pressure)
    }

    if *sd.Acceleration.X != -32.767 {
        t.Errorf("expected x acceleration -32.767, got %f", *sd.Acceleration.X)
    }

    if *sd.Acceleration.Y != -32.767 {
        t.Errorf("expected y acceleration -32.767, got %f", *sd.Acceleration.Y)
    }

    if *sd.Acceleration.Z != -32.767 {
        t.Errorf("expected z acceleration -32.767, got %f", *sd.Acceleration.Z)
    }

    if *sd.BatteryVoltage != 1.6 {
        t.Errorf("expected battery voltage 1.6, got %f", *sd.BatteryVoltage)
    }

    if *sd.TXPower != -40.0 {
        t.Errorf("expected TX power -40.0, got %f", *sd.TXPower)
    }
}

func Test_invalid_data_from_Ruuvi_project_is_recognized_as_such(t *testing.T) {
    // invalid also means unavailable here!
    v, _ := hex.DecodeString(invalidValues)

    sd, err := ruuviSensorProtocol.NewSensorData(v)
    if err != nil {
        t.Errorf("got an error initializing new sensor data from a valid byte string: %s", err)
    }

    if sd.Temperature != nil || sd.Humidity != nil || sd.Pressure != nil ||
        sd.Acceleration.X != nil || sd.Acceleration.Y != nil || sd.Acceleration.Z != nil ||
        sd.BatteryVoltage != nil || sd.TXPower != nil {
        t.Error("the invalid/unavailable sensor data should consist of only nil values but it doesn't")
    }
}
