package sdk

import (
	"time"
)

const SDK_VERSION = 123

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

const BMS_PACK_MAX = 2
const SPEED_KPH_MAX = 110
const TRIP_KM_MAX = 99999
const MESSAGE_LEN_MAX = 200

const REPORT_REALTIME_QUEUED = 3
const EEPROM_LOW_CAPACITY_PERCENT = 90
const BMS_LOW_CAPACITY_PERCENT = 20
const STACK_OVERFLOW_BYTE_MIN = 50
const NET_LOW_SIGNAL_PERCENT = 20

const (
	BATTERY_BACKUP_FULL_MV = 4400
	BATTERY_BACKUP_LOW_MV  = 3300
)

const (
	DRIVER_ID_MIN = 1
	DRIVER_ID_MAX = 5
)

const (
	MCU_DISCUR_MIN = 1
	MCU_DISCUR_MAX = 200
)

const (
	MCU_TORQUE_MIN = 1
	MCU_TORQUE_MAX = 55
)

const (
	USSD_LENGTH_MIN = 3
	USSD_LENGTH_MAX = 20
)

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

type TASK uint8

const (
	TASK_Manager TASK = iota
	TASK_Network
	TASK_Reporte
	TASK_Command
	TASK_Imu
	TASK_Remote
	TASK_Finger
	TASK_Audio
	TASK_Gate
	TASK_CanRX
	TASK_CanTX
	TASK_Limit
)
