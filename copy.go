package deep

import (
	"reflect"
	"time"
	"unsafe"
)

// Equal deeply compares the arguments y and x.
func Equal[T any](x, y T) bool {
	return reflect.DeepEqual(x, y)
}

// Copy deeply data from src to dst.
func Copy[T any](src T) (dst T) {
	dval := reflect.ValueOf(&dst)
	sval := reflect.ValueOf(&src)
	c := copier{seen: make(map[uintptr]reflect.Value)}
	c.copyValue(dval.Elem(), sval.Elem())
	return dst
}

// copier tracks pointers already copied, by address, so cyclic data
// structures terminate instead of recursing forever.
type copier struct {
	seen map[uintptr]reflect.Value
}

func (c copier) copyValue(dst, src reflect.Value) {
	if !dst.CanSet() {
		return
	}
	switch src.Type() {
	case reflect.TypeFor[time.Time]():
		dst.Set(src)
		return
	}
	switch src.Kind() {
	case reflect.Slice:
		if src.IsNil() {
			return
		}
		nval := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		for i := range src.Len() {
			c.copyValue(nval.Index(i), src.Index(i))
		}
		dst.Set(nval)
	case reflect.Array:
		nval := reflect.New(src.Type()).Elem()
		for i := range src.Len() {
			c.copyValue(nval.Index(i), src.Index(i))
		}
		dst.Set(nval)
	case reflect.Interface:
		if src.IsNil() {
			return
		}
		nval := reflect.New(src.Elem().Type()).Elem()
		c.copyValue(nval, src.Elem())
		dst.Set(nval)
	case reflect.Pointer:
		if src.IsNil() {
			return
		}
		addr := src.Pointer()
		if nval, ok := c.seen[addr]; ok {
			dst.Set(nval)
			return
		}
		nval := reflect.New(src.Elem().Type())
		c.seen[addr] = nval
		c.copyValue(nval.Elem(), src.Elem())
		dst.Set(nval)
	case reflect.Struct:
		nval := reflect.New(src.Type()).Elem()
		for i := range src.NumField() {
			df, sf := nval.Field(i), src.Field(i)
			// Unexported field: bypass reflect's read-only flag via unsafe
			if !df.CanSet() {
				df = reflect.NewAt(df.Type(), unsafe.Pointer(df.UnsafeAddr())).Elem()
				sf = reflect.NewAt(sf.Type(), unsafe.Pointer(sf.UnsafeAddr())).Elem()
			}
			c.copyValue(df, sf)
		}
		dst.Set(nval)
	case reflect.Map:
		if src.IsNil() {
			return
		}
		nval := reflect.MakeMap(src.Type())
		for _, k := range src.MapKeys() {
			sval := src.MapIndex(k)
			dval := reflect.New(sval.Type()).Elem()
			c.copyValue(dval, sval)
			nval.SetMapIndex(k, dval)
		}
		dst.Set(nval)
	default:
		dst.Set(src)
	}
}
