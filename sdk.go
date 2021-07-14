package gen_vcu_sdk

import (
	"log"

	// ttt "github.com/eclipse/paho.transport.golang"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

const (
	TOPIC_REPORT = "VCU/+/RPT"
)

type CallbackFunc func(interface{}, []byte)

type Sdk struct {
	config         transport.ClientConfig
	reportCallback CallbackFunc
}

func New(host string, port int, user, pass string) Sdk {
	return Sdk{
		config: transport.ClientConfig{
			Host:     host,
			Port:     port,
			Username: user,
			Password: pass,
			ClientId: "go_mqtt_client",
		},
	}
}

func (s *Sdk) AddReportListener(cb CallbackFunc) {
	s.reportCallback = cb
}

func (s *Sdk) ConnectAndListen() {
	t := transport.New(s.config)

	if err := t.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %s\n", err.Error())
	}

	if err := t.Subscribe(TOPIC_REPORT, s.reportHandler); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %s\n", err.Error())
	}

	util.WaitForCtrlC()
	t.Disconnect()
}

func (s *Sdk) reportHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[REPORT] %s\n", util.HexString(msg.Payload()))

	packet := &Report{Bytes: msg.Payload()}
	report, err := packet.decodeReport()
	if err != nil {
		log.Fatal(err)
	}

	// util.Debug(report)
	s.reportCallback(report, msg.Payload())
}
