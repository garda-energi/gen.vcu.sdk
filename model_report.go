package sdk

import (
	"fmt"
	"reflect"
	"time"
)

type HeaderReport struct {
	Header
	Frame Frame `type:"uint8"`
}

type ReportPacket struct {
	// name type         byte index
	Header *HeaderReport // 1 - 15
	Vcu    *Vcu          // 16 - 31
	Eeprom *Eeprom       // 32 - 33
	Gps    *Gps          // 34 - 49
	Hbar   *Hbar         // 49 - 61
	Net    *Net          // 62 - 64
	Mems   *Mems         // 65 - 90
	Remote *Remote       // 91 - 92
	Finger *Finger       // 93 - 94
	Audio  *Audio        // 95 - 98
	Hmi    *Hmi          // 99
	Bms    *Bms          // 100 - 130
	Mcu    *Mcu          // 131 - 173
	Task   *Task         // 174 - 206
}

// ValidPrefix check if r's prefix is valid
func (r *ReportPacket) ValidPrefix() bool {
	if r.Header == nil {
		return false
	}
	return r.Header.Prefix == PREFIX_REPORT
}

// Size calculate total r's size, ignoring prefix & size field
func (r *ReportPacket) Size() int {
	packetSizes := []struct {
		isNil bool // check is nil or not
		size  int  // actual size per attribute
	}{
		{r.Header == nil, 15},
		{r.Vcu == nil, 16},
		{r.Eeprom == nil, 2},
		{r.Gps == nil, 16},
		{r.Hbar == nil, 12},
		{r.Net == nil, 4},
		{r.Mems == nil, 26},
		{r.Remote == nil, 2},
		{r.Finger == nil, 2},
		{r.Audio == nil, 4},
		{r.Hmi == nil, 1},
		{r.Bms == nil, 31},
		{r.Mcu == nil, 42},
		{r.Task == nil, 32},
	}

	validSize := 0
	for _, v := range packetSizes {
		if !v.isNil {
			validSize += v.size
		}
	}

	validSize -= 3
	return validSize
}

// ValidSize check if r's size is valid
func (r *ReportPacket) ValidSize() bool {
	if r.Header == nil {
		return false
	}
	return int(r.Header.Size) == r.Size()
}

// String converts r to string.
func (r *ReportPacket) String() string {
	var out string
	rv := reflect.ValueOf(r).Elem()
	for i := 0; i < rv.NumField(); i++ {
		rvField := rv.Field(i)
		if !rvField.IsNil() {
			rvElem := rvField.Elem()
			out += fmt.Sprintf("%s => %+v\n", rvElem.Type().Name(), rvElem)
		}
	}
	return out
}

type Vcu struct {
	LogDatetime time.Time `type:"int64" len:"7"`
	State       BikeState `type:"int8"`
	Events      uint16    `type:"uint16"`
	LogBuffered uint8     `type:"uint8"`
	BatVoltage  float32   `type:"uint8" len:"1" unit:"mVolt" factor:"18.0"`
	Uptime      float32   `type:"uint32" unit:"hour" factor:"0.000277"`
}

// Events parse v's events
// func (v *Vcu) Events() VcuEvents {
// 	TODO: implement me
// }

// RealtimeData check if current report log is realtime
func (v *Vcu) RealtimeData() bool {
	if v == nil {
		return false
	}
	realtimeDuration := time.Now().Add(REPORT_REALTIME_DURATION)
	return v.LogBuffered == 0 && v.LogDatetime.After(realtimeDuration)
}

// BatteryLow check if v's backup battery voltage is low
func (v *Vcu) BatteryLow() bool {
	if v == nil {
		return false
	}
	return v.BatVoltage < BATTERY_BACKUP_LOW_MV
}

type Eeprom struct {
	Active bool  `type:"uint8"`
	Used   uint8 `type:"uint8" unit:"%"`
}

// CapacityLow check if e's storage capacity is low
func (e *Eeprom) CapacityLow() bool {
	if e == nil {
		return false
	}
	return e.Used > EEPROM_LOW_CAPACITY_PERCENT
}

type Gps struct {
	Active    bool    `type:"uint8"`
	SatInUse  uint8   `type:"uint8" unit:"Sat"`
	HDOP      float32 `type:"uint8" len:"1" factor:"0.1"`
	VDOP      float32 `type:"uint8" len:"1" factor:"0.1"`
	Speed     uint8   `type:"uint8" unit:"Kph"`
	Heading   float32 `type:"uint8" len:"1" unit:"Deg" factor:"2.0"`
	Longitude float32 `type:"int32" factor:"0.0000001"`
	Latitude  float32 `type:"int32" factor:"0.0000001"`
	Altitude  float32 `type:"uint16" len:"2" unit:"m" factor:"0.1"`
}

// ValidHorizontal check if g's horizontal section (heading, longitude, latitude) is valid
func (g *Gps) ValidHorizontal() bool {
	if g == nil {
		return false
	}
	if g.HDOP > GPS_DOP_MIN {
		return false
	}
	return (g.Longitude > GPS_LNG_MIN && g.Longitude < GPS_LNG_MAX) &&
		(g.Latitude > GPS_LAT_MIN && g.Latitude < GPS_LAT_MAX)
}

// ValidVertical check if g's vertical section (altitude) is valid
func (g *Gps) ValidVertical() bool {
	if g == nil {
		return false
	}
	return g.VDOP <= GPS_DOP_MIN
}

