package sdk

import "time"

type Header struct {
	Prefix       string    `type:"string" len:"2"`
	Size         uint8     `type:"uint8" unit:"Bytes"`
	Vin          uint32    `type:"uint32"`
	SendDatetime time.Time `type:"unix_time" len:"7"`
}
