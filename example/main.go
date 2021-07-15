package main

import (
	"fmt"

	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
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
	if online {
		fmt.Printf("%d => ONLINE\n", vin)
	} else {
		fmt.Printf("%d => OFFLINE\n", vin)
	}

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
	util.Debug(reportPacket)

	return nil
}

func reportListListener(vin int, result interface{}) error {
	frame := "SIMPLE"

	items, _ := result.(report.Items)
	if frameId, ok := items["header.frameId"]; ok {
		if frameId.Value.(report.FRAME_ID) == report.FRAME_ID_FULL {
			frame = "FULL"
		}
	}

	fmt.Printf("[%s] %d  => %+v\n", frame, vin, result)

	return nil
}
