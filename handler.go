package gen_vcu_sdk

import (
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func (s *Sdk) statusHandler(client mqtt.Client, msg mqtt.Message) {
	s.logPaylod(msg)

	online := msg.Payload()[0] == '1'

	if err := s.statusListener(getVin(msg.Topic()), online); err != nil {
		log.Fatalf("Status callback error, %v\n", err)
	}
}

func (s *Sdk) reportHandler(client mqtt.Client, msg mqtt.Message) {
	s.logPaylod(msg)

	packet := &Report{Bytes: msg.Payload()}
	report, err := packet.decodeReport()
	if err != nil {
		log.Fatalf("Can't decode report package, %v\n", err)
	}

	if err := s.reportListener(getVin(msg.Topic()), report); err != nil {
		log.Fatalf("Report callback error, %v\n", err)
	}
}

func (s *Sdk) logPaylod(msg mqtt.Message) {
	if s.logging {
		log.Printf("[%s] %s\n", msg.Topic(), util.HexString(msg.Payload()))
	}
}

func getVin(topic string) int {
	s := strings.Split(topic, "/")
	vin, _ := strconv.Atoi(s[1])

	return vin
}
