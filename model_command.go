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

type commandPacket struct {
	Header  *HeaderCommand
	Message message
}

type headerResponse struct {
	HeaderCommand
	ResCode resCode `type:"uint8"`
}

type responsePacket struct {
	Header  *headerResponse
	Message message `type:"slice"`
}

// validPrefix check if r's prefix is valid
func (r *responsePacket) validPrefix() bool {
	if r.Header == nil {
		return false
	}
	return r.Header.Prefix == PREFIX_RESPONSE
}

// size calculate r's size, ignoring prefix & size field
func (r *responsePacket) size() int {
	if r.Header == nil {
		return 0
	}
	return getPacketSize(r) - 3
}

// validSize check if r's size is valid
func (r *responsePacket) validSize() bool {
	if r.Header == nil {
		return false
	}
	return int(r.Header.Size) == r.size()
}

// belongsTo check if r is response for cmd
func (r *responsePacket) belongsTo(cmd *command) bool {
	if r.Header == nil || cmd == nil {
		return false
	}
	return r.Header.Code == cmd.code && r.Header.SubCode == cmd.subCode
}

// validCmdCode check if r's command code & subCode is valid
func (r *responsePacket) validCmdCode() bool {
	if r.Header == nil {
		return false
	}
	_, err := getCmdByCode(int(r.Header.Code), int(r.Header.SubCode))
	return err == nil
}

// validResCode check if r's response code is valid
func (r *responsePacket) validResCode() bool {
	if r.Header == nil {
		return false
	}
	for i := resCodeError; i < resCodeLimit; i++ {
		if r.Header.ResCode == i {
			return true
		}
	}
	return false
}

// hasMessage check if r has message
func (r *responsePacket) hasMessage() bool {
	return len(r.Message) > 0
}

// renderMessage subtitue BikeState to r's message
func (r *responsePacket) renderMessage() {
	if !r.hasMessage() {
		return
	}

	str := string(r.Message)
	for i := BikeStateUnknown; i < BikeStateLimit; i++ {
		old := fmt.Sprintf("{%d}", i)
		new := BikeState(i).String()
		str = strings.ReplaceAll(str, old, new)
	}
	r.Message = []byte(str)
}

// message is type for command & response message (last field)
type message []byte

// overflow check if m length is overflowed
func (m message) overflow() bool {
	return len(m) > MESSAGE_LEN_MAX
}
