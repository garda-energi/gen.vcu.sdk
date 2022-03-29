package sdk

import (
	"errors"
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
			invoker: "GenBikeState",
			arg:     BikeStateNormal,
		},
		{
			invoker: "GenLockDown",
			arg:     false,
		},
		{
			invoker: "ReportFlush",
		},
		{
			invoker: "ReportBlock",
			arg:     false,
		},
		{
			invoker: "ReportInterval",
			arg:     5 * time.Second,
		},
		{
			invoker: "ReportFrame",
			arg:     FrameFull,
		},
		{
			invoker: "RemoteSeat",
		},
		{
			invoker: "RemoteAlarm",
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
			invoker: "FotaRestart",
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
		// {
		// 	invoker: "HbarTripMeter",
		// 	arg:     uint16(4321),
		// },
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
			invoker: "McuSpeedMax",
			arg:     []uint8{90, 1},
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 50, Torque: 10}, // economy
				{DisCur: 50, Torque: 20}, // standard
				{DisCur: 50, Torque: 25}, // sport
			},
		},
		{
			invoker: "ImuAntiThief",
			arg:     true,
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
			resOut, errOut := cmder.invoke(tC.invoker, tC.arg)

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
		want    error
	}{
		{
			invoker: "GenBikeState",
			arg:     BikeStateLimit,
			want:    errInputOutOfRange("state"),
		},
		{
			invoker: "ReportInterval",
			arg:     100 * time.Hour,
			want:    errInputOutOfRange("duration"),
		},
		{
			invoker: "ReportFrame",
			arg:     FrameLimit,
			want:    errInputOutOfRange("frame"),
		},
		{
			invoker: "FingerDel",
			arg:     9,
			want:    errInputOutOfRange("id"),
		},
		{
			invoker: "NetSendUssd",
			arg:     "*123*123123*324234*423423424*4324342424234#",
			want:    errInputOutOfRange("ussd"),
		},
		{
			invoker: "NetSendUssd",
			arg:     "*#",
			want:    errInputOutOfRange("ussd"),
		},
		{
			invoker: "NetSendUssd",
			arg:     "*123*1*3*",
			want:    errors.New("invalid ussd format"),
		},
		{
			invoker: "HbarDrive",
			arg:     ModeDriveLimit,
			want:    errInputOutOfRange("drive-mode"),
		},
		{
			invoker: "HbarTrip",
			arg:     ModeTripLimit,
			want:    errInputOutOfRange("trip-mode"),
		},
		{
			invoker: "HbarAvg",
			arg:     ModeAvgLimit,
			want:    errInputOutOfRange("avg-mode"),
		},
		{
			invoker: "McuSpeedMax",
			arg:     []uint8{245, 1},
			want:    errInputOutOfRange("speed-max"),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 50, Torque: 10},                 // economy
				{DisCur: MCU_DISCUR_MIN - 1, Torque: 20}, // standard
				{DisCur: 50, Torque: 25},                 // sport
			},
			want: errInputOutOfRange(fmt.Sprint(ModeDriveStandard, ":dischare-current")),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 50, Torque: 10},                 // economy
				{DisCur: 50, Torque: 20},                 // standard
				{DisCur: 50, Torque: MCU_TORQUE_MAX + 1}, // sport
			},
			want: errInputOutOfRange(fmt.Sprint(ModeDriveSport, ":torque")),
		},
		{
			invoker: "McuTemplates",
			arg: []McuTemplate{
				{DisCur: 1, Torque: 10},  // economy
				{DisCur: 50, Torque: 20}, // standard
			},
			want: errors.New("templates should be set for all driving modes at once"),
		},
	}

	for _, tC := range testCases {
		testName := fmt.Sprint(tC.invoker, " for ", tC.want)
		t.Run(testName, func(t *testing.T) {
			cmder := newStubCommander(testVin)
			defer cmder.Destroy()

			cmderStubClient(cmder).
				mockResponse(testVin, tC.invoker, nil)

			// call related method, pass in arg, evaluate outs
			_, errOut := cmder.invoke(tC.invoker, tC.arg)

			// check output error
			if errOut == nil {
				t.Fatalf("want %s, got none", tC.want)
			}

			if err := errOut.(error); err.Error() != tC.want.Error() {
				t.Errorf("want %s, got %s", tC.want, err)
			}
		})
	}
}
