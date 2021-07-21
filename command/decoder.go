package command

import (
	"bytes"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

func (c *Command) decode(cmder *Commander, packet []byte) (*ResponsePacket, error) {
	reader := bytes.NewReader(packet)

	r := &ResponsePacket{
		Header: &HeaderResponse{},
	}

	// header
	if err := shared.Decode(reader, r.Header); err != nil {
		return nil, err
	}

	// message
	if len := reader.Len(); len > 0 {
		msg := make([]byte, len)
		reader.Read(msg)
		r.Message = msg
	}

	return r, nil
}
