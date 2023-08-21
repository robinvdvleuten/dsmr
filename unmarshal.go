package dsmr

import (
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/participle/v2"
)

// Unmarshal telegram into a Go struct.
func Unmarshal(data []byte, v any) error {
	ast, err := ParseBytes(data)
	if err != nil {
		return err
	}

	return UnmarshalAST(ast, v)
}

// UnmarshalAST unmarshalls an already parsed or constructed AST into a Go struct.
func UnmarshalAST(ast *AST, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		return fmt.Errorf("can't unmarshal into nil")
	}

	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("can only unmarshal into a pointer to a struct, not %s", rv.Type())
	}

	return unmarshalEntries(rv.Elem(), ast.entries())
}

func unmarshalEntries(v reflect.Value, entries []Entry) error {
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("%s must be a struct", v.Type())
	}

	// Collect entries from the source into a map.
	seen := map[string]Entry{}
	mentries := make(map[string][]Entry, len(entries))
	for _, entry := range entries {
		key := entry.Key()
		mentries[key] = append(mentries[key], entry)
		seen[key] = entry
	}

	// Collect the fields of the target struct.
	fields, err := flattenFields(v)
	if err != nil {
		return err
	}

	// Apply telegram entries to our fields.
	for _, field := range fields {
		tag := field.tag
		if tag.name == "" {
			continue
		}

		haventSeen := seen[tag.name] == nil
		entries := mentries[tag.name]
		if len(entries) == 0 {
			if !tag.optional && haventSeen {
				return fmt.Errorf("missing required attribute %q", tag.name)
			}

			// Apply defaults here as there's no value for this field.
			v, err := defaultValueFromTag(field, tag.defaultValue)
			if err != nil {
				return err
			}

			if v != nil {
				err = unmarshalValue(field.v, v)
				if err != nil {
					return fmt.Errorf("error applying default value to field %q, %v", field.t.Name, err)
				}
			}

			continue
		}
		delete(seen, tag.name)

		entry := entries[0]
		entries = entries[1:]
		mentries[tag.name] = entries

		// Field is a pointer, create value if necessary, then move field down.
		if field.v.Kind() == reflect.Ptr {
			if field.v.IsNil() {
				field.v.Set(reflect.New(field.v.Type().Elem()))
			}
			field.v = field.v.Elem()
			field.t.Type = field.t.Type.Elem()
		}

		switch field.v.Kind() {
		default:
			// Anything else must be a scalar value.
			if len(entries) > 0 {
				return participle.Errorf(entry.Position(), "duplicate field %q at %s", entry.Key(), entries[0].Position())
			}

			entry := entry.(*Object)
			value := entry.Value
			err = unmarshalValue(field.v, value)
			if err != nil {
				pos := entry.Position()
				if value != nil {
					pos = value.Position()
				}

				return participle.Wrapf(pos, err, "failed to unmarshal value")
			}
		}
	}

	return nil
}

var kindOfTime = reflect.TypeOf(time.Time{}).Kind()

func unmarshalValue(rv reflect.Value, v Value) error {
	k := rv.Kind()

	switch v := v.(type) {
	case *String:
		switch k {
		case reflect.String, reflect.Interface:
			rv.Set(reflect.ValueOf(v.Value))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(v.Value, 0, 64)
			if err != nil {
				return participle.Errorf(v.Position(), "error converting %q to int", v)
			}
			rv.SetInt(n)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(v.Value, 0, 64)
			if err != nil {
				return participle.Errorf(v.Position(), "error converting %q to uint", v)
			}
			rv.SetUint(n)

		case reflect.Float32, reflect.Float64:
			size := 64
			if k == reflect.Float32 {
				size = 32
			}
			n, err := strconv.ParseFloat(v.Value, size)
			if err != nil {
				return participle.Errorf(v.Position(), "error converting %q to float", v)
			}
			rv.SetFloat(n)

		default:
			return participle.Errorf(v.Position(), "unable to unmarshal string into %s", k)
		}

	case *Number:
		switch k {
		case reflect.String:
			rv.SetString(v.Value.String())

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, _ := v.Value.Int64()
			rv.SetInt(n)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, _ := v.Value.Uint64()
			rv.SetUint(n)

		case reflect.Float32, reflect.Float64, reflect.Interface:
			n, _ := v.Value.Float64()
			rv.Set(reflect.ValueOf(n))

		default:
			return participle.Errorf(v.Position(), "unable to unmarshal number into %s", k)
		}

	case *Measurement:
		switch k {
		case reflect.Slice:
			t := rv.Type().Elem()
			lv := reflect.MakeSlice(rv.Type(), 0, 2)
			for _, entry := range v.children() {
				value := reflect.New(t).Elem()
				if err := unmarshalValue(value, entry.(Value)); err != nil {
					return participle.Wrapf(entry.Position(), err, "invalid measurement value")
				}

				lv = reflect.Append(lv, value)
			}
			rv.Set(lv)

		default:
			if err := unmarshalValue(rv, v.Value); err != nil {
				return participle.Wrapf(v.Position(), err, "invalid measurement value")
			}
		}

	case *Timestamp:
		switch k {
		case reflect.String, reflect.Interface:
			rv.Set(reflect.ValueOf(v.Value))

		case kindOfTime:
			// ...

		default:
			return participle.Errorf(v.Position(), "unable to unmarshal timestamp into %s", k)
		}

	case *List:
		switch k {
		case reflect.Slice:
			t := rv.Type().Elem()
			lv := reflect.MakeSlice(rv.Type(), 0, 4)
			for _, entry := range v.children() {
				value := reflect.New(t).Elem()
				if err := unmarshalValue(value, entry.(Value)); err != nil {
					return participle.Wrapf(entry.Position(), err, "invalid measurement value")
				}

				lv = reflect.Append(lv, value)
			}
			rv.Set(lv)

		default:
			return participle.Errorf(v.Position(), "unable to unmarshal list into %s", k)
		}

	default:
		panic(v)
	}

	return nil
}

