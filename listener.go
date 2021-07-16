package gen_vcu_sdk

import (
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

type DATA_TYPE uint8

const (
	DATA_TYPE_STRUCT DATA_TYPE = iota
	DATA_TYPE_LIST
)

type StatusListenerFunc func(vin int, online bool) error
type ReportListenerFunc func(vin int, result interface{}) error

type Listener struct {
	logging    bool
	dtype      DATA_TYPE
	statusFunc StatusListenerFunc
	reportFunc ReportListenerFunc
}

func (l *Listener) status(client mqtt.Client, msg mqtt.Message) {
	l.logPaylod(msg)

	vin := parseVin(msg.Topic())
	online := isOnline(msg.Payload())
	if err := l.statusFunc(vin, online); err != nil {
		log.Fatalf("Status listener error, %v\n", err)
	}
}

func (l *Listener) report(client mqtt.Client, msg mqtt.Message) {
	l.logPaylod(msg)

	var err error
	var result interface{}

	rpt := report.New(msg.Payload())
	if l.dtype == DATA_TYPE_STRUCT {
		result, err = rpt.DecodeReportStruct()
	} else {
		result, err = rpt.DecodeReportList()
	}
	if err != nil {
		log.Fatalf("Can't decode report package, %v\n", err)
	}

	vin := parseVin(msg.Topic())
	if err := l.reportFunc(vin, result); err != nil {
		log.Fatalf("Report listener error, %v\n", err)
	}
}

func (l *Listener) logPaylod(msg mqtt.Message) {
	if l.logging {
		log.Printf("[%s] %s\n", msg.Topic(), util.HexString(msg.Payload()))
	}
}

func parseVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])

	return vin
}

func isOnline(b []byte) bool {
	return b[0] == '1'
}
