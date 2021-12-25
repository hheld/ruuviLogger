package main

import (
	"log"

	"github.com/godbus/dbus/v5"
	ble "github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
)

const (
	basement             = "C9:73:96:D3:F9:A0"
	outside              = "F2:03:48:17:D1:28"
	livingRoom           = "C5:CE:2D:AA:D3:87"
	ruuviManufacturerID  = 0x0499
	manufacturerDataProp = "ManufacturerData"
)

func main() {
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

	for device := range discovery {
		log.Printf("found device: %+v", device.Path)

		dev, err := btAdapter.GetDeviceByAddress(livingRoom)
		if err != nil {
			log.Printf("error getting device by address: %s", err)
		}

		if dev != nil {
			propsCh, err := dev.WatchProperties()

			if err == nil {
				go func() {
					for {
						select {
						case msg := <-propsCh:
							if msg.Name == manufacturerDataProp {
								if values, ok := msg.Value.(map[uint16]dbus.Variant)[ruuviManufacturerID]; ok {
									log.Printf("values: %x", values.Value().([]byte))
								}
							}
						}
					}
				}()

				break
			}
		}
	}

	select {}
}
