package sdk

import (
	"log"
	"math/rand"
	"reflect"
	"time"
)

func newFakeApi() *Sdk {
	logger := newLogger(false, "TEST")

	return &Sdk{
		logger: logger,
		client: newFakeClient(logger, false),
	}
}

func newFakeCommander(vin int) *commander {
	logger := newLogger(false, "TEST")
	client := newFakeClient(logger, true)

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

func newFakeReport(vin int) *ReportPacket {
	return &ReportPacket{
		Header: &HeaderReport{
			Header: Header{
				Prefix:       TOPIC_REPORT,
				Size:         0,
				Vin:          uint32(vin),
				SendDatetime: time.Now(),
			},
			Frame: [...]Frame{
				FrameSimple,
				FrameFull,
			}[rand.Intn(2)],
		},
		Vcu: &Vcu{
			LogDatetime: time.Now().Add(-2 * time.Second),
			State: [...]BikeState{
				BikeStateNormal,
				BikeStateReady,
				BikeStateRun,
				BikeStateStandby,
			}[rand.Intn(4)],
			Events:      uint16(rand.Uint32()),
			LogBuffered: uint8(rand.Uint32()),
			BatVoltage:  rand.Float32(),
			Uptime:      rand.Float32(),
		},
	}
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

func sdkFakeClient(api *Sdk) *fakeMqttClient {
	return api.client.Client.(*fakeMqttClient)
}

func cmderFakeClient(cmder *commander) *fakeMqttClient {
	return cmder.client.Client.(*fakeMqttClient)
}
