package sdk

import "time"

const (
	TOPIC_STATUS   = "VCU/+/STS"
	TOPIC_REPORT   = "VCU/+/RPT"
	TOPIC_COMMAND  = "VCU/+/CMD"
	TOPIC_RESPONSE = "VCU/+/RSP"
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
