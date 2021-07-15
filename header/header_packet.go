package header

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
		Src:  reflect.String,
		Len:  2,
	},
	{
		Name: "header.size",
		Src:  reflect.Uint8,
		Unit: "Bytes",
	},
	{
		Name: "header.vin",
		Src:  reflect.Uint32,
	},
	{
		Name:     "header.sendDatetime",
		Datetime: true,
	},
}

var HeaderReportPacket = append(HeaderPacket, Packet{
	Name: "header.frameID",
	Src:  reflect.Uint8,
})

var HeaderCommandPacket = append(HeaderPacket, []Packet{
	{
		Name: "header.code",
		Src:  reflect.Uint8,
	},
	{
		Name: "header.subCode",
		Src:  reflect.Uint8,
	},
}...)
