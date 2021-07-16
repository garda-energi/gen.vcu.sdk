package command

import (
	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type HeaderCmd struct {
	shared.Header
	Code    CMD_CODE    `type:"uint8"`
	SubCode CMD_SUBCODE `type:"uint8"`
}

type Cmd struct {
	HeaderCmd
	Payload [200]byte
}
