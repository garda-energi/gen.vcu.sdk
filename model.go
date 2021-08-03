package sdk

import (
	"time"
)

type Header struct {
	Prefix       string    `type:"string" len:"2"`
	Size         uint8     `type:"uint8" unit:"Bytes"`
	Vin          uint32    `type:"uint32"`
	SendDatetime time.Time `type:"unix_time" len:"7"`
}

// message is type for command & response message (last field)
type message []byte

// overflow check if m length is overflowed
func (m message) overflow() bool {
	return len(m) > MESSAGE_LEN_MAX
}

type packet []byte
type packets []packet

// online convert status payload to online status.
func (p packet) online() bool {
	return string(p) == "1"
}
