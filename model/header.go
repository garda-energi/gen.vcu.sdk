package model

type Header struct {
	Prefix       string `type:"string" len:"2"`
	Size         uint8  `type:"uint8" unit:"Bytes"`
	Vin          uint32 `type:"uint32"`
	SendDatetime int64  `type:"unix_time" len:"7"`
}

type HeaderReport struct {
	Header
	FrameID FRAME_ID `type:"uint8"`
}

type HeaderCommand struct {
	Header
	Code    CMD_CODE    `type:"uint8"`
	SubCode CMD_SUBCODE `type:"uint8"`
}
