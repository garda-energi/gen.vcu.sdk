package packet

import (
	"errors"
	"reflect"
	"strconv"
)

// type M map[string]interface{}

type Tag struct {
	Tipe      string
	Unit      string
	Scale     float32
	Chartable bool
}

func TagWalk(v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return errors.New("not a pointer value")
	}

	v = reflect.Indirect(v)

	switch v.Kind() {
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if err := TagWalk(v.Field(i).Addr()); err != nil {
				return err
			}
		}
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if err := TagWalk(v.Index(i).Addr()); err != nil {
				return err
			}
		}
	case reflect.Uint8:
		v.SetUint(10)
	case reflect.Int8:
		v.SetInt(-50)
	case reflect.String:
		v.SetString("Foo")
	case reflect.Bool:
		v.SetBool(true)

	default:
		// return errors.New("Unsupported kind: " + v.Kind().String())
	}

	// for i := 0; i < v.NumField(); i++ {
	// 	value := v.Field(i)
	// 	t := v.Type().Field(i)

	// 	// if _, ok := tipe.Tag.Lookup("type"); !ok {
	// 	if t.Type.Kind() == reflect.Struct {
	// 		tagStructWalk(value.Addr())
	// 		continue
	// 	}

	// 	tag := getTags(t)
	// 	// key := tag.Group + "." + strings.ToLower(t.Name)
	// 	val := value.Interface()

	// 	fmt.Println(t.Name, val, tag)
	// 	// switch tag.Tipe {
	// 	// case "string":
	// 	// 	val = val.([]byte)
	// 	// case "datetime":
	// 	// 	val = formatter.ToUnixTime(val.([]byte))
	// 	// }
	// 	// buf[key] = val
	// }

	return nil
}

func getTags(field reflect.StructField) Tag {
	tag := Tag{}

	// if group, ok := field.Tag.Lookup("group"); ok {
	// 	tag.Group = group
	// }
	if tipe, ok := field.Tag.Lookup("type"); ok {
		tag.Tipe = tipe
	}
	if unit, ok := field.Tag.Lookup("unit"); ok {
		tag.Unit = unit
	}
	if factor, ok := field.Tag.Lookup("factor"); ok {
		if f, err := strconv.ParseFloat(factor, 32); err == nil {
			tag.Scale = float32(f)
		}
	}
	if _, ok := field.Tag.Lookup("chartable"); ok {
		tag.Chartable = true
	}

	return tag
}
