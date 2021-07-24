package sdk

import (
	"testing"
)

const (
	resGenInfo = "534023096805001507180D19130100000156435520762E3636342C2047454E202D2032303231"
	resGenLed  = "53400E096805001507180D181F01000101"
)

func TestCommandWithNoResponse(t *testing.T) {
	cmder := newFakeCommander(nil)
	defer cmder.Destroy()

	_, err := cmder.GenInfo()

	wantErr := errPacketTimeout("ack")
	if err != wantErr {
		t.Fatalf("expected %s, got %s", wantErr, err)
	}
}
func TestCommandWithInvalidAck(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_COMMAND),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()

	wantErr := errPacketCorrupt("ack")
	if err != wantErr {
		t.Fatalf("expected %s, got %s", wantErr, err)
	}
}

func TestCommandWithValidAck(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_ACK),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()
	wantErr := errPacketTimeout("response")
	if err != wantErr {
		t.Fatalf("expected %s, got %s", wantErr, err)
	}
}

func TestCommandWithInvalidResponse(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_ACK),
		hexToByte(resGenLed),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()
	wantErr := errPacketCorrupt("response")
	if err != wantErr {
		t.Fatalf("expected %s, got %s", wantErr, err)
	}
}

func TestCommandWithValidResponse(t *testing.T) {
	cmder := newFakeCommander([][]byte{
		strToBytes(PREFIX_ACK),
		hexToByte(resGenInfo),
	})
	defer cmder.Destroy()

	_, err := cmder.GenInfo()
	if err != nil {
		t.Fatalf("unexpected error: %s\n", err)
	}
}

func newFakeCommander(responses [][]byte) *commander {
	vin := 1234
	logging := !false
	broker := &fakeBroker{
		responses: responses,
		cmdChan:   make(chan []byte),
		resChan:   make(chan struct{}),
	}

	cmder, _ := newCommander(vin, broker, logging)
	return cmder
}
