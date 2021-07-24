package sdk

import (
	"bytes"
	"encoding/binary"
	"time"
)

// encode combine command and value to bytes packet.
func encodeCommand(vin int, cmd *command, val payload) ([]byte, error) {
	if val.overflow() {
		return nil, errInputOutOfRange("payload")
	}

	var buf bytes.Buffer
	ed := binary.LittleEndian
	binary.Write(&buf, ed, reverseBytes(val))
	binary.Write(&buf, ed, cmd.subCode)
	binary.Write(&buf, ed, cmd.code)
	binary.Write(&buf, ed, reverseBytes(timeToBytes(time.Now())))
	binary.Write(&buf, binary.BigEndian, uint32(vin))
	binary.Write(&buf, ed, byte(buf.Len()))
	binary.Write(&buf, ed, []byte(PREFIX_COMMAND))
	bytes := reverseBytes(buf.Bytes())
	return bytes, nil
}

// decode extract header and message response from bytes packet.
func decodeResponse(packet []byte) (*responsePacket, error) {
	reader := bytes.NewReader(packet)

	r := &responsePacket{
		Header: &headerResponse{},
	}

	// header
	if err := decode(reader, r.Header); err != nil {
		return nil, err
	}

	// message
	if reader.Len() > 0 {
		r.Message = make(payload, reader.Len())
		reader.Read(r.Message)
	}

	return r, nil
}
