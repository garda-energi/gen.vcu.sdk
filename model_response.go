package sdk

import (
	"errors"
	"fmt"
	"strings"
)

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
	r.Message = message(str)
}

// validateResponse validate incomming response packet.
// It also render message part (subtitutes BikeState).
func (r *responsePacket) validateResponse(vin int, cmd *command) error {
	if int(r.Header.Vin) != vin {
		return errInvalidVin
	}
	if !r.belongsTo(cmd) {
		return errInvalidCmdCode
	}
	if r.Header.ResCode == resCodeOk {
		return nil
	}

	out := fmt.Sprint(r.Header.ResCode)
	// check if message is not empty
	if r.hasMessage() {
		r.renderMessage()
		out += fmt.Sprint(" ", string(r.Message))
	}
	return errors.New(out)
}
