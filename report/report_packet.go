package report

import (
	"fmt"
	"reflect"

	"github.com/pudjamansyurin/gen_vcu_sdk/header"
)

func ReportSimplePacket() []header.Packet {
	packets := header.HeaderReportPacket
	packets = append(packets, VcuPacket...)
	packets = append(packets, EepromPacket...)
	packets = append(packets, GpsPacket...)

	return packets
}

func ReportFullPacket() []header.Packet {
	packets := ReportSimplePacket()
	packets = append(packets, HbarPacket()...)
	packets = append(packets, NetPacket...)
	packets = append(packets, MemsPacket()...)
	packets = append(packets, RemotePacket...)
	packets = append(packets, FingerPacket...)
	packets = append(packets, AudioPacket...)
	packets = append(packets, HmiPacket...)
	packets = append(packets, BmsPackPacket()...)
	packets = append(packets, McuPacket()...)
	packets = append(packets, TaskPacket()...)

	return packets
}

var VcuPacket = []header.Packet{
	{
		Name:     "vcu.logDatetime",
		Datetime: true,
	},
	{
		Name: "vcu.bikeState",
		Dst:  reflect.Int8,
	},
	{
		Name: "vcu.events",
		Dst:  reflect.Int8,
	},
	{
		Name: "vcu.logBuffered",
		Dst:  reflect.Uint8,
	},
	{
		Name:   "vcu.batVoltage",
		Dst:    reflect.Uint8,
		Factor: 18.0,
		Unit:   "mVolt",
	},
	{
		Name:   "vcu.uptime",
		Src:    reflect.Uint32,
		Dst:    reflect.Float32,
		Factor: 0.0002777,
		Unit:   "hour",
	},
}

var EepromPacket = []header.Packet{
	{
		Name: "eeprom.active",
		Dst:  reflect.Bool,
	},
	{
		Name: "eeprom.used",
		Dst:  reflect.Uint8,
		Unit: "%",
	},
}

var GpsPacket = []header.Packet{
	{
		Name: "gps.active",
		Dst:  reflect.Bool,
	},
	{
		Name: "gps.satInUse",
		Dst:  reflect.Uint8,
	},
	{
		Name:   "gps.hdop",
		Dst:    reflect.Float32,
		Src:    reflect.Uint8,
		Factor: 0.1,
	},
	{
		Name:   "gps.vdop",
		Dst:    reflect.Float32,
		Src:    reflect.Uint8,
		Factor: 0.1,
	},
	{
		Name: "gps.speed",
		Dst:  reflect.Uint8,
		Unit: "Kph",
	},
	{
		Name:   "gps.heading",
		Dst:    reflect.Float32,
		Src:    reflect.Uint8,
		Factor: 2.0,
		Unit:   "Deg",
	},
	{
		Name:   "gps.longitude",
		Dst:    reflect.Float32,
		Src:    reflect.Int32,
		Factor: 0.0000001,
	},
	{
		Name:   "gps.latitude",
		Dst:    reflect.Float32,
		Src:    reflect.Int32,
		Factor: 0.0000001,
	},
	{
		Name:   "gps.altitude",
		Dst:    reflect.Float32,
		Src:    reflect.Uint16,
		Factor: 0.1,
		Unit:   "m",
	},
}

func HbarPacket() []header.Packet {
	packets := []header.Packet{
		{
			Name: "hbar.reverse",
			Dst:  reflect.Bool,
		},
	}

	for _, mode := range MODE_LIST {
		packets = append(packets, header.Packet{
			Name: fmt.Sprintf("hbar.mode.%s", mode),
			Dst:  reflect.Uint8,
		})
	}

	for _, trip := range MODE_TRIP_LIST {
		packets = append(packets, header.Packet{
			Name: fmt.Sprintf("hbar.trip.%s", trip),
			Dst:  reflect.Uint16,
			Unit: "Km",
		})
	}

	for _, avg := range MODE_AVG_LIST {
		unit := "Km"
		if avg == "efficiency" {
			unit = "Km/Kwh"
		}

		packets = append(packets, header.Packet{
			Name: fmt.Sprintf("hbar.avg.%s", avg),
			Dst:  reflect.Uint8,
			Unit: unit,
		})
	}

	return packets
}

