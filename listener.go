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
	logging    bool
	statusFunc StatusListenerFunc
	reportFunc ReportListenerFunc
}

func (l *Listener) berforeHook(msg mqtt.Message) int {
	if l.logging {
		util.LogMessage(msg)
	}
	vin := util.TopicVin(msg.Topic())
	return vin
}

func (l *Listener) status(client mqtt.Client, msg mqtt.Message) {
	vin := l.berforeHook(msg)
	online := parseOnline(msg.Payload())

	if l.statusFunc != nil {
		if err := l.statusFunc(vin, online); err != nil {
			log.Fatalf("Status listener error, %v\n", err)
		}
	}
}

func (l *Listener) report(client mqtt.Client, msg mqtt.Message) {
	vin := l.berforeHook(msg)

	result, err := report.New(msg.Payload()).Decode()
	if err != nil {
		log.Fatalf("Can't decode report package, %v\n", err)
	}

	if l.reportFunc != nil {
		if err := l.reportFunc(vin, result); err != nil {
			log.Fatalf("Report listener error, %v\n", err)
		}
	}
}

func parseOnline(b []byte) bool {
	return b[0] == '1'
}
