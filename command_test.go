package sdk

import (
	"reflect"
	"testing"
	"time"
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
