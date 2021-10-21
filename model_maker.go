package sdk

import (
	"math/rand"
	"time"
)

func makeCommandPacket(vin int, cmd *command, msg message) *commandPacket {
	return &commandPacket{
		Header: &HeaderCommand{
			Header: Header{
				Prefix:  PREFIX_COMMAND,
				Size:    0,
				Version: uint16(SDK_VERSION),
				Vin:     uint32(vin),
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
					Prefix:  PREFIX_RESPONSE,
					Size:    0,
					Version: uint16(SDK_VERSION),
					Vin:     uint32(vin),
				},
				Code:    cmd.code,
				SubCode: cmd.subCode,
			},
			ResCode: resCodeOk,
		},
		Message: msg,
	}
}

func makeReportPacket(version int, vin int, frame Frame) *ReportPacket {
	rand.Seed(time.Now().UnixNano())

	// rp := &ReportPacket{
	// };

	rp := &ReportPacket{
		Header: Header{
			Prefix:  PREFIX_REPORT,
			Size:    0,
			Version: uint16(version),
			Vin:     uint32(vin),
		},
	}

	rp.Data = PacketData{
		"Report": PacketData{
			"SendDatetime": time.Now(),
			"LogDatetime":  time.Now().Add(-2 * time.Second),
			"Frame":        frame,
		},

		"Vcu": PacketData{
			"State":       BikeState(rand.Intn(int(BikeStateLimit))),
			"Events":      uint16(rand.Uint32()),
			"LogBuffered": uint8(rand.Intn(50)),
			"BatVoltage":  randFloat(0, 4400),
			"Uptime":      randFloat(0, 1000),
			"LockDown":    randBool(),
		},

		"Eeprom": PacketData{
			"Active": randBool(),
			"Used":   uint8(rand.Intn(100)),
		},

		"Gps": PacketData{
			"Active":    randBool(),
			"SatInUse":  uint8(rand.Intn(14)),
			"HDOP":      randFloat(0, 10),
			"VDOP":      randFloat(0, 10),
			"Speed":     uint8(rand.Intn(SPEED_KPH_MAX)),
			"Heading":   randFloat(0, 360),
			"Longitude": randFloat(GPS_LNG_MIN, GPS_LNG_MAX),
			"Latitude":  randFloat(GPS_LAT_MIN, GPS_LAT_MAX),
		},

		"Hbar": PacketData{
			"Reverse": randBool(),
			"Mode": PacketData{
				"Drive": ModeDrive(rand.Intn(int(ModeDriveLimit))),
				"Trip":  ModeTrip(rand.Intn(int(ModeTripLimit))),
				"Avg":   ModeAvg(rand.Intn(int(ModeAvgLimit))),
			},
			"Trip": PacketData{
				"Odo": uint16(rand.Intn(99999)),
				"A":   uint16(rand.Intn(99999)),
				"B":   uint16(rand.Intn(99999)),
			},
			"Avg": PacketData{
				"Range":      uint8(rand.Intn(255)),
				"Efficiency": uint8(rand.Intn(255)),
			},
		},

		"Net": PacketData{
			"Signal": uint8(rand.Intn(100)),
			"State":  NetState(rand.Intn(int(NetStateLimit))),
		},

		"Imu": PacketData{
			"Active":    randBool(),
			"AntiThief": randBool(),
			"Accel": PacketData{
				"X": randFloat(0, 100),
				"Y": randFloat(0, 100),
				"Z": randFloat(0, 100),
			},
			"Gyro": PacketData{
				"X": randFloat(0, 10000),
				"Y": randFloat(0, 10000),
				"Z": randFloat(0, 10000),
			},
			"Tilt": PacketData{
				"Pitch": randFloat(0, 180),
				"Roll":  randFloat(0, 180),
			},
			"Total": PacketData{
				"Accel": randFloat(0, 100),
				"Gyro":  randFloat(0, 10000),
				"Tilt":  randFloat(0, 180),
				"Temp":  randFloat(30, 50),
			},
		},

		"Remote": PacketData{
			"Active": randBool(),
			"Nearby": randBool(),
		},

		"Finger": PacketData{
			"Active":   randBool(),
			"DriverID": uint8(rand.Intn(DRIVER_ID_MAX)),
		},

		"Audio": PacketData{
			"Active": randBool(),
			"Mute":   randBool(),
			"Volume": uint8(rand.Intn(100)),
		},

		"Hmi": PacketData{
			"Active": randBool(),
		},

		"Bms": PacketData{
			"Active": randBool(),
			"Run":    randBool(),
			"Capacity": PacketData{
				"Remaining": uint16(rand.Intn(2100)),
				"Usage":     uint16(rand.Intn(2100)),
			},
			"SOC":    uint8(rand.Intn(100)),
			"Faults": uint16(rand.Uint32()),
			"Pack": [2]PacketData{
				{
					"ID":      rand.Uint32(),
					"Fault":   uint16(rand.Uint32()),
					"Voltage": randFloat(48, 60),
					"Current": randFloat(0, 110),
					"Capacity": PacketData{
						"Remaining": uint16(rand.Intn(2100)),
						"Usage":     uint16(rand.Intn(2100)),
					},
					"SOC":  uint8(rand.Intn(100)),
					"SOH":  uint8(rand.Intn(100)),
					"Temp": uint16(randFloat(30, 50)),
				},
				{
					"Fault":   uint16(rand.Uint32()),
					"Voltage": randFloat(48, 60),
					"Current": randFloat(0, 110),
					"Capacity": PacketData{
						"Remaining": uint16(rand.Intn(2100)),
						"Usage":     uint16(rand.Intn(2100)),
					},
					"SOC":  uint8(rand.Intn(100)),
					"SOH":  uint8(rand.Intn(100)),
					"Temp": uint16(randFloat(30, 50)),
				},
			},
		},

		"Mcu": PacketData{
			"Active":    randBool(),
			"Run":       randBool(),
			"Reverse":   randBool(),
			"DriveMode": ModeDrive(rand.Intn(int(ModeDriveLimit))),
			"Speed":     uint8(rand.Intn(SPEED_KPH_MAX)),
			"RPM":       int16(rand.Intn(50000) - 25000),
			"Temp":      randFloat(30, 50),
			"Faults": PacketData{
				"Post": rand.Uint32(),
				"Run":  rand.Uint32(),
			},
			"Torque": PacketData{
				"Command":  rand.Float32(),
				"Feedback": rand.Float32(),
			},
			"DCBus": PacketData{
				"Current": rand.Float32(),
				"Voltage": rand.Float32(),
			},
			"Template": PacketData{
				"MaxRPM":   int16(rand.Intn(50000) - 25000),
				"MaxSpeed": uint8(rand.Intn(SPEED_KPH_MAX)),
				"DriveMode": [3]PacketData{
					{
						"Discur": uint8(rand.Intn(MCU_DISCUR_MAX)),
						"Torque": uint8(rand.Intn(MCU_TORQUE_MAX)),
					},
					{
						"Discur": uint8(rand.Intn(MCU_DISCUR_MAX)),
						"Torque": uint8(rand.Intn(MCU_TORQUE_MAX)),
					},
					{
						"Discur": uint8(rand.Intn(MCU_DISCUR_MAX)),
						"Torque": uint8(rand.Intn(MCU_TORQUE_MAX)),
					},
				},
			},
		},

		"Task": PacketData{
			"Stack": PacketData{
				"Manager":  uint16(rand.Intn(1000)),
				"Network":  uint16(rand.Intn(1000)),
				"Reporter": uint16(rand.Intn(1000)),
				"Command":  uint16(rand.Intn(1000)),
				"Imu":      uint16(rand.Intn(1000)),
				"Remote":   uint16(rand.Intn(1000)),
				"Finger":   uint16(rand.Intn(1000)),
				"Audio":    uint16(rand.Intn(1000)),
				"Gate":     uint16(rand.Intn(1000)),
				"CanRX":    uint16(rand.Intn(1000)),
				"CanTX":    uint16(rand.Intn(1000)),
			},
			"Wakeup": PacketData{
				"Manager":  uint8(rand.Intn(255)),
				"Network":  uint8(rand.Intn(255)),
				"Reporter": uint8(rand.Intn(255)),
				"Command":  uint8(rand.Intn(255)),
				"Imu":      uint8(rand.Intn(255)),
				"Remote":   uint8(rand.Intn(255)),
				"Finger":   uint8(rand.Intn(255)),
				"Audio":    uint8(rand.Intn(255)),
				"Gate":     uint8(rand.Intn(255)),
				"CanRX":    uint8(rand.Intn(255)),
				"CanTX":    uint8(rand.Intn(255)),
			},
		},
	}

	// if rp.Header.Frame == FrameSimple {
	// 	rp.Hbar = nil
	// 	rp.Net = nil
	// 	rp.Imu = nil
	// 	rp.Remote = nil
	// 	rp.Finger = nil
	// 	rp.Audio = nil
	// 	rp.Hmi = nil
	// 	rp.Bms = nil
	// 	rp.Mcu = nil
	// 	rp.Task = nil
	// }

	return rp
}
