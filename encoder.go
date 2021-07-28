package sdk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"reflect"
	"time"
)

// TODO: next problem in implementation.
// header has length and length get after encode.
// alternative solution : change bit #3 after encode as length of body

// // encodeResponse combine response and message to bytes packet.
// // TODO: redundant code, it is same pattern with encodeCommand()
// func encodeResponse(vin int, cmd *command, msg message) ([]byte, error) {
// 	if msg.overflow() {
// 		return nil, errInputOutOfRange("message")
// 	}

// 	rp := newResponsePacket(vin, cmd, msg)

// 	resBytes, err := encode(&rp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// change Header.Size
// 	resBytes[2] = uint8(len(resBytes) - 3)
// 	return resBytes, nil
// }

// encodeCommand combine command and message to bytes packet.
func encodeCommand(vin int, cmd *command, msg message) ([]byte, error) {
	if msg.overflow() {
		return nil, errInputOutOfRange("message")
	}

	cp := newCommandPacket(vin, cmd, msg)

	resBytes, err := encode(&cp)
	if err != nil {
		return nil, err
	}

	// change Header.Size
	resBytes[2] = uint8(len(resBytes) - 3)
	return resBytes, nil
}

// encode struct or pointer of struct to bytes
func encode(v interface{}, tags ...tagger) (resBytes []byte, resError error) {
	buf := &bytes.Buffer{}

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	// default tag
	tag := newTagger()
	// if tag is pass in argument
	if len(tags) > 0 {
		tag = tags[0]
	}

	switch rk := rv.Kind(); rk {

	case reflect.Ptr:
		if !rv.IsNil() {
			b, err := encode(rv.Interface())
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}

	case reflect.Struct:
		if rv.Type() == typeOfTime {
			t := rv.Interface().(time.Time)
			b := timeToBytes(t)
			buf.Write(b)
		} else {
			for i := 0; i < rv.NumField(); i++ {
				rvField := rv.Field(i)
				rtField := rv.Type().Field(i)

				tagField := deTag(rtField.Tag, rvField.Kind())

				b, err := encode(rvField.Addr().Interface(), tagField)
				if err != nil {
					return nil, err
				}
				buf.Write(b)
			}
		}

	case reflect.Array:
		for j := 0; j < rv.Len(); j++ {
			b, err := encode(rv.Index(j).Addr().Interface())
			if err != nil {
				return nil, err
			}
			buf.Write(b)
		}

	case reflect.Slice:
		if rv.Type() == typeOfMessage && rv.Len() > 0 {
			b := rv.Interface().(message)
			buf.Write([]byte(b))
		}

	case reflect.String:
		s := rv.String()
		b := strToBytes(s)
		buf.Write(b)

	case reflect.Bool:
		x := rv.Bool()
		b := boolToBytes(x)
		buf.Write(b)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		n := rv.Uint()
		b := uintToBytes(rk, n)[:tag.Len]
		buf.Write(b)

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		n := rv.Int()
		b := uintToBytes(rk, uint64(n))[:tag.Len]
		buf.Write(b)

	case reflect.Float32, reflect.Float64:
		var b []byte
		n := rv.Float()

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

	return buf.Bytes(), nil
}

// setVarOfTypeData create variable as typedata
func setVarOfTypeData(typedata string) reflect.Value {
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
	return rv.Elem()
}

// convertFloat64ToBytes convert float data to bytes
func convertFloat64ToBytes(typedata string, v float64) []byte {
	rv := setVarOfTypeData(typedata)
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
	dt := []int{
		t.Year() - 2000,
		int(t.Month()),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second(),
		int(t.Weekday()),
	}
	var buf bytes.Buffer
	for _, v := range dt {
		binary.Write(&buf, binary.LittleEndian, byte(v))
	}
	return buf.Bytes()
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
