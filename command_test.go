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

func TestCommandWithoutResponse(t *testing.T) {
	cmder := newFakeCommander(nil)
	defer cmder.Destroy()

	_, err := cmder.GenInfo()
	wantErr := errPacketTimeout("ack")

	if err != wantErr {
		t.Fatalf("want %s, got %s", wantErr, err)
	}
}
func TestCommandInvalidAck(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_COMMAND),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()
	wantErr := errPacketAckCorrupt

	if err != wantErr {
		t.Fatalf("want %s, got %s", wantErr, err)
	}
}

func TestCommandValidAckNoResponse(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_ACK),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()
	wantErr := errPacketTimeout("response")

	if err != wantErr {
		t.Fatalf("want %s, got %s", wantErr, err)
	}
}

func TestValidResponse(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_ACK),
		hexToByte(resGenInfo),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()

	if err != nil {
		t.Fatalf("got no error, got %s\n", err)
	}
}

func TestValidResponseOtherCommand(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_ACK),
		hexToByte(resGenLed),
	})
	defer cmder.Destroy()

	err := cmder.GenLed(true)

	if err != nil {
		t.Fatalf("want no error, got %s\n", err)
	}
}

func TestResponseInvalidPrefix(t *testing.T) {
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
}

func TestResponseInvalidSize(t *testing.T) {
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
}

func TestResponseInvalidVin(t *testing.T) {
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
}

func TestResponseInvalidCode(t *testing.T) {
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
}

func TestResponseInvalidResCode(t *testing.T) {
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
}
func TestResponseResCodeError(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_ACK),
		hexToByte(resGenInfoResCodeError),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()
	if err == nil {
		t.Fatal("want an error, got none")
	}
}

func TestResponseOverflow(t *testing.T) {
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
}

func newFakeCommander(responses [][]byte) *commander {
	vin := 354313
	logging := false
	broker := &fakeBroker{
		responses: responses,
		cmdChan:   make(chan []byte),
		resChan:   make(chan struct{}),
	}

	cmder, _ := newCommander(vin, broker, logging)
	return cmder
}
