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

var typeOfTime reflect.Type = reflect.ValueOf(time.Now()).Type()

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
					x64 = convert2Float64(tag.Tipe, x)
					// it looks like same operation ?
					if rk == reflect.Float32 {
						x64 *= tag.Factor
					} else {
						x64 *= tag.Factor
					}

				} else {
					// set sesuai biner
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

func readUint(rdr io.Reader, len int) uint64 {
	b := make([]byte, len)
	binary.Read(rdr, binary.BigEndian, &b)

	newb := make([]byte, 8)
	for i := 0; i < len; i++ {
		newb[i] = b[i]
	}

	return binary.LittleEndian.Uint64(newb)
}

func convert2Float64(typedata string, x uint64) (result float64) {
	// deklarasi
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

func parseTime(b []byte) time.Time {
	var data string
	for _, v := range b[:6] {
		data += fmt.Sprintf("%02d", uint8(v))
	}

	datetime, _ := time.Parse("060102150405", data)
	return datetime
}

func parseString(b []byte) string {
	return string(util.Reverse(b))
}