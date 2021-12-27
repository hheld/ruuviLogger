package ruuviSensorProtocol

import (
    "bytes"
    "encoding/binary"
    "errors"
)

type SensorData struct {
    Temperature    *float64
    Humidity       *float64
    Pressure       *float64
    Acceleration   AccelerationData
    BatteryVoltage *float64
    TXPower        *float64
}

type AccelerationData struct {
    X *float64
    Y *float64
    Z *float64
}

const (
    maxUint16 = ^uint16(0)
    maxInt16  = int16(maxUint16 >> 1)
    minInt16  = -maxInt16 - 1
)

func NewSensorData(data []byte) (*SensorData, error) {
    if data[0] != 0x5 {
        return nil, errors.New("data is not protocol version 5")
    }

    var sd SensorData
    var signedData int16
    var unsignedData uint16

    buf := bytes.NewReader(data[1:3])
    _ = binary.Read(buf, binary.BigEndian, &signedData)
    temperature := float64(signedData) * 0.005
    if signedData == minInt16 {
        sd.Temperature = nil
    } else {
        sd.Temperature = &temperature
    }

    buf = bytes.NewReader(data[3:5])
    _ = binary.Read(buf, binary.BigEndian, &unsignedData)
    if unsignedData == maxUint16 {
        sd.Humidity = nil
    } else {
        humidity := float64(unsignedData) * 0.0025
        sd.Humidity = &humidity
    }

    buf = bytes.NewReader(data[5:7])
    _ = binary.Read(buf, binary.BigEndian, &unsignedData)
    if unsignedData == maxUint16 {
        sd.Pressure = nil
    } else {
        pressure := (float64(unsignedData) + 50000.0) * 0.01
        sd.Pressure = &pressure
    }

    buf = bytes.NewReader(data[7:9])
    _ = binary.Read(buf, binary.BigEndian, &signedData)
    if signedData == minInt16 {
        sd.Acceleration.X = nil
    } else {
        accelerationX := float64(signedData) * 0.001
        sd.Acceleration.X = &accelerationX
    }

    buf = bytes.NewReader(data[9:11])
    _ = binary.Read(buf, binary.BigEndian, &signedData)
    if signedData == minInt16 {
        sd.Acceleration.Y = nil
    } else {
        accelerationY := float64(signedData) * 0.001
        sd.Acceleration.Y = &accelerationY
    }

    buf = bytes.NewReader(data[11:13])
    _ = binary.Read(buf, binary.BigEndian, &signedData)
    if signedData == minInt16 {
        sd.Acceleration.Z = nil
    } else {
        accelerationZ := float64(signedData) * 0.001
        sd.Acceleration.Z = &accelerationZ
    }

    buf = bytes.NewReader(data[13:15])
    _ = binary.Read(buf, binary.BigEndian, &unsignedData)
    if unsignedData == maxUint16 {
        sd.BatteryVoltage = nil
        sd.TXPower = nil
    } else {
        batteryVoltage := float64(unsignedData>>5+1600) * 0.001
        txPower := float64(2*(unsignedData&0b11111)) - 40.0
        sd.BatteryVoltage = &batteryVoltage
        sd.TXPower = &txPower
    }

    return &sd, nil
}
