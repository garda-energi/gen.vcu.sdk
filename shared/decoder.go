package shared

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

// TypeOfTime is for comparing struct type as time.Time
var TypeOfTime reflect.Type = reflect.ValueOf(time.Now()).Type()

// fix bug. in "decode func" before, It'll be error for case such as v isn't array of struct

// Decode read buffer reader than decode and set it to v.
// v is struct or pointer type that will contain decoded data
func Decode(rdr *bytes.Reader, v interface{}, tags ...Tagger) error {
	var err error

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if !rv.IsValid() || !rv.CanSet() {
		return nil
	}

	// default tag
	tag := Tagger{
		Len:    1,
		Factor: 1.0,
		Tipe:   "uint64",
	}
	// if tag is pass in argument
	if len(tags) > 0 {
		tag = tags[0]
	}

	switch rk := rv.Kind(); rk {

	case reflect.Ptr:
		rv.Set(reflect.New(rv.Type().Elem()))
		if err = Decode(rdr, rv.Interface()); err != nil {
			return err
		}

	case reflect.Struct:
		// if data type is time.Time
		if rv.Type() == TypeOfTime {
			b := make([]byte, tag.Len)
			binary.Read(rdr, binary.LittleEndian, &b)
			rv.Set(reflect.ValueOf(bytesToTime(b)))

		} else {
			for i := 0; i < rv.NumField() && rdr.Len() > 0; i++ {
				rvField := rv.Field(i)
				rtField := rv.Type().Field(i)

				tagField := DeTag(rtField.Tag, rvField.Kind())
				if err = Decode(rdr, rvField.Addr().Interface(), tagField); err != nil {
					return err
				}
			}
		}

	case reflect.Array:
		for j := 0; j < rv.Len(); j++ {
			if err = Decode(rdr, rv.Index(j).Addr().Interface()); err != nil {
				return err
			}
		}

	case reflect.String:
		x := make([]byte, tag.Len)
		binary.Read(rdr, binary.LittleEndian, &x)
		rv.SetString(bytesToStr(x))

	case reflect.Bool:
		var x bool
		binary.Read(rdr, binary.LittleEndian, &x)
		rv.SetBool(x)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		x := readUint(rdr, tag.Len)
		if !rv.OverflowUint(x) {
			rv.SetUint(x)
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		x := readUint(rdr, tag.Len)
		x_int := int64(x)
		if !rv.OverflowInt(x_int) {
			rv.SetInt(x_int)
		}

	case reflect.Float32, reflect.Float64:
		var x64 float64

		x := readUint(rdr, tag.Len)

		if tag.Factor != 1 {
			x64 = convertToFloat64(tag.Tipe, x)
			x64 *= tag.Factor

		} else {
			// set as binary
			if rk == reflect.Float32 {
				x32 := math.Float32frombits(uint32(x))
				x64 = float64(x32)
			} else {
				x64 = math.Float64frombits(x)
			}
		}

		if !rv.OverflowFloat(x64) {
			rv.SetFloat(x64)
		}

	default:
		return errors.New("unsupported kind: " + rv.Kind().String())
	}

	return err
}

// readUint read len(length) data as uint64
func readUint(rdr io.Reader, len int) uint64 {
	// sometimes, data recived in length less than 8
	b := make([]byte, len)
	binary.Read(rdr, binary.BigEndian, &b)

	newb := make([]byte, 8)
	for i := 0; i < len; i++ {
		newb[i] = b[i]
	}

	return binary.LittleEndian.Uint64(newb)
}

// convertToFloat64 convert bytes data to float64.
// data read as typedata from tag
func convertToFloat64(typedata string, x uint64) (result float64) {
	// declaration
	rv := setVarOfTypeData(typedata)

	// set sesuai memori yang dideklarasi
	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		rv.SetUint(x)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		rv.SetInt(int64(x))
	}

	// convert ke float
	switch tmpIntr := rv.Interface().(type) {
	case uint8:
		result = float64(tmpIntr)
	case uint16:
		result = float64(tmpIntr)
	case uint32:
		result = float64(tmpIntr)
	case uint64:
		result = float64(tmpIntr)
	case int8:
		result = float64(tmpIntr)
	case int16:
		result = float64(tmpIntr)
	case int32:
		result = float64(tmpIntr)
	case int64:
		result = float64(tmpIntr)
	}
	return result
}

// bytesToTime convert bytes slice (little endian) to time
// value of bytes data is :
//   1 byte of year
//   1 byte of month
//   1 byte of day
//   1 byte of hour
//   1 byte of minute
//   1 byte of second
//   1 byte of weekday (ignored)
func bytesToTime(b []byte) time.Time {
	var data string
	for _, v := range b[:6] {
		data += fmt.Sprintf("%02d", uint8(v))
	}

	datetime, _ := time.Parse("060102150405", data)
	return datetime
}

// bytesToStr convert byte slice (little endian) to string
func bytesToStr(b []byte) string {
	return string(util.Reverse(b))
}
