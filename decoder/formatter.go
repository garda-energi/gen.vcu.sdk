package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

var endian = binary.LittleEndian

func toAscii(b []byte) string {
	return string(b)
}

func toUint8(b []byte) uint8 {
	return uint8(b[0])
}

func toInt8(b []byte) int8 {
	return int8(b[0])
}

func toUint16(b []byte) uint16 {
	return endian.Uint16(b)
}

func toInt16(b []byte) int16 {
	var data int16
	buf := bytes.NewBuffer(b)
	binary.Read(buf, endian, &data)
	return data
}

func toUint32(b []byte) uint32 {
	return endian.Uint32(b)
}

func toInt32(b []byte) int32 {
	var data int32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, endian, &data)
	return data
}

func toUnixTime(b []byte) interface{} {
	var data string
	for _, v := range b {
		data += fmt.Sprintf("%d", uint8(v))
	}

	// layout := "2006-01-02T15:04:05.000Z"
	layout := "060102150405"
	datetime, _ := time.Parse(layout, data)

	return datetime.Unix()
}
