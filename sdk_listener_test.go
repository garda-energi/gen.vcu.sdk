package sdk

import (
	"fmt"
	"math/rand"
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
	vins := VinRange(5, 10)
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

	if err := api.AddListener(listener, vins...); err != nil {
		t.Error("want no error, got ", err)
	}
	defer api.RemoveListener(vins...)

	//////////////////////////////////////////////////////////

	testCases := []struct {
		desc      string
		frame     Frame
		modifier  func(rp *ReportPacket)
		validator func(rp *ReportPacket) bool
	}{
		{
			desc:  "send datetime is yesterday",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Header.SendDatetime = time.Now().UTC().Add(-24 * time.Hour)
			},
			validator: func(rp *ReportPacket) bool {
				datetime := time.Now().UTC().Add(-20 * time.Hour)
				return rp.Header.SendDatetime.Before(datetime)
			},
		},
		{
			desc:  "vcu's events has BMS_ERROR & BIKE_FALLEN",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.Events = 1<<VCU_BMS_ERROR | 1<<VCU_BIKE_FALLEN
			},
			validator: func(rp *ReportPacket) bool {
				want := VcuEvents{VCU_BIKE_FALLEN, VCU_BMS_ERROR}
				return rp.Vcu.IsEvents(want...) &&
					len(want) == len(rp.Vcu.GetEvents())
			},
		},
		{
			desc:  "vcu's events is empty",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.Events = 0
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.Vcu.GetEvents()) == 0
			},
		},
		{
			desc:  "vcu's events has invalid value",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.Events = 1 << VCU_EVENTS_MAX
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.Vcu.GetEvents()) == 0
			},
		},
		{
			desc:  "log datetime is yesterday",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.LogDatetime = time.Now().UTC().Add(-24 * time.Hour)
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Vcu.RealtimeData()
			},
		},
		{
			desc:  "log datetime is now, no buffered",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.LogBuffered = 0
				rp.Vcu.LogDatetime = time.Now().UTC()
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Vcu.RealtimeData()
			},
		},
		{
			desc:  "log is buffered",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.LogBuffered = 5
				rp.Vcu.LogDatetime = time.Now().UTC()
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Vcu.RealtimeData()
			},
		},
		{
			desc:  "backup battery medium",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.BatVoltage = BATTERY_BACKUP_FULL_MV - 300
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Vcu.BatteryLow()
			},
		},
		{
			desc:  "backup battery low",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Vcu.BatVoltage = BATTERY_BACKUP_LOW_MV - 300
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Vcu.BatteryLow()
			},
		},
		{
			desc:  "eeprom capacity medium",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Eeprom.Used = 2
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Eeprom.CapacityLow()
			},
		},
		{
			desc:  "eeprom capacity low",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Eeprom.Used = 99
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Eeprom.CapacityLow()
			},
		},
		{
			desc:  "gps dop low",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Gps.HDOP = GPS_DOP_MIN + 3
				rp.Gps.VDOP = GPS_DOP_MIN + 18
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Gps.ValidHorizontal() && !rp.Gps.ValidVertical()
			},
		},
		{
			desc:  "gps valid vdop",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Gps.VDOP = 1.5
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Gps.ValidVertical()
			},
		},
		{
			desc:  "gps valid hdop, invalid long lat",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Gps.HDOP = 2.5
				rp.Gps.Longitude = GPS_LNG_MIN - 20
				rp.Gps.Latitude = GPS_LAT_MAX + 15
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Gps.ValidHorizontal()
			},
		},
		{
			desc:  "gps valid hdop, valid long lat",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Gps.HDOP = GPS_DOP_MIN
				rp.Gps.Longitude = GPS_LNG_MIN
				rp.Gps.Latitude = GPS_LAT_MAX
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Gps.ValidHorizontal()
			},
		},
		{
			desc:  "net signal good",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Net.Signal = 75
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Net.LowSignal()
			},
		},
		{
			desc:  "net signal poor",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Net.Signal = NET_LOW_SIGNAL_PERCENT - 5
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Net.LowSignal()
			},
		},
		{
			desc:  "bms's faults has BMS_SHORT_CIRCUIT & BMS_UNBALANCE",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Bms.Faults = 1<<BMS_SHORT_CIRCUIT | 1<<BMS_UNBALANCE
			},
			validator: func(rp *ReportPacket) bool {
				want := BmsFaults{BMS_SHORT_CIRCUIT, BMS_UNBALANCE}
				return rp.Bms.IsFaults(want...) &&
					len(want) == len(rp.Bms.GetFaults())
			},
		},
		{
			desc:  "bms's faults is empty",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Bms.Faults = 0
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.Bms.GetFaults()) == 0
			},
		},
		{
			desc:  "bms's faults has invalid value",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Bms.Faults = 1 << BMS_FAULTS_MAX
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.Bms.GetFaults()) == 0
			},
		},
		{
			desc:  "bms soc full",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Bms.SOC = 100
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.Bms.LowCapacity()
			},
		},
		{
			desc:  "bms soc low",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Bms.SOC = 1
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Bms.LowCapacity()
			},
		},
		{
			desc:  "mcu's faults has MCU_POST_5V_LOW, MCU_POST_BRAKE_OPEN & MCU_RUN_ACCEL_OPEN",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Mcu.Faults.Post = 1<<MCU_POST_5V_LOW | 1<<MCU_POST_BRAKE_OPEN
				rp.Mcu.Faults.Run = 1 << MCU_RUN_ACCEL_OPEN
			},
			validator: func(rp *ReportPacket) bool {
				wantPost := []McuFaultPost{MCU_POST_5V_LOW, MCU_POST_BRAKE_OPEN}
				wantRun := []McuFaultRun{MCU_RUN_ACCEL_OPEN}
				return rp.Mcu.IsPostFaults(wantPost...) &&
					len(rp.Mcu.GetFaults().Post) == len(wantPost) &&
					rp.Mcu.IsRunFaults(wantRun...) &&
					len(rp.Mcu.GetFaults().Run) == len(wantRun)

			},
		},
		{
			desc:  "mcu's faults is empty",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Mcu.Faults.Post = 0
				rp.Mcu.Faults.Run = 0
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.Mcu.GetFaults().Post) == 0 &&
					len(rp.Mcu.GetFaults().Run) == 0
			},
		},
		{
			desc:  "some task's stack are near overflow",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Task.Stack.Manager = STACK_OVERFLOW_BYTE_MIN - 20
				rp.Task.Stack.Command = 0
			},
			validator: func(rp *ReportPacket) bool {
				return rp.Task.StackOverflow()
			},
		},
		// {
		// 	desc: "all task's stack are good",
		//	frame: FrameFull,
		// 	modifier: func(rp *ReportPacket) {
		// 		rp.Task.Stack.Manager = STACK_OVERFLOW_BYTE_MIN + uint16(rand.Intn(10))
		// 	},
		// 	validator: func(rp *ReportPacket) bool {
		// 		return rp.Task.StackOverflow()
		// 	},
		// },
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vin := vins[rand.Intn(len(vins))]

			rp := makeReportPacket(vin, tC.frame)
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
	vins := VinRange(5, 10)
	listener := Listener{
		StatusFunc: func(vin int, online bool) {
			statusChan <- &stream{
				vin:    vin,
				online: online,
			}
		},
	}

	api := newStubApi()
	api.Connect()
	defer api.Disconnect()

	if err := api.AddListener(listener, vins...); err != nil {
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
