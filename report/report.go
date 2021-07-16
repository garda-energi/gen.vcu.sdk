package report

import (
	"bytes"
	"reflect"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type Report struct {
	bytes  []byte // TODO: remove me
	reader *bytes.Reader
}

func New(raw []byte) *Report {
	return &Report{
		bytes:  raw,
		reader: bytes.NewReader(raw),
	}
}

func (r *Report) DecodeReportStruct() (interface{}, error) {
	// header := HeaderReport{}
	// if err := r.decode(&header); err != nil {
	// 	return nil, err
	// }
	// var result interface{}
	// if simpleFrame(header) {
	// 	simple := ReportSimple{}
	// 	if err := r.decode(&simple); err != nil {
	// 		return nil, err
	// 	}
	// 	result = simple
	// } else {
	// 	full := ReportFull{}
	// 	if err := r.decode(&full); err != nil {
	// 		return nil, err
	// 	}
	// 	result = full
	// }

	result := &ReportPacket{}
	if err := r.Decode(result); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Report) decode(dst interface{}) error {
	return tagWalk(bytes.NewReader(r.bytes), reflect.ValueOf(dst), "")
}

func simpleFrame(header HeaderReport) bool {
	return header.FrameID == shared.FRAME_ID_SIMPLE
}
