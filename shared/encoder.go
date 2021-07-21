package shared

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

// next problem in implementation.
// header has length and length get after encode.
// alternative solution : change bit #3 after encode as length of body

// Encode struct or pointer of struct to bytes
func Encode(v interface{}) (resBytes []byte, resError error) {
	buf := &bytes.Buffer{}

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		log.Fatal(rv.Kind())
		return nil, errors.New("type not match")
	}

	for i := 0; i < rv.NumField(); i++ {
		rvField := rv.Field(i)
		rtField := rv.Type().Field(i)

		tag := DeTag(rtField.Tag, rvField.Kind())

		switch rk := rvField.Kind(); rk {

		case reflect.Ptr:
			if !rvField.IsNil() {
				b, err := Encode(rvField.Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(b)
			}

		case reflect.Struct:
			if rvField.Type() == typeOfTime {
				t := rvField.Interface().(time.Time)
				b := timeToBytes(t)
				buf.Write(b)
			} else {
				b, err := Encode(rvField.Addr().Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(b)
			}

		case reflect.Array:
			for j := 0; j < rvField.Len(); j++ {
				b, err := Encode(rvField.Index(j).Addr().Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(b)
			}

		case reflect.String:
			s := rvField.String()
			b := util.Reverse([]byte(s))
			buf.Write(b)

		case reflect.Bool:
			x := rvField.Bool()
			var b uint8 = 0
			if x {
				b = 1
			}
			buf.Write([]byte{b})

		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			n := rvField.Uint()
			b := uintToBytes(rk, n)[:tag.Len]
			buf.Write(b)

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			n := rvField.Int()
			b := uintToBytes(rk, uint64(n))[:tag.Len]
			buf.Write(b)

		case reflect.Float32, reflect.Float64:
			var b []byte
			n := rvField.Float()

			if tag.Factor != 1 {
				n /= tag.Factor
				b = convertFloat64ToBytes(tag.Tipe, n)

			} else {
				// set sesuai biner
				if rk == reflect.Float32 {
					x32 := math.Float32bits(float32(n))
					b = uintToBytes(rk, uint64(x32))[:tag.Len]
				} else {
					x64 := math.Float64bits(n)
					b = uintToBytes(rk, x64)[:tag.Len]
				}
			}
			buf.Write(b)

		default:
			return nil, errors.New("unsupported kind: " + rv.Kind().String())
		}
	}

	return buf.Bytes(), nil
}

func uintToBytes(rk reflect.Kind, v uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	switch rk {
	case reflect.Uint8, reflect.Int8:
		return b[:1]
	case reflect.Uint16, reflect.Int16:
		return b[:2]
	case reflect.Uint32, reflect.Int32, reflect.Float32:
		return b[:4]
	default:
		return b[:8]
	}
}

// convert float data to bytes
func convertFloat64ToBytes(typedata string, v float64) []byte {
	// declaration
	// set tmp variable as datatype
	var rv reflect.Value
	switch typedata {
	case "uint8":
		var tmp uint8
		rv = reflect.ValueOf(&tmp)
	case "uint16":
		var tmp uint16
		rv = reflect.ValueOf(&tmp)
	case "uint32":
		var tmp uint32
		rv = reflect.ValueOf(&tmp)
	case "uint64", "uint":
		var tmp uint64
		rv = reflect.ValueOf(&tmp)
	case "int8":
		var tmp int8
		rv = reflect.ValueOf(&tmp)
	case "int16":
		var tmp int16
		rv = reflect.ValueOf(&tmp)
	case "int32":
		var tmp int32
		rv = reflect.ValueOf(&tmp)
	case "int64", "int":
		var tmp int64
		rv = reflect.ValueOf(&tmp)
	default:
		var tmp uint64
		rv = reflect.ValueOf(&tmp)
	}
	rv = rv.Elem()

	// convert to bytes
	b := uintToBytes(rv.Kind(), uint64(v))
	return b
}

// 1. set time to string as year-month-day-hour-minute-second
// 2. split date in string by '-'
// 3. convert to byte for each string splitted
// 4. than return bytes
func timeToBytes(t time.Time) []byte {
	dateStr := t.Format("06-01-02-15-04-05")
	datesStr := strings.Split(dateStr, "-")
	b := make([]byte, 7)
	for i, v := range datesStr {
		tmp, err := strconv.ParseUint(v, 10, 8)
		if err == nil {
			b[i] = uint8(tmp)
		}
	}
	b[6] = 1

	return b
}
