package command

import (
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

const PAYLOAD_LEN = 200

type HeaderCommand struct {
	shared.Header
	Code    uint8 `type:"uint8"`
	SubCode uint8 `type:"uint8"`
}

type CommandPacket struct {
	HeaderCommand
	Payload [PAYLOAD_LEN]byte
}

type ResponsePacket struct {
	HeaderCommand
	ResCode	RES_CODE `type:"uint8"`
	Message [PAYLOAD_LEN]byte
}

type RES_CODE uint8

const (
        RES_CODE_ERROR RES_CODE = iota
	RES_CODE_OK
        RES_CODE_INVALID
)
