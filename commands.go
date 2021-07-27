package sdk

import (
	"errors"
	"time"
)

// command store essential command informations
type command struct {
	name    string
	invoker string
	code    uint8
	subCode uint8
	timeout time.Duration
}

// cmdEvaluator is boolean evaluator for findCmd().
type cmdEvaluator func(code, subCode int, cmd *command) bool

// getCmdByName get related command by name
func getCmdByName(name string) (*command, error) {
	return findCmd(func(code, subCode int, cmd *command) bool {
		return cmd.name == name
	})
}

// getCmdByInvoker get related command by invoker
func getCmdByInvoker(invoker string) (*command, error) {
	return findCmd(func(code, subCode int, cmd *command) bool {
		return cmd.invoker == invoker
	})
}

// getCmdByCode get related command by code
func getCmdByCode(code, subCode int) (*command, error) {
	return findCmd(func(c, sc int, cmd *command) bool {
		return code == c && subCode == sc
	})
}

// findCmd find related cmd according to boolean evaluator
func findCmd(checker cmdEvaluator) (*command, error) {
	for code, subCodes := range commands {
		for subCode, cmd := range subCodes {
			if checker(code, subCode, &cmd) {
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

// commands store command name by its code & subCode as index
var commands = [][]command{
	{
		command{
			name:    "GEN_INFO",
			invoker: "GenInfo",
		},
		command{
			name:    "GEN_LED",
			invoker: "GenLed",
		},
		command{
			name:    "GEN_RTC",
			invoker: "GenRtc",
		},
		command{
			name:    "GEN_ODO",
			invoker: "GenOdo",
		},
		command{
			name:    "GEN_ANTI_THIEF",
			invoker: "GenAntiThief",
		},
		command{
			name:    "GEN_RPT_FLUSH",
			invoker: "GenReportFlush",
		},
		command{
			name:    "GEN_RPT_BLOCK",
			invoker: "GenReportBlock",
		},
	},
	{
		command{
			name:    "OVD_STATE",
			invoker: "OvdState",
		},
		command{
			name:    "OVD_RPT_INTERVAL",
			invoker: "OvdReportInterval",
		},
		command{
			name:    "OVD_RPT_FRAME",
			invoker: "OvdReportFrame",
		},
		command{
			name:    "OVD_RMT_SEAT",
			invoker: "OvdRemoteSeat",
		},
		command{
			name:    "OVD_RMT_ALARM",
			invoker: "OvdRemoteAlarm",
		},
	},
	{
		command{
			name:    "AUDIO_BEEP",
			invoker: "AudioBeep",
		},
	},
	{
		command{
			name:    "FINGER_FETCH",
			invoker: "FingerFetch",
			timeout: 15 * time.Second,
		},
		command{
			name:    "FINGER_ADD",
			invoker: "FingerAdd",
			timeout: 20 * time.Second,
		},
		command{
			name:    "FINGER_DEL",
			invoker: "FingerDel",
			timeout: 15 * time.Second,
		},
		command{
			name:    "FINGER_RST",
			invoker: "FingerRst",
			timeout: 15 * time.Second,
		},
	},
	{
		command{
			name:    "REMOTE_PAIRING",
			invoker: "RemotePairing",
			timeout: 15 * time.Second,
		},
	},
	{
		command{
			name:    "FOTA_VCU",
			invoker: "FotaVcu",
			timeout: 6 * 60 * time.Second,
		},
		command{
			name:    "FOTA_HMI",
			invoker: "FotaHmi",
			timeout: 12 * 60 * time.Second,
		},
	},
	{
		command{
			name:    "NET_SEND_USSD",
			invoker: "NetSendUssd",
		},
		command{
			name:    "NET_READ_SMS",
			invoker: "NetReadSms",
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
			name:    "HBAR_DRIVE",
			invoker: "HbarDrive",
		},
		command{
			name:    "HBAR_TRIP",
			invoker: "HbarTrip",
		},
		command{
			name:    "HBAR_AVG",
			invoker: "HbarAvg",
		},
		command{
			name:    "HBAR_REVERSE",
			invoker: "HbarReverse",
		},
	},
	{
		command{
			name:    "MCU_SPEED_MAX",
			invoker: "McuSpeedMax",
		},
		command{
			name:    "MCU_TEMPLATES",
			invoker: "McuTemplates",
		},
	},
}
