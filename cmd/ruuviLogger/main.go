package main

import (
	"github.com/muka/go-bluetooth/bluez"
	"log"
	"time"

	"github.com/godbus/dbus/v5"
	ble "github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"ruuviLogger"
	"ruuviLogger/db"
	"ruuviLogger/ruuviSensorProtocol"

	_ "ruuviLogger"
)

const (
	ruuviManufacturerID  = 0x0499
	manufacturerDataProp = "ManufacturerData"
)

type sensorStore interface {
	AddRuuvitag(name, address string) error
	GetRuuvitagID(address string) (int, error)
	AddMeasurement(data ruuviSensorProtocol.SensorData, ruuvitagID int) error
}

func discoverBluetoothDevices(ruuvitagDiscovered chan<- *device.Device1, sensorDb sensorStore) {
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
					indicesOfDiscoveredRuuvitags = append(indicesOfDiscoveredRuuvitags, i)
					if err := sensorDb.AddRuuvitag(desiredRuuvitag.Name, desiredRuuvitag.Address); err != nil {
						log.Printf("could not add newly discovered Ruuvitag to the database: %s", err)
					}
					ruuvitagDiscovered <- rt
				}
			}
		}
	}
}

func logData(rt *device.Device1, sensorDb sensorStore) {
	if rt != nil {
		sensorID, err := sensorDb.GetRuuvitagID(rt.Properties.Address)
		if err != nil {
			log.Printf("could not get the sensor ID: %s", err)
			return
		}

		propsCh, err := rt.WatchProperties()

		if err == nil {
			sensorDataCh := make(chan ruuviSensorProtocol.SensorData)

			go func(sensorData chan<- ruuviSensorProtocol.SensorData, propsCh <-chan *bluez.PropertyChanged) {
				lastSequenceNo := ^uint16(0)

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
									if sd.SequenceNo != nil {
										if *sd.SequenceNo != lastSequenceNo {
											lastSequenceNo = *sd.SequenceNo
											sensorData <- *sd
										}
									}
								}
							}
						}
					}
				}
			}(sensorDataCh, propsCh)

			go func(sensorData <-chan ruuviSensorProtocol.SensorData) {
				var lastSensorData *ruuviSensorProtocol.SensorData = nil

				for {
					select {
					case <-time.After(time.Duration(ruuviLogger.Cfg.DbWriteIntervalInSec) * time.Second):
						if lastSensorData != nil {
							if err := sensorDb.AddMeasurement(*lastSensorData, sensorID); err != nil {
								log.Printf("could not write values %s to DB for sensor %s: %s", lastSensorData.ToString(), rt.Properties.Address, err)
							}
						}
					case sd := <-sensorData:
						lastSensorData = &sd
					}
				}
			}(sensorDataCh)
		}
	}
}

func main() {
	sensorDb, err := db.ConnectToDb()
	if err != nil {
		log.Fatalf("could not connect to the database: %s", err)
	}

	ruuviTagsDiscovered := make(chan *device.Device1)

	go discoverBluetoothDevices(ruuviTagsDiscovered, sensorDb)

	for {
		select {
		case ruuvitagDevice := <-ruuviTagsDiscovered:
			log.Printf("found Ruuvitag %s", ruuvitagDevice.Properties.Address)
			logData(ruuvitagDevice, sensorDb)
		}
	}
}
