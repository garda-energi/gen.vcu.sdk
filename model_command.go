package sdk

import "time"

type HeaderCommand struct {
	Header
	Code    uint8 `type:"uint8"`
	SubCode uint8 `type:"uint8"`
}

type commandPacket struct {
	Header  *HeaderCommand
	Message message
}

// command store essential command informations
type command struct {
	name    string
	invoker string
	code    uint8
	subCode uint8
	timeout time.Duration
}

// cmdList store command name by its code & subCode as index
var cmdList = [...][]command{
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
			name:    "GEN_BIKE_STATE",
			invoker: "GenBikeState",
		},
		command{
			name:    "GEN_LOCKDOWN",
			invoker: "GenLockDown",
		},
		command{
			name:    "GEN_CAN_DEBUG",
			invoker: "GenCanDebug",
		},
	},
	{
		command{
			name:    "REPORT_FLUSH",
			invoker: "ReportFlush",
		},
		command{
			name:    "REPORT_BLOCK",
			invoker: "ReportBlock",
		},
		command{
			name:    "REPORT_INTERVAL",
			invoker: "ReportInterval",
		},
		command{
			name:    "REPORT_FRAME",
			invoker: "ReportFrame",
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
		command{
			name:    "REMOTE_SEAT",
			invoker: "RemoteSeat",
		},
		command{
			name:    "REMOTE_ALARM",
			invoker: "RemoteAlarm",
		},
	},
	{
		command{
			name:    "FOTA_RESTART",
			invoker: "FotaRestart",
			timeout: 1 * 60 * time.Second,
		},
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
			name:    "HBAR_TRIPMETER",
			invoker: "HbarTripMeter",
		},
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
	{
		command{
			name:    "IMU_ANTITHIEF",
			invoker: "ImuAntiThief",
		},
	},
}

// cmdEvaluator is boolean evaluator for findCmd().
type cmdEvaluator func(code, subCode int, cmd *command) bool

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
func findCmd(evaluator cmdEvaluator) (*command, error) {
	for code, subCodes := range cmdList {
		for subCode, cmd := range subCodes {
			if evaluator(code, subCode, &cmd) {
				cmd.code = uint8(code)
				cmd.subCode = uint8(subCode)
				if cmd.timeout == 0 {
					cmd.timeout = DEFAULT_CMD_TIMEOUT
				}
				return &cmd, nil
			}
		}
	}
	return nil, errCmdNotFound
}
