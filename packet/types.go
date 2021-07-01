package packet

// type DateTime [7]byte
// type Prefix [2]byte

const BMS_PACK_CNT = 2
const DRIVE_MODE_CNT = 3

type FrameID uint8

const (
	frINVALID FrameID = iota
	frSIMPLE
	frFULL
)

type VcuState int8

const (
	vsUNKNOWN = iota - 3
	vsLOST
	vsBACKUP
	vsNORMAL
	vsSTANDBY
	vsREADY
	vsRUN
)

type ModeDrive uint8

const (
	mdECONOMY ModeDrive = iota
	mdSTANDARD
	mdSPORT
)

type ModeTrip uint8

const (
	mtA ModeTrip = iota
	mtB
	mtODO
)

type ModePrediction uint8

const (
	mpRANGE ModePrediction = iota
	mpEFFICIENCY
)

type NetState int8

const (
	nsDOWN NetState = iota - 1
	nsREADY
	nsCONFIGURED
	nsNETWORK_ON
	nsGPRS_ON
	nsPDP_ON
	nsINTERNET_ON
	nsSERVER_ON
	nsMQTT_ON
)

type NetIpStatus int8

const (
	nipUNKNOWN NetIpStatus = iota - 1
	nipIP_INITIAL
	nipIP_START
	nipIP_CONFIG
	nipIP_GPRSACT
	nipIP_STATUS
	nipCONNECTING
	nipCONNECT_OK
	nipCLOSING
	nipCLOSED
	nipPDP_DEACT
)

type McuInvDischarge uint8

const (
	midDISABLED McuInvDischarge = iota
	midENABLED
	midCHECK
	midOCCURING
	midCOMPLETED
)
