package command

import (
	"bytes"
	"errors"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func contains(b []byte, value ...uint8) bool {
	for _, v := range value {
		if toUint8(b) == v {
			return true
		}
	}
	return false
}

func max(b []byte, max uint8) bool {
	return toUint8(b) <= max
}

func between(b []byte, min, max uint8) bool {
	return toUint8(b) >= min && toUint8(b) <= max
}

func checkAck(msg []byte) error {
	ack := util.Reverse(msg)
	if !bytes.Equal(ack, []byte(shared.PREFIX_ACK)) {
		return errors.New("ack corrupt")
	}
	return nil
}

func checkResponse(cmder *Commander, res *ResponsePacket) error {
	// check code
	if res.Header.Code != cmder.Code && res.Header.SubCode != cmder.SubCode {
		return errors.New("command & response mismatch")
	}

	// check resCode
	switch res.Header.ResCode {
	case RES_CODE_OK:
		return nil
	case RES_CODE_ERROR:
		return errors.New("response error")
	case RES_CODE_INVALID:
		return errors.New("response invalid")
	default:
		return errors.New("unknown response")
	}
}
