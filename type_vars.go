package sdk

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// typeOfTime is for comparing struct type as time.Time
var typeOfTime reflect.Type = reflect.ValueOf(time.Now()).Type()

// typeOfMessage is for comparing slice type as message ([]byte)
var typeOfMessage reflect.Type = reflect.ValueOf(message{}).Type()

var (
	errPacketAckCorrupt = errors.New("packet ack corrupt")
	errInvalidPrefix    = errors.New("prefix invalid")
	errInvalidSize      = errors.New("size invalid")
	errInvalidVin       = errors.New("vin invalid")
	errInvalidCmdCode   = errors.New("cmd code invalid")
	errInvalidResCode   = errors.New("resCode invalid")
)

type errPacketTimeout string

func (e errPacketTimeout) Error() string {
	return fmt.Sprintf("packet %s timeout", string(e))
}

type errInputOutOfRange string

func (e errInputOutOfRange) Error() string {
	return fmt.Sprintf("input %s out of range", string(e))
}

// Sleeper is building block for sleep things
type Sleeper interface {
	// Sleep pauses the current goroutine for at least the duration d.
	// A negative or zero duration causes Sleep to return immediately.
	Sleep(time.Duration)
	// After waits for the duration to elapse and then sends the current time
	// on the returned channel.
	After(d time.Duration) <-chan time.Time
}

// realSleeper implement real sleep using time module
type realSleeper struct{}

func (*realSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}
func (*realSleeper) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

var vcuStringEvents = []string{
	"NET_SOFT_RESET",
	"NET_HARD_RESET",
	"REMOTE_MISSING",
	"BIKE_FALLEN",
	"BIKE_MOVED",
	"BMS_ERROR",
	"MCU_ERROR",
}

var bmsStringFaults = []string{
	"DISCHARGE_OVER_CURRENT",
	"CHARGE_OVER_CURRENT",
	"SHORT_CIRCUIT",
	"DISCHARGE_OVER_TEMPERATURE",
	"DISCHARGE_UNDER_TEMPERATURE",
	"CHARGE_OVER_TEMPERATURE",
	"CHARGE_UNDER_TEMPERATURE",
	"UNDER_VOLTAGE",
	"OVER_VOLTAGE",
	"OVER_DISCHARGE_CAPACITY",
	"UNBALANCE",
	"SYSTEM_FAILURE",
}

var mcuStringFaults = []string{
	// Post Fault
	"HW_DESATURATION",
	"HW_OVER_CURRENT",
	"ACCEL_SHORTED",
	"ACCEL_OPEN",
	"CURRENT_LOW",
	"CURRENT_HIGH",
	"MOD_TEMP_LOW",
	"MOD_TEMP_HIGH",
	"PCB_TEMP_LOW",
	"PCB_TEMP_HIGH",
	"GATE_TEMP_LOW",
	"GATE_TEMP_HIGH",
	"5V_LOW",
	"5V_HIGH",
	"12V_LOW",
	"12V_HIGH",
	"2v5_LOW",
	"2v5_HIGH",
	"1v5_LOW",
	"1v5_HIGH",
	"DCBUS_VOLT_HIGH",
	"DCBUS_VOLT_LOW",
	"PRECHARGE_TO",
	"PRECHARGE_FAIL",
	"EE_CHECKSUM_INVALID",
	"EE_DATA_OUT_RANGE",
	"EE_UPDATE_REQ",
	"RESERVED_1",
	"RESERVED_2",
	"RESERVED_3",
	"BRAKE_SHORTED",
	"BRAKE_OPEN",

	// Run Fault
	"OVER_SPEED",
	"OVER_CURRENT",
	"OVER_VOLTAGE",
	"INV_OVER_TEMP",
	"ACCEL_SHORTED",
	"ACCEL_OPEN",
	"DIRECTION_FAULT",
	"INV_TO",
	"HW_DESATURATION",
	"HW_OVER_CURRENT",
	"UNDER_VOLTAGE",
	"CAN_LOST",
	"MOTOR_OVER_TEMP",
	"RESERVER_1",
	"RESERVER_2",
	"RESERVER_3",
	"BRAKE_SHORTED",
	"BRAKE_OPEN",
	"MODA_OVER_TEMP",
	"MODB_OVER_TEMP",
	"MODC_OVER_TEMP",
	"PCB_OVER_TEMP",
	"GATE1_OVER_TEMP",
	"GATE2_OVER_TEMP",
	"GATE3_OVER_TEMP",
	"CURRENT_FAULT",
	"RESERVER_4",
	"RESERVER_5",
	"RESERVER_6",
	"RESERVER_7",
	"RESOLVER_FAULT",
	"INV_DISCHARGE",
}
