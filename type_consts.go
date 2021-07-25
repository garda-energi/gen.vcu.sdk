package sdk

import (
	"errors"
	"fmt"
	"time"
)

const (
	TOPIC_STATUS   = "VCU/+/STS"
	TOPIC_REPORT   = "VCU/+/RPT"
	TOPIC_COMMAND  = "VCU/+/CMD"
	TOPIC_RESPONSE = "VCU/+/RSP"
)

const (
	QOS_SUB_STATUS   = 1
	QOS_SUB_REPORT   = 1
	QOS_SUB_COMMAND  = 1
	QOS_SUB_RESPONSE = 2
	QOS_PUB_COMMAND  = 2
	QOS_CMD_FLUSH    = 1
)

const (
	PREFIX_ACK      = "@A"
	PREFIX_REPORT   = "@T"
	PREFIX_COMMAND  = "@C"
	PREFIX_RESPONSE = "@S"
)

const BMS_PACK_CNT = 2
const PAYLOAD_LEN_MAX = 200
const SPEED_MAX = 110

const (
	FINGERPRINT_ID_MIN = 1
	FINGERPRINT_ID_MAX = 5
)

const (
	// TODO: check range in TATC datasheet
	MCU_DISCUR_MIN = 1
	MCU_DISCUR_MAX = 32767
)

const (
	// TODO: check range in TATC datasheet
	MCU_TORQUE_MIN = 7
	MCU_TORQUE_MAX = 3276
)

const (
	USSD_LENGTH_MIN = 3
	USSD_LENGTH_MAX = 20
)

const REPORT_REALTIME_DURATION = -5 * time.Second
const EEPROM_CRITICAL_CAPACITY_PERCENT = 90
const BATTERY_CRITICAL_MV = 3300
const NET_SIGNAL_LOW_PERCENT = 20
const BMS_LOW_CAPACITY_PERCENT = 20
const STACK_OVERFLOW_BYTE_MIN = 40

const (
	DEFAULT_CMD_TIMEOUT = 10 * time.Second
	DEFAULT_ACK_TIMEOUT = 8 * time.Second
)

const (
	REPORT_INTERVAL_MIN = time.Duration(5) * time.Second
	REPORT_INTERVAL_MAX = time.Duration(^uint16(0)) * time.Second
)

const (
	GPS_DOP_MIN = 5
	GPS_LNG_MIN = 95.011198
	GPS_LNG_MAX = 141.020354
	GPS_LAT_MIN = -11.107187
	GPS_LAT_MAX = 5.90713
)

var (
	errPacketAckCorrupt   = errors.New("packet ack corrupt")
	errInvalidPrefix      = errors.New("prefix invalid")
	errInvalidSize        = errors.New("size invalid")
	errInvalidVin         = errors.New("vin invalid")
	errInvalidCode        = errors.New("code invalid")
	errInvalidResCode     = errors.New("resCode invalid")
	errResMessageOverflow = errors.New("message overflow")
)

type errPacketTimeout string

func (e errPacketTimeout) Error() string {
	return fmt.Sprintf("packet %s timeout", string(e))
}

type errInputOutOfRange string

func (e errInputOutOfRange) Error() string {
	return fmt.Sprintf("input %s out of range", string(e))
}
