package sdk

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

const testVin = 354313

func TestCommands(t *testing.T) {
	testCases := []struct {
		cmd     string
		cmdName string
		want    interface{}
		args    interface{}
		msg     message
	}{
		{
			cmd:     "GenInfo",
			cmdName: "GEN_INFO",
			args:    nil,
			msg:     message("VCU v.664, GEN - 2021"),
		},
		{
			cmd:     "GenLed",
			cmdName: "GEN_LED",
			args:    true,
		},
		{
			cmd:     "GenRtc",
			cmdName: "GEN_RTC",
			args:    time.Now(),
		},
		{
			cmd:     "GenOdo",
			cmdName: "GEN_ODO",
			args:    uint16(4321),
		},
		{
			cmd:     "GenAntiTheaf",
			cmdName: "GEN_ANTI_THIEF",
			args:    nil,
		},
		{
			cmd:     "GenReportFlush",
			cmdName: "GEN_RPT_FLUSH",
			args:    nil,
		},
		{
			cmd:     "GenReportBlock",
			cmdName: "GEN_RPT_BLOCK",
			args:    false,
		},
		{
			cmd:     "OvdState",
			cmdName: "OVD_STATE",
			args:    BikeStateNormal,
		},
		{
			cmd:     "OvdReportInterval",
			cmdName: "OVD_RPT_INTERVAL",
			args:    5 * time.Second,
		},
		{
			cmd:     "OvdReportFrame",
			cmdName: "OVD_RPT_FRAME",
			args:    FrameFull,
		},
		{
			cmd:     "OvdRemoteSeat",
			cmdName: "OVD_RMT_SEAT",
			args:    nil,
		},
		{
			cmd:     "OvdRemoteAlarm",
			cmdName: "OVD_RMT_ALARM",
			args:    nil,
		},
		{
			cmd:     "AudioBeep",
			cmdName: "AUDIO_BEEP",
			args:    nil,
		},
		{
			cmd:     "FingerFetch",
			cmdName: "FINGER_FETCH",
			args:    nil,
			want:    []int{1, 2, 3, 4, 5},
			msg:     message([]byte("12345")),
		},
		{
			cmd:     "FingerAdd",
			cmdName: "FINGER_ADD",
			args:    nil,
			want:    3,
			msg:     message([]byte(strconv.Itoa(3))),
		},
		{
			cmd:     "FingerDel",
			cmdName: "FINGER_DEL",
			args:    3,
		},
		{
			cmd:     "FingerRst",
			cmdName: "FINGER_RST",
			args:    nil,
		},
		{
			cmd:     "RemotePairing",
			cmdName: "REMOTE_PAIRING",
			args:    nil,
		},
		{
			cmd:     "FotaVcu",
			cmdName: "FOTA_VCU",
			args:    nil,
			msg:     message("VCU upgraded v.664 -> v.665"),
		},
		{
			cmd:     "FotaHmi",
			cmdName: "FOTA_HMI",
			args:    nil,
			msg:     message("HMI upgraded v.123 -> v.124"),
		},
		{
			cmd:     "NetSendUssd",
			cmdName: "NET_SEND_USSD",
			args:    "*123*10*3#",
			msg:     message("Terima kasih, permintaan kamu akan diproses,cek SMS untuk info lengkap. Dapatkan informasi seputar kartu Tri mu di aplikasi BimaTri, download di bima.tri.co.id"),
		},
		{
			cmd:     "NetReadSms",
			cmdName: "NET_READ_SMS",
			args:    nil,
			msg:     message("Poin Bonstri kamu: 20 Sisa Kuota kamu : Kuota ++ 372 MB s.d 03/01/2031 13:30:18 Temukan beragam paket lain di bima+ https://goo.gl/RQ1DBA"),
		},
		{
			cmd:     "HbarDrive",
			cmdName: "HBAR_DRIVE",
			args:    ModeDriveEconomy,
		},
		{
			cmd:     "HbarTrip",
			cmdName: "HBAR_TRIP",
			args:    ModeTripB,
		},
		{
			cmd:     "HbarAvg",
			cmdName: "HBAR_AVG",
			args:    ModeAvgEfficiency,
		},
		{
			cmd:     "HbarReverse",
			cmdName: "HBAR_REVERSE",
			args:    false,
		},
		{
			cmd:     "McuSpeedMax",
			cmdName: "MCU_SPEED_MAX",
			args:    uint8(90),
		},
		{
			cmd:     "McuTemplates",
			cmdName: "MCU_TEMPLATES",
			args: []McuTemplate{
				{DisCur: 50, Torque: 10}, // economy
				{DisCur: 50, Torque: 20}, // standard
				{DisCur: 50, Torque: 25}, // sport
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.cmd, func(t *testing.T) {
			fakeRes := newFakeResponse(testVin, tC.cmdName)
			if tC.msg != nil {
				fakeRes.Message = tC.msg
				if tC.want == nil {
					tC.want = string(tC.msg)
				}
			}

			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				mockResponse(fakeRes),
			})
			defer cmder.Destroy()

			// call related method, pass in args, evaluate outs
			meth := reflect.ValueOf(cmder).MethodByName(tC.cmd)
			args := []reflect.Value{}
			if tC.args != nil {
				args = append(args, reflect.ValueOf(tC.args))
			}
			outs := meth.Call(args)

			outRes := outs[0]
			outError := outs[len(outs)-1]

			// check output error
			if err := outError; !err.IsNil() {
				t.Fatalf("want no error, got %s\n", err)
			}

			// check output response
			if len(outs) == 1 {
				return
			}
			if !reflect.DeepEqual(outRes.Interface(), tC.want) {
				t.Fatalf("want %s, got %s", tC.want, outRes)
			}
		})
	}
}

func newFakeCommander(responses [][]byte) *commander {
	logging := false
	client := &fakeClient{
		responses: responses,
		cmdChan:   make(chan []byte),
		resChan:   make(chan struct{}),
	}

	cmder, _ := newCommander(testVin, client, &fakeSleeper{}, logging)
	return cmder
}
