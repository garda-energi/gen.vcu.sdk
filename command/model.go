package command

import (
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type HeaderCommand struct {
	shared.Header
	Code    uint8 `type:"uint8"`
	SubCode uint8 `type:"uint8"`
}

type CommandPacket struct {
	HeaderCommand
	Payload [200]byte
}

type ResponsePacket struct {
	HeaderCommand
	ResCode	RES_CODE `type:"uint8"`
	Message [200]byte
}

type RES_CODE uint8

const (
        RES_CODE_ERROR RES_CODE = iota
	RES_CODE_OK
        RES_CODE_INVALID
)
