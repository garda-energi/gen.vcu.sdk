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
}

// func (h *Header) Validate() error {
// 	length := h.Buf.Size()
// 	minLength := int(unsafe.Sizeof(packet.HeaderPacket{}))
// 	if length < minLength {
// 		return fmt.Errorf("less header length is, %d < %d", length, minLength)
// 	}
// 	return nil
// }
