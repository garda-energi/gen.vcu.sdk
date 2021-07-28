package sdk

import (
	"fmt"
	"testing"
)

const testVin = 353313

func TestResponsePacket(t *testing.T) {
	testCases := []struct {
		desc string
		want error
		res  []byte
	}{
		{
			desc: "no packet",
			want: errPacketTimeout("ack"),
			res:  nil,
		},
		{
			desc: "invalid ack packet",
			want: errPacketAckCorrupt,
			res:  strToBytes(PREFIX_REPORT),
		},
		{
			desc: "only valid ack packet",
			want: errPacketTimeout("response"),
			res:  strToBytes(PREFIX_ACK),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmder := newStubCommander(testVin)
			defer cmder.Destroy()

			cmderStubClient(cmder).
				mockAck(testVin, tC.res)

			_, err := cmder.GenInfo()

			if err != tC.want {
				t.Errorf("want %s, got %s", tC.want, err)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	testCases := []struct {
		desc     string
		want     string
		modifier func(r *responsePacket)
	}{
		{
			desc: "invalid prefix",
			want: errInvalidPrefix.Error(),
			modifier: func(r *responsePacket) {
				r.Header.Prefix = PREFIX_REPORT
			},
		},
		{
			desc: "invalid size",
			want: errInvalidSize.Error(),
			modifier: func(r *responsePacket) {
				r.Header.Size = 55
			},
		},
		{
			desc: "invalid VIN",
			want: errInvalidVin.Error(),
			modifier: func(r *responsePacket) {
				r.Header.Vin = 12345
			},
		},
		{
			desc: "invalid cmd code, not registered",
			want: errInvalidCmdCode.Error(),
			modifier: func(r *responsePacket) {
				r.Header.Code = 9
				r.Header.SubCode = 5
			},
		},
		{
			desc: "invalid cmd code, other command",
			want: errInvalidCmdCode.Error(),
			modifier: func(r *responsePacket) {
				r.Header.Code = 1
				r.Header.SubCode = 1
			},
		},
		{
			desc: "invalid resCode",
			want: errInvalidResCode.Error(),
			modifier: func(r *responsePacket) {
				r.Header.ResCode = 99
			},
		},
		{
			desc: "message overflowed",
			want: errInvalidSize.Error(),
			modifier: func(r *responsePacket) {
				r.Message = message("Golang is very useful for writing light-weight microservices. We currently use it for generating APIs that interact with our front-end applications. If you want to build a small functional microservice quickly, then Golang is a great tool to use. It's an easy language for developers to learn quickly.")
			},
		},
		{
			desc: "simulate code error, no message",
			want: resCodeError.String(),
			modifier: func(r *responsePacket) {
				r.Header.ResCode = resCodeError
			},
		},
		{
			desc: "simulate code error, with message",
			want: fmt.Sprint(resCodeError, " State should = ", BikeStateStandby),
			modifier: func(r *responsePacket) {
				r.Header.ResCode = resCodeError
				r.Message = message("State should = {1}")
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmder := newStubCommander(testVin)
			defer cmder.Destroy()

			cmderStubClient(cmder).
				mockResponse(testVin, "GenInfo", tC.modifier)

			_, err := cmder.GenInfo()
			if err.Error() != tC.want {
				t.Errorf("want %s, got %s", tC.want, err)
			}
		})
	}
}
