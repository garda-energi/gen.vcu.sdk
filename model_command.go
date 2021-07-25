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

func (r *responsePacket) validPrefix() bool {
	if r.Header != nil {
		return r.Header.Prefix == PREFIX_RESPONSE
	}
	return false
}

func (r *responsePacket) size() int {
	if r.Header != nil {
		return 4 + 7 + 1 + 1 + 1 + len(r.Message)
	}
	return 0
}

func (r *responsePacket) validSize() bool {
	if r.Header != nil {
		return int(r.Header.Size) == r.size()
	}
	return false
}

func (r *responsePacket) matchWith(cmd *command) bool {
	if r.Header != nil && cmd != nil {
		return r.Header.Code == cmd.code && r.Header.SubCode == cmd.subCode
	}
	return false
}

func (r *responsePacket) validResCode() bool {
	if r.Header != nil {
		for i := resCodeError; i < resCodeLimit; i++ {
			if r.Header.ResCode == i {
				return true
			}
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

// renderMessage subtitue BikeState to r.Message
func (r *responsePacket) renderMessage() {
	str := string(r.Message)
	for i := BikeStateUnknown; i < BikeStateLimit; i++ {
		old := fmt.Sprintf("{%d}", i)
		new := BikeState(i).String()
		str = strings.ReplaceAll(str, old, new)
	}
	r.Message = []byte(str)
}

type message []byte

func (p message) overflow() bool {
	return len(p) > PAYLOAD_LEN_MAX
}
