package gen_vcu_sdk

import (
	"log"

	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

const (
	TOPIC_STATUS = "VCU/+/STS"
	TOPIC_REPORT = "VCU/+/RPT"
)

type StatusListenerFunc func(vin int, online bool) error
type ReportListenerFunc func(vin int, report interface{}) error

type Sdk struct {
	config     transport.ClientConfig
	logging    bool
	statusFunc StatusListenerFunc
	reportFunc ReportListenerFunc
}

func New(host string, port int, user, pass string) Sdk {
	return Sdk{
		config: transport.ClientConfig{
			Host:     host,
			Port:     port,
			Username: user,
			Password: pass,
			// ClientId: "go_mqtt_client",
		},
		logging: true,
	}
}

func (s *Sdk) ConnectAndListen() {
	t := transport.New(s.config)

	if err := t.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %v\n", err)
	}

	if err := t.Subscribe(TOPIC_STATUS, s.statusListener); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	if err := t.Subscribe(TOPIC_REPORT, s.reportListener); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	util.WaitForCtrlC()
	t.Disconnect()
}

func (s *Sdk) AddReportListener(cb ReportListenerFunc) {
	s.reportFunc = cb
}

func (s *Sdk) AddStatusListener(cb StatusListenerFunc) {
	s.statusFunc = cb
}

func (s *Sdk) Logging(enable bool) {
	s.logging = enable
}
