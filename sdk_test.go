package sdk

import (
	"log"
	"reflect"
	"testing"
)

func TestSdk(t *testing.T) {
	api := Sdk{
		logging: false,
		client:  &fakeClient{
			// responses: responses,
			// cmdChan:   make(chan []byte),
			// resChan:   make(chan struct{}),
		},
	}

	// connect to client
	if err := api.Connect(); err != nil {
		log.Fatal(err)
	}
	defer api.Disconnect()

	// // prepare the status & report listener
	// listener := Listener{
	// 	StatusFunc: func(vin int, online bool) {
	// 		status := map[bool]string{
	// 			false: "OFFLINE",
	// 			true:  "ONLINE",
	// 		}[online]
	// 		fmt.Printf("%d => %s\n", vin, status)
	// 	},
	// 	ReportFunc: func(vin int, report *ReportPacket) {
	// 		// fmt.Println(report)

	// 		// show-off all *ReportPacket methods available
	// 		// if report.Vcu.RealtimeData() {
	// 		// 	fmt.Println("Current report is realtime")
	// 		// }
	// 		// if report.Gps.ValidHorizontal() {
	// 		// 	fmt.Println("GPS longitude, latitude & heading is valid")
	// 		// }
	// 		// if report.Bms.LowCapacity() {
	// 		// 	fmt.Println("BMS need to be charged on Charging Station")
	// 		// }
	// 	},
	// }

	// // listen to report
	// // see api.Addlistener doc for usage
	// vins := VinRange(354309, 354323)
	// if err := api.AddListener(listener, vins...); err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	defer api.RemoveListener(vins...)
	// }
}

func TestSdkAddListener(t *testing.T) {
	api := Sdk{
		client: &fakeClient{},
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
			t.Fatalf("want no error, got %s", got)
		}
	})

	t.Run("use VinRange()", func(t *testing.T) {
		got := api.AddListener(Listener{
			StatusFunc: func(vin int, online bool) {},
		}, VinRange(100, 50)...)
		if got != nil {
			t.Fatalf("want no error, got %s", got)
		}
	})

	t.Run("VinRange() output", func(t *testing.T) {
		want := []int{1, 2, 3}
		got := VinRange(3, 1)
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})
}
