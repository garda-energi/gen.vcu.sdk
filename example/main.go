package main

import (
	"fmt"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
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
	fmt.Printf("%d => %s\n", vin, status)

	return nil
}

func reportListener(vin int, result interface{}) error {
	switch r := result.(type) {
	case report.ReportSimple:
		fmt.Printf("[S] %d  => %+v\n", vin, r)
	case report.ReportFull:
		fmt.Printf("[F] %d  => %+v\n", vin, r)
	}

	return nil
}
