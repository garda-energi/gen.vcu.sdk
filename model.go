package sdk

import "time"

type Header struct {
	Prefix  string `type:"string" len:"2"`
	Size    uint8  `type:"uint8" unit:"Bytes"`
	Version uint16 `type:"uint16"`
	Vin     uint32 `type:"uint32"`
}

type HeaderReport struct {
	Header
	SendDatetime time.Time `type:"unix_time" len:"7"`
	LogDatetime  time.Time `type:"int64" len:"7"`
	Frame        Frame     `type:"uint8"`
}

type HeaderCommand struct {
	Header
	Code    uint8 `type:"uint8"`
	SubCode uint8 `type:"uint8"`
}

type headerResponse struct {
	HeaderCommand
	ResCode resCode `type:"uint8"`
}

// message is type for command & response message (last field)
type message []byte

// overflow check if m length is overflowed
func (m message) overflow() bool {
	return len(m) > MESSAGE_LEN_MAX
}

type packet []byte
type packets []packet
type PacketData map[string]interface{}

// online convert status payload to online status.
func (p packet) online() bool {
	return string(p) == "1"
}

type genReportPacket struct {
	Header
	LogDatetime time.Time `type:"int64" len:"7"`
	Version     uint16    `type:"uint16"`
	Frame       Frame     `type:"uint8"`
	Payload     message
	Data        PacketData
}
