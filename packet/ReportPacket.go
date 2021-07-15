package packet

type ReportPacket struct {
	Header *HeaderPacket
	VCU    *VcuPacket
	Eeprom *EepromPacket
	GPS    *GpsPacket
	Hbar   *HbarPacket
	Net	   *NetPacket
	Mems   *MemsPacket
	Remote *RemotePacket
	Finger *FingerPacket
	Audio  *AudioPacket
	Hmi1   *Hmi1Packet
	Bms    *BmsPacket
	Mcu    *McuPacket
	Task   *TaskPacket
}

type HeaderPacket struct {
	Prefix       string `len:"2"`
	Size         uint8  `unit:"Bytes"`
	Vin          uint32 `type:"uint32"`
	SendDatetime int64  `len:"7"`
}

type VcuPacket struct {
	FrameID     uint8   ``
	LogDatetime int64   `len:"7"`
	State       int8    `type:"int8"`
	Events      uint16  ``
	LogBuffered uint8   ``
	BatVoltage  float32 `len:"1" unit:"mVolt" factor:"18.0"`
	Uptime      float32 `type:"uint32" unit:"hour" factor:"0.000277"`
}

type EepromPacket struct {
	Active bool  
	Used   uint8 `unit:"%"`
}

type GpsPacket struct {
	Active    bool    ``
	SatInUse  uint8   `unit:"Sat"`
	HDOP      float32 `len:"1" factor:"0.1"`
	VDOP      float32 `len:"1" factor:"0.1"`
	Speed     uint8   `unit:"Kph"`
	Heading   float32 `len:"1" unit:"Deg" factor:"2.0"`
	Longitude float32 `factor:"0.0000001"`
	Latitude  float32 `factor:"0.0000001"`
	Altitude  float32 `len:"2" unit:"m" factor:"0.1"`
}

type HbarPacket struct {
	Reverse bool 
	Mode    struct {
		Drive   MODE_DRIVE 
		Trip    MODE_TRIP  
		Average MODE_AVG   
	}
	Trip struct {
		A        uint16 `unit:"Km"`
		B        uint16 `unit:"Km"`
		Odometer uint16 `unit:"Km"`
	}
	Average struct {
		Range      uint8 `unit:"Km"`
		Efficiency uint8 `unit:"Km/Kwh"`
	}
}

type NetPacket struct {
	Signal   uint8        `unit:"%"`
	State    NET_STATE    `type:"int8"`
	IpStatus NET_IP_STATE `type:"int8"`
}

type MemsPacket struct {
	Active   bool 
	Detector bool 
	Accel    struct {
		X float32 `len:"2" unit:"G" factor:"0.01"`
		Y float32 `len:"2" unit:"G" factor:"0.01"`
		Z float32 `len:"2" unit:"G" factor:"0.01"`
	}
	Gyro struct {
		X float32 `len:"2" unit:"rad/s" factor:"0.1"`
		Y float32 `len:"2" unit:"rad/s" factor:"0.1"`
		Z float32 `len:"2" unit:"rad/s" factor:"0.1"`
	}
	Tilt struct {
		Pitch float32 `len:"2" unit:"Deg" factor:"0.1"`
		Roll  float32 `len:"2" unit:"Deg" factor:"0.1"`
	}
	Total struct {
		Accel       float32 `len:"2" unit:"G" factor:"0.01"`
		Gyro        float32 `len:"2" unit:"rad/s" factor:"0.1"`
		Tilt        float32 `len:"2" unit:"Deg" factor:"0.1"`
		Temperature float32 `len:"2" unit:"Celcius" factor:"0.1"`
	}
}

type RemotePacket struct {
	Active bool 
	Nearby bool 
}

type FingerPacket struct {
	Verified bool  
	DriverID uint8 
}

type AudioPacket struct {
	Active bool  
	Mute   bool  
	Volume uint8 `unit:"%"`
}
type Hmi1Packet struct {
	Active bool 
}

type BmsPacket struct {
	Active bool   
	Run    bool   
	SOC    uint8  `unit:"%"`
	Fault  uint16
	Pack   [BMS_PACK_CNT]struct {
		ID          uint32  
		Fault       uint16  
		Voltage     float32 `len:"2" unit:"Volt" factor:"0.01"`
		Current     float32 `len:"2" unit:"Ampere" factor:"0.1"`
		SOC         uint8   `unit:"%"`
		Temperature uint16  `unit:"Celcius"`
	}
}

type McuPacket struct {
	Active      bool       
	Run         bool       
	Reverse     bool       
	DriveMode   MODE_DRIVE 
	Speed       uint8      `unit:"Kph"`
	RPM         int16      `unit:"rpm"`
	Temperature float32    `len:"2" unit:"Celcius" factor:"0.1"`
	Fault       struct {
		Post uint32
		Run  uint32
	}
	Torque struct {
		Command  float32 `len:"2" unit:"Nm" factor:"0.1"`
		Feedback float32 `len:"2" unit:"Nm" factor:"0.1"`
	}
	DCBus struct {
		Current float32 `len:"2" unit:"A" factor:"0.1"`
		Voltage float32 `len:"2" unit:"V" factor:"0.1"`
	}
	Inverter struct {
		Enabled   bool              
		Lockout   bool              
		Discharge MCU_INV_DISCHARGE 
	}
	Template struct {
		MaxRPM    int16 `unit:"rpm"`
		MaxSpeed  uint8 `unit:"Kph"`
		DriveMode [DRIVE_MODE_CNT]struct {
			Discur uint16  `unit:"A"`
			Torque float32 `len:"2" unit:"Nm" factor:"0.1"`
		}
	}
}

type TaskPacket struct {
	Stack struct {
		Manager  uint16 `unit:"Bytes"`
		Network  uint16 `unit:"Bytes"`
		Reporter uint16 `unit:"Bytes"`
		Command  uint16 `unit:"Bytes"`
		Mems     uint16 `unit:"Bytes"`
		Remote   uint16 `unit:"Bytes"`
		Finger   uint16 `unit:"Bytes"`
		Audio    uint16 `unit:"Bytes"`
		Gate     uint16 `unit:"Bytes"`
		CanRX    uint16 `unit:"Bytes"`
		CanTX    uint16 `unit:"Bytes"`
	}
	Wakeup struct {
		Manager  uint8 `unit:"s"`
		Network  uint8 `unit:"s"`
		Reporter uint8 `unit:"s"`
		Command  uint8 `unit:"s"`
		Mems     uint8 `unit:"s"`
		Remote   uint8 `unit:"s"`
		Finger   uint8 `unit:"s"`
		Audio    uint8 `unit:"s"`
		Gate     uint8 `unit:"s"`
		CanRX    uint8 `unit:"s"`
		CanTX    uint8 `unit:"s"`
	}
}