var NetPacket = []header.Packet{
	{
		Name: "net.signal",
		Dst:  reflect.Uint8,
		Unit: "%",
	},
	{
		Name: "net.state",
		Dst:  reflect.Int8,
	},
	{
		Name: "net.ipStatus",
		Dst:  reflect.Int8,
	},
}

func MemsPacket() []header.Packet {
	packets := []header.Packet{
		{
			Name: "mems.active",
			Dst:  reflect.Bool,
		},
		{
			Name: "mems.motion",
			Dst:  reflect.Bool,
		},
	}

	for _, accel := range []string{"x", "y", "z"} {
		packets = append(packets, header.Packet{
			Name:   fmt.Sprintf("mems.accel.%s", accel),
			Dst:    reflect.Float32,
			Src:    reflect.Int16,
			Factor: 0.01,
			Unit:   "G",
		})
	}

	for _, gyro := range []string{"x", "y", "z"} {
		packets = append(packets, header.Packet{
			Name:   fmt.Sprintf("mems.gyro.%s", gyro),
			Dst:    reflect.Float32,
			Src:    reflect.Int16,
			Factor: 0.1,
			Unit:   "rad/s",
		})
	}

	for _, tilt := range []string{"pitch", "roll"} {
		packets = append(packets, header.Packet{
			Name:   fmt.Sprintf("mems.tilt.%s", tilt),
			Dst:    reflect.Float32,
			Src:    reflect.Int16,
			Factor: 0.1,
			Unit:   "Deg",
		})
	}

	packets = append(packets, []header.Packet{

		{
			Name:   "mems.total.accel",
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   "G",
		},
		{
			Name:   "mems.total.gyro",
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   "rad/s",
		},
		{
			Name:   "mems.total.tilt",
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   "Deg",
		},
		{
			Name:   "mems.total.temp",
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   "Celcius",
		},
	}...)

	return packets
}

var RemotePacket = []header.Packet{
	{
		Name: "remote.active",
		Dst:  reflect.Bool,
	},
	{
		Name: "remote.nearby",
		Dst:  reflect.Bool,
	},
}

var FingerPacket = []header.Packet{
	{
		Name: "finger.verified",
		Dst:  reflect.Bool,
	},
	{
		Name: "finger.driverId",
		Dst:  reflect.Uint8,
	},
}

var AudioPacket = []header.Packet{
	{
		Name: "audio.active",
		Dst:  reflect.Bool,
	},
	{
		Name: "audio.mute",
		Dst:  reflect.Bool,
	},
	{
		Name: "audio.driverId",
		Dst:  reflect.Uint8,
		Unit: "%",
	},
}

var HmiPacket = []header.Packet{
	{
		Name: "hmi.active",
		Dst:  reflect.Bool,
	},
}

func BmsPackPacket() []header.Packet {
	var packets = []header.Packet{
		{
			Name: "bms.active",
			Dst:  reflect.Bool,
		},
		{
			Name: "bms.run",
			Dst:  reflect.Bool,
		},
		{
			Name: "bms.soc",
			Dst:  reflect.Uint8,
			Unit: "%",
		},
		{
			Name: "bms.fault",
			Dst:  reflect.Uint16,
		},
	}

	for i := 0; i < BMS_PACK_CNT; i++ {
		PackPacket := []header.Packet{
			{
				Name: fmt.Sprintf("bms.pack.%d.id", i),
				Dst:  reflect.Uint32,
			},
			{
				Name: fmt.Sprintf("bms.pack.%d.fault", i),
				Dst:  reflect.Uint16,
			},
			{
				Name:   fmt.Sprintf("bms.pack.%d.voltage", i),
				Dst:    reflect.Float32,
				Src:    reflect.Uint16,
				Factor: 0.01,
				Unit:   "Volt",
			},
			{
				Name:   fmt.Sprintf("bms.pack.%d.current", i),
				Dst:    reflect.Float32,
				Src:    reflect.Uint16,
				Factor: 0.1,
				Unit:   "Ampere",
			},
			{
				Name: fmt.Sprintf("bms.pack.%d.soc", i),
				Dst:  reflect.Uint8,
				Unit: "%",
			},
			{
				Name: fmt.Sprintf("bms.pack.%d.temp", i),
				Dst:  reflect.Uint16,
				Unit: "Celcius",
			},
		}
		packets = append(packets, PackPacket...)
	}

	return packets
}

