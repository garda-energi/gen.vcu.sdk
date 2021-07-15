package main

import (
	"fmt"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	// "github.com/pudjamansyurin/gen_vcu_sdk/model"
	"github.com/pudjamansyurin/gen_vcu_sdk/packet"
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

func reportListener(vin int, report interface{}) error {
	// switch r := report.(type) {
	// case model.ReportSimple:
	// 	fmt.Printf("[S] %d  => %+v\n", vin, r)
	// case model.ReportFull:
	// 	fmt.Printf("[F] %d  => %+v\n", vin, r)
	// }
	reportPacket := report.(*packet.ReportPacket)
	fmt.Println("======= Report ========")
	fmt.Printf("[vin] %d\n", vin)
	fmt.Printf("[header]\n-prefix\t: %s\n-length\t: %d\n", reportPacket.Header.Prefix, reportPacket.Header.Size)
	if reportPacket.Mems != nil {
		fmt.Printf("[Mem]\n-active\t: %t\n-total\t: %.2f\n", reportPacket.Mems.Active, reportPacket.Mems.Total)
	}
	fmt.Println("======= ====== ========")

	return nil
}
