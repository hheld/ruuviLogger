package main

import (
	"github.com/godbus/dbus/v5"
	ble "github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"log"
	"ruuviLogger"
	"ruuviLogger/db"
	"ruuviLogger/ruuviSensorProtocol"
	"time"

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

	discovery, cancelDiscovery, err := ble.Discover(btAdapter, nil)
	if err != nil {
		log.Fatalf("error discovering Bluetooth devices: %s", err)
	}
	defer cancelDiscovery()

	var indicesOfDiscoveredRuuvitags []int

	for range discovery {
		//for {
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

		var lastSensorData *ruuviSensorProtocol.SensorData = nil

		if err == nil {
			lastSequenceNo := ^uint16(0)

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
									if sd.SequenceNo != nil {
										if *sd.SequenceNo != lastSequenceNo {
											lastSequenceNo = *sd.SequenceNo
											lastSensorData = sd
										}
									}
								}
							}
						}
					}
				}
			}()

			go func() {
				for {
					select {
					case <-time.After(time.Duration(ruuviLogger.Cfg.DbWriteIntervalInSec) * time.Second):
						if lastSensorData != nil {
							if err := sensorDb.AddMeasurement(*lastSensorData, sensorID); err != nil {
								log.Printf("could not write values %s to DB for sensor %s: %s", lastSensorData.ToString(), rt.Properties.Address, err)
							} else {
								lastSensorData = nil
							}
						}
					}
				}
			}()
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
