package packet

type ReportSimplePacket struct {
	Header HeaderPacket
	VCU    VcuPacket
	EEPROM EepromPacket
	GPS    GpsPacket
}
type ReportFullPacket struct {
	ReportSimplePacket
	HBAR   HbarPacket
	NET    NetPacket
	MEMS   MemsPacket
	Remote RemotePacket
	Finger FingerPacket
	Audio  AudioPacket
	HMI1   Hmi1Packet
	BMS    BmsPacket
	MCU    McuPacket
	TASK   TaskPacket
}

type VcuPacket struct {
	FrameID     FrameID  `type:"uint8"`
	LogDatetime int64    `type:"unix_time" len:"7"`
	State       VcuState `type:"int8"`
	Events      uint16   `type:"uint16"`
	LogBuffered uint8    `type:"uint8"`
	BatVoltage  float32  `type:"uint8" unit:"mVolt" factor:"18.0"`
	Uptime      float32  `type:"uint32" unit:"hour" factor:"0.000277"`
}

type EepromPacket struct {
	Active bool  `type:"uint8"`
	Used   uint8 `type:"uint8" unit:"%"`
}

type GpsPacket struct {
	Active    bool    `type:"uint8"`
	SatInUse  uint8   `type:"uint8" unit:"Sat"`
	HDOP      float32 `type:"uint8" factor:"0.1"`
	VDOP      float32 `type:"uint8" factor:"0.1"`
	Speed     uint8   `type:"uint8" unit:"Kph"`
	Heading   float32 `type:"uint8" unit:"Deg" factor:"2.0"`
	Longitude float32 `type:"int32" factor:"0.0000001"`
	Latitude  float32 `type:"int32" factor:"0.0000001"`
	Altitude  float32 `type:"uint16" unit:"m" factor:"0.1"`
}

type HbarPacket struct {
	Reverse bool `type:"uint8"`
	Mode    struct {
		Drive      ModeDrive      `type:"uint8"`
		Trip       ModeTrip       `type:"uint8"`
		Prediction ModePrediction `type:"uint8"`
	}
	Trip struct {
		A        uint16 `type:"uint16" unit:"Km"`
		B        uint16 `type:"uint16" unit:"Km"`
		Odometer uint16 `type:"uint16" unit:"Km"`
	}
	Prediction struct {
		Range      uint8 `type:"uint8" unit:"Km"`
		Efficiency uint8 `type:"uint8" unit:"Km/Kwh"`
	}
}

type NetPacket struct {
	Signal   uint8       `type:"uint8" unit:"%"`
	State    NetState    `type:"int8"`
	IpStatus NetIpStatus `type:"int8"`
}

type MemsPacket struct {
	Active   bool `type:"uint8"`
	Detector bool `type:"uint8"`
	Accel    struct {
		X float32 `type:"int16" unit:"G" factor:"0.01"`
		Y float32 `type:"int16" unit:"G" factor:"0.01"`
		Z float32 `type:"int16" unit:"G" factor:"0.01"`
	}
	Gyro struct {
		X float32 `type:"int16" unit:"rad/s" factor:"0.1"`
		Y float32 `type:"int16" unit:"rad/s" factor:"0.1"`
		Z float32 `type:"int16" unit:"rad/s" factor:"0.1"`
	}
	Tilt struct {
		Pitch float32 `type:"int16" unit:"Deg" factor:"0.1"`
		Roll  float32 `type:"int16" unit:"Deg" factor:"0.1"`
	}
	Total struct {
		Accel       float32 `type:"uint16" unit:"G" factor:"0.01"`
		Gyro        float32 `type:"uint16" unit:"rad/s" factor:"0.1"`
		Tilt        float32 `type:"uint16" unit:"Deg" factor:"0.1"`
		Temperature float32 `type:"uint16" unit:"Celcius" factor:"0.1"`
	}
}

type RemotePacket struct {
	Active bool `type:"uint8"`
	Nearby bool `type:"uint8"`
}

type FingerPacket struct {
	Verified bool  `type:"uint8"`
	DriverID uint8 `type:"uint8"`
}

