package sdk

import (
	"errors"
	"time"
)

type command struct {
	name    string
	method  string
	code    uint8
	subCode uint8
	timeout time.Duration
}

// getCommand get related command (code & subCode) by name
func getCommand(name string) (*command, error) {
	for code, subCodes := range commands {
		for subCode, cmd := range subCodes {
			if cmd.name == name {
				cmd.code = uint8(code)
				cmd.subCode = uint8(subCode)
				if cmd.timeout == 0 {
					cmd.timeout = DEFAULT_CMD_TIMEOUT
				}
				return &cmd, nil
			}
		}
	}
	return nil, errors.New("no command found")
}

var commands = [][]command{
	{
		command{
			name:   "GEN_INFO",
			method: "GenInfo",
		},
		command{
			name:   "GEN_LED",
			method: "GenLed",
		},
		command{
			name:   "GEN_RTC",
			method: "GenRtc",
		},
		command{
			name:   "GEN_ODO",
			method: "GenOdo",
		},
		command{
			name:   "GEN_ANTI_THIEF",
			method: "GenAntiThief",
		},
		command{
			name:   "GEN_RPT_FLUSH",
			method: "GenReportFlush",
		},
		command{
			name:   "GEN_RPT_BLOCK",
			method: "GenReportBlock",
		},
	},
	{
		command{
			name:   "OVD_STATE",
			method: "OvdState",
		},
		command{
			name:   "OVD_RPT_INTERVAL",
			method: "OvdReportInterval",
		},
		command{
			name:   "OVD_RPT_FRAME",
			method: "OvdReportFrame",
		},
		command{
			name:   "OVD_RMT_SEAT",
			method: "OvdRemoteSeat",
		},
		command{
			name:   "OVD_RMT_ALARM",
			method: "OvdRemoteAlarm",
		},
	},
	{
		command{
			name:   "AUDIO_BEEP",
			method: "AudioBeep",
		},
	},
	{
		command{
			name:    "FINGER_FETCH",
			method:  "FingerFetch",
			timeout: 15 * time.Second,
		},
		command{
			name:    "FINGER_ADD",
			method:  "FingerAdd",
			timeout: 20 * time.Second,
		},
		command{
			name:    "FINGER_DEL",
			method:  "FingerDel",
			timeout: 15 * time.Second,
		},
		command{
			name:    "FINGER_RST",
			method:  "FingerRst",
			timeout: 15 * time.Second,
		},
	},
	{
		command{
			name:    "REMOTE_PAIRING",
			method:  "RemotePairing",
			timeout: 15 * time.Second,
		},
	},
	{
		command{
			name:    "FOTA_VCU",
			method:  "FotaVcu",
			timeout: 6 * 60 * time.Second,
		},
		command{
			name:    "FOTA_HMI",
			method:  "FotaHmi",
			timeout: 12 * 60 * time.Second,
		},
	},
	{
		command{
			name:   "NET_SEND_USSD",
			method: "NetSendUssd",
		},
		command{
			name:   "NET_READ_SMS",
			method: "NetReadSms",
		},
	},
	{
		// TODO: finish CON command handler on VCU device (pending)
		command{
			name: "CON_APN",
			// desc: "Set APN connection (ex: 3gprs;3gprs;3gprs)",
			//   range: [
			//     [1, 30],
			//     [1, 30],
			//     [1, 30],
			//   ],
			//   size: 3 * 30,
			//   Tipe: "[char Name, user, pass][3]",
			//   Validator: (v) => Validator.CON(v, 3),
			//   formatCmd: (v) => AsciiToHex(v),
		},
		command{
			name: "CON_FTP",
			// desc: "Set FTP connection",
			//   range: [
			//     [1, 30],
			//     [1, 30],
			//     [1, 30],
			//   ],
			//   size: 3 * 30,
			//   Tipe: "[char host, user, pass][3]",
			//   Validator: (v) => Validator.CON(v, 3),
			//   formatCmd: (v) => AsciiToHex(v),
		},
		command{
			name: "CON_MQTT",
			// desc: "Set MQTT connection",
			//   range: [
			//     [1, 30],
			//     [1, 30],
			//     [1, 30],
			//     [1, 30],
			//   ],
			//   size: 4 * 30,
			//   Tipe: "[char host, port, user, pass][4]",
			//   Validator: (v) => Validator.CON(v, 4),
			//   formatCmd: (v) => AsciiToHex(v),
		},
	},
	{
		command{
			name:   "HBAR_DRIVE",
			method: "HbarDrive",
		},
		command{
			name:   "HBAR_TRIP",
			method: "HbarTrip",
		},
		command{
			name:   "HBAR_AVG",
			method: "HbarAvg",
		},
		command{
			name:   "HBAR_REVERSE",
			method: "HbarReverse",
		},
	},
	{
		command{
			name:   "MCU_SPEED_MAX",
			method: "McuSpeedMax",
		},
		command{
			name:   "MCU_TEMPLATES",
			method: "McuTemplates",
		},
	},
}
