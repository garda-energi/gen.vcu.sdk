package command

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"context"

	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type Command struct {
	vin int
	transport *transport.Transport
}

func New(vin int, tport *transport.Transport) *Command {
	return &Command{
		vin: vin,
		transport: tport,
	}
}

// GenInfo Gather device information
func (c *Command) GenInfo() (string, error) {
	cmder, err := getCmder("GEN_INFO")
	if err != nil {
		return "", err
	}

	if err := c.sendCommand(cmder, nil); err != nil {
		return "", err
	}

	msg, err := c.waitResponse(cmder)
	if err != nil {
		return "", err
	}

	return string(msg), nil
}

func (c *Command) sendCommand(cmder *Commander, payload []byte) error {
	packet, err := c.encode(cmder, payload)
	if err != nil {
		return err
	}

	// OnCommand[vin] = true
	topic := strings.Replace(shared.TOPIC_COMMAND, "+", fmt.Sprint(c.vin), 1)
	c.transport.Pub(topic, 1, false, packet)

	return nil
}

func (c *Command) waitResponse(cmder *Commander) ([]byte, error) {
	// wait ack
	packet, err := c.waitPacket(5*time.Second);
	if err != nil {
		return nil, err
	}
	// check ack
	if !validAck(packet) {
                return nil, errors.New("ack corrupt")
	}

	// wait response
	packet, err = c.waitPacket(cmder.Timeout);
	if err != nil {
		return nil, err
	}
	// decode & check response

	return packet, err
}

func (c *Command) waitPacket(timeout time.Duration) ([]byte, error) {
	RX.Reset(c.vin)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	data := make(chan []byte, 1)
	go func() {
	   for {
		if rx, ok := RX.Get(c.vin); ok {
                	data<- rx
                        return
                }

		select {
		case <-ctx.Done():
			return
		default:
                        time.Sleep(10 * time.Millisecond)
		}
	   }
	}()

	select {
	case <-ctx.Done():
		return nil, errors.New("packet timeout")
	case dat := <-data:
		return dat, nil
	}

	return nil, errors.New("packet error")
}

func getCmder(name string) (*Commander, error) {
	for code, subCodes := range CMDERS {
		for subCode, cmder := range subCodes {
			if cmder.Name == name {
				cmder.Code = uint8(code)
				cmder.SubCode = uint8(subCode)

 				if cmder.Timeout == 0 {
 					cmder.Timeout = 5*time.Second
 				}

				return &cmder, nil
			}
		}
	}

	return nil, errors.New("command invalid")
}