type AudioPacket struct {
	Active bool  `type:"uint8"`
	Mute   bool  `type:"uint8"`
	Volume uint8 `type:"uint8" unit:"%"`
}
type Hmi1Packet struct {
	Active bool `type:"uint8"`
}

type BmsPacket struct {
	Active bool   `type:"uint8"`
	Run    bool   `type:"uint8"`
	SOC    uint8  `type:"uint8" unit:"%"`
	Fault  uint16 `type:"uint16"`
	Pack   [BMS_PACK_CNT]struct {
		ID          uint32  `type:"uint32"`
		Fault       uint16  `type:"uint16"`
		Voltage     float32 `type:"uint16" unit:"Volt" factor:"0.01"`
		Current     float32 `type:"uint16" unit:"Ampere" factor:"0.1"`
		SOC         uint8   `type:"uint8" unit:"%"`
		Temperature uint16  `type:"uint16" unit:"Celcius"`
	}
}

type McuPacket struct {
	Active      bool      `type:"uint8"`
	Run         bool      `type:"uint8"`
	Reverse     bool      `type:"uint8"`
	DriveMode   ModeDrive `type:"uint8"`
	Speed       uint8     `type:"uint8" unit:"Kph"`
	RPM         int16     `type:"int16" unit:"rpm"`
	Temperature float32   `type:"uint16" unit:"Celcius" factor:"0.1"`
	Fault       struct {
		Post uint32 `type:"uint32"`
		Run  uint32 `type:"uint32"`
	}
	Torque struct {
		Command  float32 `type:"uint16" unit:"Nm" factor:"0.1"`
		Feedback float32 `type:"uint16" unit:"Nm" factor:"0.1"`
	}
	DCBus struct {
		Current float32 `type:"uint16" unit:"A" factor:"0.1"`
		Voltage float32 `type:"uint16" unit:"V" factor:"0.1"`
	}
	Inverter struct {
		Enabled   bool            `type:"uint8"`
		Lockout   bool            `type:"uint8"`
		Discharge McuInvDischarge `type:"uint8"`
	}
	Template struct {
		MaxRPM    int16 `type:"int16" unit:"rpm"`
		MaxSpeed  uint8 `type:"uint8" unit:"Kph"`
		DriveMode [DRIVE_MODE_CNT]struct {
			Discur uint16  `type:"uint16" unit:"A"`
			Torque float32 `type:"uint16" unit:"Nm" factor:"0.1"`
		}
	}
}

type TaskPacket struct {
	Stack struct {
		Manager  uint16 `type:"uint16" unit:"Bytes"`
		Network  uint16 `type:"uint16" unit:"Bytes"`
		Reporter uint16 `type:"uint16" unit:"Bytes"`
		Command  uint16 `type:"uint16" unit:"Bytes"`
		Mems     uint16 `type:"uint16" unit:"Bytes"`
		Remote   uint16 `type:"uint16" unit:"Bytes"`
		Finger   uint16 `type:"uint16" unit:"Bytes"`
		Audio    uint16 `type:"uint16" unit:"Bytes"`
		Gate     uint16 `type:"uint16" unit:"Bytes"`
		CanRX    uint16 `type:"uint16" unit:"Bytes"`
		CanTX    uint16 `type:"uint16" unit:"Bytes"`
	}
	Wakeup struct {
		Manager  uint8 `type:"uint8" unit:"s"`
		Network  uint8 `type:"uint8" unit:"s"`
		Reporter uint8 `type:"uint8" unit:"s"`
		Command  uint8 `type:"uint8" unit:"s"`
		Mems     uint8 `type:"uint8" unit:"s"`
		Remote   uint8 `type:"uint8" unit:"s"`
		Finger   uint8 `type:"uint8" unit:"s"`
		Audio    uint8 `type:"uint8" unit:"s"`
		Gate     uint8 `type:"uint8" unit:"s"`
		CanRX    uint8 `type:"uint8" unit:"s"`
		CanTX    uint8 `type:"uint8" unit:"s"`
	}
}
