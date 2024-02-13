package dynamicstruct

import (
	"errors"
	"reflect"
)

type (
	Writer interface {
		Write(name string, value interface{}) error
		Value() interface{}
	}

	writerImpl struct {
		fields map[string]fieldImpl
		value  interface{}
	}
)

func NewWriter(value interface{}) Writer {
	fields := map[string]fieldImpl{}

	valueOf := reflect.Indirect(reflect.ValueOf(value))
	typeOf := valueOf.Type()

	if typeOf.Kind() == reflect.Struct {
		for i := 0; i < valueOf.NumField(); i++ {
			field := typeOf.Field(i)
			fields[field.Name] = fieldImpl{
				index: i,
				field: field,
				value: valueOf.Field(i),
			}
		}
	}

	return &writerImpl{
		fields: fields,
		value:  value,
	}
}

func (w *writerImpl) Write(name string, value interface{}) error {

	valueOf := reflect.Indirect(reflect.ValueOf(value))
	typeOf := valueOf.Type()

	if typeOf.Kind() != w.fields[name].value.Type().Kind() {
		return errors.New("unassignable value")
	}

	field := w.fields[name].field
	w.fields[name] = fieldImpl{
		field: field,
		value: reflect.ValueOf(value),
	}

	structValue := reflect.Indirect(reflect.ValueOf(w.value))
	structValue.Field(int(w.fields[name].index)).Set(valueOf)

	w.value = structValue.Interface()

	return nil
}

func (w *writerImpl) Value() interface{} {
	return w.value
}
