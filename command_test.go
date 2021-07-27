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
				r.Message = message("Golang is very useful for writing light-weight microservices. We currently use it for generating APIs that interact with our front-end applications. If you want to build a small functional microservice quickly, then Golang is a great tool to use. It's an easy language for developers to learn quickly.")
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
			res := newFakeResponse(testVin, "GenInfo")
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
		invoker string
		want    interface{}
		args    interface{}
		resMsg  message
	}{
		{
			invoker: "GenInfo",
			args:    nil,
			resMsg:  message("VCU v.664, GEN - 2021"),
		},
		{
			invoker: "GenLed",
			args:    true,
		},
		{
			invoker: "GenRtc",
			args:    time.Now(),
		},
		{
			invoker: "GenOdo",
			args:    uint16(4321),
		},
		{
			invoker: "GenAntiThief",
			args:    nil,
		},
		{
			invoker: "GenReportFlush",
			args:    nil,
		},
		{
			invoker: "GenReportBlock",
			args:    false,
		},
		{
			invoker: "OvdState",
			args:    BikeStateNormal,
		},
		{
			invoker: "OvdReportInterval",
			args:    5 * time.Second,
		},
		{
			invoker: "OvdReportFrame",
			args:    FrameFull,
		},
		{
			invoker: "OvdRemoteSeat",
			args:    nil,
		},
		{
			invoker: "OvdRemoteAlarm",
			args:    nil,
		},
		{
			invoker: "AudioBeep",
			args:    nil,
		},
		{
			invoker: "FingerFetch",
			args:    nil,
			want:    []int{1, 2, 3, 4, 5},
			resMsg:  message([]byte("12345")),
		},
		{
			invoker: "FingerAdd",
			args:    nil,
			want:    3,
			resMsg:  message([]byte("3")),
		},
		{
			invoker: "FingerDel",
			args:    3,
		},
		{
			invoker: "FingerRst",
			args:    nil,
		},
		{
			invoker: "RemotePairing",
			args:    nil,
		},
		{
			invoker: "FotaVcu",
			args:    nil,
			resMsg:  message("VCU upgraded v.664 -> v.665"),
		},
		{
			invoker: "FotaHmi",
			args:    nil,
			resMsg:  message("HMI upgraded v.123 -> v.124"),
		},
		{
			invoker: "NetSendUssd",
			args:    "*123*10*3#",
			resMsg:  message("Terima kasih, permintaan kamu akan diproses,cek SMS untuk info lengkap. Dapatkan informasi seputar kartu Tri mu di aplikasi BimaTri, download di bima.tri.co.id"),
		},
		{
			invoker: "NetReadSms",
			args:    nil,
			resMsg:  message("Poin Bonstri kamu: 20 Sisa Kuota kamu : Kuota ++ 372 MB s.d 03/01/2031 13:30:18 Temukan beragam paket lain di bima+ https://goo.gl/RQ1DBA"),
		},
		{
			invoker: "HbarDrive",
			args:    ModeDriveEconomy,
		},
		{
			invoker: "HbarTrip",
			args:    ModeTripB,
		},
		{
			invoker: "HbarAvg",
			args:    ModeAvgEfficiency,
		},
		{
			invoker: "HbarReverse",
			args:    false,
		},
		{
			invoker: "McuSpeedMax",
			args:    uint8(90),
		},
		{
			invoker: "McuTemplates",
			args: []McuTemplate{
				{DisCur: 50, Torque: 10}, // economy
				{DisCur: 50, Torque: 20}, // standard
				{DisCur: 50, Torque: 25}, // sport
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.invoker, func(t *testing.T) {
			fakeRes := newFakeResponse(testVin, tC.invoker)
			if tC.resMsg != nil {
				fakeRes.Message = tC.resMsg
				if tC.want == nil {
					tC.want = string(tC.resMsg)
				}
			}

			cmder := newFakeCommander([][]byte{
				strToBytes(PREFIX_ACK),
				mockResponse(fakeRes),
			})
			defer cmder.Destroy()

			// call related method, pass in args, evaluate outs
			meth := reflect.ValueOf(cmder).MethodByName(tC.invoker)
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

func TestCommandError(t *testing.T) {
	testCases := []struct {
		desc string
	}{
		{
			desc: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

		})
	}
}
