package sdk

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// typeOfTime is for comparing struct type as time.Time
var typeOfTime reflect.Type = reflect.ValueOf(time.Now()).Type()

// typeOfMessage is for comparing slice type as message ([]byte)
var typeOfMessage reflect.Type = reflect.ValueOf(message{}).Type()

// typeOfPacketData is for comparing struct type as PacketData
var typeOfPacketData reflect.Type = reflect.ValueOf(PacketData{}).Type()

var (
	errClientDisconnected = errors.New("client disconnected")
	errCmdNotFound        = errors.New("command not found")
	errPacketAckCorrupt   = errors.New("packet ack corrupt")
	errInvalidPrefix      = errors.New("invalid prefix")
	errInvalidSize        = errors.New("invalid size")
	errInvalidVin         = errors.New("invalid vin")
	errInvalidCmdCode     = errors.New("invalid cmd code")
	errInvalidResCode     = errors.New("invalid res code")
)

type errPacketTimeout string

func (e errPacketTimeout) Error() string {
	return fmt.Sprintf("packet %s timeout", string(e))
}

type errInputOutOfRange string

func (e errInputOutOfRange) Error() string {
	return fmt.Sprintf("input %s out of range", string(e))
}

// Sleeper is building block for sleep things
type Sleeper interface {
	// Sleep pauses the current goroutine for at least the duration d.
	// A negative or zero duration causes Sleep to return immediately.
	Sleep(time.Duration)
	// After waits for the duration to elapse and then sends the current time
	// on the returned channel.
	After(d time.Duration) <-chan time.Time
}

// realSleeper implement real sleep using time module
type realSleeper struct{}

func (*realSleeper) Sleep(d time.Duration) {
	time.Sleep(d)
}
func (*realSleeper) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
