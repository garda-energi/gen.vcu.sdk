package decoder

import (
	"bytes"
	"encoding/binary"
	"reflect"

	"github.com/pudjamansyurin/gen-go-packet/packet"
	"github.com/pudjamansyurin/gen-go-packet/util"
)

type Report struct {
	Bytes []byte
	// Reader *bytes.Reader
}

func (r *Report) Decode() (interface{}, error) {
	reader := bytes.NewReader(r.Bytes)

	// var decoded interface{}
	if r.isSimpleFrame() {
		// decoded, _ = r.decode(&packet.ReportSimplePacket{})
	} else {
		// decoded, _ = r.decode(&packet.ReportFullPacket{})

		rpt := packet.ReportFullPacket{}

		if err := packet.TagWalk(reader, reflect.ValueOf(&rpt), ""); err != nil {
			return nil, err
		}

		util.Debug(rpt)
	}

	return nil, nil
}

func (r *Report) isSimpleFrame() bool {
	// TODO: check prefix
	return len(r.Bytes) == binary.Size(packet.ReportSimplePacket{})
}

// func (r *Report) Validate() error {
// 	length := len(r.Bytes)

// 	minLength := int(unsafe.Sizeof(packet.ReportPacket{}))
// 	if length < minLength {
// 		return fmt.Errorf("less report length, %d < %d", length, minLength)
// 	}
// 	return nil
// }
