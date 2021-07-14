package model

import (
	"reflect"

	"github.com/pudjamansyurin/gen_vcu_sdk/coder"
)

type Command struct {
	HeaderCommand
	Payload [200]byte
}

type ValidatorFunc func(b []byte) bool
type EncoderFunc func(b []byte) []byte

type Cmd struct {
	name       string
	desc       string
	code       CMD_CODE
	subCode    CMD_SUBCODE
	tipe       reflect.Kind
	timeoutSec int64
	validator  ValidatorFunc
	encoder    EncoderFunc
}

var CMD_LIST = []Cmd{
	{
		name:    "GEN_INFO",
		desc:    "Gather device information",
		code:    CMDC_GEN,
		subCode: CMD_SUBCODE(CMD_GEN_INFO),
	},
	{
		name:    "GEN_LED",
		desc:    "Set system led",
		code:    CMDC_GEN,
		subCode: CMD_SUBCODE(CMD_GEN_LED),
		tipe:    reflect.Bool,
	},
	// {
	//   name: "GEN_RTC",
	//   desc: "Set datetime (d[1-7])",
	//   code: CMDC_GEN,
	//   subCode: 2,
	//   size: 7,
	//   tipe: "uint8_t",
	//   range: ["YYMMDDHHmmss0d"],
	//   validator: (v) => validator.GEN.RTC(v),
	//   formatCmd: (v) => TimeStamp(v),
	// },
	{
		name:    "GEN_ODO",
		desc:    "Set odometer (km)",
		code:    CMDC_GEN,
		subCode: CMD_SUBCODE(CMD_GEN_ODO),
		tipe:    reflect.Uint16,
	},
	{
		name:    "GEN_ANTITHIEF",
		desc:    "Toggle anti-thief motion detector",
		code:    CMDC_GEN,
		subCode: CMD_SUBCODE(CMD_GEN_ANTITHIEF),
	},
	{
		name:    "GEN_RPT_FLUSH",
		desc:    "Flush report buffer",
		code:    CMDC_GEN,
		subCode: CMD_SUBCODE(CMD_GEN_RPT_FLUSH),
	},
	{
		name:    "GEN_RPT_BLOCK",
		desc:    "Block report buffer",
		code:    CMDC_GEN,
		subCode: CMD_SUBCODE(CMD_GEN_RPT_BLOCK),
		tipe:    reflect.Bool,
	},
	{
		name:    "OVD_STATE",
		desc:    "Override vehicle state",
		code:    CMDC_OVD,
		subCode: CMD_SUBCODE(CMD_OVD_STATE),
		tipe:    reflect.Uint8,
		validator: func(b []byte) bool {
			return coder.ToUint8(b) <= 3
		},
	},
	{
		name:    "OVD_RPT_INTERVAL",
		desc:    "Override report interval",
		code:    CMDC_OVD,
		subCode: CMD_SUBCODE(CMD_OVD_RPT_INTERVAL),
		tipe:    reflect.Uint16,
	},
	{
		name:    "OVD_RPT_FRAME",
		desc:    "Override report frame",
		code:    CMDC_OVD,
		subCode: CMD_SUBCODE(CMD_OVD_RPT_FRAME),
		tipe:    reflect.Uint8,
		validator: func(b []byte) bool {
			return coder.ToUint8(b) <= 2
		},
	},
	{
		name:    "OVD_RMT_SEAT",
		desc:    "Override remote seat button",
		code:    CMDC_OVD,
		subCode: CMD_SUBCODE(CMD_OVD_RMT_SEAT),
	},
	{
		name:    "OVD_RMT_ALARM",
		desc:    "Override remote alarm button",
		code:    CMDC_OVD,
		subCode: CMD_SUBCODE(CMD_OVD_RMT_ALARM),
	},
	{
		name:    "AUDIO_BEEP",
		desc:    "Beep the audio module",
		code:    CMDC_AUDIO,
		subCode: CMD_SUBCODE(CMD_AUDIO_BEEP),
	},
	{
		name:       "FINGER_FETCH",
		desc:       "Get all registered id",
		code:       CMDC_FGR,
		subCode:    CMD_SUBCODE(CMD_FGR_FETCH),
		timeoutSec: 15,
	},
	{
		name:       "FINGER_ADD",
		desc:       "Add a new fingerprint",
		code:       CMDC_FGR,
		subCode:    CMD_SUBCODE(CMD_FGR_ADD),
		timeoutSec: 20,
	},
	{
		name:       "FINGER_DEL",
		desc:       "Delete a fingerprint",
		code:       CMDC_FGR,
		subCode:    CMD_SUBCODE(CMD_FGR_DEL),
		tipe:       reflect.Uint8,
		timeoutSec: 15,
		validator: func(b []byte) bool {
			return coder.ToUint8(b) >= 1 && coder.ToUint8(b) <= FINGERPRINT_MAX
		},
	},
	{
		name:       "FINGER_RST",
		desc:       "Reset all fingerprints",
		code:       CMDC_FGR,
		subCode:    CMD_SUBCODE(CMD_FGR_RST),
		timeoutSec: 15,
	},
	{
		name:       "REMOTE_PAIRING",
		desc:       "Keyless pairing mode",
		code:       CMDC_RMT,
		subCode:    CMD_SUBCODE(CMD_RMT_PAIRING),
		timeoutSec: 15,
	},
	{
		name:       "FOTA_VCU",
		desc:       "Upgrade VCU firmware",
		code:       CMDC_FOTA,
		subCode:    CMD_SUBCODE(CMD_FOTA_VCU),
		timeoutSec: 6 * 60,
	},
	{
		name:       "FOTA_HMI",
		desc:       "Upgrade HMI firmware",
		code:       CMDC_FOTA,
		subCode:    CMD_SUBCODE(CMD_FOTA_HMI),
		timeoutSec: 12 * 60,
	},
	// {
	//   name: "NET_SEND_USSD",
	//   desc: "Send USSD (ex: *123*10*3#)",
	//   code: CMDC_NET,
	//   subCode: CMD_SUBCODE(CMD_NET_SEND_USSD),
	// 	tipe: reflect.String,
	//   size: 20,
	//   validator: (v) => validator.NET.SEND_USSD(v),
	//   formatCmd: (v) => AsciiToHex(v),
	// },
	{
		name:    "NET_READ_SMS",
		desc:    "Read last SMS",
		code:    CMDC_NET,
		subCode: CMD_SUBCODE(CMD_NET_READ_SMS),
	},
	// {
	//   name: "CON_APN",
	//   desc: "Set APN connection (ex: 3gprs;3gprs;3gprs)",
	//   code: CMDC_CON,
	//   subCode: 0,
	//   range: [
	//     [1, 30],
	//     [1, 30],
	//     [1, 30],
	//   ],
	//   size: 3 * 30,
	//   tipe: "[char name, user, pass][3]",
	//   validator: (v) => validator.CON(v, 3),
	//   formatCmd: (v) => AsciiToHex(v),
	// },
	// {
	//   name: "CON_FTP",
	//   desc: "Set FTP connection",
	//   code: CMDC_CON,
	//   subCode: 1,
	//   range: [
	//     [1, 30],
	//     [1, 30],
	//     [1, 30],
	//   ],
	//   size: 3 * 30,
	//   tipe: "[char host, user, pass][3]",
	//   validator: (v) => validator.CON(v, 3),
	//   formatCmd: (v) => AsciiToHex(v),
	// },
	// {
	//   name: "CON_MQTT",
	//   desc: "Set MQTT connection",
	//   code: CMDC_CON,
	//   subCode: 2,
	//   range: [
	//     [1, 30],
	//     [1, 30],
	//     [1, 30],
	//     [1, 30],
	//   ],
	//   size: 4 * 30,
	//   tipe: "[char host, port, user, pass][4]",
	//   validator: (v) => validator.CON(v, 4),
	//   formatCmd: (v) => AsciiToHex(v),
	// },
	{
		name:    "HBAR_DRIVE",
		desc:    "Set handlebar drive mode",
		code:    CMDC_HBAR,
		subCode: CMD_SUBCODE(CMD_HBAR_DRIVE),
		tipe:    reflect.Uint8,
		validator: func(b []byte) bool {
			return coder.ToUint8(b) <= 2
		},
	},
	{
		name:    "HBAR_TRIP",
		desc:    "Set handlebar trip mode",
		code:    CMDC_HBAR,
		subCode: CMD_SUBCODE(CMD_HBAR_TRIP),
		tipe:    reflect.Uint8,
		validator: func(b []byte) bool {
			return coder.ToUint8(b) <= 2
		},
	},
	{
		name:    "HBAR_AVG",
		desc:    "Set handlebar average mode",
		code:    CMDC_HBAR,
		subCode: CMD_SUBCODE(CMD_HBAR_AVG),
		tipe:    reflect.Uint8,
		validator: func(b []byte) bool {
			return coder.ToUint8(b) <= 1
		},
	},
	{
		name:    "HBAR_REVERSE",
		desc:    "Set handlebar reverse state",
		code:    CMDC_HBAR,
		subCode: CMD_SUBCODE(CMD_HBAR_REVERSE),
		tipe:    reflect.Uint8,
		validator: func(b []byte) bool {
			return coder.ToUint8(b) <= 1
		},
	},
	{
		name:    "MCU_SPEED_MAX",
		desc:    "Set MCU max speed",
		code:    CMDC_MCU,
		subCode: CMD_SUBCODE(CMD_MCU_SPEED_MAX),
		tipe:    reflect.Uint8,
	},
	// {
	//   name: "MCU_TEMPLATES",
	//   desc: "Set MCU templates (ex: 50,15;50,20;50,25)",
	//   code: CMDC_MCU,
	//   subCode: 1,
	//   range: [
	//     [1, 32767],
	//     [1, 3276],
	//   ],
	//   size: 4 * config.mode.drive.length,
	//   tipe: "[uint16_t discur, torque][3]",
	//   validator: (v) => validator.MCU.TEMPLATES(v),
	//   formatCmd: (v) => formatter.MCU.TEMPLATES(v),
	// },
}
