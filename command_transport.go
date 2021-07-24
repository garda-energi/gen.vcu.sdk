package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

// exec execute command and return the response.
func (c *commander) exec(cmd_name string, payload []byte) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cmd, err := getCommand(cmd_name)
	if err != nil {
		return nil, err
	}

	if err := c.sendCommand(cmd, payload); err != nil {
		return nil, err
	}

	msg, err := c.waitResponse(cmd)
	return msg, err
}

// sendCommand encode and send outgoing command.
func (c *commander) sendCommand(cmd *command, payload []byte) error {
	packet, err := encodeCommand(c.vin, cmd, payload)
	if err != nil {
		return err
	}

	c.broker.pub(setTopicToVin(TOPIC_COMMAND, c.vin), 1, true, packet)
	return nil
}

// waitResponse wait, decode and check of ACK and RESPONSE packet.
func (c *commander) waitResponse(cmd *command) ([]byte, error) {
	defer func() {
		c.flush()
		time.Sleep(1 * time.Second)
	}()

	packet, err := c.waitPacket("ack", DEFAULT_ACK_TIMEOUT)
	if err != nil {
		return nil, err
	}

	if err := validateAck(packet); err != nil {
		return nil, err
	}

	packet, err = c.waitPacket("response", cmd.timeout)
	if err != nil {
		return nil, err
	}

	res, err := decodeResponse(packet)
	if err != nil {
		return nil, err
	}

	if err := validateResponse(c.vin, cmd, res); err != nil {
		return nil, err
	}

	return res.Message, nil
}

// waitPacket wait incomming packet for current VIN.
// It throws error on timeout.
func (c *commander) waitPacket(name string, timeout time.Duration) ([]byte, error) {
	// flush channel
	for len(c.resChan) > 0 {
		<-c.resChan
	}

	select {
	case data := <-c.resChan:
		return data, nil
	case <-time.After(timeout):
		return nil, errPacketTimeout(name)
	}

}

// validateAck validate incomming ack packet.
func validateAck(msg []byte) error {
	ack := strToBytes(PREFIX_ACK)
	if !bytes.Equal(msg, ack) {
		return errPacketAckCorrupt
	}
	return nil
}

// validateResponse validate incomming response packet.
// It also parse response code and message
func validateResponse(vin int, cmd *command, res *responsePacket) error {
	if !res.validPrefix() {
		return errInvalidPrefix
	}
	if !res.validSize() {
		return errInvalidSize
	}
	if int(res.Header.Vin) != vin {
		return errInvalidVin
	}
	if !res.matchWith(cmd) {
		return errInvalidCode
	}
	if !res.validResCode() {
		return errInvalidResCode
	}
	if res.messageOverflow() {
		return errResMessageOverflow
	}

	// check resCode
	if res.Header.ResCode == resCodeOk {
		return nil
	}

	out := fmt.Sprint(res.Header.ResCode)
	// check if message is not empty
	if res.hasMessage() {
		res.renderMessage()
		out += fmt.Sprintf(", %s", res.Message)
	}
	return errors.New(out)

}
