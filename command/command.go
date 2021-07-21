package command

import (
	"strconv"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
)

type Command struct {
	vin       int
	transport *transport.Transport
}

func New(vin int, tport *transport.Transport) *Command {
	return &Command{
		vin:       vin,
		transport: tport,
	}
}

// GenInfo Gather device information
func (c *Command) GenInfo() (string, error) {
	msg, err := c.exec("GEN_INFO", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// GenLed Set buil-in led on board
func (c *Command) GenLed(on bool) error {
	_, err := c.exec("GEN_LED", makeBool(on))
	return err
}

// GenRtc Set real time clock on board
func (c *Command) GenRtc(time time.Time) error {
	_, err := c.exec("GEN_RTC", makeTime(time))
	return err
}

// FingerFetch Get all registered fingerprint ids
func (c *Command) FingerFetch() ([]int, error) {
	msg, err := c.exec("FINGER_FETCH", nil)
	if err != nil {
		return nil, err
	}

	// decode
	ids := make([]int, len(msg))
	for i := range ids {
		id, _ := strconv.Atoi(string(msg[i]))
		ids[i] = id
	}

	return ids, nil
}

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
