package config

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)

func (gos *GameObjects) LoadType(in interface{}) ([]interface{}, error) {

	t := reflect.TypeOf(in)
	t = typeElem(t)

	location := getTag(t, "location")
	if location == "" {
		return nil, errors.Errorf("location not found, %s", t.Name())
	}

	data := gos.Data(location)
	if data == "" {
		return nil, errors.Errorf("data not found, %s", location)
	}

	list, err := ParseList(location, data)
	if err != nil {
		return nil, err
	}

	out := make([]interface{}, len(list))

	for i, p := range list {

		newPtr := reflect.New(t).Interface()

		err = (*reflect_parser)(p).FillTo(newPtr)
		if err != nil {
			return nil, err
		}

		out[i] = newPtr
	}

	return out, nil
}

type reflect_parser ObjectParser

func (p *reflect_parser) FillTo(origin interface{}) error {

	v, t, err := findStructType(origin)
	if err != nil {
		return err
	}

	if !v.CanSet() {
		return errors.Errorf("value can not set, %s", origin)
	}

	for i := 0; i < t.NumField(); i++ {

		tf := t.Field(i)
		if tf.PkgPath != "" {
			continue
		}

		ignore, _ := strconv.ParseBool(tf.Tag.Get("ignore"))
		if ignore {
			continue
		}

		vf := v.Field(i)

		if !vf.IsValid() {
			return errors.Errorf("%s field is invalid", tf.Name)
		}

		if !vf.CanSet() {
			return errors.Errorf("%s field is cant set", tf.Name)
		}

		switch vf.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			err := p.setIntField(tf, vf)
			if err != nil {
				return err
			}
		case reflect.Float32, reflect.Float64:
			err := p.setFloatField(tf, vf)
			if err != nil {
				return err
			}
		case reflect.String:
			err := p.setStringField(tf, vf)
			if err != nil {
				return err
			}

		case reflect.Slice:
			err := p.setSliceField(tf, vf)
			if err != nil {
				return err
			}
		default:
			return errors.New("unsupport field kind. kind: " + vf.Kind().String())
		}

	}

	return nil
}

func (p *reflect_parser) setIntField(tf reflect.StructField, vf reflect.Value) error {
	ptr, err := p.MarshalField(tf, s2i)
	if err != nil {
		return err
	}

	i64 := int64(ptr.(int))
	if vf.OverflowInt(i64) {
		return errors.Errorf("vf.OverflowInt(i64)")
	}

	vf.SetInt(i64)
	return nil
}

func s2i(value string) (interface{}, error) {
	ii, err := strconv.Atoi(value)
	if err != nil {
		return nil, errors.Wrapf(err, "Key not found, convert default fail, %s", value)
	}

	return ii, nil
}

func (p *reflect_parser) setFloatField(tf reflect.StructField, vf reflect.Value) error {
	ptr, err := p.MarshalField(tf, s2f)
	if err != nil {
		return err
	}

	f64 := ptr.(float64)
	if vf.OverflowFloat(f64) {
		return errors.Errorf("vf.OverflowFloat(f64)")
	}

	vf.SetFloat(f64)
	return nil
}

func s2f(value string) (interface{}, error) {
	ii, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, errors.Wrapf(err, "Key not found, convert default fail, %s", value)
	}

	return ii, nil
}

func (p *reflect_parser) setStringField(tf reflect.StructField, vf reflect.Value) error {
	s, err := p.MarshalField(tf, nil)
	if err != nil {
		return err
	}

	vf.SetString(s.(string))
	return nil
}

func (p *reflect_parser) setSliceField(tf reflect.StructField, vf reflect.Value) error {

	// sep
	sep := tf.Tag.Get("sep")

	var value interface{}
	var err error
	switch tf.Type.Elem().Kind() {
	case reflect.Int:
		value = (*ObjectParser)(p).IntArray(tf.Name, sep, false)
	case reflect.Float32:
		value = F64a2F32a((*ObjectParser)(p).Float64Array(tf.Name, sep, false))
	case reflect.Float64:
		value = (*ObjectParser)(p).Float64Array(tf.Name, sep, false)
	case reflect.String:
		value = (*ObjectParser)(p).StringArray(tf.Name, sep, false)
	default:
		return errors.Errorf("unsupport slice field kind. kind: []%s", tf.Type.Elem().Kind())
	}

	if err != nil {
		return err
	}

	vf.Set(reflect.ValueOf(value))
	return nil
}

func (p *reflect_parser) MarshalField(tf reflect.StructField, parse func(string) (interface{}, error)) (interface{}, error) {

	str := ""
	v, ok := p.dataMap[strings.ToLower(tf.Name)]
	if !ok {
		v, ok := tf.Tag.Lookup("default")
		if !ok {
			return nil, errors.Errorf("key not found, %s", tf.Name)
		}
		str = v
	} else {
		if len(v) != 1 {
			return nil, errors.Errorf("key more than 1 column, key: %s, array: %s", tf.Name, v)
		}

		str = v[0]
	}

	if parse != nil {
		return parse(str)
	}

	return str, nil
}

func findStructType(in interface{}) (reflect.Value, reflect.Type, error) {
	v := reflect.ValueOf(in)
	v = valueElem(v)

	if v.Kind() != reflect.Struct {
		return reflect.Value{}, nil, errors.Errorf("value is not a struct, %s", v.Kind())
	}

	return v, v.Type(), nil
}

func valueElem(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	return v
}

func typeElem(v reflect.Type) reflect.Type {
	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	return v
}

func getTag(t reflect.Type, key string) string {
	for i := 0; i < t.NumField(); i++ {

		tf := t.Field(i)

		value, ok := tf.Tag.Lookup(key)
		if ok {
			return value
		}
	}

	return ""
}
