package gen_vcu_sdk

import (
	cmd "github.com/pudjamansyurin/gen_vcu_sdk/command"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
)

type Sdk struct {
	transport *transport.Transport
	logging  bool
}

func New(host string, port int, user, pass string, logging bool) Sdk {
	tport := transport.New(transport.Config{
                           Host: host,
                           Port: port,
                           User: user,
                           Pass: pass,
        })
	return Sdk{
		transport: tport,
		logging: logging,
	}
}

func (s *Sdk) Connect() error {
	if err := s.transport.Connect(); err != nil {
		return err
	}

	if err := s.transport.Sub(shared.TOPIC_COMMAND, cmd.CommandListener); err != nil {
        	return err
        }

	if err := s.transport.Sub(shared.TOPIC_RESPONSE, cmd.ResponseListener); err != nil {
        	return err
        }

	return nil
}

func (s *Sdk) Disconnect() {
	s.transport.Disconnect()
}

func (s *Sdk) Listen(l Listener) error {
	if l.StatusFunc != nil {
		if err := s.transport.Sub(shared.TOPIC_STATUS, StatusListener(l.StatusFunc, s.logging)); err != nil {
			return err
		}
	}

	if l.ReportFunc != nil {
		if err := s.transport.Sub(shared.TOPIC_REPORT, ReportListener(l.ReportFunc, s.logging)); err != nil {
			return err
		}
	}

	return nil
}

func (s *Sdk) NewCommand(vin int) *cmd.Command {
	return cmd.New(vin, s.transport)
}
