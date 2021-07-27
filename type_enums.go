package sdk

type Frame uint8

const (
	FrameInvalid Frame = iota
	FrameSimple
	FrameFull
)

func (m Frame) String() string {
	return [...]string{
		"INVALID",
		"SIMPLE",
		"FULL",
	}[m]
}

type BikeState int8

const (
	BikeStateUnknown BikeState = iota - 3
	BikeStateLost
	BikeStateBackup
	BikeStateNormal
	BikeStateStandby
	BikeStateReady
	BikeStateRun
	BikeStateLimit
)

func (m BikeState) String() string {
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

type ModeSub uint8

const (
	ModeSubDrive ModeSub = iota
	ModeSubTrip
	ModeSubAvg
	ModeLimit
)

func (m ModeSub) String() string {
	return [...]string{
		"DRIVE",
		"TRIP",
		"AVG",
	}[m]
}

type ModeDrive uint8

const (
	ModeDriveEconomy ModeDrive = iota
	ModeDriveStandard
	ModeDriveSport
	ModeDriveLimit
)

func (m ModeDrive) String() string {
	return [...]string{
		"ECONOMY",
		"STANDARD",
		"SPORT",
	}[m]
}

type ModeTrip uint8

const (
	ModeTripA ModeTrip = iota
	ModeTripB
	ModeTripOdo
	ModeTripLimit
)

func (m ModeTrip) String() string {
	return [...]string{
		"A",
		"B",
		"ODO",
	}[m]
}

type ModeAvg uint8

const (
	ModeAvgRange ModeAvg = iota
	ModeAvgEfficiency
	ModeAvgLimit
)

func (m ModeAvg) String() string {
	return [...]string{
		"RANGE",
		"EFFICIENCY",
	}[m]
}

func (m ModeAvg) Unit() string {
	return [...]string{
		"KM",
		"KM/KWH",
	}[m]
}

type NetState int8

const (
	NetStateDown NetState = iota - 1
	NetStateReady
	NetStateConfigured
	NetStateNetworkOn
	NetStateGprsOn
	NetStatePdpOn
	NetStateInternetOn
	NetStateServerOn
	NetStateMqttOn
)

func (m NetState) String() string {
	return [...]string{
		"DOWN",
		"READY",
		"CONFIGURED",
		"NETWORK_ON",
		"GPRS_ON",
		"PDP_ON",
		"INTERNET_ON",
		"SERVER_ON",
		"MQTT_ON",
	}[m+1]
}

type NetIpStatus int8

const (
	NetIpStatusUnknown NetIpStatus = iota - 1
	NetIpStatusInitial
	NetIpStatusStart
	NetIpStatusConfig
	NetIpStatusGprsAct
	NetIpStatusStatus
	NetIpStatusConnecting
	NetIpStatusConnectOk
	NetIpStatusClosing
	NetIpStatusClosed
	NetIpStatusPdpDeact
)

func (m NetIpStatus) String() string {
	return [...]string{
		"UNKNOWN",
		"INITIAL",
		"START",
		"CONFIG",
		"GPRSACT",
		"STATUS",
		"CONNECTING",
		"CONNECT_OK",
		"CLOSING",
		"CLOSED",
		"PDP_DEACT",
	}[m+1]
}

type McuInvDischarge uint8

const (
	McuInvDischargeDisabled McuInvDischarge = iota
	McuInvDischargeEnabled
	McuInvDischargeCheck
	McuInvDischargeOccuring
	McuInvDischargeCompleted
)

func (m McuInvDischarge) String() string {
	return [...]string{
		"DISABLED",
		"ENABLED",
		"CHECK",
		"OCCURING",
		"COMPLETED",
	}[m]
}

type resCode uint8

const (
	resCodeError resCode = iota
	resCodeOk
	resCodeInvalid
	resCodeLimit
)

func (m resCode) String() string {
	return [...]string{
		"ERROR",
		"OK",
		"INVALID",
	}[m]
}

type component string

// Component names for debug output
const (
	CMD component = "[command] "
	RPT component = "[report]  "
	CLI component = "[client]  "
)

type VcuEvent uint8
type VcuEvents []VcuEvent

const (
	VCU_NET_SOFT_RESET VcuEvent = iota
	VCU_NET_HARD_RESET
	VCU_REMOTE_MISSING
	VCU_BIKE_FALLEN
	VCU_BIKE_MOVED
	VCU_BMS_ERROR
	VCU_MCU_ERROR
	VCU_EVENTS_MAX
)

func (m VcuEvent) String() string {
	return [...]string{
		"NET_SOFT_RESET",
		"NET_HARD_RESET",
		"REMOTE_MISSING",
		"BIKE_FALLEN",
		"BIKE_MOVED",
		"BMS_ERROR",
		"MCU_ERROR",
	}[m]
}

type BmsFault uint8
type BmsFaults []BmsFault

const (
	BMS_DISCHARGE_OVER_CURRENT BmsFault = iota
	BMS_CHARGE_OVER_CURRENT
	BMS_SHORT_CIRCUIT
	BMS_DISCHARGE_OVER_TEMPERATURE
	BMS_DISCHARGE_UNDER_TEMPERATURE
	BMS_CHARGE_OVER_TEMPERATURE
	BMS_CHARGE_UNDER_TEMPERATURE
	BMS_UNDER_VOLTAGE
	BMS_OVER_VOLTAGE
	BMS_OVER_DISCHARGE_CAPACITY
	BMS_UNBALANCE
	BMS_SYSTEM_FAILURE
	BMS_FAULTS_MAX
)

func (m BmsFault) String() string {
	return [...]string{
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
	}[m]
}

type McuFaultType uint8

const (
	MCU_FAULT_POST McuFaultType = iota
	MCU_FAULT_RUN
)

type McuFault uint8
type McuFaultPost McuFault
type McuFaultRun McuFault

type McuFaults struct {
	Post []McuFaultPost
	Run  []McuFaultRun
}

const (
	MCU_POST_HW_DESATURATION McuFaultPost = iota
	MCU_POST_HW_OVER_CURRENT
	MCU_POST_ACCEL_SHORTED
	MCU_POST_ACCEL_OPEN
	MCU_POST_CURRENT_LOW
	MCU_POST_CURRENT_HIGH
	MCU_POST_MOD_TEMP_LOW
	MCU_POST_MOD_TEMP_HIGH
	MCU_POST_PCB_TEMP_LOW
	MCU_POST_PCB_TEMP_HIGH
	MCU_POST_GATE_TEMP_LOW
	MCU_POST_GATE_TEMP_HIGH
	MCU_POST_5V_LOW
	MCU_POST_5V_HIGH
	MCU_POST_12V_LOW
	MCU_POST_12V_HIGH
	MCU_POST_2v5_LOW
	MCU_POST_2v5_HIGH
	MCU_POST_1v5_LOW
	MCU_POST_1v5_HIGH
	MCU_POST_DCBUS_VOLT_HIGH
	MCU_POST_DCBUS_VOLT_LOW
	MCU_POST_PRECHARGE_TO
	MCU_POST_PRECHARGE_FAIL
	MCU_POST_EE_CHECKSUM_INVALID
	MCU_POST_EE_DATA_OUT_RANGE
	MCU_POST_EE_UPDATE_REQ
	MCU_POST_RESERVED_1
	MCU_POST_RESERVED_2
	MCU_POST_RESERVED_3
	MCU_POST_BRAKE_SHORTED
	MCU_POST_BRAKE_OPEN
	MCU_POST_FAULTS_MAX
)

func (m McuFaultPost) String() string {
	return [...]string{
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
	}[m]
}

const (
	MCU_RUN_OVER_SPEED McuFaultRun = iota
	MCU_RUN_OVER_CURRENT
	MCU_RUN_OVER_VOLTAGE
	MCU_RUN_INV_OVER_TEMP
	MCU_RUN_ACCEL_SHORTED
	MCU_RUN_ACCEL_OPEN
	MCU_RUN_DIRECTION_FAULT
	MCU_RUN_INV_TO
	MCU_RUN_HW_DESATURATION
	MCU_RUN_HW_OVER_CURRENT
	MCU_RUN_UNDER_VOLTAGE
	MCU_RUN_CAN_LOST
	MCU_RUN_MOTOR_OVER_TEMP
	MCU_RUN_RESERVER_1
	MCU_RUN_RESERVER_2
	MCU_RUN_RESERVER_3
	MCU_RUN_BRAKE_SHORTED
	MCU_RUN_BRAKE_OPEN
	MCU_RUN_MODA_OVER_TEMP
	MCU_RUN_MODB_OVER_TEMP
	MCU_RUN_MODC_OVER_TEMP
	MCU_RUN_PCB_OVER_TEMP
	MCU_RUN_GATE1_OVER_TEMP
	MCU_RUN_GATE2_OVER_TEMP
	MCU_RUN_GATE3_OVER_TEMP
	MCU_RUN_CURRENT_FAULT
	MCU_RUN_RESERVER_4
	MCU_RUN_RESERVER_5
	MCU_RUN_RESERVER_6
	MCU_RUN_RESERVER_7
	MCU_RUN_RESOLVER_FAULT
	MCU_RUN_INV_DISCHARGE
	MCU_RUN_FAULTS_MAX
)

func (m McuFaultRun) String() string {
	return [...]string{
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
	}[m]
}
