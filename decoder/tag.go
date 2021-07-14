package decoder

import (
	"bytes"
	"errors"
	"reflect"
	"strconv"
)

type Meta struct {
	Tipe   string
	Len    int
	Factor float32
	Unit   string
}

func TagWalk(rdr *bytes.Reader, v reflect.Value, t reflect.StructTag) error {
	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer value")
	}

	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag

			if err := TagWalk(rdr, v.Field(i).Addr(), tag); err != nil {
				return err
			}
		}
		return nil
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := TagWalk(rdr, v.Index(i).Addr(), ""); err != nil {
				return err
			}
		}
		return nil
	}

	err := decodePacket(rdr, v, t)
	return err
}

func decodePacket(rdr *bytes.Reader, v reflect.Value, t reflect.StructTag) error {
	if t == "" {
		return errors.New("no meta defined")
	}

	meta := decodeMeta(t)
	buf := make([]byte, meta.Len)
	rdr.Read(buf)

	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(toUint64(buf))
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var val int64 = toInt64(buf)
		if meta.Tipe == "unix_time" {
			val = toUnixTime(buf)
		}
		v.SetInt(val)
	case reflect.Float32:
		v.SetFloat(toFloat64(buf, meta.Factor))
	case reflect.Bool:
		v.SetBool(toBool(buf))
	case reflect.String:
		v.SetString(toAscii(buf))

	default:
		return errors.New("unsupported kind: " + v.Kind().String())
	}

	return nil
}

func decodeMeta(tag reflect.StructTag) Meta {
	meta := Meta{
		Factor: 1,
	}

	if tipe, ok := tag.Lookup("type"); ok {
		meta.Tipe = tipe
	}

	if len, ok := tag.Lookup("len"); ok {
		if i, err := strconv.Atoi(len); err == nil {
			meta.Len = int(i)
		}
	} else {
		switch meta.Tipe {
		case reflect.Uint8.String(), reflect.Int8.String():
			meta.Len = 1
		case reflect.Uint16.String(), reflect.Int16.String():
			meta.Len = 2
		case reflect.Uint32.String(), reflect.Int32.String():
			meta.Len = 4
		case reflect.Uint64.String(), reflect.Int64.String():
			meta.Len = 8
		}
	}

	if factor, ok := tag.Lookup("factor"); ok {
		if f, err := strconv.ParseFloat(factor, 32); err == nil {
			meta.Factor = float32(f)
		}
	}

	if unit, ok := tag.Lookup("unit"); ok {
		meta.Unit = unit
	}

	return meta
}
