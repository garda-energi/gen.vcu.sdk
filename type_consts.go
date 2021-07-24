package sdk

import (
	"errors"
	"fmt"
	"time"
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
const FINGERPRINT_MAX = 5
const SPEED_MAX = 110

const (
	DEFAULT_CMD_TIMEOUT = 10 * time.Second
	DEFAULT_ACK_TIMEOUT = 8 * time.Second
)

// var TASK_LIST = []string{
// 	"manager",
// 	"network",
// 	"reporter",
// 	"command",
// 	"mems",
// 	"remote",
// 	"finger",
// 	"audio",
// 	"gate",
// 	"canRx",
// 	"canTx",
// }
