package sdk

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type statusListener func(vin int, online bool)
type reportListener func(vin int, report *ReportPacket)

// Listener store status & report callback function
type Listener struct {
	StatusFunc statusListener
	ReportFunc reportListener
	logger     *log.Logger
}

// status is executed when got new packet on status topic.
func (ls *Listener) status() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		ls.logger.Println(debugPacket(msg))
		vin := getTopicVin(msg.Topic())
		online := parseOnline(msg.Payload())

		ls.StatusFunc(vin, online)
	}
}

// report is executed when got new packet on report topic.
func (ls *Listener) report() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		ls.logger.Println(debugPacket(msg))

		vin := getTopicVin(msg.Topic())

		result, err := newReport(msg.Payload()).decode()
		if err != nil {
			log.Fatalf("cant decode, %v\n", err)
		}

		ls.ReportFunc(vin, result)
	}
}

// parseOnline convert status payload to online status.
func parseOnline(b []byte) bool {
	return b[0] == '1'
}
