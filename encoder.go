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

// encode combine command and value to bytes packet.
// can it be replaced with encode() func bellow ?
func encodeCommand(vin int, cmd *command, val message) ([]byte, error) {
	if val.overflow() {
		return nil, errInputOutOfRange("payload")
	}

	now := time.Now()
	// var buf bytes.Buffer
	// ed := binary.LittleEndian
	// binary.Write(&buf, ed, reverseBytes(val))
	// binary.Write(&buf, ed, cmd.subCode)
	// binary.Write(&buf, ed, cmd.code)
	// binary.Write(&buf, ed, reverseBytes(timeToBytes(now)))
	// binary.Write(&buf, binary.BigEndian, uint32(vin))
	// binary.Write(&buf, ed, byte(buf.Len()))
	// binary.Write(&buf, ed, []byte(PREFIX_COMMAND))
	// bytes := reverseBytes(buf.Bytes())

	cp := &CommandPacket{
		Header: &HeaderCommand{
			Header: Header{
				Prefix:       PREFIX_COMMAND,
				Size:         0,
				Vin:          uint32(vin),
				SendDatetime: now,
			},
			Code:    cmd.code,
			SubCode: cmd.subCode,
		},
		Message: val,
	}

	resBytes, _ := encode(&cp)
	// change Header.Size
	resBytes[2] = uint8(len(resBytes) - 3)

	// compare
	// fmt.Println(resBytes)
	// fmt.Println(bytes)

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
