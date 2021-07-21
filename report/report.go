package report

import (
	"bytes"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
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
	if err := shared.Decode(r.reader, result); err != nil {
		return nil, err
	}

	return result, nil
}
