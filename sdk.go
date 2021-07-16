package gen_vcu_sdk

import (
	"log"

	"github.com/pudjamansyurin/gen_vcu_sdk/command"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

const (
	TOPIC_STATUS   = "VCU/+/STS"
	TOPIC_REPORT   = "VCU/+/RPT"
	TOPIC_COMMAND  = "VCU/+/CMD"
	TOPIC_RESPONSE = "VCU/+/RSP"
)

type Sdk struct {
	Cmd             *command.Command
	transportConfig transport.TransportConfig
	listener        Listener
}

func New(host string, port int, user, pass string) Sdk {
	return Sdk{
		transportConfig: transport.TransportConfig{
			Host: host,
			Port: port,
			User: user,
			Pass: pass,
		},
		listener: Listener{
			logging: true,
		},
	}
}

func (s *Sdk) ConnectAndListen() {
	tr := transport.New(s.transportConfig)
	if err := tr.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %v\n", err)
	}

	s.Cmd = command.New(&tr.Client)

	if err := tr.Subscribe(TOPIC_STATUS, s.listener.status); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	if err := tr.Subscribe(TOPIC_REPORT, s.listener.report); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	util.WaitForCtrlC()
	tr.Disconnect()
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
