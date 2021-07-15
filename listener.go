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
type ReportListenerFunc func(vin int, result interface{}) error

func (s *Sdk) statusListener(client mqtt.Client, msg mqtt.Message) {
	s.logPaylod(msg)

	if err := s.statusFunc(getVin(msg), isOnline(msg)); err != nil {
		log.Fatalf("Status callback error, %v\n", err)
	}
}

func (s *Sdk) reportListener(client mqtt.Client, msg mqtt.Message) {
	s.logPaylod(msg)

	rpt := report.New(msg.Payload())
	result, err := rpt.DecodeReport()
	if err != nil {
		log.Fatalf("Can't decode report package, %v\n", err)
	}

	if err := s.reportFunc(getVin(msg), result); err != nil {
		log.Fatalf("Report callback error, %v\n", err)
	}
}

func (s *Sdk) logPaylod(msg mqtt.Message) {
	if !s.logging {
		log.Printf("[%s] %s\n", msg.Topic(), util.HexString(msg.Payload()))
	}
}

func getVin(msg mqtt.Message) int {
	s := strings.Split(msg.Topic(), "/")
	vin, _ := strconv.Atoi(s[1])

	return vin
}

func isOnline(msg mqtt.Message) bool {
	return msg.Payload()[0] == '1'
}
