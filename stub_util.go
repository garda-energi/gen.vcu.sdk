package sdk

import (
	"log"
	"sync"
	"testing"
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
	stubClient :=
		&stubMqttClient{
			connected: connected,
			responses: &sync.Map{},
			ch: struct {
				res *sync.Map
				cmd *sync.Map
				rep *sync.Map
				sts *sync.Map
			}{
				res: &sync.Map{},
				cmd: &sync.Map{},
				rep: &sync.Map{},
				sts: &sync.Map{},
			},
		}

	return newFakeClient(stubClient, l)
}

func newStubCommander(vin int) *commander {
	logger := newLogger(false, "TEST")
	client := newStubClient(logger, true)
	sleeper := &stubSleeper{
		sleep: time.Millisecond,
		after: 150 * time.Millisecond,
	}

	cmder, err := newCommander(vin, client, sleeper, logger)
	if err != nil {
		log.Fatal(err)
	}
	return cmder
}

func sdkStubClient(api *Sdk) *stubMqttClient {
	return api.client.Client.(*stubMqttClient)
}

func cmderStubClient(cmder *commander) *stubMqttClient {
	return cmder.client.Client.(*stubMqttClient)
}

func assertSubscribed(t *testing.T, api *Sdk, subscribed bool, vins []int) {
	t.Helper()
	time.Sleep(time.Millisecond)

	for _, vin := range vins {
		_, found := sdkStubClient(api).ch.rep.Load(vin)
		if subscribed {
			if !found {
				t.Fatalf("%s want %v, got none", TOPIC_REPORT, vin)
			}
		} else {
			if found {
				t.Fatalf("%s want no %v, got one", TOPIC_REPORT, vin)
			}
		}
	}
}
