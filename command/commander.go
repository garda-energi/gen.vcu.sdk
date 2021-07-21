package command

import (
	"errors"
	"reflect"
	"time"
)

type ValidatorFunc func(b []byte) bool
type EncoderFunc func(b []byte) []byte

type commander struct {
	name     string
	code     uint8
	sub_code uint8
	timeout  time.Duration
	// TODO: remove if not used
	Tipe      reflect.Kind
	Validator ValidatorFunc
	Encoder   EncoderFunc
}

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
			// desc: "Set odometer (km)",
			// Tipe:    reflect.Uint16,
		},
		commander{
			name: "GEN_ANTITHIEF",
			// desc: "Toggle anti-thief motion detector",
		},
		commander{
			name: "GEN_RPT_FLUSH",
			// desc: "Flush report buffer",
		},
		commander{
			name: "GEN_RPT_BLOCK",
			// desc: "Block report buffer",
			// Tipe:    reflect.Bool,
		},
	},
	{
		commander{
			name: "OVD_STATE",
			// desc: "Override bike state",
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return between(b, uint8(shared.BIKE_STATE_NORMAL), uint8(shared.BIKE_STATE_RUN))
			// },
		},
		commander{
			name: "OVD_RPT_INTERVAL",
			// desc: "Override report interval",
			// Tipe:    reflect.Uint16,
		},
		commander{
			name: "OVD_RPT_FRAME",
			// desc: "Override report frame",
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return contains(b, uint8(shared.FRAME_ID_SIMPLE), uint8(shared.FRAME_ID_FULL))
			// },
		},
		commander{
			name: "OVD_RMT_SEAT",
			// desc: "Override remote seat button",
		},
		commander{
			name: "OVD_RMT_ALARM",
			// desc: "Override remote alarm button",
		},
	},
	{
		commander{
			name: "AUDIO_BEEP",
			// desc: "Beep the audio module",
		},
	},
	{
		commander{
			name: "FINGER_FETCH",
			// desc: "Get all registered id",
			timeout: 15 * time.Second,
		},
		commander{
			name: "FINGER_ADD",
			// desc: "Add a new fingerprint",
			timeout: 20 * time.Second,
		},
		commander{
			name: "FINGER_DEL",
			// desc: "Delete a fingerprint",
			// Tipe:    reflect.Uint8,
			timeout: 15 * time.Second,
			// Validator: func(b []byte) bool {
			// 	return between(b, 1, shared.FINGERPRINT_MAX)
			// },
		},
		commander{
			name: "FINGER_RST",
			// desc: "Reset all fingerprints",
			timeout: 15 * time.Second,
		},
	},
	{
		commander{
			name: "REMOTE_PAIRING",
			// desc: "Keyless pairing mode",
			timeout: 15 * time.Second,
		},
	},
	{
		commander{
			name: "FOTA_VCU",
			// desc: "Upgrade VCU firmware",
			timeout: 6 * 60 * time.Second,
		},
		commander{
			name: "FOTA_HMI",
			// desc: "Upgrade HMI firmware",
			timeout: 12 * 60 * time.Second,
		},
	},
	{
		commander{
			name: "NET_SEND_USSD",
			// desc: "Send USSD (ex: *123*10*3#)",
			// Tipe: reflect.String,
			//   size: 20,
			//   Validator: (v) => Validator.NET.SEND_USSD(v),
			//   formatCmd: (v) => AsciiToHex(v),
		},
		commander{
			name: "NET_READ_SMS",
			// desc: "Read last SMS",
		},
	},
	{
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
			// desc: "Set handlebar drive mode",
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, uint8(shared.MODE_DRIVE_limit)-1)
			// },
		},
		commander{
			name: "HBAR_TRIP",
			// desc: "Set handlebar trip mode",
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, uint8(shared.MODE_TRIP_limit)-1)
			// },
		},
		commander{
			name: "HBAR_AVG",
			// desc: "Set handlebar average mode",
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, uint8(shared.MODE_AVG_limit)-1)
			// },
		},
		commander{
			name: "HBAR_REVERSE",
			// desc: "Set handlebar reverse state",
			// Tipe:    reflect.Bool,
		},
	},
	{
		commander{
			name: "MCU_SPEED_MAX",
			// desc: "Set MCU max speed",
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, shared.SPEED_MAX)
			// },
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
