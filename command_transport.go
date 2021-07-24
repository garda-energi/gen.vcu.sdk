package sdk

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
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
	packet, err := c.encode(cmd, payload)
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

	packet, err := c.waitPacket(DEFAULT_ACK_TIMEOUT)
	if err != nil {
		return nil, err
	}

	if err := validateAck(packet); err != nil {
		return nil, err
	}

	packet, err = c.waitPacket(cmd.timeout)
	if err != nil {
		return nil, err
	}

	res, err := c.decode(cmd, packet)
	if err != nil {
		return nil, err
	}

	if err := validateResponse(cmd, res); err != nil {
		return nil, err
	}

	return res.Message, nil
}

// waitPacket wait incomming packet for current VIN.
// It throws error on timeout.
func (c *commander) waitPacket(timeout time.Duration) ([]byte, error) {
	// flush channel
	for len(c.resChan) > 0 {
		<-c.resChan
	}

	select {
	case data := <-c.resChan:
		return data, nil
	case <-time.After(timeout):
		return nil, errors.New("packet timeout")
	}
}

// flush clear command & response topic on broker.
// It indicates that command is done or cancelled.
func (c *commander) flush() {
	c.broker.pub(setTopicToVin(TOPIC_COMMAND, c.vin), 1, true, nil)
	c.broker.pub(setTopicToVin(TOPIC_RESPONSE, c.vin), 1, true, nil)
}

// validateAck validate incomming ack packet.
func validateAck(msg []byte) error {
	ack := strToBytes(PREFIX_ACK)
	if !bytes.Equal(msg, ack) {
		return errors.New("ack corrupt")
	}
	return nil
}

// validateResponse validate incomming response packet.
// It also parse response code and message
func validateResponse(cmd *command, res *ResponsePacket) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprint(cmd.name))
	sb.WriteString(": ")

	// check code
	if res.Header.Code != cmd.code || res.Header.SubCode != cmd.sub_code {
		sb.WriteString("response-mismatch")
		return errors.New(sb.String())
	}

	// check resCode
	resCode := &res.Header.ResCode
	if *resCode == resCodeOk {
		return nil
	}

	sb.WriteString(fmt.Sprint(*resCode))

	// check if message is empty
	if len(res.Message) > 0 {
		// subtitutes BIKE_STATE to message
		str := string(res.Message)
		for i := BikeStateUnknown; i < BikeStateLimit; i++ {
			old := fmt.Sprintf("{%d}", i)
			new := BikeState(i).String()
			str = strings.ReplaceAll(str, old, new)
		}
		res.Message = []byte(str)

		sb.WriteString(", ")
		sb.WriteString(string(res.Message))
	}

	return errors.New(sb.String())

}