func McuPacket() []header.Packet {
	packets := []header.Packet{
		{
			Name: "mcu.active",
			Dst:  reflect.Bool,
		},
		{
			Name: "mcu.run",
			Dst:  reflect.Bool,
		},
		{
			Name: "mcu.reverse",
			Dst:  reflect.Bool,
		},
		{
			Name: "mcu.driveMode",
			Dst:  reflect.Uint8,
		},
		{
			Name: "mcu.speed",
			Dst:  reflect.Uint8,
			Unit: "Kph",
		},
		{
			Name: "mcu.rpm",
			Dst:  reflect.Int16,
			Unit: "rpm",
		},
		{
			Name:   "mcu.temp",
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   "Celcius",
		},
	}

	for _, fault := range []string{"run", "post"} {
		packets = append(packets, header.Packet{
			Name: fmt.Sprintf("mcu.fault.%s", fault),
			Dst:  reflect.Uint32,
		})
	}

	for _, torque := range []string{"command", "feedback"} {
		packets = append(packets, header.Packet{
			Name:   fmt.Sprintf("mcu.torque.%s", torque),
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   "Nm",
		})
	}

	for _, dcbus := range []string{"current", "voltage"} {
		unit := "Ampere"
		if dcbus == "voltage" {
			unit = "Volt"
		}

		packets = append(packets, header.Packet{
			Name:   fmt.Sprintf("mcu.dcbus.%s", dcbus),
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   unit,
		})
	}

	packets = append(packets, []header.Packet{
		{
			Name: "mcu.inverter.enabled",
			Dst:  reflect.Bool,
		},
		{
			Name: "mcu.inverter.lockout",
			Dst:  reflect.Bool,
		},
		{
			Name: "mcu.inverter.discharge",
			Dst:  reflect.Uint8,
		},
		{
			Name: "mcu.template.maxRpm",
			Dst:  reflect.Int16,
			Unit: "rpm",
		},
		{
			Name: "mcu.template.maxSpeed",
			Dst:  reflect.Uint8,
			Unit: "Kph",
		},
	}...)

	for i := 0; i < MODE_DRIVE_CNT; i++ {
		DriveModePacket := []header.Packet{
			{
				Name: fmt.Sprintf("mcu.template.driveMode.%d.discur", i),
				Dst:  reflect.Uint16,
				Unit: "Ampere",
			},
			{
				Name:   fmt.Sprintf("mcu.template.driveMode.%d.torque", i),
				Dst:    reflect.Float32,
				Src:    reflect.Uint16,
				Factor: 0.1,
				Unit:   "Nm",
			},
		}

		packets = append(packets, DriveModePacket...)
	}

	return packets
}

func TaskPacket() []header.Packet {
	var packets []header.Packet

	for _, task := range TASK_LIST {
		packets = append(packets, header.Packet{
			Name: fmt.Sprintf("task.stack.%s", task),
			Dst:  reflect.Uint16,
			Unit: "Bytes",
		})
	}

	for _, task := range TASK_LIST {
		packets = append(packets, header.Packet{
			Name: fmt.Sprintf("task.wakeup.%s", task),
			Dst:  reflect.Uint8,
			Unit: "Seconds",
		})
	}

	return packets
}
