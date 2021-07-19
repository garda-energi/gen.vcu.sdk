package gen_vcu_sdk

import (
	"log"
	"github.com/pudjamansyurin/gen_vcu_sdk/report"
	"github.com/pudjamansyurin/gen_vcu_sdk/command"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Sdk struct {
	command  *command.Command
	transport *transport.Transport
	logging  bool
}

func New(host string, port int, user, pass string, logging bool) Sdk {
	tr := transport.New(transport.Config{
                           Host: host,
                           Port: port,
                           User: user,
                           Pass: pass,
        })
	return Sdk{
		transport: tr,
		logging: logging,
	}
}

func (s *Sdk) Connect() error {
	if err := s.transport.Connect(); err != nil {
		return err
	}

	return nil
}

func (s *Sdk) Listen(listener Listener) error {
	if listener.StatusFunc != nil {
		err := s.transport.Subscribe(shared.TOPIC_STATUS, func (client mqtt.Client, msg mqtt.Message) {
  			vin := listenerHook(msg, s.logging)
        		online := parseOnline(msg.Payload())

        		if err := listener.StatusFunc(vin, online); err != nil {
				log.Fatalf("listener callback, %v\n", err)
			}
		})
		if err != nil {
			return err
		}
	}

	if listener.ReportFunc != nil {
		err := s.transport.Subscribe(shared.TOPIC_REPORT, func (client mqtt.Client, msg mqtt.Message) {
        		vin := listenerHook(msg, s.logging)

        		result, err := report.New(msg.Payload()).Decode()
        		if err != nil {
                		log.Fatalf("cant decode, %v\n", err)
        		}

                	if err := listener.ReportFunc(vin, result); err != nil {
				log.Fatalf("listener callback %v\n", err)
			}
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Sdk) NewCommand() error {
	s.command = command.New(s.transport.Client)

        if err := s.transport.Subscribe(shared.TOPIC_COMMAND, s.command.CommandListener); err != nil {
                return err
        }

        if err := s.transport.Subscribe(shared.TOPIC_RESPONSE, s.command.ResponseListener); err != nil {
                return err
        }

	return nil
}
