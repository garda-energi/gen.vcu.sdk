package sdk

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestSdkReportListener(t *testing.T) {
	type stream struct {
		vin    int
		report *ReportPacket
	}

	reportChan := make(chan *stream)
	defer close(reportChan)

	/////////////////////// SAND BOX ////////////////////////
	api := newStubApi()
	api.Connect()
	defer api.Disconnect()

	vins := VinRange(5, 10)
	if err := api.AddListener(Listener{
		ReportFunc: func(vin int, report *ReportPacket) {
			reportChan <- &stream{
				vin:    vin,
				report: report,
			}
		},
	}, vins...); err != nil {
		t.Error("want no error, got ", err)
	}
	defer api.RemoveListener(vins...)
	//////////////////////////////////////////////////////////

	testCases := []struct {
		desc      string
		modifier  func(rp *ReportPacket)
		validator func(rp *ReportPacket) bool
	}{
		{
			desc: "send datetime is yesterday",
			modifier: func(rp *ReportPacket) {
				rp.Header.SendDatetime = time.Now().UTC().Add(-24 * time.Hour)
			},
			validator: func(rp *ReportPacket) bool {
				datetime := time.Now().UTC().Add(-20 * time.Hour)
				return rp.Header.SendDatetime.Before(datetime)
			},
		},
		{
			desc: "log datetime is yesterday",
			modifier: func(rp *ReportPacket) {
				rp.Vcu.LogDatetime = time.Now().UTC().Add(-24 * time.Hour)
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Vcu.RealtimeData()
			},
		},
		{
			desc: "log datetime is now, no buffered",
			modifier: func(rp *ReportPacket) {
				rp.Vcu.LogBuffered = 0
				rp.Vcu.LogDatetime = time.Now().UTC()
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Vcu.RealtimeData()
			},
		},
		{
			desc: "log is buffered",
			modifier: func(rp *ReportPacket) {
				rp.Vcu.LogBuffered = 5
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Vcu.RealtimeData()
			},
		},
		{
			desc: "vcu has BMS_ERROR & BIKE_FALLEN events",
			modifier: func(rp *ReportPacket) {
				rp.Vcu.Events = 1<<VCU_BMS_ERROR | 1<<VCU_BIKE_FALLEN
			},
			validator: func(rp *ReportPacket) bool {
				want := VcuEvents{VCU_BIKE_FALLEN, VCU_BMS_ERROR}
				return reflect.DeepEqual(want, rp.Vcu.GetEvents())
			},
		},
		// TODO: test all report's methods
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vin := vins[rand.Intn(len(vins))]

			rp := makeReportPacket(vin)
			tC.modifier(rp)

			sdkStubClient(api).
				mockReports(vin, []*ReportPacket{rp})

			data := <-reportChan

			if vin != data.vin {
				t.Errorf("vin want %d, got %d", vin, data.vin)
			}
			if !tC.validator(data.report) {
				t.Errorf("report received is invalid")
			}
		})
	}
}

func TestSdkStatusListener(t *testing.T) {
	type stream struct {
		vin    int
		online bool
	}

	statusChan := make(chan *stream)
	defer close(statusChan)

	/////////////////////// SAND BOX ////////////////////////
	api := newStubApi()
	api.Connect()
	defer api.Disconnect()

	vins := VinRange(5, 10)
	if err := api.AddListener(Listener{
		StatusFunc: func(vin int, online bool) {
			statusChan <- &stream{
				vin:    vin,
				online: online,
			}
		},
	}, vins...); err != nil {
		t.Error("want no error, got ", err)
	}
	defer api.RemoveListener(vins...)
	//////////////////////////////////////////////////////////

	testCases := []struct {
		desc   string
		packet packet
		online bool
		vin    int
	}{
		{
			desc:   "online status packet",
			packet: packet("1"),
			online: true,
			vin:    5,
		},
		{
			desc:   "offline status packet",
			packet: packet("0"),
			online: false,
			vin:    8,
		},
		{
			desc:   "unknown status packet",
			packet: packet("XXX"),
			online: false,
			vin:    10,
		},
	}

	for _, tC := range testCases {
		testName := fmt.Sprintf("vin %d, %s", tC.vin, tC.desc)
		t.Run(testName, func(t *testing.T) {
			sdkStubClient(api).
				mockStatus(tC.vin, tC.packet)

			data := <-statusChan

			if tC.vin != data.vin {
				t.Errorf("vin want %d, got %d", tC.vin, data.vin)
			}
			if tC.online != data.online {
				t.Errorf("status received is invalid")
			}
		})
	}
}
