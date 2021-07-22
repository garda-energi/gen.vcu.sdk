package gen_vcu_sdk

import (
	cmd "github.com/pudjamansyurin/gen_vcu_sdk/command"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
)

type Sdk struct {
	transport *transport.Transport
	logging   bool
}

// New create new instance of Sdk
func New(host string, port int, user, pass string, logging bool) Sdk {
	tport := transport.New(transport.Config{
		Host: host,
		Port: port,
		User: user,
		Pass: pass,
	})
	return Sdk{
		transport: tport,
		logging:   logging,
	}
}

// Connect open connection to mqtt broker
func (s *Sdk) Connect() error {
	return s.transport.Connect()
}

// Disconnect close connection to mqtt broker
func (s *Sdk) Disconnect() {
	s.transport.Disconnect()
}

// Listen subscribe to Status & Report topic (if callback is specified),
// it also auto subscribe to Command & Response topic
func (s *Sdk) Listen(l Listener) error {
	if l.StatusFunc != nil {
		if err := s.transport.Sub(shared.TOPIC_STATUS, 1, StatusListener(l.StatusFunc, s.logging)); err != nil {
			return err
		}
	}

	if l.ReportFunc != nil {
		if err := s.transport.Sub(shared.TOPIC_REPORT, 1, ReportListener(l.ReportFunc, s.logging)); err != nil {
			return err
		}
	}

	if err := s.transport.Sub(shared.TOPIC_COMMAND, 1, cmd.CommandListener); err != nil {
		return err
	}

	if err := s.transport.Sub(shared.TOPIC_RESPONSE, 1, cmd.ResponseListener); err != nil {
		return err
	}

	return nil
}

// NewCommand create new instance of Command
func (s *Sdk) NewCommand(vin int) *cmd.Command {
	return cmd.New(vin, s.transport)
}
