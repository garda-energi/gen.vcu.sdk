package report

import (
	"fmt"
	"reflect"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

func ReportSimplePacket() []shared.Packet {
	packets := shared.HeaderReportPacket
	packets = append(packets, VcuPacket...)
	packets = append(packets, EepromPacket...)
	packets = append(packets, GpsPacket...)

	return packets
}

func ReportFullPacket() []shared.Packet {
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

var VcuPacket = []shared.Packet{
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
		Dst:  reflect.Uint16,
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

var EepromPacket = []shared.Packet{
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

var GpsPacket = []shared.Packet{
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

func HbarPacket() []shared.Packet {
	packets := []shared.Packet{
		{
			Name: "hbar.reverse",
			Dst:  reflect.Bool,
		},
	}

	for mode := shared.MODE(0); mode < shared.MODE_SUB_limit; mode++ {
		packets = append(packets, shared.Packet{
			Name: fmt.Sprintf("hbar.mode.%s", mode),
			Dst:  reflect.Uint8,
		})
	}

	for trip := shared.MODE_TRIP(0); trip < shared.MODE_TRIP_limit; trip++ {
		packets = append(packets, shared.Packet{
			Name: fmt.Sprintf("hbar.trip.%s", trip),
			Dst:  reflect.Uint16,
			Unit: "Km",
		})
	}

	for avg := shared.MODE_AVG(0); avg < shared.MODE_AVG_limit; avg++ {
		packets = append(packets, shared.Packet{
			Name: fmt.Sprintf("hbar.avg.%s", avg),
			Dst:  reflect.Uint8,
			Unit: avg.Unit(),
		})
	}

	return packets
}

var NetPacket = []shared.Packet{
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

func MemsPacket() []shared.Packet {
	packets := []shared.Packet{
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
		packets = append(packets, shared.Packet{
			Name:   fmt.Sprintf("mems.accel.%s", accel),
			Dst:    reflect.Float32,
			Src:    reflect.Int16,
			Factor: 0.01,
			Unit:   "G",
		})
	}

	for _, gyro := range []string{"x", "y", "z"} {
		packets = append(packets, shared.Packet{
			Name:   fmt.Sprintf("mems.gyro.%s", gyro),
			Dst:    reflect.Float32,
			Src:    reflect.Int16,
			Factor: 0.1,
			Unit:   "rad/s",
		})
	}

	for _, tilt := range []string{"pitch", "roll"} {
		packets = append(packets, shared.Packet{
			Name:   fmt.Sprintf("mems.tilt.%s", tilt),
			Dst:    reflect.Float32,
			Src:    reflect.Int16,
			Factor: 0.1,
			Unit:   "Deg",
		})
	}

	packets = append(packets, []shared.Packet{

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

var RemotePacket = []shared.Packet{
	{
		Name: "remote.active",
		Dst:  reflect.Bool,
	},
	{
		Name: "remote.nearby",
		Dst:  reflect.Bool,
	},
}

var FingerPacket = []shared.Packet{
	{
		Name: "finger.verified",
		Dst:  reflect.Bool,
	},
	{
		Name: "finger.driverId",
		Dst:  reflect.Uint8,
	},
}

var AudioPacket = []shared.Packet{
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

var HmiPacket = []shared.Packet{
	{
		Name: "hmi.active",
		Dst:  reflect.Bool,
	},
}

func BmsPackPacket() []shared.Packet {
	var packets = []shared.Packet{
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

	for i := 0; i < shared.BMS_PACK_CNT; i++ {
		PackPacket := []shared.Packet{
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

func McuPacket() []shared.Packet {
	packets := []shared.Packet{
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
		packets = append(packets, shared.Packet{
			Name: fmt.Sprintf("mcu.fault.%s", fault),
			Dst:  reflect.Uint32,
		})
	}

	for _, torque := range []string{"command", "feedback"} {
		packets = append(packets, shared.Packet{
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

		packets = append(packets, shared.Packet{
			Name:   fmt.Sprintf("mcu.dcbus.%s", dcbus),
			Dst:    reflect.Float32,
			Src:    reflect.Uint16,
			Factor: 0.1,
			Unit:   unit,
		})
	}

	packets = append(packets, []shared.Packet{
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

	for i := 0; i < int(shared.MODE_DRIVE_limit); i++ {
		DriveModePacket := []shared.Packet{
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

func TaskPacket() []shared.Packet {
	var packets []shared.Packet

	for _, task := range shared.TASK_LIST {
		packets = append(packets, shared.Packet{
			Name: fmt.Sprintf("task.stack.%s", task),
			Dst:  reflect.Uint16,
			Unit: "Bytes",
		})
	}

	for _, task := range shared.TASK_LIST {
		packets = append(packets, shared.Packet{
			Name: fmt.Sprintf("task.wakeup.%s", task),
			Dst:  reflect.Uint8,
			Unit: "Seconds",
		})
	}

	return packets
}
