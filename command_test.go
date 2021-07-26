package sdk

import (
	"reflect"
	"testing"
	"time"
)

const TEST_VIN = 354313

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
		res := newFakeResponse("GEN_INFO")
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

	t.Run("invalid prefix", func(t *testing.T) {
		res := newFakeResponse("AUDIO_BEEP")
		res.Header.Prefix = PREFIX_REPORT

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		err := cmder.AudioBeep()
		wantErr := errInvalidPrefix

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}

	})

	t.Run("invalid size", func(t *testing.T) {
		res := newFakeResponse("AUDIO_BEEP")
		res.Header.Size = 55

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		err := cmder.AudioBeep()
		wantErr := errInvalidSize

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("invalid VIN", func(t *testing.T) {
		res := newFakeResponse("AUDIO_BEEP")
		res.Header.Vin = 12345

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		err := cmder.AudioBeep()
		wantErr := errInvalidVin

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("invalid code", func(t *testing.T) {
		res := newFakeResponse("AUDIO_BEEP")
		res.Header.Code = 9
		res.Header.SubCode = 5

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		err := cmder.AudioBeep()
		wantErr := errInvalidCode

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("invalid resCode", func(t *testing.T) {
		res := newFakeResponse("AUDIO_BEEP")
		res.Header.ResCode = 99

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		err := cmder.AudioBeep()
		wantErr := errInvalidResCode

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("simulate code error", func(t *testing.T) {
		res := newFakeResponse("AUDIO_BEEP")
		res.Header.ResCode = resCodeError

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		err := cmder.AudioBeep()
		if err == nil {
			t.Fatal("want an error, got none")
		}
	})

	t.Run("message overflowed", func(t *testing.T) {
		res := newFakeResponse("GEN_INFO")
		res.Message = message("'Google Go' redirects here. For the Android search app by Google, 'Google Go', for low-end Lollipop+ devices, see Android Go. For the computer program by Google to play the board game Go, see AlphaGo. For the 2003 agent-based programming language, see Go! (programming language).")
		res.Header.Size = uint8(len(res.Message))

		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			mockResponse(res),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errInvalidSize

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})
}

func TestCommands(t *testing.T) {
	testCases := []struct {
		cmd  string
		args interface{}
		res  string
	}{
		{
			cmd:  "GenInfo",
			args: nil,
			res:  "5340230968050015071A0A15080100000156435520762E3636342C2047454E202D2032303231",
		},
		{
			cmd:  "GenLed",
			args: true,
			res:  "53400E0968050015071A0A180D01000101",
		},
		{
			cmd:  "GenRtc",
			args: time.Now(),
			res:  "53400E0968050015071A0A250301000201",
		},
		{
			cmd:  "GenOdo",
			args: uint16(4321),
			res:  "53400E0968050015071A0A261801000301",
		},
		{
			cmd:  "GenAntiTheaf",
			args: nil,
			res:  "53400E0968050015071A0A272201000401",
		},
		{
			cmd:  "GenReportFlush",
			args: nil,
			res:  "53400E0968050015071A0A281C01000501",
		},
		{
			cmd:  "GenReportBlock",
			args: false,
			res:  "53400E0968050015071A0B100901000601",
		},
		{
			cmd:  "OvdState",
			args: BikeStateNormal,
			res:  "53400E0968050015071A0B111C01010001",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.cmd, func(t *testing.T) {
			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				hexToByte(tC.res),
			})
			defer cmder.Destroy()

			meth := reflect.ValueOf(cmder).MethodByName(tC.cmd)
			args := []reflect.Value{}
			if tC.args != nil {
				args = append(args, reflect.ValueOf(tC.args))
			}
			outs := meth.Call(args)

			if err := outs[len(outs)-1]; !err.IsNil() {
				t.Fatalf("want no error, got %s\n", err)
			}
		})
	}
}

func newFakeResponse(cmdName string) *responsePacket {
	cmd, _ := getCommand(cmdName)

	return &responsePacket{
		Header: &headerResponse{
			HeaderCommand: HeaderCommand{
				Header: Header{
					Prefix:       PREFIX_RESPONSE,
					Size:         0,
					Vin:          uint32(TEST_VIN),
					SendDatetime: time.Now(),
				},
				Code:    cmd.code,
				SubCode: cmd.subCode,
			},
			ResCode: resCodeOk,
		},
		Message: nil,
	}
}

// mockResponse combine response and message to bytes packet.
func mockResponse(r *responsePacket) []byte {
	resBytes, err := encode(&r)
	if err != nil {
		return nil
	}

	// change Header.Size
	if r.Header.Size == 0 {
		resBytes[2] = uint8(len(resBytes) - 3)
	}
	return resBytes
}

func newFakeCommander(responses [][]byte) *commander {
	logging := false
	broker := &fakeBroker{
		responses: responses,
		cmdChan:   make(chan []byte),
		resChan:   make(chan struct{}),
	}

	cmder, _ := newCommander(TEST_VIN, broker, &fakeSleeper{}, logging)
	return cmder
}
