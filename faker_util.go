package sdk

import (
	"log"
	"reflect"
	"time"
)

func hasVin(vins []int, vin int) bool {
	for _, v := range vins {
		if v == vin {
			return true
		}
	}
	return false
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
