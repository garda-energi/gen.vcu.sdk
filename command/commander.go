package command

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/broker"
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type Commander struct {
	vin     int
	broker  *broker.Broker
	mutex   *sync.Mutex
	resChan chan []byte
}

// New create new Commander instance and listen to command & response topic.
func New(vin int, broker *broker.Broker) (*Commander, error) {
	cmder := &Commander{
		vin:     vin,
		broker:  broker,
		mutex:   &sync.Mutex{},
		resChan: make(chan []byte, 1),
	}

	if err := cmder.listen(); err != nil {
		return nil, err
	}
	return cmder, nil
}

// GenInfo gather device information.
func (c *Commander) GenInfo() (string, error) {
	msg, err := c.exec("GEN_INFO", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// GenLed set built-in led state on device.
func (c *Commander) GenLed(on bool) error {
	_, err := c.exec("GEN_LED", shared.BoolToBytes(on))
	return err
}

// GenRtc set real time clock on device.
func (c *Commander) GenRtc(time time.Time) error {
	_, err := c.exec("GEN_RTC", shared.TimeToBytes(time))
	return err
}

// GenOdo set odometer value (in km).
func (c *Commander) GenOdo(km uint16) error {
	payload := shared.UintToBytes(reflect.Uint16, uint64(km))
	_, err := c.exec("GEN_ODO", payload)
	return err
}

// GenAntiTheaf toggle anti-thief motion detector.
func (c *Commander) GenAntiTheaf() error {
	_, err := c.exec("GEN_ANTI_THIEF", nil)
	return err
}

// GenReportFlush flush report buffer.
func (c *Commander) GenReportFlush() error {
	_, err := c.exec("GEN_RPT_FLUSH", nil)
	return err
}

// GenReportBlock block reporting mode.
func (c *Commander) GenReportBlock(on bool) error {
	_, err := c.exec("GEN_RPT_BLOCK", shared.BoolToBytes(on))
	return err
}

// OvdState override bike state.
func (c *Commander) OvdState(state shared.BIKE_STATE) error {
	min, max := shared.BIKE_STATE_NORMAL, shared.BIKE_STATE_RUN
	if state < min || state > max {
		return errors.New("state out of range")
	}

	payload := []byte{byte(state)}
	_, err := c.exec("OVD_STATE", payload)
	return err
}

// OvdReportInterval override reporting interval (in seconds).
func (c *Commander) OvdReportInterval(dur time.Duration) error {
	min, max := time.Duration(5), time.Duration(^uint16(0))
	if dur < min*time.Second || dur > max*time.Second {
		return errors.New("duration out of range")
	}
	payload := shared.UintToBytes(reflect.Uint16, uint64(dur.Seconds()))
	_, err := c.exec("OVD_RPT_INTERVAL", payload)
	return err
}

// OvdReportFrame override report frame type.
func (c *Commander) OvdReportFrame(frame shared.FRAME_ID) error {
	if frame == shared.FRAME_ID_INVALID {
		return errors.New("frame out of range")
	}

	payload := []byte{byte(frame)}
	_, err := c.exec("OVD_RPT_FRAME", payload)
	return err
}

// OvdRemoteSeat override remote seat keyless.
func (c *Commander) OvdRemoteSeat() error {
	_, err := c.exec("OVD_RMT_SEAT", nil)
	return err
}

// OvdRemoteAlarm override remote alarm keyless.
func (c *Commander) OvdRemoteAlarm() error {
	_, err := c.exec("OVD_RMT_ALARM", nil)
	return err
}

// AudioBeep beep the digital audio module.
func (c *Commander) AudioBeep() error {
	_, err := c.exec("AUDIO_BEEP", nil)
	return err
}

// FingerFetch get all registered fingerprint ids.
func (c *Commander) FingerFetch() ([]int, error) {
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
func (c *Commander) FingerAdd() (int, error) {
	msg, err := c.exec("FINGER_ADD", nil)
	if err != nil {
		return 0, err
	}

	// decode
	id, _ := strconv.Atoi(string(msg[0]))
	return id, nil
}

// FingerDel delete a fingerprint id.
func (c *Commander) FingerDel(id int) error {
	min, max := 1, FINGERPRINT_MAX
	if id < min || id > max {
		return errors.New("id out of range")
	}

	_, err := c.exec("FINGER_DEL", nil)
	return err
}

// FingerRst reset all fingerprint ids.
func (c *Commander) FingerRst() error {
	_, err := c.exec("FINGER_RST", nil)
	return err
}

// RemotePairing turn on keyless pairing mode.
func (c *Commander) RemotePairing() error {
	_, err := c.exec("REMOTE_PAIRING", nil)
	return err
}

// FotaVcu upgrade VCU (Vehicle Control Unit) firmware over the air.
func (c *Commander) FotaVcu() (string, error) {
	msg, err := c.exec("FOTA_VCU", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// FotaHmi upgrade Dashdevice/HMI (Human Machine Interface) firmware over the air.
func (c *Commander) FotaHmi() (string, error) {
	msg, err := c.exec("FOTA_HMI", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// NetSendUssd send USSD to cellular network.
// Input example: *123*10*3#
func (c *Commander) NetSendUssd(ussd string) (string, error) {
	min, max := 3, 20
	if len(ussd) < min || len(ussd) > max {
		return "", errors.New("ussd length out of range")
	}
	if !strings.HasPrefix(ussd, "*") || !strings.HasSuffix(ussd, "#") {
		return "", errors.New("ussd is invalid")
	}

	msg, err := c.exec("NET_SEND_USSD", []byte(ussd))
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// NetReadSms read latest cellular SMS inbox.
func (c *Commander) NetReadSms() (string, error) {
	msg, err := c.exec("NET_READ_SMS", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// HbarDrive set handlebar drive mode.
func (c *Commander) HbarDrive(drive shared.MODE_DRIVE) error {
	if drive == shared.MODE_DRIVE_limit {
		return errors.New("drive mode out of range")
	}

	payload := []byte{byte(drive)}
	_, err := c.exec("HBAR_DRIVE", payload)
	return err
}

// HbarTrip set handlebar trip mode.
func (c *Commander) HbarTrip(trip shared.MODE_TRIP) error {
	if trip == shared.MODE_TRIP_limit {
		return errors.New("trip mode out of range")
	}

	payload := []byte{byte(trip)}
	_, err := c.exec("HBAR_TRIP", payload)
	return err
}

// HbarAvg set handlebar average mode.
func (c *Commander) HbarAvg(avg shared.MODE_AVG) error {
	if avg == shared.MODE_AVG_limit {
		return errors.New("average mode out of range")
	}

	payload := []byte{byte(avg)}
	_, err := c.exec("HBAR_AVG", payload)
	return err
}

// HbarReverse set MCU (Motor Control Unit) reverse state.
func (c *Commander) HbarReverse(on bool) error {
	_, err := c.exec("HBAR_REVERSE", shared.BoolToBytes(on))
	return err
}

// McuSpeedMax set maximum MCU (Motor Control Unit) speed (in kph).
func (c *Commander) McuSpeedMax(kph uint8) error {
	payload := shared.UintToBytes(reflect.Uint8, uint64(kph))
	_, err := c.exec("MCU_SPEED_MAX", payload)
	return err
}

type McuTemplate struct {
	DisCur uint16
	Torque uint16
}

// McuTemplates set all MCU (Motor Control Unit)  driving mode templates.
func (c *Commander) McuTemplates(ts []McuTemplate) error {
	if len(ts) != int(shared.MODE_DRIVE_limit) {
		return errors.New("templates should be set for all driving mode at once")
	}

	var buf bytes.Buffer
	var min, maxDisCur, maxTorque uint16 = 1, 32767, 3276
	for i, t := range ts {
		driveMode := shared.MODE_DRIVE(i)
		if t.DisCur < min || t.DisCur > maxDisCur {
			return fmt.Errorf("dischare current for %s out of range", driveMode)
		}
		if t.Torque < min || t.Torque > maxTorque {
			return fmt.Errorf("torque for %s out of range", driveMode)
		}

		binary.Write(&buf, binary.LittleEndian, t.DisCur)
		binary.Write(&buf, binary.LittleEndian, t.Torque)
	}

	_, err := c.exec("MCU_TEMPLATES", buf.Bytes())
	return err
}
