package sdk

import (
	"bytes"
	"encoding/binary"
	"time"
)

// encode combine command and payload to bytes packet.
func encodeCommand(vin int, cmd *command, payload []byte) ([]byte, error) {
	if len(payload) > PAYLOAD_LEN_MAX {
		return nil, errInputOutOfRange("payload")
	}

	var buf bytes.Buffer
	ed := binary.LittleEndian
	binary.Write(&buf, ed, reverseBytes(payload))
	binary.Write(&buf, ed, cmd.sub_code)
	binary.Write(&buf, ed, cmd.code)
	binary.Write(&buf, ed, reverseBytes(timeToBytes(time.Now())))
	binary.Write(&buf, binary.BigEndian, uint32(vin))
	binary.Write(&buf, ed, byte(buf.Len()))
	binary.Write(&buf, ed, []byte(PREFIX_COMMAND))
	bytes := reverseBytes(buf.Bytes())
	return bytes, nil
}

// decode extract header and message response from bytes packet.
func decodeResponse(packet []byte) (*ResponsePacket, error) {
	reader := bytes.NewReader(packet)

	r := &ResponsePacket{
		Header: &HeaderResponse{},
	}

	// header
	if err := decode(reader, r.Header); err != nil {
		return nil, err
	}

	// message
	if len := reader.Len(); len > 0 {
		msg := make([]byte, len)
		reader.Read(msg)
		r.Message = msg
	}

	return r, nil
}
