package main

import (
	"fmt"
	"time"
	"log"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func main() {
	api := sdk.New("test.mosquitto.org", 1883, "", "", true)

	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}

	api.Listen(sdk.Listener{
		StatusFunc: statusListener,
		ReportFunc: reportListener,
	})
	time.Sleep(5 * time.Second)

	//if res, err := api.Command.GenInfo(354313); err == nil {
	//	fmt.Println(res)
	//}

	util.WaitForCtrlC()
}

func statusListener(vin int, online bool) error {
	status := map[bool]string{
		false: "OFFLINE",
		true:  "ONFLINE",
	}

	fmt.Printf("%d => %s\n", vin, status[online])
	return nil
}

func reportListener(vin int, report *report.ReportPacket) error {
	fmt.Println(report)
	return nil
}
