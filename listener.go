package gen_vcu_sdk

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type StatusListenerFunc func(vin int, online bool) error
type ReportListenerFunc func(vin int, report *report.ReportPacket) error

type Listener struct {
	StatusFunc StatusListenerFunc
	ReportFunc ReportListenerFunc
}

// StatusListener executed when got new packet on status topic.
func StatusListener(sFunc StatusListenerFunc, logging bool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if logging {
			shared.LogMessage(msg)
		}

		vin := shared.GetTopicVin(msg.Topic())
		online := parseOnline(msg.Payload())

		if err := sFunc(vin, online); err != nil {
			log.Fatalf("listener callback, %v\n", err)
		}
	}
}

// ReportListener executed when got new packet on report topic.
func ReportListener(rFunc ReportListenerFunc, logging bool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if logging {
			shared.LogMessage(msg)
		}

		vin := shared.GetTopicVin(msg.Topic())

		result, err := report.New(msg.Payload()).Decode()
		if err != nil {
			log.Fatalf("cant decode, %v\n", err)
		}

		if err := rFunc(vin, result); err != nil {
			log.Fatalf("listener callback %v\n", err)
		}
	}
}

// parseOnline convert status payload to online status.
func parseOnline(b []byte) bool {
	return b[0] == '1'
}
