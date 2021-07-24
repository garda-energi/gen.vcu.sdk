package sdk

import (
	"fmt"
	"strings"
)

type HeaderCommand struct {
	Header
	Code    uint8 `type:"uint8"`
	SubCode uint8 `type:"uint8"`
}

type CommandPacket struct {
	Header  *HeaderCommand
	Payload []byte
}

type HeaderResponse struct {
	HeaderCommand
	ResCode resCode `type:"uint8"`
}

type ResponsePacket struct {
	Header  *HeaderResponse
	Message []byte `type:"char"`
}

func (r *ResponsePacket) ValidPrefix() bool {
	return r.Header.Prefix == PREFIX_RESPONSE
}

func (r *ResponsePacket) Size() int {
	return 4 + 7 + 1 + 1 + 1 + len(r.Message)
}

func (r *ResponsePacket) AnswerFor(cmd *command) bool {
	return r.Header.Code == cmd.code && r.Header.SubCode == cmd.sub_code
}

func (r *ResponsePacket) ValidResCode() bool {
	for i := resCodeError; i < resCodeLimit; i++ {
		if r.Header.ResCode == i {
			return true
		}
	}
	return false
}

func (r *ResponsePacket) HasMessage() bool {
	return len(r.Message) > 0
}

func (r *ResponsePacket) MessageOverflow() bool {
	return len(r.Message) > PAYLOAD_LEN_MAX
}

func (r *ResponsePacket) RenderMessage() {
	str := string(r.Message)
	for i := BikeStateUnknown; i < BikeStateLimit; i++ {
		old := fmt.Sprintf("{%d}", i)
		new := BikeState(i).String()
		str = strings.ReplaceAll(str, old, new)
	}
	r.Message = []byte(str)
}
