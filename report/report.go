package report

import (
	"bytes"
)

type Report struct {
	reader *bytes.Reader
}

func New(raw []byte) *Report {
	return &Report{
		reader: bytes.NewReader(raw),
	}
}

func (r *Report) Decode() (*ReportPacket, error) {
	result := &ReportPacket{}
	if err := r.decode(result); err != nil {
		return nil, err
	}

	return result, nil
}
