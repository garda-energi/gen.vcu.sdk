package sdk

import (
	"strconv"
	"testing"
)

func TestReportMethods(t *testing.T) {
	type expectedData struct {
		vcuEvents       VcuEvents
		vcuEventsString string
		bmsFaults       BmsFaults
		bmsFaultsString string
		mcuFaults       McuFaults
		mcuFaultsString string
	}
	type tester struct {
		name string
		args struct {
			b []byte
		}
		want expectedData
	}
	testdata := make([]tester, 1)

	var resetDataTo = []struct {
		want        expectedData
		dataChanger map[string]interface{}
	}{
		{
			want: expectedData{
				vcuEvents:       VcuEvents{VCU_BMS_ERROR, VCU_MCU_ERROR},
				vcuEventsString: "[BMS_ERROR, MCU_ERROR]",
				bmsFaults:       BmsFaults{BMS_SHORT_CIRCUIT, BMS_UNDER_VOLTAGE, BMS_UNBALANCE},
				bmsFaultsString: "[SHORT_CIRCUIT, UNDER_VOLTAGE, UNBALANCE]",
				mcuFaults:       McuFaults{Post: []McuFaultPost{MCU_POST_5V_LOW}, Run: []McuFaultRun{MCU_RUN_RESERVER_1}},
				mcuFaultsString: "Post[5V_LOW]\nRun[RESERVER_1]",
			},
			dataChanger: map[string]interface{}{
				"vcuEvents":     1<<VCU_BMS_ERROR | 1<<VCU_MCU_ERROR,
				"bmsFaults":     1<<BMS_SHORT_CIRCUIT | 1<<BMS_UNDER_VOLTAGE | 1<<BMS_UNBALANCE,
				"mcuFaultsPost": 1 << MCU_POST_5V_LOW,
				"mcuFaultsRun":  1 << MCU_RUN_RESERVER_1,
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
		report, err := decodeReport(hexToByte(testDataNormal[3]))
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
				if score := compareVar(gotVcuEvent.String(), tt.want.vcuEvents.String()); score != 100 {
					t.Errorf("[VCU Event string] Got %s want %s", gotVcuEvent, tt.want.vcuEvents)
				}
				if score := compareVar(gotBmsFault, tt.want.bmsFaults); score != 100 {
					t.Errorf("[BMS Faults] Got %s want %s", gotBmsFault, tt.want.bmsFaults)
				}
				if score := compareVar(gotBmsFault.String(), tt.want.bmsFaults.String()); score != 100 {
					t.Errorf("[BMS Faults string] Got %s want %s", gotBmsFault, tt.want.bmsFaults)
				}
				if score := compareVar(gotMcuFault, tt.want.mcuFaults); score != 100 {
					t.Errorf("[MCU Faults] Got %s want %s", gotMcuFault, tt.want.mcuFaults)
				}
				if score := compareVar(gotMcuFault.String(), tt.want.mcuFaults.String()); score != 100 {
					t.Errorf("[MCU Faults string] Got %s want %s", gotMcuFault, tt.want.mcuFaults)
				}
			}
		})
	}
}
