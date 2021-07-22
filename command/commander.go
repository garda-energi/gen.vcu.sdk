package command

import (
	"errors"
	"time"
)

type commander struct {
	name     string
	code     uint8
	sub_code uint8
	timeout  time.Duration
}

// getCmder get related commander by name
func getCmder(name string) (*commander, error) {
	for code, sub_codes := range commands {
		for sub_code, cmder := range sub_codes {
			if cmder.name == name {
				cmder.code = uint8(code)
				cmder.sub_code = uint8(sub_code)

				if cmder.timeout == 0 {
					cmder.timeout = DEFAULT_CMD_TIMEOUT
				}

				return &cmder, nil
			}
		}
	}

	return nil, errors.New("command invalid")
}

var commands = [][]commander{
	{
		commander{
			name: "GEN_INFO",
		},
		commander{
			name: "GEN_LED",
		},
		commander{
			name: "GEN_RTC",
		},
		commander{
			name: "GEN_ODO",
		},
		commander{
			name: "GEN_ANTI_THIEF",
		},
		commander{
			name: "GEN_RPT_FLUSH",
		},
		commander{
			name: "GEN_RPT_BLOCK",
		},
	},
	{
		commander{
			name: "OVD_STATE",
		},
		commander{
			name: "OVD_RPT_INTERVAL",
		},
		commander{
			name: "OVD_RPT_FRAME",
		},
		commander{
			name: "OVD_RMT_SEAT",
		},
		commander{
			name: "OVD_RMT_ALARM",
		},
	},
	{
		commander{
			name: "AUDIO_BEEP",
		},
	},
	{
		commander{
			name:    "FINGER_FETCH",
			timeout: 15 * time.Second,
		},
		commander{
			name:    "FINGER_ADD",
			timeout: 20 * time.Second,
		},
		commander{
			name:    "FINGER_DEL",
			timeout: 15 * time.Second,
		},
		commander{
			name:    "FINGER_RST",
			timeout: 15 * time.Second,
		},
	},
	{
		commander{
			name:    "REMOTE_PAIRING",
			timeout: 15 * time.Second,
		},
	},
	{
		commander{
			name:    "FOTA_VCU",
			timeout: 6 * 60 * time.Second,
		},
		commander{
			name:    "FOTA_HMI",
			timeout: 12 * 60 * time.Second,
		},
	},
	{
		commander{
			name: "NET_SEND_USSD",
		},
		commander{
			name: "NET_READ_SMS",
		},
	},
	{
		// TODO: finish CON command handler on VCU device (pending)
		commander{
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
		commander{
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
		commander{
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
		commander{
			name: "HBAR_DRIVE",
		},
		commander{
			name: "HBAR_TRIP",
		},
		commander{
			name: "HBAR_AVG",
		},
		commander{
			name: "HBAR_REVERSE",
		},
	},
	{
		commander{
			name: "MCU_SPEED_MAX",
		},
		commander{
			name: "MCU_TEMPLATES",
			// desc: "Set MCU templates (ex: 50,15;50,20;50,25)",
			//   range: [
			//     [1, 32767],
			//     [1, 3276],
			//   ],
			//   size: 4 * config.mode.drive.length,
			//   Tipe: "[uint16_t discur, torque][3]",
			//   Validator: (v) => Validator.MCU.TEMPLATES(v),
			//   formatCmd: (v) => formatter.MCU.TEMPLATES(v),
		},
	},
}
