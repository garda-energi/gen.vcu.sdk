package sdk

import (
	"reflect"
	"strconv"
)

type tagger struct {
	Tipe   string
	Len    int
	Factor float64
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
		t.Tipe = tipe
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
