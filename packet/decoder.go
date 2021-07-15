package packet

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
	"reflect"
	"strconv"
)

func Decode(buf io.Reader, v interface{}, buflen int) (int, error) {
	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return 0, errors.New("Type not match")
	}

	var bufRead = 0
	for i := 0; i < rv.NumField() && bufRead < buflen; i++ {
		var len int = 0
		var factor float64 = 1
		rvField := rv.Field(i)
		typeField := rv.Type().Field(i)
		tag := typeField.Tag
		if tagLen := tag.Get("len"); tagLen != "" {
			tmplen, _ := strconv.ParseInt(tagLen, 10, 64)
			len = int(tmplen)
		}
		if tagF := tag.Get("factor"); tagF != "" {
			factor, _ = strconv.ParseFloat(tagF, 64)
		}

		// fmt.Println(buflen, typeField.Name, rvField.Kind())
		if rvField.IsValid() {
			if rvField.CanSet() {
				switch rk := rvField.Kind(); rk {

				case reflect.Ptr:
					rvField.Set(reflect.New(rvField.Type().Elem()))
					readLen, err := Decode(buf, rvField.Interface(), (buflen - bufRead))
					bufRead += readLen
					if err != nil { return bufRead, err }

				case reflect.Struct:
					readLen, err := Decode(buf, rvField.Addr().Interface(), (buflen - bufRead))
					bufRead += readLen
					if err != nil { return bufRead, err }

				case reflect.Array:
					for j := 0; j < rvField.Len(); j++ {					
						readLen, err := Decode(buf, rvField.Index(j).Addr().Interface(), (buflen - bufRead))
						bufRead += readLen
						if err != nil { return bufRead, err }
					}

				case reflect.String:
					x := make([]byte, len)
					binary.Read(buf, binary.LittleEndian, &x)
					rvField.SetString(string(x))
					bufRead += len

				case reflect.Bool:
					var x bool
					binary.Read(buf, binary.LittleEndian, &x)
					rvField.SetBool(x)
					bufRead += 1

				case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
					x, readLen := readUint(buf, rk, len)
					if !rvField.OverflowUint(x) {
						rvField.SetUint(x)
					}
					bufRead += readLen

				case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
					x, readLen := readUint(buf, rk, len)
					x_int := int64(x)
					if !rvField.OverflowInt(x_int) {
						rvField.SetInt(x_int)
					}
					bufRead += readLen

				case reflect.Float32:
					var x32 float32
					var x64 float64
					x, readLen := readUint(buf, rk, len)
					n := uint32(x)
					if factor != 1 {
						x64 = float64(n) * factor

					} else {
						x32 = math.Float32frombits(n)
						x64 = float64(x32)
					}
					if !rvField.OverflowFloat(x64) {
						rvField.SetFloat(x64)
					}
					bufRead += readLen

				case reflect.Float64:
					var x64 float64
					x, readLen := readUint(buf, rk, len)
					if factor == 1 {
						x64 = float64(x) * factor

					} else {
						x64 = math.Float64frombits(x)
					}
					if !rvField.OverflowFloat(x64) {
						rvField.SetFloat(x64)
					}
					bufRead += readLen

				default:
				}
			}
		}
	}
	return bufRead, nil
}

func readUint(buf io.Reader, rk reflect.Kind, len int) (uint64, int) {
	var maxlen int = 1
	switch rk {
	case reflect.Uint8, reflect.Int8:
		maxlen = 1
	case reflect.Uint16, reflect.Int16:
		maxlen = 2
	case reflect.Uint32, reflect.Int32, reflect.Float32:
		maxlen = 4
	case reflect.Uint64, reflect.Uint, reflect.Int64, reflect.Int, reflect.Float64:
		maxlen = 8
	}
	if len == 0 {
		len = maxlen
	}
	b := make([]byte, len)
	newb := make([]byte, 8)
	binary.Read(buf, binary.BigEndian, &b)
	for i := 0; i < len; i++ {
		newb[i] = b[i]
	}
	return binary.LittleEndian.Uint64(newb), len
}