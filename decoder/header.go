package decoder

import (
	"bufio"
	"encoding/binary"
	"errors"

	"github.com/pudjamansyurin/gen-go-packet/packet"
)

type Header struct {
	Buf *bufio.Reader
}

func (h *Header) Decode() (packet.HeaderPacket, error) {
	var data packet.HeaderPacket

	if err := binary.Read(h.Buf, binary.LittleEndian, &data); err != nil {
		return packet.HeaderPacket{}, errors.New("cant decode header")
	}
	return data, nil

	// data := make(M)
	// for _, v := range packet.Header {
	// 	buf := make([]byte, v.Size)
	// 	n, err := reader.Read(buf)
	// 	if n != v.Size || err != nil {
	// 		return M{}, errors.New("packet corrupted")
	// 	}

	// 	switch v.Format.DataType {
	// 	case reflect.String:
	// 		data[v.Field] = toAscii(buf)
	// 	case reflect.Uint8:
	// 		data[v.Field] = toUint8(buf)
	// 	case reflect.Uint32:
	// 		data[v.Field] = toUint32(buf)
	// 	case reflect.Interface:
	// 		if v.Format.InterfaceName == "datetime" {
	// 			data[v.Field] = toUnixTime(buf)
	// 		}
	// 	}
	// }
	// return data, nil
}

// func (h *Header) Validate() error {
// 	length := h.Buf.Size()
// 	minLength := int(unsafe.Sizeof(packet.HeaderPacket{}))
// 	if length < minLength {
// 		return fmt.Errorf("less header length is, %d < %d", length, minLength)
// 	}
// 	return nil
// }
