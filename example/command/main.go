package main

import (
	"fmt"
	"log"

	sdk "github.com/pudjamansyurin/gen.vcu.sdk"
)

func main() {
	stopChan := sdk.SetupGracefulShutdown() // optional code

	api := sdk.New(sdk.ClientConfig{
		Host: "test.mosquitto.org",
		Port: 1884,
		User: "rw",
		Pass: "readwrite",
	}, true)

	// connect to client
	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}
	defer api.Disconnect()

	// listen to commands & response
	if dev354313, err := api.NewCommander(354313); err != nil {
		fmt.Println(err)
	} else {
		defer dev354313.Destroy()

		// expose all commands available
		if info, err := dev354313.GenInfo(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(info)
		}

		if err := dev354313.GenLed(false); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Device led (on-board) was turned-off")
		}

		// rtc := time.Now()
		// if err := dev354313.GenRtc(rtc); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("RTC synced to", rtc)
		// }

		// bikeState := sdk.BikeStateNormal
		// if err := dev354313.GenBikeState(bikeState); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Bike state is changed to", bikeState)
		// }

		// if err := dev354313.GenLockDown(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Lock-down mode is disabled")
		// }

		// if err := dev354313.GenRestart(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Main chip was restarted")
		// }

		// if err := dev354313.ReportFlush(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Report buffer was flushed")
		// }

		// if err := dev354313.ReportBlock(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Report is unblocked")
		// }

		// reportInterval := 5 * time.Second
		// if err := dev354313.ReportInterval(reportInterval); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Report interval changed to", reportInterval)
		// }

		// frame := sdk.FrameFull
		// if err := dev354313.ReportFrame(frame); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Report frame changed to", frame)
		// }

		// if err := dev354313.AudioBeep(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Beep sound has been generated")
		// }

		// if ids, err := dev354313.FingerFetch(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Registered driverID are :", ids)
		// }

		// if id, err := dev354313.FingerAdd(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("New driverID registered as", id)
		// }

		// driverId := 1
		// if err := dev354313.FingerDel(driverId); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("DriverID deleted for", driverId)
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

		// if err := dev354313.RemoteSeat(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Remote button seat was toggled")
		// }

		// if err := dev354313.RemoteAlarm(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Remote button alarm was toggled")
		// }

		// if res, err := dev354313.FotaVcu(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("VCU firmware is updgraded:", res)
		// }

		// if res, err := dev354313.FotaHmi(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("HMI firmware is updgraded:", res)
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

		// trip := sdk.ModeTripOdo
		// km := uint16(54321)
		// if err := dev354313.HbarTripMeter(trip, km); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Trip mode %s changed to %d km", trip, km)
		// }

		// driveMode := sdk.ModeDriveStandard
		// if err := dev354313.HbarDrive(driveMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Drive mode changed to", driveMode)
		// }

		// tripMode := sdk.ModeTripOdo
		// if err := dev354313.HbarTrip(tripMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Trip mode changed to", tripMode)
		// }

		// avgMode := sdk.ModeAvgEfficiency
		// if err := dev354313.HbarAvg(avgMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Average mode changed to", avgMode)
		// }

		// kph := uint8(100)
		// if err := dev354313.McuSpeedMax(kph); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Motor speed is limited to", kph, "kph")
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
		// 		fmt.Println("Motor template for", sdk.ModeDrive(i), "changed to", t)
		// 	}
		// }

		// if err := dev354313.ImuAntiThief(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("Anti-theaf detector was disabled")
		// }
	}

	<-stopChan
}
