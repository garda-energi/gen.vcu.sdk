package command

import (
	"errors"
	"strconv"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
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

// GenOdo Set odometer value in km
func (c *Command) GenOdo(km uint16) error {
	_, err := c.exec("GEN_ODO", makeU16(km))
	return err
}

// GenAntiTheaf Toggle anti-thief motion detector
func (c *Command) GenAntiTheaf() error {
	_, err := c.exec("GEN_ANTI_THIEF", nil)
	return err
}

// GenReportFlush Flush report buffer
func (c *Command) GenReportFlush() error {
	_, err := c.exec("GEN_RPT_FLUSH", nil)
	return err
}

// GenReportBlock Block reporting mode
func (c *Command) GenReportBlock(on bool) error {
	_, err := c.exec("GEN_RPT_BLOCK", makeBool(on))
	return err
}

// OvdState Override bike state
func (c *Command) OvdState(state shared.BIKE_STATE) error {
	min, max := shared.BIKE_STATE_NORMAL, shared.BIKE_STATE_RUN
	if state < min || state > max {
		return errors.New("state out of range")
	}

	payload := []byte{byte(state)}
	_, err := c.exec("OVD_STATE", payload)
	return err
}

// OvdReportInterval Override reporting interval in seconds
func (c *Command) OvdReportInterval(dur time.Duration) error {
	min, max := time.Duration(5), time.Duration(^uint16(0))
	if dur < min*time.Second || dur > max*time.Second {
		return errors.New("duration out of range")
	}

	payload := makeU16(uint16(dur.Seconds()))
	_, err := c.exec("OVD_RPT_INTERVAL", payload)
	return err
}

// OvdReportFrame Override report frame type
func (c *Command) OvdReportFrame(frame shared.FRAME_ID) error {
	min, max := shared.FRAME_ID_SIMPLE, shared.FRAME_ID_FULL
	if frame < min || frame > max {
		return errors.New("frame out of range")
	}

	payload := []byte{byte(frame)}
	_, err := c.exec("OVD_RPT_FRAME", payload)
	return err
}

// OvdRemoteSeat Override remote seat fob
func (c *Command) OvdRemoteSeat() error {
	_, err := c.exec("OVD_RMT_SEAT", nil)
	return err
}

// OvdRemoteAlarm Override remote alarm fob
func (c *Command) OvdRemoteAlarm() error {
	_, err := c.exec("OVD_RMT_ALARM", nil)
	return err
}

// AudioBeep Beep the digital audio module
func (c *Command) AudioBeep() error {
	_, err := c.exec("AUDIO_BEEP", nil)
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

// FingerAdd Add a new fingerprint id
func (c *Command) FingerAdd() (int, error) {
	msg, err := c.exec("FINGER_ADD", nil)
	if err != nil {
		return 0, err
	}

	// decode
	id, _ := strconv.Atoi(string(msg[0]))
	return id, nil
}

// FingerDel Delete a fingerprint id
func (c *Command) FingerDel(id int) error {
	min, max := 1, FINGERPRINT_MAX
	if id < min || id > max {
		return errors.New("id out of range")
	}

	_, err := c.exec("FINGER_DEL", nil)
	return err
}

// FingerRst Reset all fingerprint ids
func (c *Command) FingerRst() error {
	_, err := c.exec("FINGER_RST", nil)
	return err
}

// RemotePairing Enter keyless pairing mode
func (c *Command) RemotePairing() error {
	_, err := c.exec("REMOTE_PAIRING", nil)
	return err
}

// FotaVcu Upgrade VCU firmware over the air
func (c *Command) FotaVcu() (string, error) {
	msg, err := c.exec("FOTA_VCU", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// FotaHmi Upgrade HMI firmware over the air
func (c *Command) FotaHmi() (string, error) {
	msg, err := c.exec("FOTA_HMI", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
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
