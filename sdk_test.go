package sdk

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSdk(t *testing.T) {
	api := Sdk{
		logger: newLogger(false, "TEST"),
		client: newFakeClient(false, nil),
	}

	// prepare the status & report listener
	listener := Listener{
		StatusFunc: func(vin int, online bool) {
			fmt.Println(vin, " => ", online)
		},
		ReportFunc: func(vin int, report *ReportPacket) {
			fmt.Println(report)
		},
	}

	t.Run("with disconnected client", func(t *testing.T) {
		want := errClientDisconnected
		got := api.AddListener(listener, 100)
		if want != got {
			t.Fatalf("want %s, got %s", want, got)
		}
	})

	// connect to client
	api.Connect()
	defer api.Disconnect()

	t.Run("with connected client, check the subscribed & unsubscribed vins", func(t *testing.T) {
		vins := VinRange(5, 10)
		got := api.AddListener(listener, vins...)

		if got != nil {
			t.Fatal("want no error, got ", got)
		}

		subscribed := api.client.(*fakeClient).subscribed
		for _, topic := range []string{TOPIC_STATUS, TOPIC_REPORT} {
			gotVins := subscribed[topic]
			if !reflect.DeepEqual(vins, gotVins) {
				t.Fatalf("%s want %v, got %v", topic, vins, gotVins)
			}
		}

		api.RemoveListener(vins...)
		for _, topic := range []string{TOPIC_STATUS, TOPIC_REPORT} {
			gotLen := len(subscribed[topic])
			if gotLen > 0 {
				t.Fatalf("want 0 vins, got %d", gotLen)
			}
		}
	})
}

func TestSdkAddListener(t *testing.T) {
	api := Sdk{
		client: newFakeClient(true, nil),
	}

	t.Run("without listener", func(t *testing.T) {
		want := "at least 1 listener supplied"
		got := api.AddListener(Listener{}, 123)
		if want != got.Error() {
			t.Fatalf("want %s, got %s", want, got)
		}
	})

	t.Run("without vin args", func(t *testing.T) {
		want := "at least 1 vin supplied"
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		})
		if want != got.Error() {
			t.Fatalf("want %s, got %s", want, got)
		}
	})

	t.Run("with only 1 listener", func(t *testing.T) {
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		}, 123)
		if got != nil {
			t.Fatal("want no error, got ", got)
		}
	})

	t.Run("with only 2 listener", func(t *testing.T) {
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
			ReportFunc: func(vin int, report *ReportPacket) {},
		}, 123)
		if got != nil {
			t.Fatal("want no error, got ", got)
		}
	})

	t.Run("use VinRange() as input", func(t *testing.T) {
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		}, VinRange(1, 20)...)
		if got != nil {
			t.Fatal("want no error, got ", got)
		}
	})

	t.Run("check VinRange() output", func(t *testing.T) {
		want := []int{1, 2, 3, 4}
		got := VinRange(4, 1)
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})
}
