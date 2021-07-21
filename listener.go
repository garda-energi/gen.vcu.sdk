package gen_vcu_sdk

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

type StatusListenerFunc func(vin int, online bool) error
type ReportListenerFunc func(vin int, report *report.ReportPacket) error

type Listener struct {
	StatusFunc StatusListenerFunc
	ReportFunc ReportListenerFunc
}

func StatusListener(sFunc StatusListenerFunc, logging bool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if logging {
			util.LogMessage(msg)
		}

		vin := util.TopicVin(msg.Topic())
		online := parseOnline(msg.Payload())

		if err := sFunc(vin, online); err != nil {
			log.Fatalf("listener callback, %v\n", err)
		}
	}
}

func ReportListener(rFunc ReportListenerFunc, logging bool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if logging {
			util.LogMessage(msg)
		}

		vin := util.TopicVin(msg.Topic())

		result, err := report.New(msg.Payload()).Decode()
		if err != nil {
			log.Fatalf("cant decode, %v\n", err)
		}

		if err := rFunc(vin, result); err != nil {
			log.Fatalf("listener callback %v\n", err)
		}
	}
}

func parseOnline(b []byte) bool {
	return b[0] == '1'
}
