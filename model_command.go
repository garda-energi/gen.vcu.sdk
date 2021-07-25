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

// type CommandPacket struct {
// 	Header  *HeaderCommand
// 	Message message
// }

type headerResponse struct {
	HeaderCommand
	ResCode resCode `type:"uint8"`
}

type responsePacket struct {
	Header  *headerResponse
	Message message `type:"slice"`
}

type message []byte

func (p message) overflow() bool {
	return len(p) > PAYLOAD_LEN_MAX
}

func (r *responsePacket) validPrefix() bool {
	return r.Header.Prefix == PREFIX_RESPONSE
}

func (r *responsePacket) size() int {
	return 4 + 7 + 1 + 1 + 1 + len(r.Message)
}

func (r *responsePacket) validSize() bool {
	return int(r.Header.Size) == r.size()
}

func (r *responsePacket) matchWith(cmd *command) bool {
	return r.Header.Code == cmd.code && r.Header.SubCode == cmd.subCode
}

func (r *responsePacket) validResCode() bool {
	for i := resCodeError; i < resCodeLimit; i++ {
		if r.Header.ResCode == i {
			return true
		}
	}
	return false
}

func (r *responsePacket) hasMessage() bool {
	return len(r.Message) > 0
}

func (r *responsePacket) messageOverflow() bool {
	return len(r.Message) > PAYLOAD_LEN_MAX
}

func (r *responsePacket) renderMessage() {
	str := string(r.Message)
	for i := BikeStateUnknown; i < BikeStateLimit; i++ {
		old := fmt.Sprintf("{%d}", i)
		new := BikeState(i).String()
		str = strings.ReplaceAll(str, old, new)
	}
	r.Message = []byte(str)
}
