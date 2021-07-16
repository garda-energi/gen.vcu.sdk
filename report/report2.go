package report

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/pudjamansyurin/gen_vcu_sdk/shared"
)

type Items map[string]Item

type Item struct {
	Value interface{}
	Unit  string
}

func (r *Report) DecodeReportList() (Items, error) {
	var result = Items{}

	rdr := bytes.NewReader(r.bytes)
	for _, packet := range ReportFullPacket() {
		buf := make([]byte, getLen(packet))
		if n, err := rdr.Read(buf); err != nil && n == 0 {
			break
			// return Items{}, fmt.Errorf("failed to read buffer, %v", err)
		}

		item, err := getItem(packet, buf)
		if err != nil {
			return Items{}, fmt.Errorf("decoding error, %v\b", err)
		}

		result[packet.Name] = item
	}

	return result, nil
}

func getItem(packet shared.Packet, buf []byte) (Item, error) {
	item := Item{
		Unit: packet.Unit,
	}

	switch packet.Dst {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := toUint64(buf)
		if packet.Factor != 0 {
			v *= uint64(packet.Factor)
		}
		item.Value = v
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := toInt64(buf)
		if packet.Factor != 0 {
			v *= int64(packet.Factor)
		}
		item.Value = v
	case reflect.Float32:
		item.Value = toFloat64(buf, float64(packet.Factor))
	case reflect.Bool:
		item.Value = toBool(buf)
	case reflect.String:
		item.Value = toAscii(buf)

	default:
		if packet.Datetime {
			item.Value = toUnixTime(buf)
		} else {
			return Item{}, errors.New("unsupported kind: " + packet.Dst.String())
		}
	}

	return item, nil
}

func getLen(packet shared.Packet) int {
	var len int

	if packet.Len != 0 {
		len = packet.Len
	} else {
		src := packet.Src
		if src == reflect.Invalid {
			src = packet.Dst
		}

		switch src {
		case reflect.Bool, reflect.Uint8, reflect.Int8:
			len = 1
		case reflect.Uint16, reflect.Int16:
			len = 2
		case reflect.Uint32, reflect.Int32:
			len = 4
		case reflect.Uint64, reflect.Int64:
			len = 8
		default:
			if packet.Datetime {
				len = 7
			}
		}
	}

	return len
}
