package keypath

import (
	"github.com/go-courier/reflectx"
	"github.com/go-courier/reflectx/typesutil"
	"go/ast"
	"reflect"
	"strings"
)

func NewKeyPathDecoder(values map[string]string) *KeyPathDecoder {
	return &KeyPathDecoder{
		values: values,
	}
}

type KeyPathDecoder struct {
	values map[string]string
}

func (d *KeyPathDecoder) Decode(v interface{}) error {
	walker := NewPathWalker()
	rv, ok := v.(reflect.Value)
	if !ok {
		rv = reflect.ValueOf(v)
	}
	return d.scanAndSetValue(walker, rv)
}

func (d *KeyPathDecoder) scanAndSetValue(walker *PathWalker, rv reflect.Value) error {
	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			p := walker.String()

			// only sub path new empty
			for key := range d.values {
				if strings.HasPrefix(key, p) {
					rv.Set(reflectx.New(rv.Type()))
					return d.scanAndSetValue(walker, rv.Elem())
				}
			}
			return nil
		}
		return d.scanAndSetValue(walker, rv.Elem())
	case reflect.Func, reflect.Interface, reflect.Chan:
		// skip
	default:
		typ := rv.Type()
		if _, ok := typesutil.EncodingTextMarshalerTypeReplacer(typesutil.FromRType(typ)); ok {
			if v, ok := d.values[walker.String()]; ok {
				if err := reflectx.UnmarshalText(rv, []byte(v)); err != nil {
					return err
				}
			}
			return nil
		}

		switch rv.Kind() {
		case reflect.Map:
			p := walker.String()

			for key := range d.values {
				if strings.HasPrefix(key, p) {
					k := strings.Split(strings.TrimLeft(key, p+"."), ".")[0]

					rvKey := reflect.ValueOf(k)

					for _, kk := range rv.MapKeys() {
						if k == kk.String() {
							rvKey = kk
						}
					}

					v := reflectx.New(rv.Type().Elem())

					mV := rv.MapIndex(rvKey)

					if mV.IsValid() {
						v.Set(mV)
					}

					walker.Enter(k)
					err := d.scanAndSetValue(walker, v)
					if err != nil {
						return err
					}
					walker.Exit()
					rv.SetMapIndex(rvKey, v)
				}
			}

		case reflect.Array, reflect.Slice:
			for i := 0; i < rv.Len(); i++ {
				walker.Enter(i)
				if err := d.scanAndSetValue(walker, rv.Index(i)); err != nil {
					return err
				}
				walker.Exit()
			}

		case reflect.Struct:
			tpe := rv.Type()
			for i := 0; i < rv.NumField(); i++ {
				field := tpe.Field(i)

				flags := (map[string]bool)(nil)
				name := field.Name

				if !ast.IsExported(name) {
					continue
				}

				if tag, ok := field.Tag.Lookup("json"); ok {
					n, fs := tagValueAndFlags(tag)
					if n == "-" {
						continue
					}
					if n != "" {
						name = n
					}
					flags = fs
				}

				inline := flags == nil && reflectx.Deref(field.Type).Kind() == reflect.Struct && field.Anonymous

				if !inline {
					walker.Enter(name)
				}

				if err := d.scanAndSetValue(walker, rv.Field(i)); err != nil {
					return err
				}

				if !inline {
					walker.Exit()
				}
			}
		default:
			if v, ok := d.values[walker.String()]; ok {
				if err := reflectx.UnmarshalText(rv, []byte(v)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func tagValueAndFlags(tagString string) (string, map[string]bool) {
	valueAndFlags := strings.Split(tagString, ",")
	v := valueAndFlags[0]
	tagFlags := map[string]bool{}
	if len(valueAndFlags) > 1 {
		for _, flag := range valueAndFlags[1:] {
			tagFlags[flag] = true
		}
	}
	return v, tagFlags
}
