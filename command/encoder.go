package command

import (
	"encoding/binary"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

func toUint8(b []byte) uint8 {
	return uint8(b[0])
}

func (c *Command) encode(cmd Commander, payload []byte) []byte {
	var sb strings.Builder

	sb.Write(payload)
	sb.WriteByte(byte(cmd.SubCode))
	sb.WriteByte(byte(cmd.Code))
	sb.Write(buildTime(time.Now()))

	vin32 := make([]byte, 4)
	binary.BigEndian.PutUint32(vin32, uint32(c.vin))
	sb.Write(vin32)

	sb.WriteByte(byte(sb.Len()))
	sb.WriteString(shared.PREFIX_COMMAND)

	bytes := builderToPacket(sb)
	// util.Debug(util.HexString(bytes))
	return bytes
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

	bytes := builderToPacket(sb)
	return bytes
}

func builderToPacket(sb strings.Builder) []byte {
	return util.Reverse([]byte(sb.String()))
}
