package sdk

import (
	"bytes"
)

type report struct {
	reader *bytes.Reader
}

// newReport create new instance of *report
func newReport(raw []byte) *report {
	return &report{
		reader: bytes.NewReader(raw),
	}
}

// decode report packet from reader
func (r *report) decode() (*ReportPacket, error) {
	result := &ReportPacket{}
	if err := decode(r.reader, result); err != nil {
		return nil, err
	}
	return result, nil
}
