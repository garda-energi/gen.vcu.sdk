package sdk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"math"
	"reflect"
	"time"
)

// TODO: next problem in implementation.
// header has length and length get after encode.
// alternative solution : change bit #3 after encode as length of body

// encode struct or pointer of struct to bytes
func encode(v interface{}) (resBytes []byte, resError error) {
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

		tag := deTag(rtField.Tag, rvField.Kind())

		switch rk := rvField.Kind(); rk {

		case reflect.Ptr:
			if !rvField.IsNil() {
				b, err := encode(rvField.Interface())
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
				b, err := encode(rvField.Addr().Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(b)
			}

		case reflect.Array:
			for j := 0; j < rvField.Len(); j++ {
				b, err := encode(rvField.Index(j).Addr().Interface())
				if err != nil {
					return nil, err
				}
				buf.Write(b)
			}

		case reflect.String:
			s := rvField.String()
			b := strToBytes(s)
			buf.Write(b)

		case reflect.Bool:
			x := rvField.Bool()
			b := boolToBytes(x)
			buf.Write(b)

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

				// i don't know but it's work
				if rk == reflect.Float32 {
					n = float64(float32(n))
				}
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

// setVarOfTypeData create variable as typedata
func setVarOfTypeData(typedata string) reflect.Value {
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

	return rv
}

// convertFloat64ToBytes convert float data to bytes
func convertFloat64ToBytes(typedata string, v float64) []byte {
	// declaration
	rv := setVarOfTypeData(typedata)

	// convert to bytes
	b := uintToBytes(rv.Kind(), uint64(v))
	return b
}

// uintToBytes convert uint category type to byte slice (little endian)
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

// timeToBytes convert time to slice byte (big endian)
func timeToBytes(t time.Time) []byte {
	var buf bytes.Buffer
	ed := binary.LittleEndian
	binary.Write(&buf, ed, byte(t.Year()-2000))
	binary.Write(&buf, ed, byte(t.Month()))
	binary.Write(&buf, ed, byte(t.Day()))
	binary.Write(&buf, ed, byte(t.Hour()))
	binary.Write(&buf, ed, byte(t.Minute()))
	binary.Write(&buf, ed, byte(t.Second()))
	binary.Write(&buf, ed, byte(t.Weekday()))
	bytes := buf.Bytes()

	return bytes
}

// boolToBytes convert bool to byte slice.
func boolToBytes(d bool) []byte {
	var b uint8 = 0
	if d {
		b = 1
	}
	return []byte{b}
}

// strToBytes convert string to byte slice (little endian)
func strToBytes(d string) []byte {
	return reverseBytes([]byte(d))
}