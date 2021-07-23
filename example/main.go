package main

import (
	"fmt"
	"log"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	// "github.com/pudjamansyurin/gen_vcu_sdk/command"
)

func main() {
	api := sdk.New("test.mosquitto.org", 1883, "", "", true)

	// connect to broker
	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}
	defer api.Disconnect()

	// prepare the status & report listener
	listener := &sdk.Listener{
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
	}

	// listen to report
	// see api.Addlistener doc for usage
	reportVins := sdk.VinRange(354309, 354323)
	if err := api.AddListener(reportVins, listener); err != nil {
		fmt.Println(err)
	} else {
		defer api.RemoveListener(reportVins)
	}

	// listen to commands & response
	if dev354313, err := api.NewCommand(354313); err != nil {
		fmt.Println(err)
	} else {
		defer dev354313.Destroy()

		info, err := dev354313.GenInfo()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(info)
		}

		// if err := dev354313.GenLed(false); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println("On-board led was turned-off")
		// }

		// rtc := time.Now()
		// if err := dev354313.GenRtc(rtc); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("RTC synced to %s\n", rtc)
		// }

		// km := uint16(61234)
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

		// bikeState := shared.BIKE_STATE_NORMAL
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

		// frame := shared.FRAME_ID_FULL
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

		// ids, err := dev354313.FingerFetch()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Registered driverID are : %v\n", ids)
		// }

		// id, err := dev354313.FingerAdd()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("New driverID registered = %d\n", id)
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

		// vcuRes, err := dev354313.FotaVcu()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("VCU fw is updgraded, %s\n", vcuRes)
		// }

		// hmiRes, err := dev354313.FotaHmi()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("HMI fw is updgraded, %s\n", hmiRes)
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
		// 	fmt.Printf("Drive mode changed to %s\n", driveMode)
		// }

		// tripMode := shared.MODE_TRIP_ODO
		// if err := dev354313.HbarTrip(tripMode); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Printf("Trip mode changed to %s\n", tripMode)
		// }

		// avgMode := shared.MODE_AVG_EFFICIENCY
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

		// templates := []command.McuTemplate{
		// 	{DischargeCurrent: 50, Torque: 10}, // economy
		// 	{DischargeCurrent: 50, Torque: 20}, // standard
		// 	{DischargeCurrent: 50, Torque: 25}, // sport
		// }
		// if err := dev354313.McuTemplates(templates); err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	for i, t := range templates {
		// 		fmt.Printf("Motor template for %s changed to %+v\n", shared.MODE_DRIVE(i), t)
		// 	}
		// }
	}

	shared.WaitForCtrlC()
}
