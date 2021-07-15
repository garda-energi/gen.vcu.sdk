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
	transport transport.Transport
	listener  Listener
}

func New(host string, port int, user, pass string) Sdk {
	return Sdk{
		transport: transport.New(host, port, user, pass),
		listener: Listener{
			logging: true,
		},
	}
}

func (s *Sdk) ConnectAndListen() {
	if err := s.transport.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %v\n", err)
	}

	if err := s.transport.Subscribe(TOPIC_STATUS, s.listener.status); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	if err := s.transport.Subscribe(TOPIC_REPORT, s.listener.report); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	util.WaitForCtrlC()
	s.transport.Disconnect()
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
