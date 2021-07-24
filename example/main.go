package main

import (
	"fmt"
	"log"

	sdk "github.com/pudjamansyurin/gen.vcu.sdk"
)

func main() {
	stopChan := sdk.SetupGracefulShutdown() // optional code

	api := sdk.New(sdk.BrokerConfig{
		Host: "test.mosquitto.org",
		Port: 1883,
		User: "",
		Pass: "",
	}, true)

	// connect to broker
	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}
	defer api.Disconnect()

	// prepare the status & report listener
	listener := &sdk.Listener{
		StatusFunc: func(vin int, online bool) {
			status := map[bool]string{
				false: "OFFLINE",
				true:  "ONLINE",
			}[online]

			fmt.Printf("%d => %s\n", vin, status)
		},
		ReportFunc: func(vin int, report *sdk.ReportPacket) {
			fmt.Println(report)
		},
	}

	// listen to report
	// see api.Addlistener doc for usage
	vins := sdk.VinRange(354309, 354323)
	if err := api.AddListener(vins, listener); err != nil {
		fmt.Println(err)
	} else {
		defer api.RemoveListener(vins)
	}

	// listen to commands & response
	if dev354313, err := api.NewCommander(354313); err != nil {
		fmt.Println(err)
	} else {
		defer dev354313.Destroy()

		if info, err := dev354313.GenInfo(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(info)
		}

		// if err := dev354313.GenLed(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Device led (on-board) was turned-off")
		// }

		// rtc := time.Now()
		// if err := dev354313.GenRtc(rtc); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("RTC synced to %s\n", rtc)
		// }

		// km := uint16(54321)
		// if err := dev354313.GenOdo(km); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Odometer changed to %d km\n", km)
		// }

		// if err := dev354313.GenAntiTheaf(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Anti-theaf detector was toggled")
		// }

		// if err := dev354313.GenReportFlush(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Report buffer was flushed")
		// }

		// if err := dev354313.GenReportBlock(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Report is unblocked")
		// }

		// bikeState := sdk.BikeStateNormal
		// if err := dev354313.OvdState(bikeState); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Bike state is changed to %s\n", bikeState)
		// }

		// reportInterval := 5 * time.Second
		// if err := dev354313.OvdReportInterval(reportInterval); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Report interval changed to %s\n", reportInterval)
		// }

		// frame := sdk.FrameFull
		// if err := dev354313.OvdReportFrame(frame); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Report frame changed to %s\n", frame)
		// }

		// if err := dev354313.OvdRemoteSeat(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Remote button seat was toggled")
		// }

		// if err := dev354313.OvdRemoteAlarm(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Remote button alarm was toggled")
		// }

		// if err := dev354313.AudioBeep(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Beep sound has been generated")
		// }

		// if ids, err := dev354313.FingerFetch(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Registered driverID are : %v\n", ids)
		// }

		// if id, err := dev354313.FingerAdd(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("New driverID registered as %d\n", id)
		// }

		// driverId := 1
		// if err := dev354313.FingerDel(driverId); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("DriverID %d deleted\n", driverId)
		// }

		// if err := dev354313.FingerRst(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("All driverID are deleted")
		// }

		// if err := dev354313.RemotePairing(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Success pairing with new keyless/fob")
		// }

		// if res, err := dev354313.FotaVcu(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("VCU firmware is updgraded, %s\n", res)
		// }

		// if res, err := dev354313.FotaHmi(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("HMI firmware is updgraded, %s\n", res)
		// }

		// if res, err := dev354313.NetSendUssd("*123*10*3#"); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(res)
		// }

		// if res, err := dev354313.NetReadSms(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(res)
		// }

		// driveMode := sdk.ModeDriveStandard
		// if err := dev354313.HbarDrive(driveMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Drive mode changed to %s\n", driveMode)
		// }

		// tripMode := sdk.ModeTripOdo
		// if err := dev354313.HbarTrip(tripMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Trip mode changed to %s\n", tripMode)
		// }

		// avgMode := sdk.ModeAvgEfficiency
		// if err := dev354313.HbarAvg(avgMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Average mode changed to %s\n", avgMode)
		// }

		// if err := dev354313.HbarReverse(true); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Motor now in reverse direction")
		// }

		// kph := uint8(100)
		// if err := dev354313.McuSpeedMax(kph); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Motor speed is limited to %d kph\n", kph)
		// }

		// templates := []sdk.McuTemplate{
		// 	{DisCur: 50, Torque: 10}, // economy
		// 	{DisCur: 50, Torque: 20}, // standard
		// 	{DisCur: 50, Torque: 25}, // sport
		// }
		// if err := dev354313.McuTemplates(templates); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	for i, t := range templates {
		// 		fmt.Printf("Motor template for %s changed to %+v\n", sdk.ModeDrive(i), t)
		// 	}
		// }
	}

	<-stopChan
}
