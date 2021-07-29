package sdk

import (
	"strconv"
	"testing"
	"time"
)

func TestReportEventsAndFaults(t *testing.T) {
	type expectedData struct {
		vcuEvents       VcuEvents
		vcuEventsString string
		bmsFaults       BmsFaults
		bmsFaultsString string
		mcuFaults       McuFaults
		mcuFaultsString string

		isDataRealtime       bool
		isBateryLow          bool
		isEepromCapacityLow  bool
		isGpsValidHorizontal bool
		isGpsValidVertical   bool
		isNetLowSignal       bool
		isBmsCapacityLow     bool
	}
	type tester struct {
		name string
		args struct {
			b packet
		}
		want expectedData
	}
	testdata := make([]tester, 1)

	var resetDataTo = []struct {
		packet      packet
		want        expectedData
		dataChanger map[string]interface{}
	}{
		{
			packet: hexToByte(testDataNormal[3]),
			want: expectedData{
				// error or fault
				vcuEvents:       VcuEvents{VCU_BMS_ERROR, VCU_MCU_ERROR},
				vcuEventsString: "[BMS_ERROR, MCU_ERROR]",
				bmsFaults:       BmsFaults{BMS_SHORT_CIRCUIT, BMS_UNDER_VOLTAGE, BMS_UNBALANCE},
				bmsFaultsString: "[SHORT_CIRCUIT, UNDER_VOLTAGE, UNBALANCE]",
				mcuFaults:       McuFaults{Post: []McuFaultPost{MCU_POST_5V_LOW}, Run: []McuFaultRun{MCU_RUN_RESERVER_1}},
				mcuFaultsString: "Post[5V_LOW]\nRun[RESERVER_1]",

				// vcu
				isDataRealtime:       true,
				isBateryLow:          false,
				isEepromCapacityLow:  false,
				isGpsValidHorizontal: true,
				isGpsValidVertical:   true,
				isNetLowSignal:       true,
				isBmsCapacityLow:     false,
			},
			dataChanger: map[string]interface{}{
				"vcuEvents":     1<<VCU_BMS_ERROR | 1<<VCU_MCU_ERROR,
				"bmsFaults":     1<<BMS_SHORT_CIRCUIT | 1<<BMS_UNDER_VOLTAGE | 1<<BMS_UNBALANCE,
				"mcuFaultsPost": 1 << MCU_POST_5V_LOW,
				"mcuFaultsRun":  1 << MCU_RUN_RESERVER_1,

				"vcuLogBuffred":  0,
				"vcuLogDatetime": time.Now(),
				"vcuBatVoltage":  4572,
				"eepromUsed":     10,
				"netSignal":      9,
				"bmsSoc":         90,
			},
		},
	}

	for i := range testdata {
		// set data
		// 1. decode from hex than return report
		// 2. change report data
		// 3. encode
		// 3. save to testdata. it will be tested

		// 1
		report, err := decodeReport(resetDataTo[i].packet)
		if err != nil {
			t.Error("Create Dataset error: ", err)
		}

		// 2
		if vcuEvents, ok := resetDataTo[i].dataChanger["vcuEvents"]; ok {
			report.Vcu.Events = uint16(vcuEvents.(int))
		}
		if bmsFaults, ok := resetDataTo[i].dataChanger["bmsFaults"]; ok {
			report.Bms.Faults = uint16(bmsFaults.(int))
		}
		if mcuFaultsPost, ok := resetDataTo[i].dataChanger["mcuFaultsPost"]; ok {
			report.Mcu.Faults.Post = uint32(mcuFaultsPost.(int))
		}
		if mcuFaultsRun, ok := resetDataTo[i].dataChanger["mcuFaultsRun"]; ok {
			report.Mcu.Faults.Run = uint32(mcuFaultsRun.(int))
		}

		if vcuLogBuffred, ok := resetDataTo[i].dataChanger["vcuLogBuffred"]; ok {
			report.Vcu.LogBuffered = uint8(vcuLogBuffred.(int))
		}
		if vcuLogDatetime, ok := resetDataTo[i].dataChanger["vcuLogDatetime"]; ok {
			report.Vcu.LogDatetime = vcuLogDatetime.(time.Time)
		}
		if vcuBatVoltage, ok := resetDataTo[i].dataChanger["vcuBatVoltage"]; ok {
			report.Vcu.BatVoltage = float32(vcuBatVoltage.(int))
		}
		if eepromUsed, ok := resetDataTo[i].dataChanger["eepromUsed"]; ok {
			report.Eeprom.Used = uint8(eepromUsed.(int))
		}
		if netSignal, ok := resetDataTo[i].dataChanger["netSignal"]; ok && report.Net != nil {
			report.Net.Signal = uint8(netSignal.(int))
		}

		// 3
		b, err := encode(report)
		if err != nil {
			t.Errorf("Create Dataset error: %s", err)
		}

		testdata[i].name = "data #" + strconv.Itoa(i)
		testdata[i].args.b = b
		testdata[i].want = resetDataTo[i].want
	}

	for _, tt := range testdata {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeReport(tt.args.b)
			if err != nil {
				t.Error("Error Decode: ", err)
			} else {
				gotVcuEvent := got.Vcu.GetEvents()
				gotBmsFault := got.Bms.GetFaults()
				gotMcuFault := got.Mcu.GetFaults()

				if score := compareVar(gotVcuEvent, tt.want.vcuEvents); score != 100 {
					t.Errorf("[VCU Event] Got %s want %s", gotVcuEvent, tt.want.vcuEvents)
				}
				if score := compareVar(gotVcuEvent.String(), tt.want.vcuEventsString); score != 100 {
					t.Errorf("[VCU Event string] Got %s want %s", gotVcuEvent, tt.want.vcuEventsString)
				}
				if len(tt.want.vcuEvents) > 0 {
					isEvent := got.Vcu.IsEvent(tt.want.vcuEvents[0])
					if !isEvent {
						t.Errorf("[VCU Event IsEvent] Got %t want %t", isEvent, true)
					}
				}

				if score := compareVar(gotBmsFault, tt.want.bmsFaults); score != 100 {
					t.Errorf("[BMS Faults] Got %s want %s", gotBmsFault, tt.want.bmsFaults)
				}
				if score := compareVar(gotBmsFault.String(), tt.want.bmsFaultsString); score != 100 {
					t.Errorf("[BMS Faults string] Got %s want %s", gotBmsFault, tt.want.bmsFaultsString)
				}
				if len(tt.want.bmsFaults) > 0 {
					isFault := got.Bms.IsFault(tt.want.bmsFaults[0])
					if !isFault {
						t.Errorf("[BMS Faults isFault] Got %t want %t", isFault, true)
					}
				}

				if score := compareVar(gotMcuFault, tt.want.mcuFaults); score != 100 {
					t.Errorf("[MCU Faults] Got %s want %s", gotMcuFault, tt.want.mcuFaults)
				}
				if score := compareVar(gotMcuFault.String(), tt.want.mcuFaultsString); score != 100 {
					t.Errorf("[MCU Faults string] Got %s want %s", gotMcuFault, tt.want.mcuFaultsString)
				}
				if len(tt.want.mcuFaults.Post) > 0 {
					isFault := got.Mcu.IsFaultPost(tt.want.mcuFaults.Post[0])
					if !isFault {
						t.Errorf("[MCU Faults isFault] Got %t want %t", isFault, true)
					}
				}

				if score := compareVar(got.Vcu.RealtimeData(), tt.want.isDataRealtime); score != 100 {
					t.Errorf("[VCU RealtimeData] Got %t want %t", got.Vcu.RealtimeData(), tt.want.isDataRealtime)
				}
				if score := compareVar(got.Vcu.BatteryLow(), tt.want.isBateryLow); score != 100 {
					t.Errorf("[VCU BatteryLow] Got %t (%f) want %t", got.Vcu.BatteryLow(), got.Vcu.BatVoltage, tt.want.isBateryLow)
				}
				if score := compareVar(got.Eeprom.CapacityLow(), tt.want.isEepromCapacityLow); score != 100 {
					t.Errorf("[Eeprom CapacityLow] Got %t want %t", got.Eeprom.CapacityLow(), tt.want.isEepromCapacityLow)
				}
				if score := compareVar(got.Gps.ValidHorizontal(), tt.want.isGpsValidHorizontal); score != 100 {
					t.Errorf("[GPS ValidHorizontal] Got %t want %t", got.Gps.ValidHorizontal(), tt.want.isGpsValidHorizontal)
				}
				if score := compareVar(got.Gps.ValidVertical(), tt.want.isGpsValidVertical); score != 100 {
					t.Errorf("[GPS ValidVertical] Got %t want %t", got.Gps.ValidVertical(), tt.want.isGpsValidVertical)
				}
				if score := compareVar(got.Net.LowSignal(), tt.want.isNetLowSignal); score != 100 {
					t.Errorf("[Net LowSignal] Got %t want %t", got.Net.LowSignal(), tt.want.isNetLowSignal)
				}
				if score := compareVar(got.Bms.LowCapacity(), tt.want.isBmsCapacityLow); score != 100 {
					t.Errorf("[Bms LowCapacity] Got %t want %t", got.Bms.LowCapacity(), tt.want.isBmsCapacityLow)
				}
			}
		})
	}
}
