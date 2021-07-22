package command

import (
	"encoding/binary"
	"errors"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

// encode combine command and payload to bytes packet.
func (c *Command) encode(cmder *commander, payload []byte) ([]byte, error) {
	if len(payload) > PAYLOAD_LEN {
		return nil, errors.New("payload overload")
	}

	var sb strings.Builder
	sb.Write(payload)
	sb.WriteByte(byte(cmder.sub_code))
	sb.WriteByte(byte(cmder.code))
	sb.Write(shared.TimeToBytes(time.Now()))

	vin32 := make([]byte, 4)
	binary.BigEndian.PutUint32(vin32, uint32(c.vin))
	sb.Write(vin32)

	sb.WriteByte(byte(sb.Len()))
	sb.WriteString(shared.PREFIX_COMMAND)

	bytes := shared.StrToBytes(sb.String())
	return bytes, nil
}
