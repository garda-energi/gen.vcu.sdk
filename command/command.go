package command

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
	"log"
	"context"

	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
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

	packet := c.encode(cmder, nil)
	c.exec(packet)

	if err := c.waitResponse(cmder); err != nil {
		return "", err
	}

	return "", nil
}

func (c *Command) exec(packet []byte) {
	// OnCommand[vin] = true
	topic := strings.Replace(shared.TOPIC_COMMAND, "+", fmt.Sprint(c.vin), 1)
	c.transport.Pub(topic, 1, false, packet)
}

func (c *Command) waitResponse(cmder *Commander) error {
	// wait ack
	packet, err := c.waitPacket(5*time.Second);
	if err != nil {
		return err
	}
	// check ack
	ack := util.Reverse([]byte(shared.PREFIX_ACK))
	if !bytes.Equal(packet, ack) {
                return errors.New("ack corrupt")
	}

	// wait response
	packet, err = c.waitPacket(cmder.Timeout);
	if err != nil {
		return err
	}
	// decode response
	log.Println(packet)

	return nil
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
