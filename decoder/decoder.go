package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

var endian = binary.LittleEndian

func toUnixTime(b []byte) int64 {
	var data string
	for _, v := range b {
		data += fmt.Sprintf("%d", uint8(v))
	}

	datetime, _ := time.Parse("060102150405", data)
	return datetime.Unix()
}

func toBool(b []byte) bool {
	return b[0] == 1
}

func toAscii(b []byte) string {
	return reverse(string(b))
}

func toFloat64(b []byte, factor float32) float64 {
	var data float64

	switch len(b) {
	case 1:
		data = float64(int8(b[0]))
	case 2:
		var d int16
		binary.Read(bytes.NewBuffer(b), endian, &d)
		data = float64(d)
	case 4:
		var d int32
		binary.Read(bytes.NewBuffer(b), endian, &d)
		data = float64(d)
	case 8:
		binary.Read(bytes.NewBuffer(b), endian, &data)
	}

	return data * float64(factor)
}

func toUint64(b []byte) uint64 {
	var data uint64

	switch len(b) {
	case 1:
		data = uint64(uint8(b[0]))
	case 2:
		data = uint64(endian.Uint16(b))
	case 4:
		data = uint64(endian.Uint32(b))
	case 8:
		data = endian.Uint64(b)
	}

	return data
}

func toInt64(b []byte) int64 {
	var data int64

	switch len(b) {
	case 1:
		data = int64(int8(b[0]))
	case 2:
		var d int16
		binary.Read(bytes.NewBuffer(b), endian, &d)
		data = int64(d)
	case 4:
		var d int32
		binary.Read(bytes.NewBuffer(b), endian, &d)
		data = int64(d)
	case 8:
		binary.Read(bytes.NewBuffer(b), endian, &data)
	}

	return data
}

// reverse returns a string with the bytes of s in reverse order.
func reverse(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := len(s) - 1; i >= 0; i-- {
		b.WriteByte(s[i])
	}
	return b.String()
}
