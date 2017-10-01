package export

import (
	"reflect"
	"strconv"
)

type stringer interface {
	String() string
}

func structToJson(d []byte, indent []byte, rv reflect.Value) []byte {
	d = append(d, '{')
	if rv.NumField() > 0 {
		indent = append(indent, "    "...)
		d = append(d, '\n')
		rt := rv.Type()
		for i := 0; i < rv.NumField(); i++ {
			if i != 0 {
				d = append(d, ',', '\n')
			}
			d = append(d, indent...)
			d = strconv.AppendQuote(d, rt.Field(i).Name)
			d = append(d, ':', ' ')
			d = valueToJson(d, indent, rv.Field(i))
		}
		d = append(d, '\n')
		d = append(d, indent[:len(indent)-4]...)
	}
	d = append(d, '}')
	return d
}

func arrayToJson(d []byte, indent []byte, rv reflect.Value) []byte {
	d = append(d, '[')
	if rv.Len() > 0 {
		indent = append(indent, "    "...)
		d = append(d, '\n')
		d = append(d, indent...)
		for i := 0; i < rv.Len(); i++ {
			if i != 0 {
				d = append(d, ',', ' ')
			}
			d = valueToJson(d, indent, rv.Index(i))
		}
		d = append(d, '\n')
		d = append(d, indent[:len(indent)-4]...)
	}
	d = append(d, ']')
	return d
}

func ifaceToJson(d []byte, indent []byte, e interface{}) ([]byte, bool) {
	if s, ok := e.(stringer); ok {
		return strconv.AppendQuote(d, s.String()), true
	}
	switch v := e.(type) {
	case int64:
		d = strconv.AppendInt(d, v, 10)
	case uint64:
		d = strconv.AppendUint(d, v, 10)
	case uint32:
		d = strconv.AppendUint(d, uint64(v), 10)
	case float64:
		d = strconv.AppendFloat(d, v, 'g', -1, 64)
	case bool:
		d = strconv.AppendBool(d, v)
	default:
		return d, false
	}
	return d, true
}

func valueToJson(d []byte, indent []byte, rv reflect.Value) []byte {
	if rv.CanInterface() {
		var ok bool
		d, ok = ifaceToJson(d, indent, rv.Interface())
		if ok {
			return d
		}
	}
	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			d = append(d, `nil`...)
		} else {
			d = strconv.AppendUint(d, uint64(rv.Pointer()), 10)
		}
	case reflect.Struct:
		d = structToJson(d, indent, rv)
	case reflect.Slice, reflect.Array:
		d = arrayToJson(d, indent, rv)
	default:
		d = append(d, `"kind: `...)
		d = strconv.AppendUint(d, uint64(rv.Kind()), 10)
		d = append(d, ` type:`...)
		d = append(d, rv.Type().String()...)
		d = append(d, '"')
	}
	return d
}
