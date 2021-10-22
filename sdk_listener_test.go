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
				rp.Data["Report"].(PacketData)["SendDatetime"] = time.Now().UTC().Add(-24 * time.Hour)
			},
			validator: func(rp *ReportPacket) bool {
				datetime := time.Now().UTC().Add(-20 * time.Hour)
				return rp.Data["Report"].(PacketData)["SendDatetime"].(time.Time).Before(datetime)
			},
		},
		{
			desc:  "log is not queued",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Report"].(PacketData)["Queued"] = 0
			},
			validator: func(rp *ReportPacket) bool {
				return rp.RealtimeData()
			},
		},
		{
			desc:  "log is queued",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Report"].(PacketData)["LogDatetime"] = time.Now().UTC()
				rp.Data["Report"].(PacketData)["Queued"] = 5
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.RealtimeData()
			},
		},
		{
			desc:  "vcu's events has BMS_ERROR & BIKE_FALLEN",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Vcu"].(PacketData)["Events"] = 1<<VCU_BMS_ERROR | 1<<VCU_BIKE_FALLEN
			},
			validator: func(rp *ReportPacket) bool {
				want := VcuEvents{VCU_BIKE_FALLEN, VCU_BMS_ERROR}
				return rp.VcuIsEvents(want...) &&
					len(want) == len(rp.VcuGetEvents())
			},
		},
		{
			desc:  "vcu's events is empty",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Vcu"].(PacketData)["Events"] = 0
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.VcuGetEvents()) == 0
			},
		},
		{
			desc:  "vcu's events has invalid value",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Vcu"].(PacketData)["Events"] = 1 << VCU_EVENTS_MAX
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.VcuGetEvents()) == 0
			},
		},
		{
			desc:  "backup battery medium",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Vcu"].(PacketData)["BatVoltage"] = float32(BATTERY_BACKUP_FULL_MV - 300)
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.VcuBatteryLow()
			},
		},
		{
			desc:  "backup battery low",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Vcu"].(PacketData)["BatVoltage"] = float32(BATTERY_BACKUP_LOW_MV - 300)
			},
			validator: func(rp *ReportPacket) bool {
				return rp.VcuBatteryLow()
			},
		},
		{
			desc:  "eeprom capacity medium",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Eeprom"].(PacketData)["Used"] = 2
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.EepromCapacityLow()
			},
		},
		{
			desc:  "eeprom capacity low",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Eeprom"].(PacketData)["Used"] = 99
			},
			validator: func(rp *ReportPacket) bool {
				return rp.EepromCapacityLow()
			},
		},
		{
			desc:  "gps dop low",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Gps"].(PacketData)["HDOP"] = float32(GPS_DOP_MIN + 3)
				rp.Data["Gps"].(PacketData)["VDOP"] = float32(GPS_DOP_MIN + 18)
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.GpsValidHorizontal() && !rp.GpsValidVertical()
			},
		},
		{
			desc:  "gps valid vdop",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Gps"].(PacketData)["VDOP"] = 1.5
			},
			validator: func(rp *ReportPacket) bool {
				return rp.GpsValidVertical()
			},
		},
		{
			desc:  "gps valid hdop, invalid long lat",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Gps"].(PacketData)["HDOP"] = 2.5
				rp.Data["Gps"].(PacketData)["Longitude"] = GPS_LNG_MIN - 20
				rp.Data["Gps"].(PacketData)["Latitude"] = GPS_LAT_MAX - 20
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.GpsValidHorizontal()
			},
		},
		{
			desc:  "gps valid hdop, valid long lat",
			frame: FrameSimple,
			modifier: func(rp *ReportPacket) {
				rp.Data["Gps"].(PacketData)["HDOP"] = float32(GPS_DOP_MIN)
				rp.Data["Gps"].(PacketData)["Longitude"] = float32(GPS_LNG_MIN)
				rp.Data["Gps"].(PacketData)["Latitude"] = float32(GPS_LAT_MAX)
			},
			validator: func(rp *ReportPacket) bool {
				return rp.GpsValidHorizontal()
			},
		},
		{
			desc:  "net signal good",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Net"].(PacketData)["Signal"] = 75
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.NetLowSignal()
			},
		},
		{
			desc:  "net signal poor",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Net"].(PacketData)["Signal"] = NET_LOW_SIGNAL_PERCENT - 5
			},
			validator: func(rp *ReportPacket) bool {
				return rp.NetLowSignal()
			},
		},
		{
			desc:  "bms's faults has BMS_SHORT_CIRCUIT & BMS_UNBALANCE",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Bms"].(PacketData)["Faults"] = 1<<BMS_SHORT_CIRCUIT | 1<<BMS_UNBALANCE
			},
			validator: func(rp *ReportPacket) bool {
				want := BmsFaults{BMS_SHORT_CIRCUIT, BMS_UNBALANCE}
				return rp.BmsIsFaults(want...) &&
					len(want) == len(rp.BmsGetFaults())
			},
		},
		{
			desc:  "bms's faults is empty",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Bms"].(PacketData)["Faults"] = 0
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.BmsGetFaults()) == 0
			},
		},
		{
			desc:  "bms's faults has invalid value",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Bms"].(PacketData)["Faults"] = 1 << BMS_FAULTS_MAX
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.BmsGetFaults()) == 0
			},
		},
		{
			desc:  "bms soc full",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Bms"].(PacketData)["SOC"] = 100
			},
			validator: func(rp *ReportPacket) bool {
				return !rp.BmsLowCapacity()
			},
		},
		{
			desc:  "bms soc low",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Bms"].(PacketData)["SOC"] = 1
			},
			validator: func(rp *ReportPacket) bool {
				return rp.BmsLowCapacity()
			},
		},
		{
			desc:  "mcu's faults has MCU_POST_5V_LOW, MCU_POST_BRAKE_OPEN & MCU_RUN_ACCEL_OPEN",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Mcu"].(PacketData)["Faults"].(PacketData)["Post"] = uint32(1<<MCU_POST_5V_LOW | 1<<MCU_POST_BRAKE_OPEN)
				rp.Data["Mcu"].(PacketData)["Faults"].(PacketData)["Run"] = uint32(1 << MCU_RUN_ACCEL_OPEN)
			},
			validator: func(rp *ReportPacket) bool {
				wantPost := []McuFaultPost{MCU_POST_5V_LOW, MCU_POST_BRAKE_OPEN}
				wantRun := []McuFaultRun{MCU_RUN_ACCEL_OPEN}
				return rp.McuIsPostFaults(wantPost...) &&
					len(rp.McuGetFaults().Post) == len(wantPost) &&
					rp.McuIsRunFaults(wantRun...) &&
					len(rp.McuGetFaults().Run) == len(wantRun)

			},
		},
		{
			desc:  "mcu's faults is empty",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Mcu"].(PacketData)["Faults"].(PacketData)["Post"] = uint32(0)
				rp.Data["Mcu"].(PacketData)["Faults"].(PacketData)["Run"] = uint32(0)
			},
			validator: func(rp *ReportPacket) bool {
				return len(rp.McuGetFaults().Post) == 0 &&
					len(rp.McuGetFaults().Run) == 0
			},
		},
		{
			desc:  "some task's stack are near overflow",
			frame: FrameFull,
			modifier: func(rp *ReportPacket) {
				rp.Data["Task"].(PacketData)["Stack"].(PacketData)["Manager"] = uint16(STACK_OVERFLOW_BYTE_MIN - 20)
				rp.Data["Task"].(PacketData)["Stack"].(PacketData)["Command"] = uint16(0)
			},
			validator: func(rp *ReportPacket) bool {
				return rp.TaskStackOverflow()
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
			version := 1

			rp := makeReportPacket(version, vin, tC.frame)
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
