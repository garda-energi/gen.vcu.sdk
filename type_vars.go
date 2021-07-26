package sdk

import (
	"errors"
	"fmt"
	"time"
)

var (
	errPacketAckCorrupt   = errors.New("packet ack corrupt")
	errInvalidPrefix      = errors.New("prefix invalid")
	errInvalidSize        = errors.New("size invalid")
	errInvalidVin         = errors.New("vin invalid")
	errInvalidCode        = errors.New("code invalid")
	errInvalidResCode     = errors.New("resCode invalid")
	errResMessageOverflow = errors.New("message overflow")
)

type errPacketTimeout string

func (e errPacketTimeout) Error() string {
	return fmt.Sprintf("packet %s timeout", string(e))
}

type errInputOutOfRange string

func (e errInputOutOfRange) Error() string {
	return fmt.Sprintf("input %s out of range", string(e))
}

// Sleeper is building block for sleep function
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
