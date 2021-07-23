package command

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

// encode combine command and payload to bytes packet.
func (c *Commander) encode(cmd *command, payload []byte) ([]byte, error) {
	if len(payload) > PAYLOAD_LEN {
		return nil, errors.New("payload overload")
	}

	var buf bytes.Buffer
	ed := binary.LittleEndian
	binary.Write(&buf, ed, shared.Reverse(payload))
	binary.Write(&buf, ed, cmd.sub_code)
	binary.Write(&buf, ed, cmd.code)
	binary.Write(&buf, ed, shared.Reverse(shared.TimeToBytes(time.Now())))
	binary.Write(&buf, binary.BigEndian, uint32(c.vin))
	binary.Write(&buf, ed, byte(buf.Len()))
	binary.Write(&buf, ed, []byte(shared.PREFIX_COMMAND))
	bytes := shared.Reverse(buf.Bytes())
	return bytes, nil
}

// decode extract response header and message from bytes packet.
func (c *Commander) decode(cmd *command, packet []byte) (*ResponsePacket, error) {
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
