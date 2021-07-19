package command

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

type Command struct {
	transport *transport.Transport
	cmd    *CommandList
}

func New(tr *transport.Transport) *Command {
	return &Command{
		transport: tr,
		cmd      : NewCommandList(),
	}
}

// GenInfo Gather device information
func (c *Command) GenInfo(vin int) (string, error) {
	cmd := c.cmd.GEN.INFO

	packet := c.encode(vin, cmd, nil)
	c.publish(vin, packet)

	if err := waitAck(vin, cmd); err != nil {
		return "", err
	}

	res, err := waitResponse(vin, cmd)
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (c *Command) publish(vin int, packet []byte) {
	// OnCommand[vin] = true
	c.transport.Pub(cmdTopic(vin), 1, false, packet)
}

func waitAck(vin int, cmd Commander) error {
	RX.Reset(vin)

	ack := util.Reverse([]byte(shared.PREFIX_ACK))
	done := make(chan struct{})
	go func() {
		for {
			if rx, ok := RX.Get(vin); ok {
				if bytes.Equal(rx, ack) {
					done <- struct{}{}
				}
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(5 * time.Second):
	}

	return errors.New("ack timeout")
}

func waitResponse(vin int, cmd Commander) ([]byte, error) {
	RX.Reset(vin)

	done := make(chan []byte)
	go func() {
		for {
			if rx, ok := RX.Get(vin); ok {
				done <- rx
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	select {
	case res := <-done:
		return res, nil
	case <-time.After(getTimeout(cmd)):
	}

	return nil, errors.New("response timeout")
}

func getTimeout(cmd Commander) time.Duration {
	timeout := cmd.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return timeout
}

func cmdTopic(vin int) string {
	return strings.Replace(shared.TOPIC_COMMAND, "+", fmt.Sprint(vin), 1)
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
