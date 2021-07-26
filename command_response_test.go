package sdk

import (
	"errors"
	"testing"
)

func TestResponse(t *testing.T) {
	t.Run("no packet", func(t *testing.T) {
		cmder := newFakeCommander(nil)
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errPacketTimeout("ack")

		if err != wantErr {
			t.Errorf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("invalid ack packet", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_COMMAND),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errPacketAckCorrupt

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("only valid ack packet", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errPacketTimeout("response")

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("valid packet", func(t *testing.T) {
		wantMsg := "VCU v.664, GEN - 2021"
		res := newFakeResponse(testVin, "GEN_INFO")
		res.Message = message(wantMsg)

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		msg, err := cmder.GenInfo()

		if err != nil {
			t.Fatalf("want no error, got %s\n", err)
		}
		if msg != wantMsg {
			t.Fatalf("want %s, got %s", wantMsg, msg)
		}
	})
}

func TestResponseError(t *testing.T) {
	testCases := []struct {
		desc      string
		wantErr   error
		formatter func(r *responsePacket)
	}{
		{
			desc:    "invalid prefix",
			wantErr: errInvalidPrefix,
			formatter: func(r *responsePacket) {
				r.Header.Prefix = PREFIX_REPORT
			},
		},
		{
			desc:    "invalid size",
			wantErr: errInvalidSize,
			formatter: func(r *responsePacket) {
				r.Header.Size = 55
			},
		},
		{
			desc:    "invalid VIN",
			wantErr: errInvalidVin,
			formatter: func(r *responsePacket) {
				r.Header.Vin = 12345
			},
		},
		{
			desc:    "invalid cmd code",
			wantErr: errInvalidCmdCode,
			formatter: func(r *responsePacket) {
				r.Header.Code = 9
				r.Header.SubCode = 5
			},
		},
		{
			desc:    "invalid resCode",
			wantErr: errInvalidResCode,
			formatter: func(r *responsePacket) {
				r.Header.ResCode = 99
			},
		},
		{
			desc:    "message overflowed",
			wantErr: errInvalidSize,
			formatter: func(r *responsePacket) {
				r.Message = message("'Google Go' redirects here. For the Android search app by Google, 'Google Go', for low-end Lollipop+ devices, see Android Go. For the computer program by Google to play the board game Go, see AlphaGo. For the 2003 agent-based programming language, see Go! (programming language).")
				r.Header.Size = uint8(len(r.Message))
			},
		},
		{
			desc:    "simulate code error",
			wantErr: errors.New(resCodeError.String()),
			formatter: func(r *responsePacket) {
				r.Header.ResCode = resCodeError
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := newFakeResponse(testVin, "GEN_INFO")
			tC.formatter(res)

			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				mockResponse(res),
			})
			defer cmder.Destroy()

			_, err := cmder.GenInfo()
			if err.Error() != tC.wantErr.Error() {
				t.Fatalf("want %s, got %s", tC.wantErr, err)
			}
		})
	}
}
