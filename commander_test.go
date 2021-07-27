package sdk

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func callCmd(cmder *commander, invoker string, arg interface{}) (res, err interface{}) {
	method := reflect.ValueOf(cmder).MethodByName(invoker)
	ins := []reflect.Value{}
	if arg != nil {
		ins = append(ins, reflect.ValueOf(arg))
	}
	outs := method.Call(ins)

	err = outs[len(outs)-1].Interface()
	if len(outs) > 1 {
		res = outs[0].Interface()
	}
	return
}

func TestCommander(t *testing.T) {
	testCases := []struct {
		invoker string
		arg     interface{}
		resMsg  message
		wantOut interface{}
	}{
		{
			invoker: "GenInfo",
			resMsg:  message("VCU v.664, GEN - 2021"),
		},
		{
			invoker: "GenLed",
			arg:     true,
		},
		{
			invoker: "GenRtc",
			arg:     time.Now(),
		},
		{
			invoker: "GenOdo",
			arg:     uint16(4321),
		},
		{
			invoker: "GenAntiThief",
		},
		{
			invoker: "GenReportFlush",
		},
		{
			invoker: "GenReportBlock",
			arg:     false,
		},
		{
			invoker: "OvdState",
			arg:     BikeStateNormal,
		},
		{
			invoker: "OvdReportInterval",
			arg:     5 * time.Second,
		},
		{
			invoker: "OvdReportFrame",
			arg:     FrameFull,
		},
		{
			invoker: "OvdRemoteSeat",
		},
		{
			invoker: "OvdRemoteAlarm",
		},
		{
			invoker: "AudioBeep",
		},
		{
			invoker: "FingerFetch",
			resMsg:  message([]byte("12345")),
			wantOut: []int{1, 2, 3, 4, 5},
		},
		{
			invoker: "FingerAdd",
			resMsg:  message([]byte("3")),
			wantOut: 3,
		},
		{
			invoker: "FingerDel",
			arg:     3,
		},
		{
			invoker: "FingerRst",
		},
		{
			invoker: "RemotePairing",
		},
		{
			invoker: "FotaVcu",
			resMsg:  message("VCU upgraded v.664 -> v.665"),
		},
		{
			invoker: "FotaHmi",
			resMsg:  message("HMI upgraded v.123 -> v.124"),
		},
		{
			invoker: "NetSendUssd",
			arg:     "*123*10*3#",
			resMsg:  message("Terima kasih, permintaan kamu akan diproses,cek SMS untuk info lengkap. Dapatkan informasi seputar kartu Tri mu di aplikasi BimaTri, download di bima.tri.co.id"),
		},
		{
			invoker: "NetReadSms",
			resMsg:  message("Poin Bonstri kamu: 20 Sisa Kuota kamu : Kuota ++ 372 MB s.d 03/01/2031 13:30:18 Temukan beragam paket lain di bima+ https://goo.gl/RQ1DBA"),
		},
		{
			invoker: "HbarDrive",
			arg:     ModeDriveEconomy,
		},
		{
			invoker: "HbarTrip",
			arg:     ModeTripB,
		},
		{
			invoker: "HbarAvg",
			arg:     ModeAvgEfficiency,
		},
		{
			invoker: "HbarReverse",
			arg:     false,
		},
		{
			invoker: "McuSpeedMax",
			arg:     uint8(90),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 50, Torque: 10}, // economy
				{DisCur: 50, Torque: 20}, // standard
				{DisCur: 50, Torque: 25}, // sport
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.invoker, func(t *testing.T) {
			// generate fake response
			fakeRes := newFakeResponse(testVin, tC.invoker)
			if tC.resMsg != nil {
				fakeRes.Message = tC.resMsg
				if tC.wantOut == nil {
					tC.wantOut = string(tC.resMsg)
				}
			}

			// initialize fake commander
			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				mockResponse(fakeRes),
			})
			defer cmder.Destroy()

			// call related method, pass in arg, evaluate outs
			resOut, errOut := callCmd(cmder, tC.invoker, tC.arg)

			// check output error
			if errOut != nil {
				if err := errOut.(error); err != nil {
					t.Fatalf("want no error, got %s", err)
				}
			}

			// check output response
			if resOut != nil {
				if !reflect.DeepEqual(resOut, tC.wantOut) {
					t.Fatalf("want %s, got %s", tC.wantOut, resOut)
				}
			}
		})
	}
}
func TestCommanderInvalidInput(t *testing.T) {
	testCases := []struct {
		invoker string
		arg     interface{}
		wantErr string
	}{
		{
			invoker: "OvdState",
			arg:     BikeStateLimit,
			wantErr: errInputOutOfRange("state").Error(),
		},
		{
			invoker: "OvdReportInterval",
			arg:     100 * time.Hour,
			wantErr: errInputOutOfRange("duration").Error(),
		},
		{
			invoker: "OvdReportFrame",
			arg:     FrameInvalid,
			wantErr: errInputOutOfRange("frame").Error(),
		},
		{
			invoker: "FingerDel",
			arg:     9,
			wantErr: errInputOutOfRange("id").Error(),
		},
		{
			invoker: "NetSendUssd",
			arg:     "*123*123123*324234*423423424*4324342424234#",
			wantErr: errInputOutOfRange("ussd").Error(),
		},
		{
			invoker: "NetSendUssd",
			arg:     "*#",
			wantErr: errInputOutOfRange("ussd").Error(),
		},
		{
			invoker: "NetSendUssd",
			arg:     "*123*1*3*",
			wantErr: "invalid ussd format",
		},
		{
			invoker: "HbarDrive",
			arg:     ModeDriveLimit,
			wantErr: errInputOutOfRange("drive-mode").Error(),
		},
		{
			invoker: "HbarTrip",
			arg:     ModeTripLimit,
			wantErr: errInputOutOfRange("trip-mode").Error(),
		},
		{
			invoker: "HbarAvg",
			arg:     ModeAvgLimit,
			wantErr: errInputOutOfRange("avg-mode").Error(),
		},
		{
			invoker: "McuSpeedMax",
			arg:     uint8(245),
			wantErr: errInputOutOfRange("speed-max").Error(),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 50, Torque: 10},                 // economy
				{DisCur: MCU_DISCUR_MIN - 1, Torque: 20}, // standard
				{DisCur: 50, Torque: 25},                 // sport
			},
			wantErr: errInputOutOfRange(fmt.Sprintf("%s:dischare-current", ModeDriveStandard)).Error(),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 50, Torque: 10},                 // economy
				{DisCur: 50, Torque: 20},                 // standard
				{DisCur: 50, Torque: MCU_TORQUE_MAX + 1}, // sport
			},
			wantErr: errInputOutOfRange(fmt.Sprintf("%s:torque", ModeDriveSport)).Error(),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 0, Torque: 10},  // economy
				{DisCur: 50, Torque: 20}, // standard
			},
			wantErr: "templates should be set for all driving mode at once",
		},
	}

	for _, tC := range testCases {
		testName := fmt.Sprintf("%s, error %s", tC.invoker, tC.wantErr)
		t.Run(testName, func(t *testing.T) {
			// initialize fake commander
			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				mockResponse(newFakeResponse(testVin, tC.invoker)),
			})
			defer cmder.Destroy()

			// call related method, pass in arg, evaluate outs
			_, errOut := callCmd(cmder, tC.invoker, tC.arg)

			// check output error
			if err := errOut.(error).Error(); err != tC.wantErr {
				t.Fatalf("want %s, got %s", tC.wantErr, err)
			}
		})
	}
}
