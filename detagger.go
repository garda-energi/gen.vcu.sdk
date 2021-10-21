package sdk

import (
	"reflect"
	"strconv"
)

type tagger struct {
	Name         string
	Tipe         VarDataType
	Len          int
	Factor       float64
	UnfactorType VarDataType
	Sub          []tagger
}

func newTagger() tagger {
	return tagger{
		Len:    1,
		Factor: 1.0,
		Tipe:   "uint64",
	}
}

// deTag decode tagger from struct field.
func deTag(tag reflect.StructTag, rk reflect.Kind) tagger {
	t := newTagger()

	if factor, ok := tag.Lookup("factor"); ok {
		v, _ := strconv.ParseFloat(factor, 64)
		t.Factor = v
	}

	if tipe, ok := tag.Lookup("type"); ok {
		t.Tipe = VarDataType(tipe)
	}

	if len, ok := tag.Lookup("len"); ok {
		v, _ := strconv.ParseInt(len, 10, 64)
		t.Len = int(v)
	} else {
		switch rk {
		case reflect.Uint8, reflect.Int8:
			t.Len = 1
		case reflect.Uint16, reflect.Int16:
			t.Len = 2
		case reflect.Uint32, reflect.Int32, reflect.Float32:
			t.Len = 4
		case reflect.Uint64, reflect.Uint, reflect.Int64, reflect.Int, reflect.Float64:
			t.Len = 8
		}
	}

	return t
}

// normalize normalization tagger from uncomplete tag
func (tag tagger) normalize() tagger {

	if tag.Factor == 0 {
		tag.Factor = 1.0
	}

	if tag.Len == 0 {
		switch tag.Tipe {
		case Uint8_t, Int8_t, Boolean_t:
			tag.Len = 1
		case Uint16_t, Int16_t:
			tag.Len = 2
		case Uint32_t, Int32_t, Float_t:
			tag.Len = 4
		case Uint64_t, Int64_t:
			tag.Len = 8
		}
	}

	return tag
}

// getSize calculate tagger and sub tagger size
func (tag tagger) getSize() int {
	length := 0
	switch tag.Tipe {
	case Struct_t:
		for _, sub := range tag.Sub {
			length += sub.getSize()
		}
	case Array_t:
		if len(tag.Sub) != 0 {
			length += tag.Len * tag.Sub[0].getSize()
		}
	case Float_t:
		if tag.Len != 0 {
			length += tag.Len
		} else {
			length += 4
		}
	case Time_t:
		length += 7
	case Boolean_t, Uint8_t, Int8_t:
		length += 1
	case Uint16_t, Int16_t:
		length += 2
	case Uint32_t, Int32_t:
		length += 4
	case Uint64_t, Int64_t:
		length += 8
	}
	return length
}

// getPacketSize calculate packet size
func getPacketSize(v interface{}) int {
	size := 0
	rv := reflect.ValueOf(v)

	for rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}

	switch rk := rv.Kind(); rk {

	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			rvField := rv.Field(i)
			rtField := rv.Type().Field(i)

			tagField := deTag(rtField.Tag, rvField.Kind())
			if rvField.Type() == typeOfTime {
				size += tagField.Len
			} else if rk := rvField.Kind(); rk == reflect.Struct || rk == reflect.Array || rk == reflect.Ptr || rk == reflect.Slice {
				size += getPacketSize(rvField.Addr().Interface())
			} else {
				size += tagField.Len
			}
		}

	case reflect.Array, reflect.Slice:
		if rv.Type() == typeOfMessage {
			size += rv.Len()
		} else {
			for i := 0; i < rv.Len(); i++ {
				size += getPacketSize(rv.Index(i).Addr().Interface())
			}
		}

	default:
	}

	return size
}
