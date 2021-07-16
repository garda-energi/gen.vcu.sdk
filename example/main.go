package main

import (
	"fmt"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

func main() {
	api := sdk.New("test.mosquitto.org", 1883, "", "")

	api.AddStatusListener(statusListener)

	// only choose one data type, not both
	dtype := sdk.DATA_TYPE_STRUCT
	api.SetDataType(dtype)
	if dtype == sdk.DATA_TYPE_LIST {
		api.AddReportListener(reportListListener)
	} else {
		api.AddReportListener(reportStructListener)
	}

	api.Logging(false)
	api.ConnectAndListen()
}

func statusListener(vin int, online bool) error {
	status := map[bool]string{
		false: "OFFLINE",
		true:  "ONFLINE",
	}

	fmt.Printf("%d => %s\n", vin, status[online])
	return nil
}

func reportStructListener(vin int, result interface{}) error {
	// switch r := result.(type) {
	// case report.ReportSimple:
	// 	fmt.Printf("[SIMPLE] %d  => %+v\n", vin, r)
	// case report.ReportFull:
	// 	fmt.Printf("[FULL] %d  => %+v\n", vin, r)
	// }

	reportPacket := result.(*report.ReportPacket)
	fmt.Printf("%+v\n", reportPacket)
	return nil
}

func reportListListener(vin int, result interface{}) error {
	frame := "SIMPLE"

	items, _ := result.(report.Items)
	if frameId, ok := items["header.frameId"]; ok {
		if frameId.Value.(shared.FRAME_ID) == shared.FRAME_ID_FULL {
			frame = "FULL"
		}
	}

	fmt.Printf("[%s] %d  => %+v\n", frame, vin, result)
	return nil
}
