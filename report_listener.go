package sdk

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StatusListenerFunc func(vin int, online bool) error
type ReportListenerFunc func(vin int, report *ReportPacket) error

// Listener store status & report callback function
type Listener struct {
	StatusFunc StatusListenerFunc
	ReportFunc ReportListenerFunc
}

// statusListener is executed when got new packet on status topic.
func statusListener(sFunc StatusListenerFunc, logging bool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if logging {
			logPacket(msg)
		}

		vin := getTopicVin(msg.Topic())
		online := parseOnline(msg.Payload())

		if err := sFunc(vin, online); err != nil {
			log.Fatalf("listener callback, %v\n", err)
		}
	}
}

// reportListener is executed when got new packet on report topic.
func reportListener(rFunc ReportListenerFunc, logging bool) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		if logging {
			logPacket(msg)
		}

		vin := getTopicVin(msg.Topic())

		result, err := newReport(msg.Payload()).decode()
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
