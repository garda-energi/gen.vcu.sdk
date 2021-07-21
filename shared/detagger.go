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

func DeTag(tag reflect.StructTag, rk reflect.Kind) Tagger {
	t := Tagger{
		Factor: 1.0,
		Tipe:   "uint64",
	}

	t.Len = DeTagLen(tag, rk)

	if factor, ok := tag.Lookup("factor"); ok {
		v, _ := strconv.ParseFloat(factor, 64)
		t.Factor = v
	}

	if tipe, ok := tag.Lookup("type"); ok {
		t.Tipe = tipe
	}

	return t
}

func DeTagLen(tag reflect.StructTag, rk reflect.Kind) int {
	len := 1

	if l, ok := tag.Lookup("len"); ok {
		v, _ := strconv.ParseInt(l, 10, 64)
		len = int(v)
	} else {
		switch rk {
		case reflect.Uint8, reflect.Int8:
			len = 1
		case reflect.Uint16, reflect.Int16:
			len = 2
		case reflect.Uint32, reflect.Int32, reflect.Float32:
			len = 4
		case reflect.Uint64, reflect.Uint, reflect.Int64, reflect.Int, reflect.Float64:
			len = 8
		}
	}

	return len
}