type field struct {
	t   reflect.StructField
	v   reflect.Value
	tag tag
}

func flattenFields(v reflect.Value) ([]field, error) {
	out := make([]field, 0, v.NumField())
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		ft := t.Field(i)
		if ft.Anonymous {
			if f.Kind() != reflect.Struct {
				return nil, fmt.Errorf("%s: anonymous field must be a struct", ft.Name)
			}
			sub, err := flattenFields(f)
			if err != nil {
				return nil, fmt.Errorf("%s: %s", ft.Name, err)
			}

			out = append(out, sub...)
		} else {
			tag := parseTag(v.Type(), ft)
			out = append(out, field{ft, f, tag})
		}
	}

	return out, nil
}

type tag struct {
	name         string
	optional     bool
	defaultValue string
}

func parseTag(parent reflect.Type, t reflect.StructField) tag {
	defaultValue := t.Tag.Get("default")
	s, _ := t.Tag.Lookup("dsmr")

	parts := strings.Split(s, ",")
	name := parts[0]
	if name == "-" {
		return tag{}
	}
	if name == "" {
		name = t.Name
	}
	if len(parts) == 1 {
		return tag{name: name, defaultValue: defaultValue, optional: defaultValue != ""}
	}

	id := fieldID(parent, t)
	option := parts[1]
	switch option {
	case "optional", "omitempty":
		return tag{name: name, optional: true, defaultValue: defaultValue}
	default:
		panic("invalid DSMR tag option " + option + " on " + id)
	}
}

func fieldID(parent reflect.Type, t reflect.StructField) string {
	return fmt.Sprintf("%s.%s.%s", parent.PkgPath(), parent.Name(), t.Name)
}

func defaultValueFromTag(f field, defaultValue string) (Value, error) {
	v, err := valueFromTag(f, defaultValue)
	if err != nil {
		return nil, fmt.Errorf("error parsing default value: %v", err)
	}
	return v, nil
}

func valueFromTag(f field, defaultValue string) (Value, error) {
	if defaultValue == "" {
		return nil, nil // nolint: nilnil
	}

	k := f.v.Kind()
	if k == reflect.Ptr {
		k = f.v.Type().Elem().Kind()
	}

	switch k {
	case reflect.String:
		return &String{Value: defaultValue}, nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(defaultValue, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting %q to int", defaultValue)
		}
		return &Number{Value: big.NewFloat(0).SetInt64(n)}, nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := strconv.ParseUint(defaultValue, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("error converting %q to uint", defaultValue)
		}
		return &Number{Value: big.NewFloat(0).SetUint64(n)}, nil

	case reflect.Float32, reflect.Float64:
		size := 64
		if k == reflect.Float32 {
			size = 32
		}
		n, err := strconv.ParseFloat(defaultValue, size)
		if err != nil {
			return nil, fmt.Errorf("error converting %q to float", defaultValue)
		}
		return &Number{Value: big.NewFloat(n)}, nil

	default:
		return nil, fmt.Errorf("only primitive types, map & slices can have tag value, not %q", f.v.Type())
	}
}
