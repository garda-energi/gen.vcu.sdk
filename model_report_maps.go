package sdk

// version : structure
var ReportPacketStructures = map[int]tagger{
	1: {
		Tipe: Struct_t, Sub: []tagger{
			{Name: "Report", Tipe: Struct_t, Sub: []tagger{
				{Name: "SendDatetime", Tipe: Time_t, Len: 7},
				{Name: "LogDatetime", Tipe: Time_t, Len: 7},
				{Name: "Frame", Tipe: Uint8_t},
				{Name: "Queued", Tipe: Uint8_t},
			}},
			{Name: "Vcu", Tipe: Struct_t, Sub: []tagger{
				{Name: "State", Tipe: Int8_t},
				{Name: "Events", Tipe: Uint16_t},
				{Name: "Version", Tipe: Uint16_t},
				{Name: "BatVoltage", Tipe: Float_t, Len: 1, Factor: 18.0},
				{Name: "Uptime", Tipe: Float_t, Len: 4, Factor: 0.000277},
				{Name: "LockDown", Tipe: Boolean_t},
				{Name: "CANDebug", Tipe: Boolean_t},
			}},
			{Name: "Eeprom", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Used", Tipe: Uint8_t},
			}},
			{Name: "Gps", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "SatInUse", Tipe: Uint8_t},
				{Name: "HDOP", Tipe: Float_t, Len: 1, Factor: 0.1},
				{Name: "VDOP", Tipe: Float_t, Len: 1, Factor: 0.1},
				{Name: "Speed", Tipe: Uint8_t},
				{Name: "Heading", Tipe: Float_t, Len: 1, Factor: 2.0},
				{Name: "Longitude", Tipe: Float_t, UnfactorType: Int32_t, Len: 4, Factor: 0.0000001},
				{Name: "Latitude", Tipe: Float_t, UnfactorType: Int32_t, Len: 4, Factor: 0.0000001},
				{Name: "Altitude", Tipe: Float_t, UnfactorType: Int16_t, Len: 2, Factor: 0.1},
			}},
			{Name: "Net", Tipe: Struct_t, Sub: []tagger{
				{Name: "Signal", Tipe: Uint8_t},
				{Name: "State", Tipe: Int8_t},
			}},
			{Name: "Imu", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "AntiThief", Tipe: Boolean_t},
				{Name: "Tilt", Tipe: Struct_t, Sub: []tagger{
					{Name: "Pitch", Tipe: Float_t, UnfactorType: Int16_t, Len: 2, Factor: 0.1},
					{Name: "Roll", Tipe: Float_t, UnfactorType: Int16_t, Len: 2, Factor: 0.1},
				}},
				{Name: "Total", Tipe: Struct_t, Sub: []tagger{
					{Name: "Accel", Tipe: Float_t, Len: 2, Factor: 0.01},
					{Name: "Gyro", Tipe: Float_t, Len: 2, Factor: 0.1},
					{Name: "Tilt", Tipe: Float_t, Len: 2, Factor: 0.1},
					{Name: "Temperature", Tipe: Int8_t},
				}},
			}},
			{Name: "Remote", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Nearby", Tipe: Boolean_t},
			}},
			{Name: "Finger", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "DriverID", Tipe: Uint8_t},
			}},
			{Name: "Audio", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Mute", Tipe: Uint8_t},
				{Name: "Volume", Tipe: Uint8_t},
			}},
			{Name: "Hmi", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Version", Tipe: Uint16_t},
			}},
			{Name: "Bms", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Run", Tipe: Boolean_t},
				{Name: "Faults", Tipe: Uint16_t},
				{Name: "SOC", Tipe: Uint8_t},
				{Name: "Pack", Tipe: Array_t, Len: 2, Sub: []tagger{
					{Tipe: Struct_t, Sub: []tagger{
						{Name: "ID", Tipe: Uint32_t},
						{Name: "Faults", Tipe: Uint16_t},
						{Name: "Voltage", Tipe: Float_t, Len: 2, Factor: 0.01},
						{Name: "Current", Tipe: Float_t, Len: 2, Factor: 0.1},
						{Name: "Capacity", Tipe: Struct_t, Sub: []tagger{
							{Name: "Remaining", Tipe: Uint16_t},
							{Name: "Usage", Tipe: Uint16_t},
						}},
						{Name: "SOC", Tipe: Uint8_t},
						{Name: "SOH", Tipe: Uint8_t},
						{Name: "Temperature", Tipe: Int8_t},
					}},
				}},
			}},
			{Name: "Hbar", Tipe: Struct_t, Sub: []tagger{
				{Name: "Reverse", Tipe: Boolean_t},
				{Name: "Mode", Tipe: Struct_t, Sub: []tagger{
					{Name: "Drive", Tipe: Uint8_t},
					{Name: "Trip", Tipe: Uint8_t},
					{Name: "Avg", Tipe: Uint8_t},
				}},
				{Name: "Trip", Tipe: Struct_t, Sub: []tagger{
					{Name: "Odo", Tipe: Uint16_t},
					{Name: "A", Tipe: Uint16_t},
					{Name: "B", Tipe: Uint16_t},
				}},
				{Name: "Avg", Tipe: Struct_t, Sub: []tagger{
					{Name: "Range", Tipe: Uint8_t},
					{Name: "Efficiency", Tipe: Uint8_t},
				}},
			}},
			{Name: "Mcu", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Run", Tipe: Boolean_t},
				{Name: "Reverse", Tipe: Boolean_t},
				{Name: "DriveMode", Tipe: Uint8_t},
				{Name: "Speed", Tipe: Uint8_t},
				{Name: "RPM", Tipe: Int16_t},
				{Name: "Temperature", Tipe: Int8_t},
				{Name: "Faults", Tipe: Struct_t, Sub: []tagger{
					{Name: "Post", Tipe: Uint32_t},
					{Name: "Run", Tipe: Uint32_t},
				}},
				{Name: "Torque", Tipe: Struct_t, Sub: []tagger{
					{Name: "Commanded", Tipe: Float_t, Len: 2, Factor: 0.1},
					{Name: "Feedback", Tipe: Float_t, Len: 2, Factor: 0.1},
				}},
				{Name: "DCBus", Tipe: Struct_t, Sub: []tagger{
					{Name: "Current", Tipe: Float_t, Len: 2, Factor: 0.1},
					{Name: "Voltage", Tipe: Float_t, Len: 2, Factor: 0.1},
				}},
				{Name: "Template", Tipe: Struct_t, Sub: []tagger{
					{Name: "MaxRPM", Tipe: Int16_t},
					{Name: "MaxSpeed", Tipe: Uint8_t},
					{Name: "DriveMode", Tipe: Array_t, Len: 3, Sub: []tagger{
						{Tipe: Struct_t, Sub: []tagger{
							{Name: "Discur", Tipe: Uint8_t},
							{Name: "Torque", Tipe: Uint8_t},
						}},
					}},
				}},
			}},
			{Name: "Task", Tipe: Struct_t, Sub: []tagger{
				{Name: "Stack", Tipe: Struct_t, Sub: []tagger{
					{Name: "Manager", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Network", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Reporter", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Command", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Imu", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Remote", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Finger", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Audio", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Gate", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanRX", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanTX", Tipe: Uint8_t, Sub: []tagger{}},
				}},
				{Name: "Wakeup", Tipe: Struct_t, Sub: []tagger{
					{Name: "Manager", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Network", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Reporter", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Command", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Imu", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Remote", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Finger", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Audio", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Gate", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanRX", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanTX", Tipe: Uint8_t, Sub: []tagger{}},
				}},
			}},
		},
	},
	2: {
		Tipe: Struct_t, Sub: []tagger{
			{Name: "Report", Tipe: Struct_t, Sub: []tagger{
				{Name: "SendDatetime", Tipe: Time_t, Len: 7},
				{Name: "LogDatetime", Tipe: Time_t, Len: 7},
				{Name: "Frame", Tipe: Uint8_t},
				{Name: "Queued", Tipe: Uint8_t},
			}},
			{Name: "Vcu", Tipe: Struct_t, Sub: []tagger{
				{Name: "State", Tipe: Int8_t},
				{Name: "Events", Tipe: Uint16_t},
				{Name: "Version", Tipe: Uint16_t},
				{Name: "BatVoltage", Tipe: Float_t, Len: 1, Factor: 18.0},
				{Name: "Uptime", Tipe: Float_t, Len: 4, Factor: 0.000277},
				{Name: "LockDown", Tipe: Boolean_t},
				{Name: "CANDebug", Tipe: Boolean_t},
			}},
			{Name: "Eeprom", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Used", Tipe: Uint8_t},
			}},
			{Name: "Gps", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "SatInUse", Tipe: Uint8_t},
				{Name: "HDOP", Tipe: Float_t, Len: 1, Factor: 0.1},
				{Name: "VDOP", Tipe: Float_t, Len: 1, Factor: 0.1},
				{Name: "Speed", Tipe: Uint8_t},
				{Name: "Heading", Tipe: Float_t, Len: 1, Factor: 2.0},
				{Name: "Longitude", Tipe: Float_t, UnfactorType: Int32_t, Len: 4, Factor: 0.0000001},
				{Name: "Latitude", Tipe: Float_t, UnfactorType: Int32_t, Len: 4, Factor: 0.0000001},
				{Name: "Altitude", Tipe: Float_t, UnfactorType: Int16_t, Len: 2, Factor: 0.1},
			}},
			{Name: "Net", Tipe: Struct_t, Sub: []tagger{
				{Name: "Signal", Tipe: Uint8_t},
				{Name: "State", Tipe: Int8_t},
			}},
			{Name: "Imu", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "AntiThief", Tipe: Boolean_t},
				{Name: "Tilt", Tipe: Struct_t, Sub: []tagger{
					{Name: "Pitch", Tipe: Float_t, UnfactorType: Int16_t, Len: 2, Factor: 0.1},
					{Name: "Roll", Tipe: Float_t, UnfactorType: Int16_t, Len: 2, Factor: 0.1},
				}},
				{Name: "Total", Tipe: Struct_t, Sub: []tagger{
					{Name: "Accel", Tipe: Float_t, Len: 2, Factor: 0.01},
					{Name: "Gyro", Tipe: Float_t, Len: 2, Factor: 0.1},
					{Name: "Tilt", Tipe: Float_t, Len: 2, Factor: 0.1},
					{Name: "Temperature", Tipe: Int8_t},
				}},
			}},
			{Name: "Remote", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Nearby", Tipe: Boolean_t},
			}},
			{Name: "Finger", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "DriverID", Tipe: Uint8_t},
			}},
			{Name: "Audio", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Mute", Tipe: Uint8_t},
				{Name: "Volume", Tipe: Uint8_t},
			}},
			{Name: "Hmi", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Version", Tipe: Uint16_t},
			}},
			{Name: "Bms", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Run", Tipe: Boolean_t},
				{Name: "Faults", Tipe: Uint16_t},
				{Name: "SOC", Tipe: Uint8_t},
				{Name: "Capacity", Tipe: Struct_t, Sub: []tagger{
					{Name: "Remaining", Tipe: Uint16_t},
					{Name: "Usage", Tipe: Uint16_t},
				}},
				{Name: "Pack", Tipe: Array_t, Len: 2, Sub: []tagger{
					{Tipe: Struct_t, Sub: []tagger{
						{Name: "ID", Tipe: Uint32_t},
						{Name: "Faults", Tipe: Uint16_t},
						{Name: "Voltage", Tipe: Float_t, Len: 2, Factor: 0.01},
						{Name: "Current", Tipe: Float_t, Len: 2, Factor: 0.1},
						{Name: "Capacity", Tipe: Struct_t, Sub: []tagger{
							{Name: "Remaining", Tipe: Uint16_t},
							{Name: "Usage", Tipe: Uint16_t},
						}},
						{Name: "SOC", Tipe: Uint8_t},
						{Name: "SOH", Tipe: Uint8_t},
						{Name: "Temperature", Tipe: Int8_t},
					}},
				}},
			}},
			{Name: "Hbar", Tipe: Struct_t, Sub: []tagger{
				{Name: "Reverse", Tipe: Boolean_t},
				{Name: "Mode", Tipe: Struct_t, Sub: []tagger{
					{Name: "Drive", Tipe: Uint8_t},
					{Name: "Trip", Tipe: Uint8_t},
					{Name: "Avg", Tipe: Uint8_t},
				}},
				{Name: "Trip", Tipe: Struct_t, Sub: []tagger{
					{Name: "Odo", Tipe: Uint16_t},
					{Name: "A", Tipe: Uint16_t},
					{Name: "B", Tipe: Uint16_t},
				}},
				{Name: "Avg", Tipe: Struct_t, Sub: []tagger{
					{Name: "Range", Tipe: Uint8_t},
					{Name: "Efficiency", Tipe: Uint8_t},
				}},
			}},
			{Name: "Mcu", Tipe: Struct_t, Sub: []tagger{
				{Name: "Active", Tipe: Boolean_t},
				{Name: "Run", Tipe: Boolean_t},
				{Name: "Reverse", Tipe: Boolean_t},
				{Name: "DriveMode", Tipe: Uint8_t},
				{Name: "Speed", Tipe: Uint8_t},
				{Name: "RPM", Tipe: Int16_t},
				{Name: "Temperature", Tipe: Int8_t},
				{Name: "Error", Tipe: Array_t, Len: 8, Sub: []tagger{
					{Name: "err", Tipe: Uint8_t, Len: 1},
				}},
				{Name: "DCBus", Tipe: Struct_t, Sub: []tagger{
					{Name: "Current", Tipe: Float_t, Len: 2, Factor: 0.1},
					{Name: "Voltage", Tipe: Float_t, Len: 2, Factor: 0.1},
				}},
			}},
			{Name: "Task", Tipe: Struct_t, Sub: []tagger{
				{Name: "Stack", Tipe: Struct_t, Sub: []tagger{
					{Name: "Manager", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Network", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Reporter", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Command", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Imu", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Remote", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Finger", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Audio", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Gate", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanRX", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanTX", Tipe: Uint8_t, Sub: []tagger{}},
				}},
				{Name: "Wakeup", Tipe: Struct_t, Sub: []tagger{
					{Name: "Manager", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Network", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Reporter", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Command", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Imu", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Remote", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Finger", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Audio", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "Gate", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanRX", Tipe: Uint8_t, Sub: []tagger{}},
					{Name: "CanTX", Tipe: Uint8_t, Sub: []tagger{}},
				}},
			}},
		},
	},
}
