package sdk

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

const testVin = 354313

func TestResponsePacket(t *testing.T) {
	testCases := []struct {
		desc      string
		wantErr   error
		responses [][]byte
	}{
		{
			desc:      "no packet",
			wantErr:   errPacketTimeout("ack"),
			responses: nil,
		},
		{
			desc:      "invalid ack packet",
			wantErr:   errPacketAckCorrupt,
			responses: [][]byte{strToBytes(PREFIX_REPORT)},
		},
		{
			desc:      "only valid ack packet",
			wantErr:   errPacketTimeout("response"),
			responses: [][]byte{strToBytes(PREFIX_ACK)},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cmder := newFakeCommander(tC.responses)
			defer cmder.Destroy()

			_, err := cmder.GenInfo()

			if err != tC.wantErr {
				t.Fatalf("want %s, got %s", tC.wantErr, err)
			}
		})
	}
}

func TestResponseError(t *testing.T) {
	testCases := []struct {
		desc      string
		wantErr   error
		formatter func(r *responsePacket)
	}{
		{
			desc:    "invalid prefix",
			wantErr: errInvalidPrefix,
			formatter: func(r *responsePacket) {
				r.Header.Prefix = PREFIX_REPORT
			},
		},
		{
			desc:    "invalid size",
			wantErr: errInvalidSize,
			formatter: func(r *responsePacket) {
				r.Header.Size = 55
			},
		},
		{
			desc:    "invalid VIN",
			wantErr: errInvalidVin,
			formatter: func(r *responsePacket) {
				r.Header.Vin = 12345
			},
		},
		{
			desc:    "invalid cmd code",
			wantErr: errInvalidCmdCode,
			formatter: func(r *responsePacket) {
				r.Header.Code = 9
				r.Header.SubCode = 5
			},
		},
		{
			desc:    "invalid resCode",
			wantErr: errInvalidResCode,
			formatter: func(r *responsePacket) {
				r.Header.ResCode = 99
			},
		},
		{
			desc:    "message overflowed",
			wantErr: errInvalidSize,
			formatter: func(r *responsePacket) {
				r.Message = message("'Google Go' redirects here. For the Android search app by Google, 'Google Go', for low-end Lollipop+ devices, see Android Go. For the computer program by Google to play the board game Go, see AlphaGo. For the 2003 agent-based programming language, see Go! (programming language).")
				r.Header.Size = uint8(len(r.Message))
			},
		},
		{
			desc:    "simulate code error",
			wantErr: errors.New(resCodeError.String()),
			formatter: func(r *responsePacket) {
				r.Header.ResCode = resCodeError
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res := newFakeResponse(testVin, "GEN_INFO")
			tC.formatter(res)

			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				mockResponse(res),
			})
			defer cmder.Destroy()

			_, err := cmder.GenInfo()
			if err.Error() != tC.wantErr.Error() {
				t.Fatalf("want %s, got %s", tC.wantErr, err)
			}
		})
	}
}

func TestCommands(t *testing.T) {
	testCases := []struct {
		cmd  string
		want interface{}
		args interface{}
		msg  message
	}{
		{
			cmd:  "GEN_INFO",
			args: nil,
			msg:  message("VCU v.664, GEN - 2021"),
		},
		{
			cmd:  "GEN_LED",
			args: true,
		},
		{
			cmd:  "GEN_RTC",
			args: time.Now(),
		},
		{
			cmd:  "GEN_ODO",
			args: uint16(4321),
		},
		{
			cmd:  "GEN_ANTI_THIEF",
			args: nil,
		},
		{
			cmd:  "GEN_RPT_FLUSH",
			args: nil,
		},
		{
			cmd:  "GEN_RPT_BLOCK",
			args: false,
		},
		{
			cmd:  "OVD_STATE",
			args: BikeStateNormal,
		},
		{
			cmd:  "OVD_RPT_INTERVAL",
			args: 5 * time.Second,
		},
		{
			cmd:  "OVD_RPT_FRAME",
			args: FrameFull,
		},
		{
			cmd:  "OVD_RMT_SEAT",
			args: nil,
		},
		{
			cmd:  "OVD_RMT_ALARM",
			args: nil,
		},
		{
			cmd:  "AUDIO_BEEP",
			args: nil,
		},
		{
			cmd:  "FINGER_FETCH",
			args: nil,
			want: []int{1, 2, 3, 4, 5},
			msg:  message([]byte("12345")),
		},
		{
			cmd:  "FINGER_ADD",
			args: nil,
			want: 3,
			msg:  message([]byte("3")),
		},
		{
			cmd:  "FINGER_DEL",
			args: 3,
		},
		{
			cmd:  "FINGER_RST",
			args: nil,
		},
		{
			cmd:  "REMOTE_PAIRING",
			args: nil,
		},
		{
			cmd:  "FOTA_VCU",
			args: nil,
			msg:  message("VCU upgraded v.664 -> v.665"),
		},
		{
			cmd:  "FOTA_HMI",
			args: nil,
			msg:  message("HMI upgraded v.123 -> v.124"),
		},
		{
			cmd:  "NET_SEND_USSD",
			args: "*123*10*3#",
			msg:  message("Terima kasih, permintaan kamu akan diproses,cek SMS untuk info lengkap. Dapatkan informasi seputar kartu Tri mu di aplikasi BimaTri, download di bima.tri.co.id"),
		},
		{
			cmd:  "NET_READ_SMS",
			args: nil,
			msg:  message("Poin Bonstri kamu: 20 Sisa Kuota kamu : Kuota ++ 372 MB s.d 03/01/2031 13:30:18 Temukan beragam paket lain di bima+ https://goo.gl/RQ1DBA"),
		},
		{
			cmd:  "HBAR_DRIVE",
			args: ModeDriveEconomy,
		},
		{
			cmd:  "HBAR_TRIP",
			args: ModeTripB,
		},
		{
			cmd:  "HBAR_AVG",
			args: ModeAvgEfficiency,
		},
		{
			cmd:  "HBAR_REVERSE",
			args: false,
		},
		{
			cmd:  "MCU_SPEED_MAX",
			args: uint8(90),
		},
		{
			cmd: "MCU_TEMPLATES",
			args: []McuTemplate{
				{DisCur: 50, Torque: 10}, // economy
				{DisCur: 50, Torque: 20}, // standard
				{DisCur: 50, Torque: 25}, // sport
			},
		},
	}

	for _, tC := range testCases {
		cmd, _ := getCmdByName(tC.cmd)
		t.Run(cmd.invoker, func(t *testing.T) {
			fakeRes := newFakeResponse(testVin, tC.cmd)
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
			meth := reflect.ValueOf(cmder).MethodByName(cmd.invoker)
			args := []reflect.Value{}
			if tC.args != nil {
				args = append(args, reflect.ValueOf(tC.args))
			}
			outs := meth.Call(args)

			// check output error
			outError := outs[len(outs)-1]
			if !outError.IsNil() {
				t.Fatalf("want no error, got %s\n", outError)
			}

			// check output response
			if len(outs) < 2 {
				return
			}

			outRes := outs[0]
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
	sleeper := &fakeSleeper{
		sleep: time.Millisecond,
		after: 125 * time.Millisecond,
	}
	cmder, _ := newCommander(testVin, client, sleeper, logging)
	return cmder
}