type Hbar struct {
	Reverse bool `type:"uint8"`
	Mode    struct {
		Drive ModeDrive `type:"uint8"`
		Trip  ModeTrip  `type:"uint8"`
		Avg   ModeAvg   `type:"uint8"`
	}
	Trip struct {
		A        uint16 `type:"uint16" unit:"Km"`
		B        uint16 `type:"uint16" unit:"Km"`
		Odometer uint16 `type:"uint16" unit:"Km"`
	}
	Avg struct {
		Range      uint8 `type:"uint8" unit:"Km"`
		Efficiency uint8 `type:"uint8" unit:"Km/Kwh"`
	}
}

type Net struct {
	Signal   uint8       `type:"uint8" unit:"%"`
	State    NetState    `type:"int8"`
	IpStatus NetIpStatus `type:"int8"`
}

// LowSignal check if n's signal is low
func (n *Net) LowSignal() bool {
	if n == nil {
		return false
	}
	return n.Signal <= NET_LOW_SIGNAL_PERCENT
}

type Mems struct {
	Active bool `type:"uint8"`
	Motion bool `type:"uint8"`
	Accel  struct {
		X float32 `type:"int16" len:"2" unit:"G" factor:"0.01"`
		Y float32 `type:"int16" len:"2" unit:"G" factor:"0.01"`
		Z float32 `type:"int16" len:"2" unit:"G" factor:"0.01"`
	}
	Gyro struct {
		X float32 `type:"int16" len:"2" unit:"rad/s" factor:"0.1"`
		Y float32 `type:"int16" len:"2" unit:"rad/s" factor:"0.1"`
		Z float32 `type:"int16" len:"2" unit:"rad/s" factor:"0.1"`
	}
	Tilt struct {
		Pitch float32 `type:"int16" len:"2" unit:"Deg" factor:"0.1"`
		Roll  float32 `type:"int16" len:"2" unit:"Deg" factor:"0.1"`
	}
	Total struct {
		Accel float32 `type:"uint16" len:"2" unit:"G" factor:"0.01"`
		Gyro  float32 `type:"uint16" len:"2" unit:"rad/s" factor:"0.1"`
		Tilt  float32 `type:"uint16" len:"2" unit:"Deg" factor:"0.1"`
		Temp  float32 `type:"uint16" len:"2" unit:"Celcius" factor:"0.1"`
	}
}

type Remote struct {
	Active bool `type:"uint8"`
	Nearby bool `type:"uint8"`
}

type Finger struct {
	Verified bool  `type:"uint8"`
	DriverID uint8 `type:"uint8"`
}

type Audio struct {
	Active bool  `type:"uint8"`
	Mute   bool  `type:"uint8"`
	Volume uint8 `type:"uint8" unit:"%"`
}
type Hmi struct {
	Active bool `type:"uint8"`
}

type Bms struct {
	Active bool   `type:"uint8"`
	Run    bool   `type:"uint8"`
	SOC    uint8  `type:"uint8" unit:"%"`
	Fault  uint16 `type:"uint16"`
	Pack   [BMS_PACK_MAX]struct {
		ID      uint32  `type:"uint32"`
		Fault   uint16  `type:"uint16"`
		Voltage float32 `type:"uint16" len:"2" unit:"Volt" factor:"0.01"`
		Current float32 `type:"uint16" len:"2" unit:"Ampere" factor:"0.1"`
		SOC     uint8   `type:"uint8" unit:"%"`
		Temp    uint16  `type:"uint16" unit:"Celcius"`
	}
}

// Faults parse b's fault field
// func (b *Bms) Faults() BmsFault {
// 	TODO: implement me
// }

// LowCapacity check if b's SoC (State of Charge) is low
func (b *Bms) LowCapacity() bool {
	if b == nil {
		return false
	}
	return b.SOC < BMS_LOW_CAPACITY_PERCENT
}

type Mcu struct {
	Active    bool      `type:"uint8"`
	Run       bool      `type:"uint8"`
	Reverse   bool      `type:"uint8"`
	DriveMode ModeDrive `type:"uint8"`
	Speed     uint8     `type:"uint8" unit:"Kph"`
	RPM       int16     `type:"int16" unit:"rpm"`
	Temp      float32   `type:"uint16" len:"2" unit:"Celcius" factor:"0.1"`
	Fault     struct {
		Post uint32 `type:"uint32"`
		Run  uint32 `type:"uint32"`
	}
	Torque struct {
		Command  float32 `type:"uint16" len:"2" unit:"Nm" factor:"0.1"`
		Feedback float32 `type:"uint16" len:"2" unit:"Nm" factor:"0.1"`
	}
	DCBus struct {
		Current float32 `type:"uint16" len:"2" unit:"A" factor:"0.1"`
		Voltage float32 `type:"uint16" len:"2" unit:"V" factor:"0.1"`
	}
	Inverter struct {
		Enabled   bool            `type:"uint8"`
		Lockout   bool            `type:"uint8"`
		Discharge McuInvDischarge `type:"uint8"`
	}
	Template struct {
		MaxRPM    int16 `type:"int16" unit:"rpm"`
		MaxSpeed  uint8 `type:"uint8" unit:"Kph"`
		DriveMode [ModeDriveLimit]struct {
			Discur uint16  `type:"uint16" unit:"A"`
			Torque float32 `type:"uint16" len:"2" unit:"Nm" factor:"0.1"`
		}
	}
}

// Fault parse m's fault field
// func (m *Mcu) Faults() McuFault {
// 	TODO: implement me
// }

type Task struct {
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

// StackOverflow check if t's stack is near overflow fault
func (t *Task) StackOverflow() bool {
	if t == nil {
		return false
	}
	rv := reflect.ValueOf(t.Stack)
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).Uint() < STACK_OVERFLOW_BYTE_MIN {
			return true
		}
	}
	return false
}
