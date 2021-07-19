package command

import (
	"encoding/binary"
	"strings"
	"time"
	"errors"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func toUint8(b []byte) uint8 {
	return uint8(b[0])
}

func (c *Command) encode(cmder *Commander, payload []byte) ([]byte, error) {
	var sb strings.Builder

	if len(payload) > PAYLOAD_LEN {
		return nil, errors.New("payload overload")
	}
	sb.Write(payload)
	sb.WriteByte(byte(cmder.SubCode))
	sb.WriteByte(byte(cmder.Code))
	sb.Write(buildTime(time.Now()))

	vin32 := make([]byte, 4)
	binary.BigEndian.PutUint32(vin32, uint32(c.vin))
	sb.Write(vin32)

	sb.WriteByte(byte(sb.Len()))
	sb.WriteString(shared.PREFIX_COMMAND)

	bytes := sbToPacket(sb)
	// util.Debug(util.HexString(bytes))

	return bytes, nil
}

func buildTime(t time.Time) []byte {
	var sb strings.Builder

	sb.WriteByte(byte(t.Year() - 2000))
	sb.WriteByte(byte(t.Month()))
	sb.WriteByte(byte(t.Day()))
	sb.WriteByte(byte(t.Hour()))
	sb.WriteByte(byte(t.Minute()))
	sb.WriteByte(byte(t.Second()))
	sb.WriteByte(byte(t.Weekday()))

	bytes := sbToPacket(sb)
	return bytes
}

func sbToPacket(sb strings.Builder) []byte {
	return util.Reverse([]byte(sb.String()))
}
