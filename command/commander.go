package command

import (
	"reflect"
	"time"
)

type ValidatorFunc func(b []byte) bool
type EncoderFunc func(b []byte) []byte

type Commander struct {
	Name      string
	Desc      string
	Code      uint8
	SubCode   uint8
	Tipe      reflect.Kind
	Timeout   time.Duration
	Validator ValidatorFunc
	Encoder   EncoderFunc
}

var CMDERS = [][]Commander{
	{
		Commander{
			Name: "GEN_INFO",
		},
		Commander{
			Name: "GEN_LED",
			Desc: "Set system led",
			// Code:    CMDC_GEN,
			// SubCode: CMD_SUBCODE(CMD_GEN_LED),
			// Tipe:    reflect.Bool,
		},
		Commander{
			Name: "GEN_RTC",
			Desc: "Set datetime (d[1-7])",
			//   Code: CMDC_GEN,
			//   SubCode: 2,
			//   size: 7,
			// Tipe: time.Time,
			//   range: ["YYMMDDHHmmss0d"],
			//   Validator: (v) => Validator.GEN.RTC(v),
			//   formatCmd: (v) => TimeStamp(v),
		},
		Commander{
			Name: "GEN_ODO",
			Desc: "Set odometer (km)",
			// Code:    CMDC_GEN,
			// SubCode: CMD_SUBCODE(CMD_GEN_ODO),
			// Tipe:    reflect.Uint16,
		},
		Commander{
			Name: "GEN_ANTITHIEF",
			Desc: "Toggle anti-thief motion detector",
			// Code:    CMDC_GEN,
			// SubCode: CMD_SUBCODE(CMD_GEN_ANTITHIEF),
		},
		Commander{
			Name: "GEN_RPT_FLUSH",
			Desc: "Flush report buffer",
			// Code:    CMDC_GEN,
			// SubCode: CMD_SUBCODE(CMD_GEN_RPT_FLUSH),
		},
		Commander{
			Name: "GEN_RPT_BLOCK",
			Desc: "Block report buffer",
			// Code:    CMDC_GEN,
			// SubCode: CMD_SUBCODE(CMD_GEN_RPT_BLOCK),
			// Tipe:    reflect.Bool,
		},
	},
	{
		Commander{
			Name: "OVD_STATE",
			Desc: "Override bike state",
			// Code:    CMDC_OVD,
			// SubCode: CMD_SUBCODE(CMD_OVD_STATE),
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return between(b, uint8(shared.BIKE_STATE_NORMAL), uint8(shared.BIKE_STATE_RUN))
			// },
		},
		Commander{
			Name: "OVD_RPT_INTERVAL",
			Desc: "Override report interval",
			// Code:    CMDC_OVD,
			// SubCode: CMD_SUBCODE(CMD_OVD_RPT_INTERVAL),
			// Tipe:    reflect.Uint16,
		},
		Commander{
			Name: "OVD_RPT_FRAME",
			Desc: "Override report frame",
			// Code:    CMDC_OVD,
			// SubCode: CMD_SUBCODE(CMD_OVD_RPT_FRAME),
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return contains(b, uint8(shared.FRAME_ID_SIMPLE), uint8(shared.FRAME_ID_FULL))
			// },
		},
		Commander{
			Name: "OVD_RMT_SEAT",
			Desc: "Override remote seat button",
			// Code:    CMDC_OVD,
			// SubCode: CMD_SUBCODE(CMD_OVD_RMT_SEAT),
		},
		Commander{
			Name: "OVD_RMT_ALARM",
			Desc: "Override remote alarm button",
			// Code:    CMDC_OVD,
			// SubCode: CMD_SUBCODE(CMD_OVD_RMT_ALARM),
		},
	},
	{
		Commander{
			Name: "AUDIO_BEEP",
			Desc: "Beep the audio module",
			// Code:    CMDC_AUDIO,
			// SubCode: CMD_SUBCODE(CMD_AUDIO_BEEP),
		},
	},
	{
		Commander{
			Name: "FINGER_FETCH",
			Desc: "Get all registered id",
			// Code:    CMDC_FGR,
			// SubCode: CMD_SUBCODE(CMD_FGR_FETCH),
			Timeout: 15 * time.Second,
		},
		Commander{
			Name: "FINGER_ADD",
			Desc: "Add a new fingerprint",
			// Code:    CMDC_FGR,
			// SubCode: CMD_SUBCODE(CMD_FGR_ADD),
			Timeout: 20 * time.Second,
		},
		Commander{
			Name: "FINGER_DEL",
			Desc: "Delete a fingerprint",
			// Code:    CMDC_FGR,
			// SubCode: CMD_SUBCODE(CMD_FGR_DEL),
			// Tipe:    reflect.Uint8,
			Timeout: 15 * time.Second,
			// Validator: func(b []byte) bool {
			// 	return between(b, 1, shared.FINGERPRINT_MAX)
			// },
		},
		Commander{
			Name: "FINGER_RST",
			Desc: "Reset all fingerprints",
			// Code:    CMDC_FGR,
			// SubCode: CMD_SUBCODE(CMD_FGR_RST),
			Timeout: 15 * time.Second,
		},
	},
	{
		Commander{
			Name: "REMOTE_PAIRING",
			Desc: "Keyless pairing mode",
			// Code:    CMDC_RMT,
			// SubCode: CMD_SUBCODE(CMD_RMT_PAIRING),
			Timeout: 15 * time.Second,
		},
	},
	{
		Commander{
			Name: "FOTA_VCU",
			Desc: "Upgrade VCU firmware",
			// Code:    CMDC_FOTA,
			// SubCode: CMD_SUBCODE(CMD_FOTA_VCU),
			Timeout: 6 * 60 * time.Second,
		},
		Commander{
			Name: "FOTA_HMI",
			Desc: "Upgrade HMI firmware",
			// Code:    CMDC_FOTA,
			// SubCode: CMD_SUBCODE(CMD_FOTA_HMI),
			Timeout: 12 * 60 * time.Second,
		},
	},
	{
		Commander{
			Name: "NET_SEND_USSD",
			Desc: "Send USSD (ex: *123*10*3#)",
			//   Code: CMDC_NET,
			//   SubCode: CMD_SUBCODE(CMD_NET_SEND_USSD),
			// Tipe: reflect.String,
			//   size: 20,
			//   Validator: (v) => Validator.NET.SEND_USSD(v),
			//   formatCmd: (v) => AsciiToHex(v),
		},
		Commander{
			Name: "NET_READ_SMS",
			Desc: "Read last SMS",
			// Code:    CMDC_NET,
			// SubCode: CMD_SUBCODE(CMD_NET_READ_SMS),
		},
	},
	{
		Commander{
			Name: "CON_APN",
			Desc: "Set APN connection (ex: 3gprs;3gprs;3gprs)",
			//   Code: CMDC_CON,
			//   SubCode: 0,
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
		Commander{
			Name: "CON_FTP",
			Desc: "Set FTP connection",
			//   Code: CMDC_CON,
			//   SubCode: 1,
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
		Commander{
			Name: "CON_MQTT",
			Desc: "Set MQTT connection",
			//   Code: CMDC_CON,
			//   SubCode: 2,
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
		Commander{
			Name: "HBAR_DRIVE",
			Desc: "Set handlebar drive mode",
			// Code:    CMDC_HBAR,
			// SubCode: CMD_SUBCODE(CMD_HBAR_DRIVE),
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, uint8(shared.MODE_DRIVE_limit)-1)
			// },
		},
		Commander{
			Name: "HBAR_TRIP",
			Desc: "Set handlebar trip mode",
			// Code:    CMDC_HBAR,
			// SubCode: CMD_SUBCODE(CMD_HBAR_TRIP),
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, uint8(shared.MODE_TRIP_limit)-1)
			// },
		},
		Commander{
			Name: "HBAR_AVG",
			Desc: "Set handlebar average mode",
			// Code:    CMDC_HBAR,
			// SubCode: CMD_SUBCODE(CMD_HBAR_AVG),
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, uint8(shared.MODE_AVG_limit)-1)
			// },
		},
		Commander{
			Name: "HBAR_REVERSE",
			Desc: "Set handlebar reverse state",
			// Code:    CMDC_HBAR,
			// SubCode: CMD_SUBCODE(CMD_HBAR_REVERSE),
			// Tipe:    reflect.Bool,
		},
	},
	{
		Commander{
			Name: "MCU_SPEED_MAX",
			Desc: "Set MCU max speed",
			// Code:    CMDC_MCU,
			// SubCode: CMD_SUBCODE(CMD_MCU_SPEED_MAX),
			// Tipe:    reflect.Uint8,
			// Validator: func(b []byte) bool {
			// 	return max(b, shared.SPEED_MAX)
			// },
		},
		Commander{
			Name: "MCU_TEMPLATES",
			Desc: "Set MCU templates (ex: 50,15;50,20;50,25)",
			//   Code: CMDC_MCU,
			//   SubCode: 1,
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
