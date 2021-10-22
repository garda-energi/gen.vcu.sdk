package sdk

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"
)

// decodeResponse extract header and message response from bytes packet.
func decodeResponse(packet packet) (*responsePacket, error) {
	// TODO: redundant with decodeReport logic, implement DRY
	reader := bytes.NewReader(packet)
	result := &responsePacket{}
	if err := decode(reader, result); err != nil {
		return nil, err
	}

	if !result.validPrefix() {
		return nil, errInvalidPrefix
	}
	if !result.validSize() {
		return nil, errInvalidSize
	}
	if reader.Len() != 0 {
		return nil, errInvalidSize
	}

	if !result.validCmdCode() {
		return nil, errInvalidCmdCode
	}
	if !result.validResCode() {
		return nil, errInvalidResCode
	}
	return result, nil
}

// decodeReport extract report from bytes packet.
func decodeReport(packet packet) (*ReportPacket, error) {
	reportPacket := &ReportPacket{}

	// get version
	reader := bytes.NewReader(packet)
	decode(reader, reportPacket)
	rpStructure, isGot := ReportPacketStructures[int(reportPacket.Header.Version)]
	if !isGot {
		return nil, errors.New(fmt.Sprintf("Cannot handle report version %d", reportPacket.Header.Version))
	}

	// Check validity
	if !reportPacket.ValidPrefix() {
		return nil, errInvalidPrefix
	}

	// decode payload
	payloadReader := bytes.NewReader(reportPacket.Payload)
	if err := decode(payloadReader, &reportPacket.Data, rpStructure); err != nil {
		return nil, err
	}

	// check length
	if payloadReader.Len() != 0 {
		return nil, errInvalidSize
	}

	return reportPacket, nil
}

