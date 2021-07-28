package sdk

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSdk(t *testing.T) {
	listener := Listener{
		StatusFunc: func(vin int, online bool) {
			fmt.Println(vin, "=>", online)
		},
		ReportFunc: func(vin int, report *ReportPacket) {
			fmt.Println(report)
		},
	}

	t.Run("with dis/connected client", func(t *testing.T) {
		api := newStubApi()

		vin := 100

		got := api.AddListener(listener, vin)
		want := errClientDisconnected
		if want != got {
			t.Errorf("want %s, got %s", want, got)
		}

		api.Connect()
		got = api.AddListener(listener, vin)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("check the un/subscribed vins", func(t *testing.T) {
		api := newStubApi()
		api.Connect()

		vins := VinRange(5, 10)

		_ = api.AddListener(listener, vins...)
		assertSubscribed(t, api, true, TOPIC_STATUS, vins)
		assertSubscribed(t, api, true, TOPIC_REPORT, vins)

		addVins := []int{13, 15}
		curVins := append(vins, addVins...)
		_ = api.AddListener(listener, addVins...)
		assertSubscribed(t, api, true, TOPIC_STATUS, curVins)
		assertSubscribed(t, api, true, TOPIC_REPORT, curVins)

		delVins := []int{4, 5, 6, 15}
		curVins = []int{7, 8, 9, 10, 13}
		api.RemoveListener(delVins...)
		assertSubscribed(t, api, false, TOPIC_STATUS, delVins)
		assertSubscribed(t, api, false, TOPIC_REPORT, delVins)
		assertSubscribed(t, api, true, TOPIC_STATUS, curVins)
		assertSubscribed(t, api, true, TOPIC_REPORT, curVins)
	})
}

func TestSdkAddListener(t *testing.T) {
	t.Run("without listener", func(t *testing.T) {
		api := newStubApi()
		api.Connect()

		want := "at least 1 listener supplied"
		got := api.AddListener(Listener{}, 123)
		if want != got.Error() {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("without vin args", func(t *testing.T) {
		api := newStubApi()
		api.Connect()

		want := "at least 1 vin supplied"
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		})
		if want != got.Error() {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("with only 1 listener, 1 vins", func(t *testing.T) {
		api := newStubApi()
		api.Connect()

		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		}, 123)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("with 2 listener, 1 vins", func(t *testing.T) {
		api := newStubApi()
		api.Connect()

		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
			ReportFunc: func(vin int, report *ReportPacket) {},
		}, 123)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("use VinRange() as input", func(t *testing.T) {
		api := newStubApi()
		api.Connect()

		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		}, VinRange(1, 20)...)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("check VinRange() output", func(t *testing.T) {
		api := newStubApi()
		api.Connect()

		want := []int{1, 2, 3, 4}
		got := VinRange(4, 1)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

func assertSubscribed(t *testing.T, api *Sdk, subscribed bool, topic string, vins []int) {
	t.Helper()

	fc := sdkStubClient(api)
	for _, vin := range vins {
		_, found := fc.vins[vin][topic]
		if subscribed {
			if !found {
				t.Fatalf("%s want %v, got none", topic, vin)
			}
		} else {
			if found {
				t.Fatalf("%s want no %v, got one", topic, vin)
			}
		}
	}
}
