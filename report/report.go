package report

import (
	"bytes"
	"reflect"
)

type Report struct {
	bytes []byte
}

func New(raw []byte) *Report {
	return &Report{
		bytes: raw,
	}
}

// func (r *Report) DecodeReportPacket() (interface{}, error) {

// }

func (r *Report) DecodeReport() (interface{}, error) {
	header := HeaderReport{}
	if err := r.decode(&header); err != nil {
		return nil, err
	}

	var result interface{}
	if simpleFrame(header) {
		simple := ReportSimple{}
		if err := r.decode(&simple); err != nil {
			return nil, err
		}
		result = simple
	} else {
		full := ReportFull{}
		if err := r.decode(&full); err != nil {
			return nil, err
		}
		result = full
	}

	return result, nil
}

func (r *Report) decode(dst interface{}) error {
	return tagWalk(bytes.NewReader(r.bytes), reflect.ValueOf(dst), "")
}

func simpleFrame(header HeaderReport) bool {
	return header.FrameID == FRAME_ID_SIMPLE
}
