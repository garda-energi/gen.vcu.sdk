package command

import (
	"errors"
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
	c.transport.Pub(shared.SetTopicToVin(shared.TOPIC_COMMAND, c.vin), 1, true, packet)
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

	checkResponseTicker := time.NewTicker(10 * time.Millisecond)
	defer checkResponseTicker.Stop()

	dataChan := make(chan []byte, 1)
	go func() {
		for range checkResponseTicker.C {
			if data, ok := responses.get(c.vin); ok {
				dataChan <- data
				break
			}
		}
	}()

	select {
	case data := <-dataChan:
		return data, nil
	case <-time.After(timeout):
		return nil, errors.New("packet timeout")
	}
}

// flush clear command & response topic on broker.
// It indicates that command is done or cancelled.
func (c *Command) flush() {
	c.transport.Pub(shared.SetTopicToVin(shared.TOPIC_COMMAND, c.vin), 1, true, nil)
	c.transport.Pub(shared.SetTopicToVin(shared.TOPIC_RESPONSE, c.vin), 1, true, nil)
}
