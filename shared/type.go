package shared

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

var TASK_LIST = []string{
	"manager",
	"network",
	"reporter",
	"command",
	"mems",
	"remote",
	"finger",
	"audio",
	"gate",
	"canRx",
	"canTx",
}

type FRAME_ID uint8

const (
	FRAME_ID_INVALID FRAME_ID = iota
	FRAME_ID_SIMPLE
	FRAME_ID_FULL
)

func (m FRAME_ID) String() string {
	return [...]string{
		"INVALID",
		"SIMPLE",
		"FULL",
	}[m]
}

type BIKE_STATE int8

const (
	BIKE_STATE_UNKNOWN BIKE_STATE = iota - 3
	BIKE_STATE_LOST
	BIKE_STATE_BACKUP
	BIKE_STATE_NORMAL
	BIKE_STATE_STANDBY
	BIKE_STATE_READY
	BIKE_STATE_RUN
	BIKE_STATE_limit
)

func (m BIKE_STATE) String() string {
	return [...]string{
		"UNKNOWN",
		"LOST",
		"BACKUP",
		"NORMAL",
		"STANDBY",
		"READY",
		"RUN",
	}[m+3]
}

type MODE uint8

const (
	MODE_SUB_DRIVE MODE = iota
	MODE_SUB_TRIP
	MODE_SUB_AVG
	MODE_SUB_limit
)

func (m MODE) String() string {
	return [...]string{
		"DRIVE",
		"TRIP",
		"AVG",
	}[m]
}

type MODE_DRIVE uint8

const (
	MODE_DRIVE_ECONOMY MODE_DRIVE = iota
	MODE_DRIVE_STANDARD
	MODE_DRIVE_SPORT
	MODE_DRIVE_limit
)

func (m MODE_DRIVE) String() string {
	return [...]string{
		"ECONOMY",
		"STANDARD",
		"SPORT",
	}[m]
}

type MODE_TRIP uint8

const (
	MODE_TRIP_A MODE_TRIP = iota
	MODE_TRIP_B
	MODE_TRIP_ODO
	MODE_TRIP_limit
)

func (m MODE_TRIP) String() string {
	return [...]string{
		"A",
		"B",
		"ODO",
	}[m]
}

type MODE_AVG uint8

const (
	MODE_AVG_RANGE MODE_AVG = iota
	MODE_AVG_EFFICIENCY
	MODE_AVG_limit
)

func (m MODE_AVG) String() string {
	return [...]string{
		"RANGE",
		"EFFICIENCY",
	}[m]
}

func (m MODE_AVG) Unit() string {
	return [...]string{
		"KM",
		"KM/KWH",
	}[m]
}

type NET_STATE int8

const (
	NET_STATE_DOWN NET_STATE = iota - 1
	NET_STATE_READY
	NET_STATE_CONFIGURED
	NET_STATE_NETWORK_ON
	NET_STATE_GPRS_ON
	NET_STATE_PDP_ON
	NET_STATE_INTERNET_ON
	NET_STATE_SERVER_ON
	NET_STATE_MQTT_ON
)

type NET_IP_STATE int8

const (
	NET_IP_STATE_UNKNOWN NET_IP_STATE = iota - 1
	NET_IP_STATE_INITIAL
	NET_IP_STATE_START
	NET_IP_STATE_CONFIG
	NET_IP_STATE_GPRSACT
	NET_IP_STATE_STATUS
	NET_IP_STATE_CONNECTING
	NET_IP_STATE_CONNECT_OK
	NET_IP_STATE_CLOSING
	NET_IP_STATE_CLOSED
	NET_IP_STATE_PDP_DEACT
)

type MCU_INV_DISCHARGE uint8

const (
	MCU_INV_DISCHARGE_DISABLED MCU_INV_DISCHARGE = iota
	MCU_INV_DISCHARGE_ENABLED
	MCU_INV_DISCHARGE_CHECK
	MCU_INV_DISCHARGE_OCCURING
	MCU_INV_DISCHARGE_COMPLETED
)
