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
			"Queued":       uint8(rand.Intn(50)),
		},

		"Vcu": PacketData{
			"State":      BikeState(rand.Intn(int(BikeStateLimit))),
			"Events":     uint16(rand.Uint32()),
			"Version":    uint16(rand.Uint32()),
			"BatVoltage": randFloat(0, 4400),
			"Uptime":     randFloat(0, 1000),
			"LockDown":   randBool(),
			"CANDebug":   uint8(rand.Intn(100)),
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
			"Altitude":  randFloat(0, 10),
		},

		"Net": PacketData{
			"Signal": uint8(rand.Intn(100)),
			"State":  NetState(rand.Intn(int(NetStateLimit))),
		},

		"Imu": PacketData{
			"Active":    randBool(),
			"AntiThief": randBool(),
			"Tilt": PacketData{
				"Pitch": randFloat(0, 180),
				"Roll":  randFloat(0, 180),
			},
			"Total": PacketData{
				"Accel":       randFloat(0, 100),
				"Gyro":        randFloat(0, 10000),
				"Tilt":        randFloat(0, 180),
				"Temperature": int8(rand.Intn(128)),
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
			"Active":  randBool(),
			"Version": uint16(rand.Uint32()),
		},

		"Bms": PacketData{
			"Active": randBool(),
			"Run":    randBool(),
			"Faults": uint16(rand.Uint32()),
			"SOC":    uint8(rand.Intn(100)),
			"Pack": [2]PacketData{
				{
					"ID":      rand.Uint32(),
					"Faults":  uint16(rand.Uint32()),
					"Voltage": randFloat(48, 60),
					"Current": randFloat(0, 110),
					"Capacity": PacketData{
						"Remaining": uint16(rand.Intn(2100)),
						"Usage":     uint16(rand.Intn(2100)),
					},
					"SOC":         uint8(rand.Intn(100)),
					"SOH":         uint8(rand.Intn(100)),
					"Temperature": int8(rand.Intn(128)),
				},
				{
					"ID":      rand.Uint32(),
					"Faults":  uint16(rand.Uint32()),
					"Voltage": randFloat(48, 60),
					"Current": randFloat(0, 110),
					"Capacity": PacketData{
						"Remaining": uint16(rand.Intn(2100)),
						"Usage":     uint16(rand.Intn(2100)),
					},
					"SOC":         uint8(rand.Intn(100)),
					"SOH":         uint8(rand.Intn(100)),
					"Temperature": int8(rand.Intn(128)),
				},
			},
		},

		"Hbar": PacketData{
			"Reverse": randBool(),
			"Mode": PacketData{
				"Drive": ModeDrive(rand.Intn(int(ModeDriveLimit))),
				"Trip":  ModeTrip(rand.Intn(int(ModeTripLimit))),
				"Avg":   ModeAvg(rand.Intn(int(ModeAvgLimit))),
			},
			"Trip": PacketData{
				"Odo": uint16(rand.Intn(TRIP_KM_MAX)),
				"A":   uint16(rand.Intn(TRIP_KM_MAX)),
				"B":   uint16(rand.Intn(TRIP_KM_MAX)),
			},
			"Avg": PacketData{
				"Range":      uint8(rand.Intn(255)),
				"Efficiency": uint8(rand.Intn(255)),
			},
		},

		"Mcu": PacketData{
			"Active":      randBool(),
			"Run":         randBool(),
			"Reverse":     randBool(),
			"DriveMode":   ModeDrive(rand.Intn(int(ModeDriveLimit))),
			"Speed":       uint8(rand.Intn(SPEED_KPH_MAX)),
			"RPM":         int16(rand.Intn(50000) - 25000),
			"Temperature": int8(rand.Intn(128)),
			"Faults": PacketData{
				"Post": rand.Uint32(),
				"Run":  rand.Uint32(),
			},
			"Torque": PacketData{
				"Commanded": rand.Float32(),
				"Feedback":  rand.Float32(),
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
