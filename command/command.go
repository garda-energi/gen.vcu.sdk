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

var CMD_LIST = NewCommandList()

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
	cmd := CMD_LIST.GEN.INFO

	packet := c.encode(cmd, nil)
	c.exec(packet)

	if err := c.waitResponse(cmd); err != nil {
		return "", err
	}

	return "", nil
}

func (c *Command) exec(packet []byte) {
	// OnCommand[vin] = true
	topic := strings.Replace(shared.TOPIC_COMMAND, "+", fmt.Sprint(c.vin), 1)
	c.transport.Pub(topic, 1, false, packet)
}

func (c *Command) waitResponse(cmd Commander) error {
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
	timeout := cmd.Timeout
	if timeout == 0 {
		timeout = 5*time.Second
	}
	packet, err = c.waitPacket(timeout);
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

// func getCmdPacket(code CMD_CODE, subCode CMD_SUBCODE) (CmdPacket, error) {
// 	// for _, cmd := range CMD_LIST {
// 	// 	if cmd.Code == code && cmd.SubCode == subCode {
// 	// 		return cmd, nil
// 	// 	}
// 	// }

// 	cmd, ok := CmdList[code][subCode]
// 	if ok {
// 		return cmd, nil
// 	}
// 	return CmdPacket{}, errors.New("command code not found")
// }
