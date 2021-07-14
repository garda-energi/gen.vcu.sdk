package main

import (
	"github.com/davecgh/go-spew/spew"
	sdk "github.com/pudjamansyurin/gen_vcu_sdk"
)

func main() {
	sdk := sdk.New("test.mosquitto.org", 1883, "", "")

	sdk.AddReportListener(reportListener)
	sdk.ConnectAndListen()
}

func reportListener(report interface{}, bytes []byte) {
	spew.Dump(report)
}
