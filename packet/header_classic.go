package packet

import (
	"reflect"
)

var Header = []Packet{
	{
		Group:    "packet",
		Field:    "prefix",
		Title:    "Prefix",
		Header:   true,
		Required: true,
		Size:     2,
		Format:   Format{DataType: reflect.String},
		// Format: (v) => HexToAscii(cend(v)),
		// display: (vf) => vf,
		// FormatCmd: (_) => cend(AsciiToHex(config.prefix.command)),
	},
	{
		Group:     "packet",
		Field:     "size",
		Title:     "Size",
		Header:    true,
		Required:  true,
		Chartable: true,
		Unit:      "Bytes",
		Size:      1,
		Format:    Format{DataType: reflect.Uint8},
		// Format: (v) => HexToUnsignedInt(cend(v)),
		// display: (vf) => Dot(vf),
		// FormatCmd: (hex) => cend(IntToHex(hex.length / 2, 1 * 2)),
	},
	{
		Group:    "packet",
		Field:    "vin",
		Title:    "VIN",
		Header:   true,
		Required: true,
		Size:     4,
		Format:   Format{DataType: reflect.Uint32},
		// Format: (v) => HexToUnsignedInt(cend(v)),
		// display: (vf) => vf,
		// FormatCmd: (v) => cend(IntToHex(v, 4 * 2)),
	},
	{
		Group:     "packet.datetime",
		Field:     "sendDatetime",
		Title:     "Send Datetime",
		Header:    true,
		Required:  true,
		Chartable: true,
		Size:      7,
		Format:    Format{DataType: reflect.Interface, InterfaceName: "datetime"},
		// Format: (v) => Number(dayjs(parseDatetime(v), "YYMMDDHHmmss0d").unix()),
		// display: (vf) => dayjs.unix(vf).Format("ddd, DD-MM-YY HH:mm:ss"),
		// FormatCmd: (unix) =>
		//   buildTimestamp(dayjs.unix(unix).Format("YYMMDDHHmmss0d")),
	},
}

var CommandHeader = append(Header, []Packet{
	{
		Group:    "command",
		Field:    "code",
		Title:    "Code",
		Header:   true,
		Required: true,
		Size:     1,
		Format:   Format{DataType: reflect.Uint8},
		// Format: (v) => HexToUnsignedInt(v),
		// display: (vf) => vf,
		// FormatCmd: (v) => cend(IntToHex(v, 1 * 2)),
	},
	{
		Group:    "command",
		Field:    "subCode",
		Title:    "Sub Code",
		Header:   true,
		Required: true,
		Size:     1,
		Format:   Format{DataType: reflect.Uint8},
		// Format: (v) => HexToUnsignedInt(v),
		// display: (vf) => vf,
		// FormatCmd: (v) => cend(IntToHex(v, 1 * 2)),
	},
}...)
