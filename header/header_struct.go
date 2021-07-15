package header

type Header struct {
	Prefix       string `type:"string" len:"2"`
	Size         uint8  `type:"uint8" unit:"Bytes"`
	Vin          uint32 `type:"uint32"`
	SendDatetime int64  `type:"unix_time" len:"7"`
}
