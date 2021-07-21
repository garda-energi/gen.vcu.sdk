package command

import (
	"encoding/binary"
	"errors"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func toUint8(b []byte) uint8 {
	return uint8(b[0])
}

func (c *Command) encode(cmder *commander, payload []byte) ([]byte, error) {
	if len(payload) > PAYLOAD_LEN {
		return nil, errors.New("payload overload")
	}

	var sb strings.Builder
	sb.Write(payload)
	sb.WriteByte(byte(cmder.sub_code))
	sb.WriteByte(byte(cmder.code))
	sb.Write(makeTime(time.Now()))

	vin32 := make([]byte, 4)
	binary.BigEndian.PutUint32(vin32, uint32(c.vin))
	sb.Write(vin32)

	sb.WriteByte(byte(sb.Len()))
	sb.WriteString(shared.PREFIX_COMMAND)

	bytes := util.Reverse([]byte(sb.String()))
	return bytes, nil
}

func makeTime(t time.Time) []byte {
	var sb strings.Builder

	sb.WriteByte(byte(t.Year() - 2000))
	sb.WriteByte(byte(t.Month()))
	sb.WriteByte(byte(t.Day()))
	sb.WriteByte(byte(t.Hour()))
	sb.WriteByte(byte(t.Minute()))
	sb.WriteByte(byte(t.Second()))
	sb.WriteByte(byte(t.Weekday()))

	return util.Reverse([]byte(sb.String()))
}

func makeBool(b bool) []byte {
	enc := []byte{0}
	if b {
		enc[0] = 1
	}
	return enc
}
