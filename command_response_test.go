package sdk

import (
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

	t.Run("invalid prefix", func(t *testing.T) {
		res := newFakeResponse(testVin, "AUDIO_BEEP")
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
		res := newFakeResponse(testVin, "AUDIO_BEEP")
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
		res := newFakeResponse(testVin, "AUDIO_BEEP")
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
		res := newFakeResponse(testVin, "AUDIO_BEEP")
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
		res := newFakeResponse(testVin, "AUDIO_BEEP")
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
		res := newFakeResponse(testVin, "AUDIO_BEEP")
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
		res := newFakeResponse(testVin, "GEN_INFO")
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
