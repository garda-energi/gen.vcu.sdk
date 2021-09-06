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
	mutex   *sync.Mutex
	resChan chan packet
	client  *client
	sleeper Sleeper
}

// newCommander create new *commander instance and listen to command & response topic.
func newCommander(vin int, c *client, s Sleeper, l *log.Logger) (*commander, error) {
	cmder := &commander{
		vin:     vin,
		logger:  l,
		mutex:   &sync.Mutex{},
		resChan: make(chan packet, 1),
		client:  c,
		sleeper: s,
	}

	if err := cmder.listen(); err != nil {
		return nil, err
	}
	return cmder, nil
}

// GenInfo gather device information.
func (c *commander) GenInfo() (string, error) {
	msg, err := c.exec("GenInfo", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// GenLed set built-in led state on device.
func (c *commander) GenLed(on bool) error {
	_, err := c.exec("GenLed", boolToBytes(on))
	return err
}

// GenRtc set real time clock on device.
func (c *commander) GenRtc(time time.Time) error {
	_, err := c.exec("GenRtc", timeToBytes(time))
	return err
}

// GenBikeState override bike state.
func (c *commander) GenBikeState(state BikeState) error {
	min, max := BikeStateNormal, BikeStateRun
	if state < min || state > max {
		return errInputOutOfRange("state")
	}

	_, err := c.exec("GenBikeState", message{byte(state)})
	return err
}

// ReportFlush flush pending report in device buffer.
func (c *commander) ReportFlush() error {
	_, err := c.exec("ReportFlush", nil)
	return err
}

// ReportBlock stop device reporting mode.
func (c *commander) ReportBlock(on bool) error {
	_, err := c.exec("ReportBlock", boolToBytes(on))
	return err
}

// ReportInterval override reporting interval.
func (c *commander) ReportInterval(dur time.Duration) error {
	if dur < REPORT_INTERVAL_MIN || dur > REPORT_INTERVAL_MAX {
		return errInputOutOfRange("duration")
	}
	msg := uintToBytes(reflect.Uint16, uint64(dur.Seconds()))
	_, err := c.exec("ReportInterval", msg)
	return err
}

// ReportFrame override report frame type.
func (c *commander) ReportFrame(frame Frame) error {
	if frame == FrameLimit {
		return errInputOutOfRange("frame")
	}

	_, err := c.exec("ReportFrame", message{byte(frame)})
	return err
}

// AudioBeep beep the digital audio module.
func (c *commander) AudioBeep() error {
	_, err := c.exec("AudioBeep", nil)
	return err
}

// FingerFetch get all registered fingerprint ids.
func (c *commander) FingerFetch() ([]int, error) {
	msg, err := c.exec("FingerFetch", nil)
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
	msg, err := c.exec("FingerAdd", nil)
	if err != nil {
		return 0, err
	}

	// decode
	id, _ := strconv.Atoi(string(msg[0]))
	return id, nil
}

// FingerDel delete a fingerprint id.
func (c *commander) FingerDel(id int) error {
	if id < DRIVER_ID_MIN || id > DRIVER_ID_MAX {
		return errInputOutOfRange("id")
	}

	_, err := c.exec("FingerDel", nil)
	return err
}

// FingerRst reset all fingerprint ids.
func (c *commander) FingerRst() error {
	_, err := c.exec("FingerRst", nil)
	return err
}

// RemotePairing turn on keyless pairing mode.
func (c *commander) RemotePairing() error {
	_, err := c.exec("RemotePairing", nil)
	return err
}

// RemoteSeat override seat button on remote/keyless.
func (c *commander) RemoteSeat() error {
	_, err := c.exec("RemoteSeat", nil)
	return err
}

// RemoteAlarm override alarm button on remote/keyless.
func (c *commander) RemoteAlarm() error {
	_, err := c.exec("RemoteAlarm", nil)
	return err
}

// FotaVcu upgrade VCU (Vehicle Control Unit) firmware over the air.
func (c *commander) FotaVcu() (string, error) {
	msg, err := c.exec("FotaVcu", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// FotaHmi upgrade Dashboard/HMI (Human Machine Interface) firmware over the air.
func (c *commander) FotaHmi() (string, error) {
	msg, err := c.exec("FotaHmi", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// NetSendUssd send USSD to cellular network.
// Input example: *123*10*3#
func (c *commander) NetSendUssd(ussd string) (string, error) {
	if len(ussd) < USSD_LENGTH_MIN || len(ussd) > USSD_LENGTH_MAX {
		return "", errInputOutOfRange("ussd")
	}
	if !strings.HasPrefix(ussd, "*") || !strings.HasSuffix(ussd, "#") {
		return "", errors.New("invalid ussd format")
	}

	msg, err := c.exec("NetSendUssd", message(ussd))
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// NetReadSms read latest cellular SMS inbox.
func (c *commander) NetReadSms() (string, error) {
	msg, err := c.exec("NetReadSms", nil)
	if err != nil {
		return "", err
	}
	return string(msg), nil
}

// HbarTripMeter set trip meter value (in km).
func (c *commander) HbarTripMeter(trip ModeTrip, km uint16) error {
	if trip == ModeTripLimit {
		return errInputOutOfRange("trip-mode")
	}

	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, byte(trip))
	binary.Write(&buf, binary.LittleEndian, uintToBytes(reflect.Uint16, uint64(km)))

	_, err := c.exec("HbarTripMeter", buf.Bytes())
	return err
}

// HbarDrive set handlebar drive mode.
func (c *commander) HbarDrive(drive ModeDrive) error {
	if drive == ModeDriveLimit {
		return errInputOutOfRange("drive-mode")
	}

	_, err := c.exec("HbarDrive", message{byte(drive)})
	return err
}

// HbarTrip set handlebar trip mode.
func (c *commander) HbarTrip(trip ModeTrip) error {
	if trip == ModeTripLimit {
		return errInputOutOfRange("trip-mode")
	}

	_, err := c.exec("HbarTrip", message{byte(trip)})
	return err
}

// HbarAvg set handlebar average mode.
func (c *commander) HbarAvg(avg ModeAvg) error {
	if avg == ModeAvgLimit {
		return errInputOutOfRange("avg-mode")
	}

	_, err := c.exec("HbarAvg", message{byte(avg)})
	return err
}


// McuSpeedMax set maximum MCU (Motor Control Unit) speed (in kph).
func (c *commander) McuSpeedMax(kph uint8) error {
	if kph > SPEED_KPH_MAX {
		return errInputOutOfRange("speed-max")
	}

	msg := uintToBytes(reflect.Uint8, uint64(kph))
	_, err := c.exec("McuSpeedMax", msg)
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
	for i, t := range ts {
		driveMode := ModeDrive(i)
		if t.DisCur < MCU_DISCUR_MIN || t.DisCur > MCU_DISCUR_MAX {
			return errInputOutOfRange(fmt.Sprint(driveMode, ":dischare-current"))
		}
		if t.Torque < MCU_TORQUE_MIN || t.Torque > MCU_TORQUE_MAX {
			return errInputOutOfRange(fmt.Sprint(driveMode, ":torque"))
		}

		binary.Write(&buf, binary.LittleEndian, t.DisCur)
		binary.Write(&buf, binary.LittleEndian, t.Torque)
	}

	_, err := c.exec("McuTemplates", buf.Bytes())
	return err
}

// ImuAntiThief set anti-thief motion detector.
func (c *commander) ImuAntiThief(on bool) error {
	_, err := c.exec("ImuAntiThief", boolToBytes(on))
	return err
}
