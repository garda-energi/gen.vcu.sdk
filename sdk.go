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

type Sdk struct {
	config   transport.ClientConfig
	listener Listener
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
		listener: Listener{
			logging: true,
		},
	}
}

func (s *Sdk) ConnectAndListen() {
	t := transport.New(s.config)

	if err := t.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %v\n", err)
	}

	if err := t.Subscribe(TOPIC_STATUS, s.listener.status); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	if err := t.Subscribe(TOPIC_REPORT, s.listener.report); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	util.WaitForCtrlC()
	t.Disconnect()
}

func (s *Sdk) AddStatusListener(cb StatusListenerFunc) {
	s.listener.statusFunc = cb
}

func (s *Sdk) AddReportListener(cb ReportListenerFunc) {
	s.listener.reportFunc = cb
}

func (s *Sdk) Logging(enable bool) {
	s.listener.logging = enable
}
