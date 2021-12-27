package main

import (
    "log"

    "github.com/godbus/dbus/v5"
    ble "github.com/muka/go-bluetooth/api"
    "github.com/muka/go-bluetooth/bluez/profile/adapter"
    "github.com/muka/go-bluetooth/bluez/profile/device"
    "ruuviLogger"
    "ruuviLogger/ruuviSensorProtocol"

    _ "ruuviLogger"
)

const (
    ruuviManufacturerID  = 0x0499
    manufacturerDataProp = "ManufacturerData"
)

func discoverBluetoothDevices(ruuvitagDiscovered chan<- *device.Device1) {
    btAdapter, err := ble.GetDefaultAdapter()
    if err != nil {
        log.Fatalf("could not find a Bluetooth adapter: %s", err)
    }

    discoveryOptions := adapter.DiscoveryFilter{
        UUIDs:         []string{},
        RSSI:          0,
        Pathloss:      0,
        Transport:     "le",
        DuplicateData: true,
    }

    discovery, cancelDiscovery, err := ble.Discover(btAdapter, &discoveryOptions)
    if err != nil {
        log.Fatalf("error discovering Bluetooth devices: %s", err)
    }
    defer cancelDiscovery()

    var indicesOfDiscoveredRuuvitags []int

    for range discovery {
        for i, desiredRuuvitag := range ruuviLogger.Settings.Ruuvitags {
            rt, err := btAdapter.GetDeviceByAddress(desiredRuuvitag.Address)
            if err != nil {
                log.Printf("error getting device by address: %s", err)
            }

            if rt != nil {
                // check if we have already found this one earlier
                foundBefore := false

                for _, idxFound := range indicesOfDiscoveredRuuvitags {
                    if idxFound == i {
                        foundBefore = true
                        break
                    }
                }

                if !foundBefore {
                    ruuvitagDiscovered <- rt
                    indicesOfDiscoveredRuuvitags = append(indicesOfDiscoveredRuuvitags, i)
                }
            }
        }
    }
}

func logData(rt *device.Device1) {
    if rt != nil {
        propsCh, err := rt.WatchProperties()

        if err == nil {
            go func() {
                for {
                    select {
                    case msg := <-propsCh:
                        if msg.Name == manufacturerDataProp {
                            if values, ok := msg.Value.(map[uint16]dbus.Variant)[ruuviManufacturerID]; ok {
                                sd, err := ruuviSensorProtocol.NewSensorData(values.Value().([]byte))
                                if err != nil {
                                    log.Printf("error interpreting sensor data from device %s from message %x: %s",
                                        rt.Properties.Address, values.Value().([]byte), err)
                                } else {
                                    log.Printf("values: %s from device %s", sd.ToString(), rt.Properties.Address)
                                }
                            }
                        }
                    }
                }
            }()
        }
    }
}

func main() {
    ruuviTagsDiscovered := make(chan *device.Device1)

    go discoverBluetoothDevices(ruuviTagsDiscovered)

    for {
        select {
        case ruuvitagDevice := <-ruuviTagsDiscovered:
            log.Printf("found Ruuvitag %s", ruuvitagDevice.Properties.Address)
            logData(ruuvitagDevice)
        }
    }
}
