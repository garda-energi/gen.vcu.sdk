package main

import (
	"fmt"
	"log"
	"time"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func main() {
	api := sdk.New("test.mosquitto.org", 1883, "", "", true)

	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}
	defer api.Disconnect()

	api.Listen(sdk.Listener{
		StatusFunc: func(vin int, online bool) error {
			status := map[bool]string{
				false: "OFFLINE",
				true:  "ONFLINE",
			}[online]

			fmt.Printf("%d => %s\n", vin, status)
			return nil
		},
		ReportFunc: func(vin int, report *report.ReportPacket) error {
			// fmt.Println(report)
			return nil
		},
	})

	{
		dev354313 := api.NewCommand(354313)

		info, err := dev354313.GenInfo()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(info)
		}

		// if err := dev354313.GenLed(true); err != nil {
		// 	fmt.Println(err)
		// }

		rtc := time.Now()
		if err := dev354313.GenRtc(rtc); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("RTC synced to %s\n", rtc)
		}

		// if err := dev354313.GenOdo(0); err != nil {
		// 	fmt.Println(err)
		// }

		// if err := dev354313.OvdState(shared.BIKE_STATE_NORMAL); err != nil {
		// 	fmt.Println(err)
		// }

		// ids, err := dev354313.FingerFetch()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(ids)
		// }
	}

	util.WaitForCtrlC()
}
