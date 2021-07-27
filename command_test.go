package sdk

import (
	"fmt"
	"testing"
	"time"
)

const testVin = 354313

func newFakeCommander(responses [][]byte) *commander {
	logging := false
	client := &fakeClient{
		responses: responses,
		cmdChan:   make(chan []byte),
		resChan:   make(chan struct{}),
	}
	sleeper := &fakeSleeper{
		sleep: time.Millisecond,
		after: 125 * time.Millisecond,
	}
	cmder, _ := newCommander(testVin, client, sleeper, logging)
	return cmder
}
func TestResponsePacket(t *testing.T) {
	testCases := []struct {
		desc      string
		wantErr   error
		responses [][]byte
	}{
		{
			desc:      "no packet",
			wantErr:   errPacketTimeout("ack"),
			responses: nil,
		},
		{
			desc:      "invalid ack packet",
			wantErr:   errPacketAckCorrupt,
			responses: [][]byte{strToBytes(PREFIX_REPORT)},
		},
		{
			desc:      "only valid ack packet",
			wantErr:   errPacketTimeout("response"),
			responses: [][]byte{strToBytes(PREFIX_ACK)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmder := newFakeCommander(tC.responses)
			defer cmder.Destroy()

			_, err := cmder.GenInfo()

			if err != tC.wantErr {
				t.Fatalf("want %s, got %s", tC.wantErr, err)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	testCases := []struct {
		desc      string
		wantErr   string
		formatter func(r *responsePacket)
	}{
		{
			desc:    "invalid prefix",
			wantErr: errInvalidPrefix.Error(),
			formatter: func(r *responsePacket) {
				r.Header.Prefix = PREFIX_REPORT
			},
		},
		{
			desc:    "invalid size",
			wantErr: errInvalidSize.Error(),
			formatter: func(r *responsePacket) {
				r.Header.Size = 55
			},
		},
		{
			desc:    "invalid VIN",
			wantErr: errInvalidVin.Error(),
			formatter: func(r *responsePacket) {
				r.Header.Vin = 12345
			},
		},
		{
			desc:    "invalid cmd code, not registered",
			wantErr: errInvalidCmdCode.Error(),
			formatter: func(r *responsePacket) {
				r.Header.Code = 9
				r.Header.SubCode = 5
			},
		},
		{
			desc:    "invalid cmd code, other command",
			wantErr: errInvalidCmdCode.Error(),
			formatter: func(r *responsePacket) {
				r.Header.Code = 1
				r.Header.SubCode = 1
			},
		},
		{
			desc:    "invalid resCode",
			wantErr: errInvalidResCode.Error(),
			formatter: func(r *responsePacket) {
				r.Header.ResCode = 99
			},
		},
		{
			desc:    "message overflowed",
			wantErr: errInvalidSize.Error(),
			formatter: func(r *responsePacket) {
				r.Message = message("Golang is very useful for writing light-weight microservices. We currently use it for generating APIs that interact with our front-end applications. If you want to build a small functional microservice quickly, then Golang is a great tool to use. It's an easy language for developers to learn quickly.")
			},
		},
		{
			desc:    "simulate code error, no message",
			wantErr: resCodeError.String(),
			formatter: func(r *responsePacket) {
				r.Header.ResCode = resCodeError
			},
		},
		{
			desc:    "simulate code error, with message",
			wantErr: fmt.Sprintf("%s, State should = %s.", resCodeError, BikeStateStandby),
			formatter: func(r *responsePacket) {
				r.Header.ResCode = resCodeError
				r.Message = message("State should = {1}.")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := fakeResponse(testVin, "GenInfo")
			tC.formatter(res)

			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				mockResponse(res),
			})
			defer cmder.Destroy()

			_, err := cmder.GenInfo()
			if err.Error() != tC.wantErr {
				t.Fatalf("want %s, got %s", tC.wantErr, err)
			}
		})
	}
}
