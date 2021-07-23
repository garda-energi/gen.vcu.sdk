package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	// "github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

func main() {
	api := sdk.New("test.mosquitto.org", 1883, "", "", true)

	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}
	defer api.Disconnect()

	// api.Listen(sdk.Listener{
	// 	StatusFunc: func(vin int, online bool) error {
	// 		status := map[bool]string{
	// 			false: "OFFLINE",
	// 			true:  "ONLINE",
	// 		}[online]

	// 		fmt.Printf("%d => %s\n", vin, status)
	// 		return nil
	// 	},
	// 	ReportFunc: func(vin int, report *report.ReportPacket) error {
	// 		// fmt.Println(report)
	// 		return nil
	// 	},
	// })

	//
	// listen by list
	// api.AddListener([1, 2 ,3], .....)
	//
	// listen one spesific vin
	// api.AddListener([2341], .....)
	//
	// listen by range
	// api.AddListener(sdk.VinRange(min, max), .....)
	api.AddListener(sdk.VinRange(354309, 354323), sdk.Listener{
		StatusFunc: func(vin int, online bool) error {
			status := map[bool]string{
				false: "OFFLINE",
				true:  "ONLINE",
			}[online]

			fmt.Printf("%d => %s\n", vin, status)
			return nil
		},
		ReportFunc: func(vin int, report *report.ReportPacket) error {
			// fmt.Println(report)
			return nil
		},
	})
	fmt.Println("Listening")

	// Try to remove listener by range
	go func() {
		time.Sleep(20 * time.Second)
		api.RemoveListener(sdk.VinRange(354309, 354323))
		fmt.Println("Listener was removed")
	}()

	{
		dev354313 := api.NewCommand(354313)

		info, err := dev354313.GenInfo()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(info)
		}

		// if err := dev354313.GenLed(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("On-board led was turned-off")
		// }

		// rtc := time.Now()
		// if err := dev354313.GenRtc(rtc); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("RTC synced to %s\n", rtc)
		// }

		// km := uint16(5)
		// if err := dev354313.GenOdo(km); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Odometer changed to %d km", km)
		// }

		// if err := dev354313.GenAntiTheaf(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Anti-theaf detector was toggled")
		// }

		// if err := dev354313.GenReportFlush(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Report buffer was flushed")
		// }

		// if err := dev354313.GenReportBlock(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Report is unblocked")
		// }

		// bikeState := shared.BIKE_STATE_NORMAL
		// if err := dev354313.OvdState(bikeState); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Bike state is changed to %s", bikeState)
		// }

		// reportInterval := 5 * time.Second
		// if err := dev354313.OvdReportInterval(reportInterval); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Report interval changed to %s", reportInterval)
		// }

		// frame := shared.FRAME_ID_FULL
		// if err := dev354313.OvdReportFrame(frame); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Report frame changed to %s", frame)
		// }

		// if err := dev354313.OvdRemoteSeat(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Remote button seat was toggled")
		// }

		// if err := dev354313.OvdRemoteAlarm(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Remote button alarm was toggled")
		// }

		// if err := dev354313.AudioBeep(); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Beep sound has been generated")
		// }

		// ids, err := dev354313.FingerFetch()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Registered driverID are : %v", ids)
		// }

		// id, err := dev354313.FingerAdd()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("New driverID registered = %d", id)
		// }

		// driverId := 1
		// if err := dev354313.FingerDel(driverId); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("DriverID %d deleted", driverId)
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

		// vcuRes, err := dev354313.FotaVcu()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("VCU fw is updgraded, %s", vcuRes)
		// }

		// hmiRes, err := dev354313.FotaHmi()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("HMI fw is updgraded, %s", hmiRes)
		// }

		// res, err := dev354313.NetSendUssd("*123*10*3#")
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(res)
		// }

		// sms, err := dev354313.NetReadSms()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(sms)
		// }

		// driveMode := shared.MODE_DRIVE_STANDARD
		// if err := dev354313.HbarDrive(driveMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Drive mode changed to %s", driveMode)
		// }

		// tripMode := shared.MODE_TRIP_ODO
		// if err := dev354313.HbarTrip(tripMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Trip mode changed to %s", tripMode)
		// }

		// avgMode := shared.MODE_AVG_EFFICIENCY
		// if err := dev354313.HbarAvg(avgMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Average mode changed to %s", avgMode)
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
		// 	fmt.Printf("Motor speed is limited to %d kph", kph)
		// }
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
