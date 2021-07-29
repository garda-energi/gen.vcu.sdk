package sdk

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestCommandHandler(t *testing.T) {
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
			cmder := newStubCommander(testVin)
			defer cmder.Destroy()

			cmderStubClient(cmder).
				mockResponse(testVin, tC.invoker, func(rp *responsePacket) {
					if tC.resMsg != nil {
						rp.Message = tC.resMsg
					}
				})

			// call related method, pass in arg, evaluate outs
			resOut, errOut := callCommand(cmder, tC.invoker, tC.arg)

			// check output error
			if errOut != nil {
				if err := errOut.(error); err != nil {
					t.Error("want no error, got ", err)
				}
			}

			// check output response
			if tC.resMsg != nil && tC.wantOut == nil {
				tC.wantOut = string(tC.resMsg)
			}

			if resOut != nil {
				if !reflect.DeepEqual(resOut, tC.wantOut) {
					t.Errorf("want %s, got %s", tC.wantOut, resOut)
				}
			}
		})
	}
}
func TestCommandInvalidInputHandler(t *testing.T) {
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
			arg:     FrameLimit,
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
			wantErr: errInputOutOfRange(fmt.Sprint(ModeDriveStandard, ":dischare-current")).Error(),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 50, Torque: 10},                 // economy
				{DisCur: 50, Torque: 20},                 // standard
				{DisCur: 50, Torque: MCU_TORQUE_MAX + 1}, // sport
			},
			wantErr: errInputOutOfRange(fmt.Sprint(ModeDriveSport, ":torque")).Error(),
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
		testName := fmt.Sprint(tC.invoker, " for ", tC.wantErr)
		t.Run(testName, func(t *testing.T) {
			cmder := newStubCommander(testVin)
			defer cmder.Destroy()

			cmderStubClient(cmder).
				mockResponse(testVin, tC.invoker, nil)

			// call related method, pass in arg, evaluate outs
			_, errOut := callCommand(cmder, tC.invoker, tC.arg)

			// check output error
			if err := errOut.(error).Error(); err != tC.wantErr {
				t.Errorf("want %s, got %s", tC.wantErr, err)
			}
		})
	}
}
