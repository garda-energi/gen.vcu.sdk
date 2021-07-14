package main

import (
	"fmt"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
)

func main() {
	sdk := sdk.New("test.mosquitto.org", 1883, "", "")

	sdk.AddStatusListener(statusListener)
	sdk.AddReportListener(reportListener)

	sdk.Logging(false)
	sdk.ConnectAndListen()
}

func statusListener(vin int, online bool) error {
	// why go not support ternary operation ?
	status := "OFFLINE"
	if online {
		status = "ONLINE"
	}
	fmt.Printf("%d is %s\n", vin, status)

	return nil
}

func reportListener(vin int, report interface{}) error {
	fmt.Printf("%+v\n", report)

	return nil
}
