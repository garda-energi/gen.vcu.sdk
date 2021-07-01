package packet

type ReportSimplePacket struct {
	Header HeaderPacket
	VCU    VcuPacket
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
	// MCU McuPacket
	TASK TaskPacket
}

type VcuPacket struct {
	FrameID     FrameID  `type:"uint8" chartable:""`
	LogDatetime int64    `type:"byte" len:"7" chartable:""`
	State       VcuState `type:"int8" chartable:""`
	Events      uint16   `type:"uint16" chartable:""`
	LogBuffered uint8    `type:"uint8" chartable:""`
	BatVoltage  float32  `type:"uint8" unit:"mVolt" factor:"18.0" chartable:""`
	Uptime      float32  `type:"uint32" unit:"hour" factor:"0.000277" chartable:""`
}

type GpsPacket struct {
	Active    bool    `type:"uint8" chartable:""`
	SatInUse  uint8   `type:"uint8" unit:"Sat" chartable:""`
	HDOP      float32 `type:"uint8" factor:"0.1" chartable:""`
	VDOP      float32 `type:"uint8" factor:"0.1" chartable:""`
	Speed     uint8   `type:"uint8" unit:"Kph" chartable:""`
	Heading   float32 `type:"uint8" unit:"Deg" factor:"2.0" chartable:""`
	Longitude float32 `type:"int32" factor:"0.0000001" chartable:""`
	Latitude  float32 `type:"int32" factor:"0.0000001" chartable:""`
	Altitude  float32 `type:"uint16" unit:"m" factor:"0.1" chartable:""`
}

type HbarPacket struct {
	Reverse bool `type:"uint8" chartable:""`
	Mode    struct {
		Drive      ModeDrive      `type:"uint8" chartable:""`
		Trip       ModeTrip       `type:"uint8" chartable:""`
		Prediction ModePrediction `type:"uint8" chartable:""`
	}
	Trip struct {
		A        uint16 `type:"uint16" unit:"Km" chartable:""`
		B        uint16 `type:"uint16" unit:"Km" chartable:""`
		Odometer uint16 `type:"uint16" unit:"Km" chartable:""`
	}
	Prediction struct {
		Range      uint8 `type:"uint8" unit:"Km" chartable:""`
		Efficiency uint8 `type:"uint8" unit:"Km/Kwh" chartable:""`
	}
}

type NetPacket struct {
	Signal   uint8       `type:"uint8" unit:"%" chartable:""`
	State    NetState    `type:"int8" chartable:""`
	IpStatus NetIpStatus `type:"int8" chartable:""`
}

type MemsPacket struct {
	Active   bool `type:"uint8" chartable:""`
	Detector bool `type:"uint8" chartable:""`
	Accel    struct {
		X float32 `type:"int16" unit:"G" factor:"0.01" chartable:""`
		Y float32 `type:"int16" unit:"G" factor:"0.01" chartable:""`
		Z float32 `type:"int16" unit:"G" factor:"0.01" chartable:""`
	}
	Gyro struct {
		X float32 `type:"int16" unit:"rad/s" factor:"0.1" chartable:""`
		Y float32 `type:"int16" unit:"rad/s" factor:"0.1" chartable:""`
		Z float32 `type:"int16" unit:"rad/s" factor:"0.1" chartable:""`
	}
	Tilt struct {
		Pitch float32 `type:"int16" unit:"Deg" factor:"0.1" chartable:""`
		Roll  float32 `type:"int16" unit:"Deg" factor:"0.1" chartable:""`
	}
	Total struct {
		Accel       float32 `type:"uint16" unit:"G" factor:"0.01" chartable:""`
		Gyro        float32 `type:"uint16" unit:"rad/s" factor:"0.1" chartable:""`
		Tilt        float32 `type:"uint16" unit:"Deg" factor:"0.1" chartable:""`
		Temperature float32 `type:"uint16" unit:"Celcius" factor:"0.1" chartable:""`
	}
}

type RemotePacket struct {
	Active bool `type:"uint8" chartable:""`
	Nearby bool `type:"uint8" chartable:""`
}

type FingerPacket struct {
	Verified bool  `type:"uint8" chartable:""`
	DriverID uint8 `type:"uint8" chartable:""`
}

type AudioPacket struct {
	Active bool  `type:"uint8" chartable:""`
	Mute   bool  `type:"uint8" chartable:""`
	Volume uint8 `type:"uint8" unit:"%" chartable:""`
}
type Hmi1Packet struct {
	Active bool `type:"uint8" chartable:""`
}

