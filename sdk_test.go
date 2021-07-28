package sdk

import (
	"fmt"
	"reflect"
	"testing"
)

var logger = newLogger(false, "TEST")

func TestSdk(t *testing.T) {
	api := Sdk{
		logger: logger,
		client: newFakeClient(logger, false, nil),
	}

	// prepare the status & report listener
	listener := Listener{
		StatusFunc: func(vin int, online bool) {
			fmt.Println(vin, "=>", online)
		},
		ReportFunc: func(vin int, report *ReportPacket) {
			fmt.Println(report)
		},
	}

	t.Run("with disconnected client", func(t *testing.T) {
		api.Disconnect()
		want := errClientDisconnected

		got := api.AddListener(listener, 100)
		if want != got {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("with connected client", func(t *testing.T) {
		api.Connect()
		vins := VinRange(5, 10)

		got := api.AddListener(listener, vins...)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("check the subscribed vins", func(t *testing.T) {
		api.Connect()
		vins := VinRange(5, 10)

		_ = api.AddListener(listener, vins...)
		assertSubscribed(t, &api, TOPIC_STATUS, vins)
		assertSubscribed(t, &api, TOPIC_REPORT, vins)

		moreVins := []int{13, 15}
		wantVins := append(vins, moreVins...)
		_ = api.AddListener(listener, moreVins...)
		assertSubscribed(t, &api, TOPIC_STATUS, wantVins)
		assertSubscribed(t, &api, TOPIC_REPORT, wantVins)
	})

	// t.Run("check the subscribed & unsubscribed vins", func(t *testing.T) {
	// 	api.Connect()
	// 	vins := VinRange(5, 10)

	// 	_ = api.AddListener(listener, vins...)
	// 	subscribed := api.client.Client.(*fakeMqttClient).subscribed
	// 	for _, topic := range []string{TOPIC_STATUS, TOPIC_REPORT} {
	// 		gotVins := subscribed[topic]
	// 		sort.Ints(gotVins)
	// 		if !reflect.DeepEqual(vins, gotVins) {
	// 			t.Errorf("%s want %v, got %v", topic, vins, gotVins)
	// 		}
	// 	}

	// 	api.RemoveListener(vins...)
	// 	for _, topic := range []string{TOPIC_STATUS, TOPIC_REPORT} {
	// 		want := 0
	// 		got := len(subscribed[topic])
	// 		if want != got {
	// 			t.Errorf("want %d vins, got %d", want, got)
	// 		}
	// 	}
	// })
}

func TestSdkAddListener(t *testing.T) {
	api := Sdk{
		client: newFakeClient(logger, true, nil),
	}

	t.Run("without listener", func(t *testing.T) {
		want := "at least 1 listener supplied"
		got := api.AddListener(Listener{}, 123)
		if want != got.Error() {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("without vin args", func(t *testing.T) {
		want := "at least 1 vin supplied"
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		})
		if want != got.Error() {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("with only 1 listener, 1 vins", func(t *testing.T) {
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		}, 123)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("with 2 listener, 1 vins", func(t *testing.T) {
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
			ReportFunc: func(vin int, report *ReportPacket) {},
		}, 123)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("use VinRange() as input", func(t *testing.T) {
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		}, VinRange(1, 20)...)
		if got != nil {
			t.Error("want no error, got ", got)
		}
	})

	t.Run("check VinRange() output", func(t *testing.T) {
		want := []int{1, 2, 3, 4}
		got := VinRange(4, 1)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v, got %v", want, got)
		}
	})
}

func assertSubscribed(t *testing.T, api *Sdk, topic string, vins []int) {
	t.Helper()
	subs := api.client.Client.(*fakeMqttClient).subscribed

	gotVins := subs[topic]
	for _, vin := range vins {
		if !hasVin(gotVins, vin) {
			t.Fatalf("%s want %v, got %v", topic, vins, gotVins)
		}
		// fmt.Println(vin, gotVins)
	}
}
