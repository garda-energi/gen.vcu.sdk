package command

import (
	"encoding/binary"
	"errors"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

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

	bytes := util.Reverse([]byte(sb.String()))
	return bytes, nil
}

func makeBool(d bool) []byte {
	b := []byte{0}
	if d {
		b[0] = 1
	}
	return b
}