type BmsPacket struct {
	Active bool   `type:"uint8" chartable:""`
	Run    bool   `type:"uint8" chartable:""`
	SOC    uint8  `type:"uint8" unit:"%" chartable:""`
	Fault  uint16 `type:"uint16" chartable:""`
	Pack   [BMS_PACK_CNT]struct {
		ID          uint32  `type:"uint32"`
		Fault       uint16  `type:"uint16" chartable:""`
		Voltage     float32 `type:"uint16" unit:"Volt" factor:"0.01" chartable:""`
		Current     float32 `type:"uint16" unit:"Ampere" factor:"0.1" chartable:""`
		SOC         uint8   `type:"uint8" unit:"%" chartable:""`
		Temperature uint16  `type:"uint16" unit:"Celcius" chartable:""`
		// One struct {
		// 	ID          uint32  `type:"uint32"`
		// 	Fault       uint16  `type:"uint16" chartable:""`
		// 	Voltage     float32 `type:"uint16" unit:"Volt" factor:"0.01" chartable:""`
		// 	Current     float32 `type:"uint16" unit:"Ampere" factor:"0.1" chartable:""`
		// 	SOC         uint8   `type:"uint8" unit:"%" chartable:""`
		// 	Temperature uint16  `type:"uint16" unit:"Celcius" chartable:""`
		// }
		// Two struct {
		// 	ID          uint32  `type:"uint32"`
		// 	Fault       uint16  `type:"uint16" chartable:""`
		// 	Voltage     float32 `type:"uint16" unit:"Volt" factor:"0.01" chartable:""`
		// 	Current     float32 `type:"uint16" unit:"Ampere" factor:"0.1" chartable:""`
		// 	SOC         uint8   `type:"uint8" unit:"%" chartable:""`
		// 	Temperature uint16  `type:"uint16" unit:"Celcius" chartable:""`
		// }
	}
}

type McuPacket struct {
	Active      bool      `type:"uint8" chartable:""`
	Run         bool      `type:"uint8" chartable:""`
	Reverse     bool      `type:"uint8" chartable:""`
	DriveMode   ModeDrive `type:"uint8" chartable:""`
	Speed       uint8     `type:"uint8" unit:"Kph" chartable:""`
	RPM         int16     `type:"int16" unit:"rpm" chartable:""`
	Temperature float32   `type:"uint16" unit:"Celcius" factor:"0.1" chartable:""`
	Fault       struct {
		Post uint32 `type:"uint32" chartable:""`
		Run  uint32 `type:"uint32" chartable:""`
	}
	Torque struct {
		Command  float32 `type:"uint16" unit:"Nm" factor:"0.1" chartable:""`
		Feedback float32 `type:"uint16" unit:"Nm" factor:"0.1" chartable:""`
	}
	DCBus struct {
		Current float32 `type:"uint16" unit:"A" factor:"0.1" chartable:""`
		Voltage float32 `type:"uint16" unit:"V" factor:"0.1" chartable:""`
	}
	Inverter struct {
		Enabled   bool            `type:"uint8" chartable:""`
		Lockout   bool            `type:"uint8" chartable:""`
		Discharge McuInvDischarge `type:"uint8" chartable:""`
	}
	Template struct {
		MaxRPM    int16 `type:"int16" unit:"rpm" chartable:""`
		MaxSpeed  uint8 `type:"uint8" unit:"Kph" chartable:""`
		DriveMode [DRIVE_MODE_CNT]struct {
			Discur uint16  `type:"uint16" unit:"A" chartable:""`
			Torque float32 `type:"uint16" unit:"Nm" factor:"0.1" chartable:""`
			// Economy struct {
			// 	Discur uint16  `type:"uint16" unit:"A" chartable:""`
			// 	Torque float32 `type:"uint16" unit:"Nm" factor:"0.1" chartable:""`
			// }
			// Standard struct {
			// 	Discur uint16  `type:"uint16" unit:"A" chartable:""`
			// 	Torque float32 `type:"uint16" unit:"Nm" factor:"0.1" chartable:""`
			// }
			// Sport struct {
			// 	Discur uint16  `type:"uint16" unit:"A" chartable:""`
			// 	Torque float32 `type:"uint16" unit:"Nm" factor:"0.1" chartable:""`
			// }
		}
	}
}

type TaskPacket struct {
	Stack struct {
		Manager  uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Network  uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Reporter uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Command  uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Mems     uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Remote   uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Finger   uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Audio    uint16 `type:"uint16" unit:"Bytes" chartable:""`
		Gate     uint16 `type:"uint16" unit:"Bytes" chartable:""`
		CanRX    uint16 `type:"uint16" unit:"Bytes" chartable:""`
		CanTX    uint16 `type:"uint16" unit:"Bytes" chartable:""`
	}
	Wakeup struct {
		Manager  uint8 `type:"uint8" unit:"s" chartable:""`
		Network  uint8 `type:"uint8" unit:"s" chartable:""`
		Reporter uint8 `type:"uint8" unit:"s" chartable:""`
		Command  uint8 `type:"uint8" unit:"s" chartable:""`
		Mems     uint8 `type:"uint8" unit:"s" chartable:""`
		Remote   uint8 `type:"uint8" unit:"s" chartable:""`
		Finger   uint8 `type:"uint8" unit:"s" chartable:""`
		Audio    uint8 `type:"uint8" unit:"s" chartable:""`
		Gate     uint8 `type:"uint8" unit:"s" chartable:""`
		CanRX    uint8 `type:"uint8" unit:"s" chartable:""`
		CanTX    uint8 `type:"uint8" unit:"s" chartable:""`
	}
}
