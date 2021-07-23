package command

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

// encode combine command and payload to bytes packet.
func (c *Command) encode(cmder *commander, payload []byte) ([]byte, error) {
	if len(payload) > PAYLOAD_LEN {
		return nil, errors.New("payload overload")
	}

	var buf bytes.Buffer
	ed := binary.LittleEndian
	binary.Write(&buf, ed, shared.Reverse(payload))
	binary.Write(&buf, ed, cmder.sub_code)
	binary.Write(&buf, ed, cmder.code)
	binary.Write(&buf, ed, shared.Reverse(shared.TimeToBytes(time.Now())))
	binary.Write(&buf, binary.BigEndian, uint32(c.vin))
	binary.Write(&buf, ed, byte(buf.Len()))
	binary.Write(&buf, ed, []byte(shared.PREFIX_COMMAND))
	bytes := shared.Reverse(buf.Bytes())
	return bytes, nil
}
