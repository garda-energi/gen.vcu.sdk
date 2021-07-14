package main

import (
	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func main() {
	sdk := sdk.New("test.mosquitto.org", 1883, "", "")

	sdk.AddReportListener(reportListener)
	sdk.ConnectAndListen()
}

func reportListener(report interface{}, bytes []byte) {
	util.Debug(report)
}
