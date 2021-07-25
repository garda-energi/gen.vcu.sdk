package sdk

import (
	"errors"
	"time"
)

type command struct {
	name    string
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
			name: "GEN_INFO",
		},
		command{
			name: "GEN_LED",
		},
		command{
			name: "GEN_RTC",
		},
		command{
			name: "GEN_ODO",
		},
		command{
			name: "GEN_ANTI_THIEF",
		},
		command{
			name: "GEN_RPT_FLUSH",
		},
		command{
			name: "GEN_RPT_BLOCK",
		},
	},
	{
		command{
			name: "OVD_STATE",
		},
		command{
			name: "OVD_RPT_INTERVAL",
		},
		command{
			name: "OVD_RPT_FRAME",
		},
		command{
			name: "OVD_RMT_SEAT",
		},
		command{
			name: "OVD_RMT_ALARM",
		},
	},
	{
		command{
			name: "AUDIO_BEEP",
		},
	},
	{
		command{
			name:    "FINGER_FETCH",
			timeout: 15 * time.Second,
		},
		command{
			name:    "FINGER_ADD",
			timeout: 20 * time.Second,
		},
		command{
			name:    "FINGER_DEL",
			timeout: 15 * time.Second,
		},
		command{
			name:    "FINGER_RST",
			timeout: 15 * time.Second,
		},
	},
	{
		command{
			name:    "REMOTE_PAIRING",
			timeout: 15 * time.Second,
		},
	},
	{
		command{
			name:    "FOTA_VCU",
			timeout: 6 * 60 * time.Second,
		},
		command{
			name:    "FOTA_HMI",
			timeout: 12 * 60 * time.Second,
		},
	},
	{
		command{
			name: "NET_SEND_USSD",
		},
		command{
			name: "NET_READ_SMS",
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
			name: "HBAR_DRIVE",
		},
		command{
			name: "HBAR_TRIP",
		},
		command{
			name: "HBAR_AVG",
		},
		command{
			name: "HBAR_REVERSE",
		},
	},
	{
		command{
			name: "MCU_SPEED_MAX",
		},
		command{
			name: "MCU_TEMPLATES",
		},
	},
}
