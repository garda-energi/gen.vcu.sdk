package shared

import "reflect"

type Packet struct {
	Name     string
	Src      reflect.Kind
	Dst      reflect.Kind
	Len      int
	Unit     string
	Factor   float32
	Datetime bool
}

var HeaderPacket = []Packet{
	{
		Name: "header.prefix",
		Dst:  reflect.String,
		Len:  2,
	},
	{
		Name: "header.size",
		Dst:  reflect.Uint8,
		Unit: "Bytes",
	},
	{
		Name: "header.vin",
		Dst:  reflect.Uint32,
	},
	{
		Name:     "header.sendDatetime",
		Datetime: true,
	},
}

var HeaderReportPacket = append(HeaderPacket, Packet{
	Name: "header.frameID",
	Dst:  reflect.Uint8,
})

var HeaderCommandPacket = append(HeaderPacket, []Packet{
	{
		Name: "header.code",
		Dst:  reflect.Uint8,
	},
	{
		Name: "header.subCode",
		Dst:  reflect.Uint8,
	},
}...)
