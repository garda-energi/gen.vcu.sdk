package sdk

import (
	"strconv"
	"testing"
)

type testDataChanger struct {
	byteIdx int    // byte index
	newByte []byte // byte which replace to
}

// change byte for test
func (dc testDataChanger) changeByte(b []byte) {
	i := dc.byteIdx
	for _, v := range dc.newByte {
		b[i] = v
		i++
	}
}

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
		dataChanger []testDataChanger
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
			dataChanger: []testDataChanger{
				// dont use byteIdx, because the report struct field is dynamic (not stable)
				// how about this: reportStruct (modified) -> encode (hexstring) -> decode reportStruct (validate)
				{byteIdx: 23, newByte: []byte{0x60, 0x00}},              // change vcu event data (byte index 23)
				{byteIdx: 101, newByte: []byte{0x84, 0x04}},             // change bms fault data (byte index 101)
				{byteIdx: 138, newByte: []byte{0x00, 0x10, 0x00, 0x00}}, // change mcu post fault data (byte index 138)
				{byteIdx: 142, newByte: []byte{0x00, 0x20, 0x00, 0x00}}, // change mcu run fault data (byte index 142)
			},
		},
	}

	for i := range testdata {
		testdata[i].name = "data #" + strconv.Itoa(i)
		testdata[i].args.b = hexToByte(testDataNormal[3])
		testdata[i].want = resetDataTo[i].want

		for _, dc := range resetDataTo[i].dataChanger {
			dc.changeByte(testdata[i].args.b)
		}
	}

	for _, tt := range testdata {
		if tt.name == "" {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeReport(tt.args.b)
			if err != nil {
				t.Errorf("Error Decode: %s", err)
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
