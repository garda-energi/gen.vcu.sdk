package command

import (
	"errors"
	"reflect"
	"strconv"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/transport"
)

type Command struct {
	vin       int
	transport *transport.Transport
}

// New create new Command instance.
func New(vin int, tport *transport.Transport) *Command {
	return &Command{
		vin:       vin,
		transport: tport,
	}
}

// GenInfo gather device information.
func (c *Command) GenInfo() (string, error) {
	msg, err := c.exec("GEN_INFO", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// GenLed set built-in led state on board.
func (c *Command) GenLed(on bool) error {
	_, err := c.exec("GEN_LED", boolToBytes(on))
	return err
}

// GenRtc set real time clock on board.
func (c *Command) GenRtc(time time.Time) error {
	_, err := c.exec("GEN_RTC", shared.TimeToBytes(time))
	return err
}

// GenOdo set odometer value (in km).
func (c *Command) GenOdo(km uint16) error {
	payload := shared.UintToBytes(reflect.Uint16, uint64(km))
	_, err := c.exec("GEN_ODO", payload)
	return err
}

// GenAntiTheaf toggle anti-thief motion detector.
func (c *Command) GenAntiTheaf() error {
	_, err := c.exec("GEN_ANTI_THIEF", nil)
	return err
}

// GenReportFlush flush report buffer.
func (c *Command) GenReportFlush() error {
	_, err := c.exec("GEN_RPT_FLUSH", nil)
	return err
}

// GenReportBlock block reporting mode.
func (c *Command) GenReportBlock(on bool) error {
	_, err := c.exec("GEN_RPT_BLOCK", boolToBytes(on))
	return err
}

// OvdState override bike state.
func (c *Command) OvdState(state shared.BIKE_STATE) error {
	min, max := shared.BIKE_STATE_NORMAL, shared.BIKE_STATE_RUN
	if state < min || state > max {
		return errors.New("state out of range")
	}

	payload := []byte{byte(state)}
	_, err := c.exec("OVD_STATE", payload)
	return err
}

// OvdReportInterval override reporting interval (in seconds).
func (c *Command) OvdReportInterval(dur time.Duration) error {
	min, max := time.Duration(5), time.Duration(^uint16(0))
	if dur < min*time.Second || dur > max*time.Second {
		return errors.New("duration out of range")
	}
	payload := shared.UintToBytes(reflect.Uint16, uint64(dur.Seconds()))
	_, err := c.exec("OVD_RPT_INTERVAL", payload)
	return err
}

// OvdReportFrame override report frame type.
func (c *Command) OvdReportFrame(frame shared.FRAME_ID) error {
	min, max := shared.FRAME_ID_SIMPLE, shared.FRAME_ID_FULL
	if frame < min || frame > max {
		return errors.New("frame out of range")
	}

	payload := []byte{byte(frame)}
	_, err := c.exec("OVD_RPT_FRAME", payload)
	return err
}

// OvdRemoteSeat override remote seat keyless.
func (c *Command) OvdRemoteSeat() error {
	_, err := c.exec("OVD_RMT_SEAT", nil)
	return err
}

// OvdRemoteAlarm override remote alarm keyless.
func (c *Command) OvdRemoteAlarm() error {
	_, err := c.exec("OVD_RMT_ALARM", nil)
	return err
}

// AudioBeep beep the digital audio module.
func (c *Command) AudioBeep() error {
	_, err := c.exec("AUDIO_BEEP", nil)
	return err
}

// FingerFetch get all registered fingerprint ids.
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

// FingerAdd add a new fingerprint id.
func (c *Command) FingerAdd() (int, error) {
	msg, err := c.exec("FINGER_ADD", nil)
	if err != nil {
		return 0, err
	}

	// decode
	id, _ := strconv.Atoi(string(msg[0]))
	return id, nil
}

// FingerDel delete a fingerprint id.
func (c *Command) FingerDel(id int) error {
	min, max := 1, FINGERPRINT_MAX
	if id < min || id > max {
		return errors.New("id out of range")
	}

	_, err := c.exec("FINGER_DEL", nil)
	return err
}

// FingerRst reset all fingerprint ids.
func (c *Command) FingerRst() error {
	_, err := c.exec("FINGER_RST", nil)
	return err
}

// RemotePairing turn on keyless pairing mode.
func (c *Command) RemotePairing() error {
	_, err := c.exec("REMOTE_PAIRING", nil)
	return err
}

// FotaVcu upgrade VCU firmware over the air.
func (c *Command) FotaVcu() (string, error) {
	msg, err := c.exec("FOTA_VCU", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// FotaHmi upgrade HMI firmware over the air.
func (c *Command) FotaHmi() (string, error) {
	msg, err := c.exec("FOTA_HMI", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}
