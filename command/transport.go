package command

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

// exec execute command and return the response.
func (c *Command) exec(cmd_name string, payload []byte) ([]byte, error) {
	cmder, err := getCmder(cmd_name)
	if err != nil {
		return nil, err
	}

	if err := c.sendCommand(cmder, payload); err != nil {
		return nil, err
	}

	msg, err := c.waitResponse(cmder)
	return msg, err
}

// sendCommand encode and send outgoing command.
func (c *Command) sendCommand(cmder *commander, payload []byte) error {
	packet, err := c.encode(cmder, payload)
	if err != nil {
		return err
	}

	// TODO: monitor outgoing command in-memory buffer
	// OnCommand[vin] = true
	c.transport.Pub(c.topic(shared.TOPIC_COMMAND), 1, true, packet)
	return nil
}

// waitResponse wait, decode and check of ACK and RESPONSE.
func (c *Command) waitResponse(cmder *commander) ([]byte, error) {
	defer func() {
		c.flush()
		time.Sleep(1 * time.Second)
	}()

	packet, err := c.waitPacket(DEFAULT_ACK_TIMEOUT)
	if err != nil {
		return nil, err
	}

	if err := checkAck(packet); err != nil {
		return nil, err
	}

	packet, err = c.waitPacket(cmder.timeout)
	if err != nil {
		return nil, err
	}

	res, err := c.decode(cmder, packet)
	if err != nil {
		return nil, err
	}

	if err := checkResponse(cmder, res); err != nil {
		return nil, err
	}

	return res.Message, nil
}

// waitPacket wait incomming packet to related VIN.
// It throws error on timeout.
func (c *Command) waitPacket(timeout time.Duration) ([]byte, error) {
	responses.reset(c.vin)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	data := make(chan []byte, 1)
	go func() {
		for {
			if rx, ok := responses.get(c.vin); ok {
				data <- rx
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
}

// flush clear command & response topic on broker.
// It indicates that command is done or cancelled.
func (c *Command) flush() {
	c.transport.Pub(c.topic(shared.TOPIC_COMMAND), 1, true, nil)
	c.transport.Pub(c.topic(shared.TOPIC_RESPONSE), 1, true, nil)
}

// topic subtitutes VIN into topic.
func (c *Command) topic(topic_pattern string) string {
	return strings.Replace(topic_pattern, "+", fmt.Sprint(c.vin), 1)
}
