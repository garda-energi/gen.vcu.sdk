package sdk

import (
	"log"
	"reflect"
	"time"
)

func newApi() *Sdk {
	logger := newLogger(false, "TEST")

	return &Sdk{
		logger: logger,
		client: newFakeClient(logger, false, nil),
	}
}

func newFakeResponse(vin int, invoker string, modifier func(*responsePacket)) [][]byte {
	cmd, err := getCmdByInvoker(invoker)
	if err != nil {
		log.Fatal(err)
	}

	// get default rp, and modify it
	rp := newResponsePacket(vin, cmd, nil)
	if modifier != nil {
		modifier(rp)
	}

	// encode
	resBytes, err := encode(rp)
	if err != nil {
		log.Fatal(err)
	}
	if rp.Header.Size == 0 {
		resBytes[2] = uint8(len(resBytes) - 3)
	}

	return [][]byte{
		strToBytes(PREFIX_ACK),
		resBytes,
	}
}

func newFakeCommander(vin int, responses [][]byte) *commander {
	logger := newLogger(false, "TEST")
	client := newFakeClient(logger, true, responses)

	sleeper := &fakeSleeper{
		sleep: time.Millisecond,
		after: 125 * time.Millisecond,
	}

	cmder, err := newCommander(vin, client, sleeper, logger)
	if err != nil {
		log.Fatal(err)
	}
	return cmder
}

func callFakeCmd(cmder *commander, invoker string, arg interface{}) (res, err interface{}) {
	method := reflect.ValueOf(cmder).MethodByName(invoker)
	ins := []reflect.Value{}
	if arg != nil {
		ins = append(ins, reflect.ValueOf(arg))
	}
	outs := method.Call(ins)

	err = outs[len(outs)-1].Interface()
	if len(outs) > 1 {
		res = outs[0].Interface()
	}
	return
}

func findVinIn(vins []int, vin int) int {
	for i, v := range vins {
		if v == vin {
			return i
		}
	}
	return -1
}
