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

type HeaderResponse struct {
	HeaderCommand
	ResCode RES_CODE `type:"uint8"`
}

type CommandPacket struct {
	Header  *HeaderCommand
	Payload []byte
}

type ResponsePacket struct {
	Header  *HeaderResponse
	Message []byte `type:"char"`
}

type RES_CODE uint8

const (
	RES_CODE_ERROR RES_CODE = iota
	RES_CODE_OK
	RES_CODE_INVALID
)
