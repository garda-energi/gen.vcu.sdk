package shared

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"reflect"
	"time"

	"github.com/pudjamansyurin/gen_vcu_sdk/util"
)

// this variable is for comparing struct type as time.Time
var typeOfTime reflect.Type = reflect.ValueOf(time.Now()).Type()

// read buffer reader than decode and set it to v.
// v is struct or pointer type that will contain decoded data
func Decode(rdr *bytes.Reader, v interface{}) error {
	var err error

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		log.Fatal(rv.Kind())
		return errors.New("type not match")
	}

	for i := 0; i < rv.NumField() && rdr.Len() > 0; i++ {
		rvField := rv.Field(i)
		rtField := rv.Type().Field(i)

		if rvField.IsValid() && rvField.CanSet() {
			tag := DeTag(rtField.Tag, rvField.Kind())

			switch rk := rvField.Kind(); rk {

			case reflect.Ptr:
				rvField.Set(reflect.New(rvField.Type().Elem()))
				if err = Decode(rdr, rvField.Interface()); err != nil {
					return err
				}

			case reflect.Struct:
				// if data type is time.Time
				if rvField.Type() == typeOfTime {
					b := make([]byte, tag.Len)
					binary.Read(rdr, binary.LittleEndian, &b)
					rvField.Set(reflect.ValueOf(parseTime(b)))
				} else {
					if err = Decode(rdr, rvField.Addr().Interface()); err != nil {
						return err
					}
				}

			case reflect.Array:
				for j := 0; j < rvField.Len(); j++ {
					if err = Decode(rdr, rvField.Index(j).Addr().Interface()); err != nil {
						return err
					}
				}

			case reflect.String:
				x := make([]byte, tag.Len)
				binary.Read(rdr, binary.LittleEndian, &x)
				rvField.SetString(parseString(x))

			case reflect.Bool:
				var x bool
				binary.Read(rdr, binary.LittleEndian, &x)
				rvField.SetBool(x)

			case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
				x := readUint(rdr, tag.Len)
				if !rvField.OverflowUint(x) {
					rvField.SetUint(x)
				}

			case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
				x := readUint(rdr, tag.Len)
				x_int := int64(x)
				if !rvField.OverflowInt(x_int) {
					rvField.SetInt(x_int)
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

				if !rvField.OverflowFloat(x64) {
					rvField.SetFloat(x64)
				}

			default:
				return errors.New("unsupported kind: " + rv.Kind().String())
			}
		}
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

// parseTime convert bytes slice (little endian) to time
// value of bytes data is :
//   1 byte of year
//   1 byte of month
//   1 byte of day
//   1 byte of hour
//   1 byte of minute
//   1 byte of second
//   1 byte of weekday (ignored)
func parseTime(b []byte) time.Time {
	var data string
	for _, v := range b[:6] {
		data += fmt.Sprintf("%02d", uint8(v))
	}

	datetime, _ := time.Parse("060102150405", data)
	return datetime
}

// parseString convert byte slice (little endian) to string
func parseString(b []byte) string {
	return string(util.Reverse(b))
}
