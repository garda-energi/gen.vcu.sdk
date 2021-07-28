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

// status is executed when received new packet on status topic.
func (ls *Listener) status() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		ls.logger.Println(RPT, debugPacket(msg))
		vin := getTopicVin(msg.Topic())
		online := online(msg.Payload())

		ls.StatusFunc(vin, online)
	}
}

// report is executed when received new packet on report topic.
func (ls *Listener) report() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		ls.logger.Println(RPT, debugPacket(msg))

		vin := getTopicVin(msg.Topic())

		result, err := decodeReport(msg.Payload())
		if err != nil {
			log.Fatalln("cant decode", err)
		}

		ls.ReportFunc(vin, result)
	}
}

// online convert status payload to online status.
func online(b []byte) bool {
	return b[0] == '1'
}
