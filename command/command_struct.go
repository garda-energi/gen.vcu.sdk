package command

import (
	"github.com/pudjamansyurin/gen_vcu_sdk/header"
)

type HeaderCommand struct {
	header.Header
	Code    CMD_CODE    `type:"uint8"`
	SubCode CMD_SUBCODE `type:"uint8"`
}

type Command struct {
	HeaderCommand
	Payload [200]byte
}
