package sdk

import (
	"math/rand"
	"time"
)

func makeCommandPacket(vin int, cmd *command, msg message) *commandPacket {
	return &commandPacket{
		Header: &HeaderCommand{
			Header: Header{
				Prefix:       PREFIX_COMMAND,
				Size:         0,
				Vin:          uint32(vin),
				SendDatetime: time.Now(),
			},
			Code:    cmd.code,
			SubCode: cmd.subCode,
		},
		Message: msg,
	}
}

func makeResponsePacket(vin int, cmd *command, msg message) *responsePacket {
	return &responsePacket{
		Header: &headerResponse{
			HeaderCommand: HeaderCommand{
				Header: Header{
					Prefix:       PREFIX_RESPONSE,
					Size:         0,
					Vin:          uint32(vin),
					SendDatetime: time.Now(),
				},
				Code:    cmd.code,
				SubCode: cmd.subCode,
			},
			ResCode: resCodeOk,
		},
		Message: msg,
	}
}

func makeReportPacket(vin int, full Frame) *ReportPacket {
	rand.Seed(time.Now().UnixNano())

	rp := &ReportPacket{
		Header: &HeaderReport{
			Header: Header{
				Prefix:       PREFIX_REPORT,
				Size:         0,
				Vin:          uint32(vin),
				SendDatetime: time.Now(),
			},
			Frame: full, // Frame(rand.Intn(int(FrameLimit))),
		},
		Vcu: &Vcu{
			LogDatetime: time.Now().Add(-2 * time.Second),
			State:       BikeState(rand.Intn(int(BikeStateLimit))),
			Events:      uint16(rand.Uint32()),
			LogBuffered: uint8(rand.Intn(50)),
			BatVoltage:  randFloat(0, 4400),
			Uptime:      randFloat(0, 1000),
		},
		Eeprom: &Eeprom{
			Active: randBool(),
			Used:   uint8(rand.Intn(100)),
		},
		Gps: &Gps{
			Active:    randBool(),
			SatInUse:  uint8(rand.Intn(14)),
			HDOP:      randFloat(0, 10),
			VDOP:      randFloat(0, 10),
			Speed:     uint8(rand.Intn(SPEED_KPH_MAX)),
			Heading:   randFloat(0, 360),
			Longitude: randFloat(GPS_LNG_MIN, GPS_LNG_MAX),
			Latitude:  randFloat(GPS_LAT_MIN, GPS_LAT_MAX),
		},
		Hbar: &Hbar{
			Reverse: randBool(),
			Mode: struct {
				Drive ModeDrive "type:\"uint8\""
				Trip  ModeTrip  "type:\"uint8\""
				Avg   ModeAvg   "type:\"uint8\""
			}{
				Drive: ModeDrive(rand.Intn(int(ModeDriveLimit))),
				Trip:  ModeTrip(rand.Intn(int(ModeTripLimit))),
				Avg:   ModeAvg(rand.Intn(int(ModeAvgLimit))),
			},
			Trip: struct {
				Odo uint16 "type:\"uint16\" unit:\"Km\""
				A   uint16 "type:\"uint16\" unit:\"Km\""
				B   uint16 "type:\"uint16\" unit:\"Km\""
			}{
				Odo: uint16(rand.Intn(99999)),
				A:   uint16(rand.Intn(99999)),
				B:   uint16(rand.Intn(99999)),
			},
			Avg: struct {
				Range      uint8 "type:\"uint8\" unit:\"Km\""
				Efficiency uint8 "type:\"uint8\" unit:\"Km/Kwh\""
			}{
				Range:      uint8(rand.Intn(255)),
				Efficiency: uint8(rand.Intn(255)),
			},
		},
		Net: &Net{
			Signal:   uint8(rand.Intn(100)),
			State:    NetState(rand.Intn(int(NetStateLimit))),
			IpStatus: NetIpStatus(rand.Intn(int(NetIpStatusLimit))),
		},
		Mems: &Mems{
			Active:    randBool(),
			AntiThief: randBool(),
			Accel: struct {
				X float32 "type:\"int16\" len:\"2\" unit:\"G\" factor:\"0.01\""
				Y float32 "type:\"int16\" len:\"2\" unit:\"G\" factor:\"0.01\""
				Z float32 "type:\"int16\" len:\"2\" unit:\"G\" factor:\"0.01\""
			}{
				X: randFloat(0, 100),
				Y: randFloat(0, 100),
				Z: randFloat(0, 100),
			},
			Gyro: struct {
				X float32 "type:\"int16\" len:\"2\" unit:\"rad/s\" factor:\"0.1\""
				Y float32 "type:\"int16\" len:\"2\" unit:\"rad/s\" factor:\"0.1\""
				Z float32 "type:\"int16\" len:\"2\" unit:\"rad/s\" factor:\"0.1\""
			}{
				X: randFloat(0, 10000),
				Y: randFloat(0, 10000),
				Z: randFloat(0, 10000),
			},
			Tilt: struct {
				Pitch float32 "type:\"int16\" len:\"2\" unit:\"Deg\" factor:\"0.1\""
				Roll  float32 "type:\"int16\" len:\"2\" unit:\"Deg\" factor:\"0.1\""
			}{
				Pitch: randFloat(0, 180),
				Roll:  randFloat(0, 180),
			},
			Total: struct {
				Accel float32 "type:\"uint16\" len:\"2\" unit:\"G\" factor:\"0.01\""
				Gyro  float32 "type:\"uint16\" len:\"2\" unit:\"rad/s\" factor:\"0.1\""
				Tilt  float32 "type:\"uint16\" len:\"2\" unit:\"Deg\" factor:\"0.1\""
				Temp  float32 "type:\"uint16\" len:\"2\" unit:\"Celcius\" factor:\"0.1\""
			}{
				Accel: randFloat(0, 100),
				Gyro:  randFloat(0, 10000),
				Tilt:  randFloat(0, 180),
				Temp:  randFloat(30, 50),
			},
		},
		Remote: &Remote{
			Active: randBool(),
			Nearby: randBool(),
		},
		Finger: &Finger{
			Active:   randBool(),
			DriverID: uint8(rand.Intn(DRIVER_ID_MAX)),
		},
		Audio: &Audio{
			Active: randBool(),
			Mute:   randBool(),
			Volume: uint8(rand.Intn(100)),
		},
		Hmi: &Hmi{
			Active: randBool(),
		},
		Bms: &Bms{
			Active: randBool(),
			Run:    randBool(),
			SOC:    uint8(rand.Intn(100)),
			Faults: uint16(rand.Uint32()),
			Pack: [2]struct {
				ID      uint32  "type:\"uint32\""
				Fault   uint16  "type:\"uint16\""
				Voltage float32 "type:\"uint16\" len:\"2\" unit:\"Volt\" factor:\"0.01\""
				Current float32 "type:\"uint16\" len:\"2\" unit:\"Ampere\" factor:\"0.1\""
				SOC     uint8   "type:\"uint8\" unit:\"%\""
				Temp    uint16  "type:\"uint16\" unit:\"Celcius\""
			}{
				{
					ID:      rand.Uint32(),
					Fault:   uint16(rand.Uint32()),
					Voltage: randFloat(48, 60),
					Current: randFloat(0, 110),
					SOC:     uint8(rand.Intn(100)),
					Temp:    uint16(randFloat(30, 50)),
				},
				{
					Fault:   uint16(rand.Uint32()),
					Voltage: randFloat(48, 60),
					Current: randFloat(0, 110),
					SOC:     uint8(rand.Intn(100)),
					Temp:    uint16(randFloat(30, 50)),
				},
			},
		},
		Mcu: &Mcu{
			Active:    randBool(),
			Run:       randBool(),
			Reverse:   randBool(),
			DriveMode: ModeDrive(rand.Intn(int(ModeDriveLimit))),
			Speed:     uint8(rand.Intn(SPEED_KPH_MAX)),
			RPM:       int16(rand.Intn(50000) - 25000),
			Temp:      randFloat(30, 50),
			Faults: struct {
				Post uint32 "type:\"uint32\""
				Run  uint32 "type:\"uint32\""
			}{
				Post: rand.Uint32(),
				Run:  rand.Uint32(),
			},
			Torque: struct {
				Command  float32 "type:\"uint16\" len:\"2\" unit:\"Nm\" factor:\"0.1\""
				Feedback float32 "type:\"uint16\" len:\"2\" unit:\"Nm\" factor:\"0.1\""
			}{
				Command:  rand.Float32(),
				Feedback: rand.Float32(),
			},
			DCBus: struct {
				Current float32 "type:\"uint16\" len:\"2\" unit:\"A\" factor:\"0.1\""
				Voltage float32 "type:\"uint16\" len:\"2\" unit:\"V\" factor:\"0.1\""
			}{
				Current: rand.Float32(),
				Voltage: rand.Float32(),
			},
			Inverter: struct {
				Enabled   bool            "type:\"uint8\""
				Lockout   bool            "type:\"uint8\""
				Discharge McuInvDischarge "type:\"uint8\""
			}{
				Enabled:   randBool(),
				Lockout:   randBool(),
				Discharge: McuInvDischarge(rand.Intn(int(McuInvDischargeLimit))),
			},
			Template: struct {
				MaxRPM    int16 "type:\"int16\" unit:\"rpm\""
				MaxSpeed  uint8 "type:\"uint8\" unit:\"Kph\""
				DriveMode [3]struct {
					Discur uint16  "type:\"uint16\" unit:\"A\""
					Torque float32 "type:\"uint16\" len:\"2\" unit:\"Nm\" factor:\"0.1\""
				}
			}{
				MaxRPM:   int16(rand.Intn(50000) - 25000),
				MaxSpeed: uint8(rand.Intn(SPEED_KPH_MAX)),
				DriveMode: [3]struct {
					Discur uint16  "type:\"uint16\" unit:\"A\""
					Torque float32 "type:\"uint16\" len:\"2\" unit:\"Nm\" factor:\"0.1\""
				}{
					{
						Discur: uint16(rand.Float32()),
						Torque: rand.Float32(),
					},
					{
						Discur: uint16(rand.Float32()),
						Torque: rand.Float32(),
					},
					{
						Discur: uint16(rand.Float32()),
						Torque: rand.Float32(),
					},
				},
			},
		},
		Task: &Task{
			Stack: struct {
				Manager  uint16 "type:\"uint16\" unit:\"Bytes\""
				Network  uint16 "type:\"uint16\" unit:\"Bytes\""
				Reporter uint16 "type:\"uint16\" unit:\"Bytes\""
				Command  uint16 "type:\"uint16\" unit:\"Bytes\""
				Mems     uint16 "type:\"uint16\" unit:\"Bytes\""
				Remote   uint16 "type:\"uint16\" unit:\"Bytes\""
				Finger   uint16 "type:\"uint16\" unit:\"Bytes\""
				Audio    uint16 "type:\"uint16\" unit:\"Bytes\""
				Gate     uint16 "type:\"uint16\" unit:\"Bytes\""
				CanRX    uint16 "type:\"uint16\" unit:\"Bytes\""
				CanTX    uint16 "type:\"uint16\" unit:\"Bytes\""
			}{
				Manager:  uint16(rand.Intn(1000)),
				Network:  uint16(rand.Intn(1000)),
				Reporter: uint16(rand.Intn(1000)),
				Command:  uint16(rand.Intn(1000)),
				Mems:     uint16(rand.Intn(1000)),
				Remote:   uint16(rand.Intn(1000)),
				Finger:   uint16(rand.Intn(1000)),
				Audio:    uint16(rand.Intn(1000)),
				Gate:     uint16(rand.Intn(1000)),
				CanRX:    uint16(rand.Intn(1000)),
				CanTX:    uint16(rand.Intn(1000)),
			},
			Wakeup: struct {
				Manager  uint8 "type:\"uint8\" unit:\"s\""
				Network  uint8 "type:\"uint8\" unit:\"s\""
				Reporter uint8 "type:\"uint8\" unit:\"s\""
				Command  uint8 "type:\"uint8\" unit:\"s\""
				Mems     uint8 "type:\"uint8\" unit:\"s\""
				Remote   uint8 "type:\"uint8\" unit:\"s\""
				Finger   uint8 "type:\"uint8\" unit:\"s\""
				Audio    uint8 "type:\"uint8\" unit:\"s\""
				Gate     uint8 "type:\"uint8\" unit:\"s\""
				CanRX    uint8 "type:\"uint8\" unit:\"s\""
				CanTX    uint8 "type:\"uint8\" unit:\"s\""
			}{
				Manager:  uint8(rand.Intn(255)),
				Network:  uint8(rand.Intn(255)),
				Reporter: uint8(rand.Intn(255)),
				Command:  uint8(rand.Intn(255)),
				Mems:     uint8(rand.Intn(255)),
				Remote:   uint8(rand.Intn(255)),
				Finger:   uint8(rand.Intn(255)),
				Audio:    uint8(rand.Intn(255)),
				Gate:     uint8(rand.Intn(255)),
				CanRX:    uint8(rand.Intn(255)),
				CanTX:    uint8(rand.Intn(255)),
			},
		},
	}

	if rp.Header.Frame == FrameSimple {
		rp.Hbar = nil
		rp.Net = nil
		rp.Mems = nil
		rp.Remote = nil
		rp.Finger = nil
		rp.Audio = nil
		rp.Hmi = nil
		rp.Bms = nil
		rp.Mcu = nil
		rp.Task = nil
	}

	return rp
}
