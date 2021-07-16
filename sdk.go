package gen_vcu_sdk

import (
	"log"

	"github.com/pudjamansyurin/gen_vcu_sdk/command"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

type Sdk struct {
	Command  *command.Command
	config   transport.Config
	listener Listener
}

func New(host string, port int, user, pass string) Sdk {
	return Sdk{
		config: transport.Config{
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
	tr := transport.New(s.config)
	if err := tr.Connect(); err != nil {
		log.Fatalf("[MQTT] Failed to connect, %v\n", err)
	}

	s.Command = command.New(tr.Client)

	if err := tr.Subscribe(shared.TOPIC_COMMAND, s.Command.CommandListener); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	if err := tr.Subscribe(shared.TOPIC_RESPONSE, s.Command.ResponseListener); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	if err := tr.Subscribe(shared.TOPIC_STATUS, s.listener.status); err != nil {
		log.Fatalf("[MQTT] Failed to subscribe, %v\n", err)
	}

	if err := tr.Subscribe(shared.TOPIC_REPORT, s.listener.report); err != nil {
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
