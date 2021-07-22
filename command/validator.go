package command

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

// func contains(b []byte, value ...uint8) bool {
// 	for _, v := range value {
// 		if toUint8(b) == v {
// 			return true
// 		}
// 	}
// 	return false
// }

// func max(b []byte, max uint8) bool {
// 	return toUint8(b) <= max
// }

// func between(b []byte, min, max uint8) bool {
// 	return toUint8(b) >= min && toUint8(b) <= max
// }

// checkAck validate incomming ack packet
func checkAck(msg []byte) error {
	ack := util.Reverse(msg)
	if !bytes.Equal(ack, []byte(shared.PREFIX_ACK)) {
		return errors.New("ack corrupt")
	}
	return nil
}

// checkResponse validate incomming response packet,
// it also parse response code and message
func checkResponse(cmder *commander, res *ResponsePacket) error {
	// check code
	if res.Header.Code != cmder.code && res.Header.SubCode != cmder.sub_code {
		return errors.New("response-mismatch")
	}

	// check resCode
	var err string
	switch res.Header.ResCode {
	case RES_CODE_OK:
		return nil
	case RES_CODE_ERROR:
		err = "response-error"
	case RES_CODE_INVALID:
		err = "response-invalid"
	default:
		err = "response-unknown"
	}

	if len(res.Message) == 0 {
		return errors.New(err)
	}

	// TODO: subtitutes BIKE_STATE to message

	return fmt.Errorf("%s, %s", err, res.Message)
}
