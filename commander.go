package sdk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

type commander struct {
	vin     int
	logger  *log.Logger
	broker  Broker
	mutex   *sync.Mutex
	resChan chan []byte
}

// newCommander create new *commander instance and listen to command & response topic.
func newCommander(vin int, broker Broker, logging bool) (*commander, error) {
	cmder := &commander{
		vin:     vin,
		logger:  newLogger(logging, "COMMAND"),
		broker:  broker,
		mutex:   &sync.Mutex{},
		resChan: make(chan []byte, 1),
	}

	if err := cmder.listenResponse(); err != nil {
		return nil, err
	}
	return cmder, nil
}

// GenInfo gather device information.
func (c *commander) GenInfo() (string, error) {
	msg, err := c.exec("GEN_INFO", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// GenLed set built-in led state on device.
func (c *commander) GenLed(on bool) error {
	_, err := c.exec("GEN_LED", boolToBytes(on))
	return err
}

// GenRtc set real time clock on device.
func (c *commander) GenRtc(time time.Time) error {
	_, err := c.exec("GEN_RTC", timeToBytes(time))
	return err
}

// GenOdo set odometer value (in km).
func (c *commander) GenOdo(km uint16) error {
	payload := uintToBytes(reflect.Uint16, uint64(km))
	_, err := c.exec("GEN_ODO", payload)
	return err
}

// GenAntiTheaf toggle anti-thief motion detector.
func (c *commander) GenAntiTheaf() error {
	_, err := c.exec("GEN_ANTI_THIEF", nil)
	return err
}

// GenReportFlush flush pending report in device buffer.
func (c *commander) GenReportFlush() error {
	_, err := c.exec("GEN_RPT_FLUSH", nil)
	return err
}

// GenReportBlock stop device reporting mode.
func (c *commander) GenReportBlock(on bool) error {
	_, err := c.exec("GEN_RPT_BLOCK", boolToBytes(on))
	return err
}

// OvdState override bike state.
func (c *commander) OvdState(state BikeState) error {
	min, max := BikeStateNormal, BikeStateRun
	if state < min || state > max {
		return errInputOutOfRange("state")
	}

	payload := []byte{byte(state)}
	_, err := c.exec("OVD_STATE", payload)
	return err
}

// OvdReportInterval override reporting interval.
func (c *commander) OvdReportInterval(dur time.Duration) error {
	min, max := time.Duration(5), time.Duration(^uint16(0))
	if dur < min*time.Second || dur > max*time.Second {
		return errInputOutOfRange("duration")
	}
	payload := uintToBytes(reflect.Uint16, uint64(dur.Seconds()))
	_, err := c.exec("OVD_RPT_INTERVAL", payload)
	return err
}

// OvdReportFrame override report frame type.
func (c *commander) OvdReportFrame(frame Frame) error {
	if frame == FrameInvalid {
		return errInputOutOfRange("frame")
	}

	payload := []byte{byte(frame)}
	_, err := c.exec("OVD_RPT_FRAME", payload)
	return err
}

// OvdRemoteSeat override seat button on remote/keyless.
func (c *commander) OvdRemoteSeat() error {
	_, err := c.exec("OVD_RMT_SEAT", nil)
	return err
}

// OvdRemoteAlarm override alarm button on remote/keyless.
func (c *commander) OvdRemoteAlarm() error {
	_, err := c.exec("OVD_RMT_ALARM", nil)
	return err
}

// AudioBeep beep the digital audio module.
func (c *commander) AudioBeep() error {
	_, err := c.exec("AUDIO_BEEP", nil)
	return err
}

// FingerFetch get all registered fingerprint ids.
func (c *commander) FingerFetch() ([]int, error) {
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
func (c *commander) FingerAdd() (int, error) {
	msg, err := c.exec("FINGER_ADD", nil)
	if err != nil {
		return 0, err
	}

	// decode
	id, _ := strconv.Atoi(string(msg[0]))
	return id, nil
}

// FingerDel delete a fingerprint id.
func (c *commander) FingerDel(id int) error {
	min, max := 1, FINGERPRINT_MAX
	if id < min || id > max {
		return errInputOutOfRange("id")
	}

	_, err := c.exec("FINGER_DEL", nil)
	return err
}

// FingerRst reset all fingerprint ids.
func (c *commander) FingerRst() error {
	_, err := c.exec("FINGER_RST", nil)
	return err
}

// RemotePairing turn on keyless pairing mode.
func (c *commander) RemotePairing() error {
	_, err := c.exec("REMOTE_PAIRING", nil)
	return err
}

// FotaVcu upgrade VCU (Vehicle Control Unit) firmware over the air.
func (c *commander) FotaVcu() (string, error) {
	msg, err := c.exec("FOTA_VCU", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// FotaHmi upgrade Dashboard/HMI (Human Machine Interface) firmware over the air.
func (c *commander) FotaHmi() (string, error) {
	msg, err := c.exec("FOTA_HMI", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// NetSendUssd send USSD to cellular network.
// Input example: *123*10*3#
func (c *commander) NetSendUssd(ussd string) (string, error) {
	min, max := 3, 20
	if len(ussd) < min || len(ussd) > max {
		return "", errInputOutOfRange("ussd")
	}
	if !strings.HasPrefix(ussd, "*") || !strings.HasSuffix(ussd, "#") {
		return "", errors.New("invalid ussd format")
	}

	msg, err := c.exec("NET_SEND_USSD", []byte(ussd))
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// NetReadSms read latest cellular SMS inbox.
func (c *commander) NetReadSms() (string, error) {
	msg, err := c.exec("NET_READ_SMS", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// HbarDrive set handlebar drive mode.
func (c *commander) HbarDrive(drive ModeDrive) error {
	if drive == ModeDriveLimit {
		return errInputOutOfRange("drive-mode")
	}

	payload := []byte{byte(drive)}
	_, err := c.exec("HBAR_DRIVE", payload)
	return err
}

// HbarTrip set handlebar trip mode.
func (c *commander) HbarTrip(trip ModeTrip) error {
	if trip == ModeTripLimit {
		return errInputOutOfRange("trip-mode")
	}

	payload := []byte{byte(trip)}
	_, err := c.exec("HBAR_TRIP", payload)
	return err
}

// HbarAvg set handlebar average mode.
func (c *commander) HbarAvg(avg ModeAvg) error {
	if avg == ModeAvgLimit {
		return errInputOutOfRange("avg-mode")
	}

	payload := []byte{byte(avg)}
	_, err := c.exec("HBAR_AVG", payload)
	return err
}

// HbarReverse set MCU (Motor Control Unit) reverse state.
func (c *commander) HbarReverse(on bool) error {
	_, err := c.exec("HBAR_REVERSE", boolToBytes(on))
	return err
}

// McuSpeedMax set maximum MCU (Motor Control Unit) speed (in kph).
func (c *commander) McuSpeedMax(kph uint8) error {
	payload := uintToBytes(reflect.Uint8, uint64(kph))
	_, err := c.exec("MCU_SPEED_MAX", payload)
	return err
}

type McuTemplate struct {
	DisCur uint16
	Torque uint16
}

// McuTemplates set all MCU (Motor Control Unit) driving mode templates.
func (c *commander) McuTemplates(ts []McuTemplate) error {
	if len(ts) != int(ModeDriveLimit) {
		return errors.New("templates should be set for all driving mode at once")
	}

	var buf bytes.Buffer
	var min, maxDisCur, maxTorque uint16 = 1, 32767, 3276
	for i, t := range ts {
		driveMode := ModeDrive(i)
		if t.DisCur < min || t.DisCur > maxDisCur {
			return errInputOutOfRange(fmt.Sprintf("%s:dischare-current", driveMode))
		}
		if t.Torque < min || t.Torque > maxTorque {
			return errInputOutOfRange(fmt.Sprintf("%s:torque", driveMode))
		}

		binary.Write(&buf, binary.LittleEndian, t.DisCur)
		binary.Write(&buf, binary.LittleEndian, t.Torque)
	}

	_, err := c.exec("MCU_TEMPLATES", buf.Bytes())
	return err
}
