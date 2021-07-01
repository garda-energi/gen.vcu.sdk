package packet

type HeaderPacket struct {
	Prefix       string `type:"byte" len:"2"`
	Size         uint8  `type:"uint8" unit:"Bytes" chartable:""`
	Vin          uint32 `type:"uint32"`
	SendDatetime int64  `type:"byte" len:"7" chartable:""`
}
