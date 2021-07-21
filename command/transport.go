package command

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

func (c *Command) sendCommand(cmder *commander, payload []byte) error {
	packet, err := c.encode(cmder, payload)
	if err != nil {
		return err
	}

	// OnCommand[vin] = true
	c.transport.Pub(c.topic(shared.TOPIC_COMMAND), 1, true, packet)
	return nil
}

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

func (c *Command) flush() {
	c.transport.Pub(c.topic(shared.TOPIC_COMMAND), 1, true, nil)
	c.transport.Pub(c.topic(shared.TOPIC_RESPONSE), 1, true, nil)
}

func (c *Command) topic(topic_pattern string) string {
	return strings.Replace(topic_pattern, "+", fmt.Sprint(c.vin), 1)
}
