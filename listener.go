package gen_vcu_sdk

import (
	"log"
	"strconv"
	"strings"

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

func (l *Listener) berforHook(msg mqtt.Message) int {
	if l.logging {
		logPaylod(msg)
	}
	vin := parseVin(msg.Topic())
	return vin
}

func (l *Listener) status(client mqtt.Client, msg mqtt.Message) {
	vin := l.berforHook(msg)
	online := parseOnline(msg.Payload())

	if l.statusFunc != nil {
		if err := l.statusFunc(vin, online); err != nil {
			log.Fatalf("Status listener error, %v\n", err)
		}
	}
}

func (l *Listener) report(client mqtt.Client, msg mqtt.Message) {
	vin := l.berforHook(msg)

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

func logPaylod(msg mqtt.Message) {
	log.Printf("[%s] %s\n", msg.Topic(), util.HexString(msg.Payload()))
}

func parseVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])

	return vin
}

func parseOnline(b []byte) bool {
	return b[0] == '1'
}
