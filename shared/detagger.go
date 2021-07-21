package shared

import (
	"reflect"
	"strconv"
)

type Tagger struct {
	Tipe   string
	Len    int
	Factor float64
}

func deTag(tag reflect.StructTag, rk reflect.Kind) Tagger {
	t := Tagger{
		Len:    1,
		Factor: 1.0,
		Tipe:   "uint64",
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

	if factor, ok := tag.Lookup("factor"); ok {
		v, _ := strconv.ParseFloat(factor, 64)
		t.Factor = v
	}

	if tipe, ok := tag.Lookup("type"); ok {
		t.Tipe = tipe
	}

	return t
}
