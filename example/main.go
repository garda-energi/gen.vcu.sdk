package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pudjamansyurin/gen_vcu_sdk/api"
)

func main() {
	api := api.New("test.mosquitto.org", 1883, "", "")

	api.AddReportListener(reportListener)
	api.ConnectAndListen()
}

func reportListener(report interface{}, bytes []byte) {
	spew.Dump(report)
}
