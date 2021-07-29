package sdk

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestSdkReportMethods(t *testing.T) {
	type stream struct {
		vin    int
		report *ReportPacket
	}

	reportChan := make(chan *stream)
	defer close(reportChan)

	/////////////////////// SAND BOX ////////////////////////
	listener := Listener{
		ReportFunc: func(vin int, report *ReportPacket) {
			reportChan <- &stream{
				vin:    vin,
				report: report,
			}
		},
	}

	api := newStubApi()
	api.Connect()
	defer api.Disconnect()

	vins := VinRange(5, 10)
	if err := api.AddListener(listener, vins...); err != nil {
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
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vin := vins[rand.Intn(len(vins))]

			rp := makeReportPacket(vin)
			tC.modifier(rp)

			sdkStubClient(api).
				mockReport(vin, []*ReportPacket{rp})

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
