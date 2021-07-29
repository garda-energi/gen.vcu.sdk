package sdk

import (
	"log"
	"reflect"
	"sync"
	"time"
)

func newStubApi() *Sdk {
	logger := newLogger(false, "TEST")
	return &Sdk{
		logger: logger,
		client: newStubClient(logger, false),
	}
}

func newStubClient(l *log.Logger, connected bool) *client {
	_ = newClientOptions(&ClientConfig{}, l)
	return &client{
		Client: &stubMqttClient{
			connected: connected,

			vinsMutex: &sync.RWMutex{},
			vins:      make(map[int]map[string]packets),

			cmdChan: make(map[int]chan []byte),
			resChan: make(map[int]chan struct{}),
		},
		logger: l,
	}
}

func newStubCommander(vin int) *commander {
	logger := newLogger(false, "TEST")
	client := newStubClient(logger, true)
	sleeper := &stubSleeper{
		sleep: time.Millisecond,
		after: 125 * time.Millisecond,
	}

	cmder, err := newCommander(vin, client, sleeper, logger)
	if err != nil {
		log.Fatal(err)
	}
	return cmder
}

// func sdkStubClient(api *Sdk) *stubMqttClient {
// 	return api.client.Client.(*stubMqttClient)
// }

func cmderStubClient(cmder *commander) *stubMqttClient {
	return cmder.client.Client.(*stubMqttClient)
}

func callCommand(cmder *commander, invoker string, arg interface{}) (res, err interface{}) {
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
