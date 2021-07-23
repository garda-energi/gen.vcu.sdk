package command

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

// exec execute command and return the response.
func (c *Command) exec(cmd_name string, payload []byte) ([]byte, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	cmder, err := getCommander(cmd_name)
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

	c.broker.Pub(shared.SetTopicToVin(shared.TOPIC_COMMAND, c.vin), 1, true, packet)
	return nil
}

// waitResponse wait, decode and check of ACK and RESPONSE packet.
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
func (c *Command) flush() {
	c.broker.Pub(shared.SetTopicToVin(shared.TOPIC_COMMAND, c.vin), 1, true, nil)
	c.broker.Pub(shared.SetTopicToVin(shared.TOPIC_RESPONSE, c.vin), 1, true, nil)
}

// checkAck validate incomming ack packet.
func checkAck(msg []byte) error {
	ack := shared.StrToBytes(shared.PREFIX_ACK)
	if !bytes.Equal(msg, ack) {
		return errors.New("ack corrupt")
	}
	return nil
}

// checkResponse validate incomming response packet.
// It also parse response code and message
func checkResponse(cmder *commander, res *ResponsePacket) error {
	// check code
	if res.Header.Code != cmder.code || res.Header.SubCode != cmder.sub_code {
		return errors.New("response-mismatch")
	}

	// check resCode
	resCode := &res.Header.ResCode
	if *resCode == RES_CODE_OK {
		return nil
	}

	// check if message is empty
	if len(res.Message) == 0 {
		return fmt.Errorf("%s", *resCode)
	}

	// subtitutes BIKE_STATE to message
	str := string(res.Message)
	for i := shared.BIKE_STATE_UNKNOWN; i < shared.BIKE_STATE_limit; i++ {
		old := fmt.Sprintf("{%d}", i)
		new := shared.BIKE_STATE(i).String()
		str = strings.ReplaceAll(str, old, new)
	}
	res.Message = []byte(str)

	return fmt.Errorf("%s, %s", *resCode, res.Message)
}
