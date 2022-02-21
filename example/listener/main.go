package main

import (
	"fmt"
	"log"

	sdk "github.com/garda-energi/gen.vcu.sdk"
)

func main() {
	stopChan := sdk.SetupGracefulShutdown() // optional code

	api := sdk.New(sdk.ClientConfig{
		Host: "broker.emqx.io",
		Port: 1883,
		// Host: "test.mosquitto.org",
		// Port: 1884,
		User: "rw",
		Pass: "readwrite",
	}, false)

	// connect to client
	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}
	defer api.Disconnect()

	// prepare the status & report listener
	listener := sdk.Listener{
		StatusFunc: func(vin int, online bool) {
			status := map[bool]string{
				false: "OFFLINE",
				true:  "ONLINE",
			}[online]
			fmt.Println(vin, "=>", status)
		},
		ReportFunc: func(vin int, report *sdk.ReportPacket) {
			fmt.Println(vin, "=>", report)

			// Get Error
			bikeErr := report.GetBikeError()
			if bikeErr != sdk.BIKE_NOERROR {
				fmt.Println("Code:", byte(bikeErr))
				fmt.Println("Error:", bikeErr.Error())
			}

			// fmt.Println(string(report.Json()))

			// fmt.Println(report.GetValue("Bms.Pack.[0].Capacity.Remaining"))
			// fmt.Println(report.GetType("Bms.Pack.[0].Capacity.Remaining"))

			// vcuEvent, isOK := report.GetValue("Vcu.Events").(uint16)
			// if isOK {
			// 	fmt.Println(vcuEvent)
			// }
			// fmt.Println(report.GetType("Vcu.Events"))

			// expose all *ReportPacket methods available
			// if report.Vcu.RealtimeData() {
			// 	fmt.Println("Current report is realtime")
			// }
			// if report.Gps.ValidHorizontal() {
			// 	fmt.Println("GPS longitude, latitude & heading is valid")
			// }
			// if report.Bms.LowCapacity() {
			// 	fmt.Println("BMS need to be charged on Charging Station")
			// }
		},
	}

	// listen to all vins
	if err := api.AddListener(listener); err != nil {
		fmt.Println(err)
	} else {
		defer api.RemoveListener()
	}

	// listen to range of vins
	// see api.Addlistener doc for usage
	// vins := sdk.VinRange(354309, 354323)
	// if err := api.AddListener(listener, vins...); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	defer api.RemoveListener(vins...)
	// }

	<-stopChan
}