// decode read buffer reader than decode and set it to v.
// v is struct or pointer type that will contain decoded data.
// fix bug. in "decode func" before, It'll be error for case such as v isn't array of struct.
func decode(rdr *bytes.Reader, v interface{}, tags ...tagger) error {
	var err error
	var isVMaps bool = false
	var mapElm reflect.Value // map rv
	var rv reflect.Value

	rv = reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if !rv.IsValid() || (rv.Kind() != reflect.Map && !rv.CanSet()) {
		return nil
	}

	// default tag
	tag := newTagger()
	// if tag is pass in argument
	if len(tags) > 0 {
		tag = tags[0]
	}

	// if v is map
	if rv.Kind() == reflect.Map && rv.Type() == typeOfPacketData {
		if tag.Tipe == "" {
			return errors.New("Decode to map: need tag tipe.")
		}

		mapElm = rv

		switch tag.Tipe {
		case Struct_t:
			if mapElm.IsNil() && mapElm.CanAddr() {
				mapElm.Set(reflect.MakeMap(typeOfPacketData))
			}
			break
		case Array_t:
			if len(tag.Sub) == 0 {
				return errors.New("Tag (" + tag.Name + "): Tipe Array_t cannot be 0")
			}
			_, contentType := createZeroElm(tag.Sub[0].Tipe)
			rv = reflect.MakeSlice(reflect.SliceOf(contentType), tag.Len, tag.Len)
			mapElm.SetMapIndex(reflect.ValueOf(tag.Name), rv)
			break
		default:
			rv, _ = createZeroElm(tag.Tipe)
		}
		isVMaps = true
	}

	switch rk := rv.Kind(); rk {
	case reflect.Ptr:
		rv.Set(reflect.New(rv.Type().Elem()))
		if err = decode(rdr, rv.Interface()); err != nil {
			return err
		}

	case reflect.Struct:
		// if data type is time.Time
		if rv.Type() == typeOfTime {
			b := make([]byte, tag.Len)
			if err = binary.Read(rdr, binary.LittleEndian, &b); err != nil {
				return err
			}
			if isVMaps {
				mapElm.SetMapIndex(reflect.ValueOf(tag.Name), reflect.ValueOf(bytesToTime(b)))
			} else {
				rv.Set(reflect.ValueOf(bytesToTime(b)))
			}

		} else {
			for i := 0; i < rv.NumField() && rdr.Len() > 0; i++ {
				rvField := rv.Field(i)
				rtField := rv.Type().Field(i)

				tagField := deTag(rtField.Tag, rvField.Kind())
				if err = decode(rdr, rvField.Addr().Interface(), tagField); err != nil {
					return err
				}
			}
		}

	case reflect.Array:
		for j := 0; j < rv.Len(); j++ {
			if err = decode(rdr, rv.Index(j).Addr().Interface()); err != nil {
				return err
			}
		}

	case reflect.Slice:
		if isVMaps {
			for j := 0; j < rv.Len(); j++ {
				if err = decode(rdr, rv.Index(j).Addr().Interface(), tag.Sub[0]); err != nil {
					return err
				}
			}

		} else {
			if rv.Type() == typeOfMessage {
				x := make(message, rdr.Len())
				if err = binary.Read(rdr, binary.LittleEndian, &x); err != nil {
					return err
				}
				rv.Set(reflect.ValueOf(x))
			}
		}

	case reflect.String:
		x := make([]byte, tag.Len)
		if err = binary.Read(rdr, binary.LittleEndian, &x); err != nil {
			return err
		}
		if isVMaps {
			mapElm.SetMapIndex(reflect.ValueOf(tag.Name), reflect.ValueOf(bytesToStr(x)))
		} else {
			rv.SetString(bytesToStr(x))
		}

	case reflect.Bool:
		var x bool
		if err = binary.Read(rdr, binary.LittleEndian, &x); err != nil {
			return err
		}
		if isVMaps {
			mapElm.SetMapIndex(reflect.ValueOf(tag.Name), reflect.ValueOf(x))
		} else {
			rv.SetBool(x)
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		x, err := readUint(rdr, tag.Len)
		if err != nil {
			return err
		}
		if !rv.OverflowUint(x) {
			if isVMaps {
				var rvx reflect.Value // reflect value of x
				switch rk {
				case reflect.Uint8:
					rvx = reflect.ValueOf(uint8(x))
				case reflect.Uint16:
					rvx = reflect.ValueOf(uint16(x))
				case reflect.Uint32:
					rvx = reflect.ValueOf(uint32(x))
				case reflect.Uint64:
					rvx = reflect.ValueOf(uint64(x))
				case reflect.Uint:
					rvx = reflect.ValueOf(uint(x))
				}
				mapElm.SetMapIndex(reflect.ValueOf(tag.Name), rvx)
			} else {
				rv.SetUint(x)
			}
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		x, err := readUint(rdr, tag.Len)
		if err != nil {
			return err
		}
		if !rv.OverflowInt(int64(x)) {
			if isVMaps {
				var rvx reflect.Value // reflect value of x
				switch rk {
				case reflect.Int8:
					rvx = reflect.ValueOf(int8(x))
				case reflect.Int16:
					rvx = reflect.ValueOf(int16(x))
				case reflect.Int32:
					rvx = reflect.ValueOf(int32(x))
				case reflect.Int64:
					rvx = reflect.ValueOf(int64(x))
				case reflect.Int:
					rvx = reflect.ValueOf(int(x))
				}
				mapElm.SetMapIndex(reflect.ValueOf(tag.Name), rvx)
			} else {
				rv.SetInt(int64(x))
			}
		}

	case reflect.Float32, reflect.Float64:
		var x64 float64
		var x_int int64
		var x_uint uint64
		var err error

		switch tag.UnfactorType {
		case Int8_t, Int16_t, Int32_t, Int64_t:
			x_int, err = readInt(rdr, tag.Len)
			x_uint = uint64(x_int)
		default:
			x_uint, err = readUint(rdr, tag.Len)
		}

		if err != nil {
			return err
		}

		if tag.Factor != 1 {
			tagType := tag.Tipe
			if tag.UnfactorType != "" {
				tagType = tag.UnfactorType
			}
			x64 = convertToFloat64(tagType, x_uint)
			x64 *= tag.Factor

		} else {
			// set as binary
			if rk == reflect.Float32 {
				x32 := math.Float32frombits(uint32(x_uint))
				x64 = float64(x32)
			} else {
				x64 = math.Float64frombits(x_uint)
			}
		}

		if !rv.OverflowFloat(x64) {
			if isVMaps {
				mapElm.SetMapIndex(reflect.ValueOf(tag.Name), reflect.ValueOf(float32(x64)))
			} else {
				rv.SetFloat(x64)
			}
		}

	case reflect.Map:
		if isVMaps {
			for _, tagSub := range tag.Sub {
				tagSub = tagSub.normalize()
				if rdr.Len() == 0 {
					break
				}
				if tagSub.Tipe == Struct_t {
					rv = reflect.MakeMap(typeOfPacketData)
					mapElm.SetMapIndex(reflect.ValueOf(tagSub.Name), rv)
					if err = decode(rdr, rv.Interface(), tagSub); err != nil {
						return err
					}
				} else {
					if err = decode(rdr, mapElm.Interface(), tagSub); err != nil {
						return err
					}
				}
			}
		}

	case reflect.Interface:
		if err = decode(rdr, rv.Elem(), tag); err != nil {
			return err
		}

	default:
		return errors.New("unsupported kind: " + rv.Kind().String())
	}

	return err
}

// createZeroElm create zero/nil variable by VarDataType
func createZeroElm(dt VarDataType) (reflect.Value, reflect.Type) {
	var rt reflect.Type

	switch dt {
	case Boolean_t:
		rt = reflect.TypeOf(false)
	case Uint8_t:
		rt = reflect.TypeOf(uint8(0))
	case Uint16_t:
		rt = reflect.TypeOf(uint16(0))
	case Uint32_t:
		rt = reflect.TypeOf(uint32(0))
	case Uint64_t:
		rt = reflect.TypeOf(uint64(0))
	case Int8_t:
		rt = reflect.TypeOf(int8(0))
	case Int16_t:
		rt = reflect.TypeOf(int16(0))
	case Int32_t:
		rt = reflect.TypeOf(int32(0))
	case Int64_t:
		rt = reflect.TypeOf(int64(0))
	case Float_t:
		rt = reflect.TypeOf(float32(0))
	case Time_t:
		rt = typeOfTime
	case Struct_t:
		rt = typeOfPacketData
	}

	return reflect.Zero(rt), rt
}

// readUint read len(length) data as uint64
func readUint(rdr io.Reader, len int) (uint64, error) {
	// sometimes, data recived in length less than 8
	b := make([]byte, len)
	err := binary.Read(rdr, binary.BigEndian, &b)
	if err != nil {
		return 0, err
	}

	newb := make([]byte, 8)
	for i := 0; i < len; i++ {
		newb[i] = b[i]
	}
	return binary.LittleEndian.Uint64(newb), nil
}

// readInt read len(length) data as int64 (signed int)
func readInt(rdr io.Reader, len int) (int64, error) {
	// sometimes, data recived in length less than 8
	b := make([]byte, len)
	err := binary.Read(rdr, binary.BigEndian, &b)
	if err != nil {
		return 0, err
	}

	var result int64 = 0

	switch len {
	case 1:
		result = int64(int8(b[0]))
	case 2:
		result = int64(int16(binary.LittleEndian.Uint16(b)))
	case 3, 4:
		result = int64(int32(binary.LittleEndian.Uint32(b)))
	default:
		result = int64(binary.LittleEndian.Uint64(b))
	}

	return result, nil
}

// convertToFloat64 convert bytes data to float64.
// data is read as typedata from tag.
func convertToFloat64(typedata VarDataType, x uint64) (result float64) {
	rv := setVarOfTypeData(typedata)

	// set ad declared memory size
	switch rv.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		rv.SetUint(x)
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		rv.SetInt(int64(x))
	}

	// convert to float
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
	return string(reverseBytes(b))
}
