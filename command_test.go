package sdk

import (
	"testing"
)

const (
	resGenLed                = "53400E096805001507180D181F01000101"
	resGenInfo               = "534023096805001507180D19130100000156435520762E3636342C2047454E202D2032303231"
	resGenInfoInvalidPrefix  = "123423096805001507180D19130100000156435520762E3636342C2047454E202D2032303231"
	resGenInfoInvalidSize    = "534023096805001507180D19130100000156435520762E3636342C204"
	resGenInfoInvalidVin     = "534023123456781507180D19130100000156435520762E3636342C2047454E202D2032303231"
	resGenInfoInvalidCode    = "534023096805001507180D19130167890156435520762E3636342C2047454E202D2032303231"
	resGenInfoInvalidResCode = "534023096805001507180D19130100009956435520762E3636342C2047454E202D2032303231"
	resGenInfoResCodeError   = "534023096805001507180D19130100000056435520762E3636342C2047454E202D2032303231"
	resGenInfoOverflow       = "5340E0096805001507180D19130100000156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D203230323156435520762E3636342C2047454E202D2032303231"
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
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfo),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()

		if err != nil {
			t.Fatalf("want no error, got %s\n", err)
		}
	})

	t.Run("different command, valid packet", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenLed),
		})
		defer cmder.Destroy()

		err := cmder.GenLed(true)

		if err != nil {
			t.Fatalf("want no error, got %s\n", err)
		}
	})

	t.Run("invalid prefix", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfoInvalidPrefix),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errInvalidPrefix

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}

	})

	t.Run("invalid size", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfoInvalidSize),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errInvalidSize

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("invalid VIN", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfoInvalidVin),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errInvalidVin

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("invalid code", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfoInvalidCode),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errInvalidCode

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("invalid resCode", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfoInvalidResCode),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errInvalidResCode

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})

	t.Run("simulate code error", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfoResCodeError),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		if err == nil {
			t.Fatal("want an error, got none")
		}
	})

	t.Run("message overflowed", func(t *testing.T) {
		cmder := newFakeCommander([][]byte{
			strToBytes(PREFIX_ACK),
			hexToByte(resGenInfoOverflow),
		})
		defer cmder.Destroy()

		_, err := cmder.GenInfo()
		wantErr := errResMessageOverflow

		if err != wantErr {
			t.Fatalf("want %s, got %s", wantErr, err)
		}
	})
}

func newFakeCommander(responses [][]byte) *commander {
	vin := 354313
	logging := false
	broker := &fakeBroker{
		responses: responses,
		cmdChan:   make(chan []byte),
		resChan:   make(chan struct{}),
	}

	cmder, _ := newCommander(vin, broker, &fakeSleeper{}, logging)
	return cmder
}
