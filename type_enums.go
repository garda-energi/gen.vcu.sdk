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
