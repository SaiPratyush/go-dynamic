package dynamicstruct

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

type (
	Writer interface {
		HasField(name string) bool
		Field(name string) Field
		Fields() []Field
		ToStruct(value interface{}) error
		ToSliceOfReaders() []Writer
		ToMapReaderOfReaders() map[interface{}]Writer
		Value() interface{}
	}

	Field interface {
		Name() string

		// Getters
		PointerInt() *int
		Int() int
		PointerInt8() *int8
		Int8() int8
		PointerInt16() *int16
		Int16() int16
		PointerInt32() *int32
		Int32() int32
		PointerInt64() *int64
		Int64() int64
		PointerUint() *uint
		Uint() uint
		PointerUint8() *uint8
		Uint8() uint8
		PointerUint16() *uint16
		Uint16() uint16
		PointerUint32() *uint32
		Uint32() uint32
		PointerUint64() *uint64
		Uint64() uint64
		PointerFloat32() *float32
		Float32() float32
		PointerFloat64() *float64
		Float64() float64
		PointerString() *string
		String() string
		PointerBool() *bool
		Bool() bool
		PointerTime() *time.Time
		Time() time.Time
		Interface() interface{}

		// Setters
		SetPointerInt(value *int)
		SetInt(value int)
		SetPointerInt8(value *int8)
		SetInt8(value int8)
		SetPointerInt16(value *int16)
		SetInt16(value int16)
		SetPointerInt32(value *int32)
		SetInt32(value int32)
		SetPointerInt64(value *int64)
		SetInt64(value int64)
		SetPointerUint(value *uint)
		SetUint(value uint)
		SetPointerUint8(value *uint8)
		SetUint8(value uint8)
		SetPointerUint16(value *uint16)
		SetUint16(value uint16)
		SetPointerUint32(value *uint32)
		SetUint32(value uint32)
		SetPointerUint64(value *uint64)
		SetUint64(value uint64)
		SetPointerFloat32(value *float32)
		SetFloat32(value float32)
		SetPointerFloat64(value *float64)
		SetFloat64(value float64)
		SetPointerString(value *string)
		SetString(value string)
		SetPointerBool(value *bool)
		SetBool(value bool)
		SetPointerTime(value *time.Time)
		SetTime(value time.Time)
		SetInterface(value interface{})
	}

	readImpl struct {
		fields map[string]fieldImpl
		value  interface{}
	}

	fieldImpl struct {
		index int
		field reflect.StructField
		value reflect.Value
	}
)

func NewReader(value interface{}) Writer {
	fields := map[string]fieldImpl{}

	valueOf := reflect.Indirect(reflect.ValueOf(value))
	typeOf := valueOf.Type()

	if typeOf.Kind() == reflect.Struct {
		for i := 0; i < valueOf.NumField(); i++ {
			field := typeOf.Field(i)
			fields[field.Name] = fieldImpl{
				field: field,
				value: valueOf.Field(i),
			}
		}
	}

	return readImpl{
		fields: fields,
		value:  value,
	}
}

func (r readImpl) HasField(name string) bool {
	_, ok := r.fields[name]
	return ok
}

func (r readImpl) Field(name string) Field {
	if !r.HasField(name) {
		return nil
	}
	return r.fields[name]
}

func (r readImpl) Fields() []Field {
	var fields []Field

	for _, field := range r.fields {
		fields = append(fields, field)
	}

	return fields
}

func (r readImpl) ToStruct(value interface{}) error {
	valueOf := reflect.ValueOf(value)

	if valueOf.Kind() != reflect.Ptr || valueOf.IsNil() {
		return errors.New("ToStruct: expected a pointer as an argument")
	}

	valueOf = valueOf.Elem()
	typeOf := valueOf.Type()

	if valueOf.Kind() != reflect.Struct {
		return errors.New("ToStruct: expected a pointer to struct as an argument")
	}

	for i := 0; i < valueOf.NumField(); i++ {
		fieldType := typeOf.Field(i)
		fieldValue := valueOf.Field(i)

		original, ok := r.fields[fieldType.Name]
		if !ok {
			continue
		}

		if fieldValue.CanSet() && r.haveSameTypes(original.value.Type(), fieldValue.Type()) {
			fieldValue.Set(original.value)
		}
	}

	return nil
}

func (r readImpl) ToSliceOfReaders() []Writer {
	valueOf := reflect.Indirect(reflect.ValueOf(r.value))
	typeOf := valueOf.Type()

	if typeOf.Kind() != reflect.Slice && typeOf.Kind() != reflect.Array {
		return nil
	}

	var readers []Writer

	for i := 0; i < valueOf.Len(); i++ {
		readers = append(readers, NewReader(valueOf.Index(i).Interface()))
	}

	return readers
}

func (r readImpl) ToMapReaderOfReaders() map[interface{}]Writer {
	valueOf := reflect.Indirect(reflect.ValueOf(r.value))
	typeOf := valueOf.Type()

	if typeOf.Kind() != reflect.Map {
		return nil
	}

	readers := map[interface{}]Writer{}

	for _, keyValue := range valueOf.MapKeys() {
		readers[keyValue.Interface()] = NewReader(valueOf.MapIndex(keyValue).Interface())
	}

	return readers
}

func (r readImpl) Value() interface{} {
	return r.value
}

func (r readImpl) haveSameTypes(first reflect.Type, second reflect.Type) bool {
	if first.Kind() != second.Kind() {
		return false
	}

	switch first.Kind() {
	case reflect.Ptr:
		return r.haveSameTypes(first.Elem(), second.Elem())
	case reflect.Struct:
		return first.PkgPath() == second.PkgPath() && first.Name() == second.Name()
	case reflect.Slice:
		return r.haveSameTypes(first.Elem(), second.Elem())
	case reflect.Map:
		return r.haveSameTypes(first.Elem(), second.Elem()) && r.haveSameTypes(first.Key(), second.Key())
	default:
		return first.Kind() == second.Kind()
	}
}

func (f fieldImpl) Name() string {
	return f.field.Name
}

func (f fieldImpl) PointerInt() *int {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int()
	return &value
}

func (f fieldImpl) Int() int {
	return int(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt8() *int8 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int8()
	return &value
}

func (f fieldImpl) Int8() int8 {
	return int8(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt16() *int16 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int16()
	return &value
}

func (f fieldImpl) Int16() int16 {
	return int16(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt32() *int32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int32()
	return &value
}

func (f fieldImpl) Int32() int32 {
	return int32(reflect.Indirect(f.value).Int())
}

func (f fieldImpl) PointerInt64() *int64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Int64()
	return &value
}

func (f fieldImpl) Int64() int64 {
	return reflect.Indirect(f.value).Int()
}

func (f fieldImpl) PointerUint() *uint {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint()
	return &value
}

func (f fieldImpl) Uint() uint {
	return uint(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint8() *uint8 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint8()
	return &value
}

func (f fieldImpl) Uint8() uint8 {
	return uint8(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint16() *uint16 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint16()
	return &value
}

func (f fieldImpl) Uint16() uint16 {
	return uint16(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint32() *uint32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint32()
	return &value
}

func (f fieldImpl) Uint32() uint32 {
	return uint32(reflect.Indirect(f.value).Uint())
}

func (f fieldImpl) PointerUint64() *uint64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Uint64()
	return &value
}

func (f fieldImpl) Uint64() uint64 {
	return reflect.Indirect(f.value).Uint()
}

func (f fieldImpl) PointerFloat32() *float32 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Float32()
	return &value
}

func (f fieldImpl) Float32() float32 {
	return float32(reflect.Indirect(f.value).Float())
}

func (f fieldImpl) PointerFloat64() *float64 {
	if f.value.IsNil() {
		return nil
	}
	value := f.Float64()
	return &value
}

func (f fieldImpl) Float64() float64 {
	return reflect.Indirect(f.value).Float()
}

func (f fieldImpl) PointerString() *string {
	if f.value.IsNil() {
		return nil
	}
	value := f.String()
	return &value
}

func (f fieldImpl) String() string {
	return reflect.Indirect(f.value).String()
}

func (f fieldImpl) PointerBool() *bool {
	if f.value.IsNil() {
		return nil
	}
	value := f.Bool()
	return &value
}

func (f fieldImpl) Bool() bool {
	return reflect.Indirect(f.value).Bool()
}

func (f fieldImpl) PointerTime() *time.Time {
	if f.value.IsNil() {
		return nil
	}
	value := f.Time()
	return &value
}

func (f fieldImpl) Time() time.Time {
	value, ok := reflect.Indirect(f.value).Interface().(time.Time)
	if !ok {
		panic(fmt.Sprintf(`field "%s" is not instance of time.Time`, f.field.Name))
	}

	return value
}

func (f fieldImpl) Interface() interface{} {
	return f.value.Interface()
}

func (f fieldImpl) haveSameTypes(first reflect.Type, second reflect.Type) bool {
	if first.Kind() != second.Kind() {
		return false
	}

	switch first.Kind() {
	case reflect.Ptr:
		return f.haveSameTypes(first.Elem(), second.Elem())
	case reflect.Struct:
		return first.PkgPath() == second.PkgPath() && first.Name() == second.Name()
	case reflect.Slice:
		return f.haveSameTypes(first.Elem(), second.Elem())
	case reflect.Map:
		return f.haveSameTypes(first.Elem(), second.Elem()) && f.haveSameTypes(first.Key(), second.Key())
	default:
		return first.Kind() == second.Kind()
	}
}

func (f fieldImpl) SetPointerInt(value *int) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetInt(*value)
}

func (f fieldImpl) SetInt(value int) {
	f.value.SetInt(int64(value))
}

func (f fieldImpl) SetPointerInt8(value *int8) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetInt8(*value)
}

func (f fieldImpl) SetInt8(value int8) {
	f.value.SetInt(int64(value))
}

func (f fieldImpl) SetPointerInt16(value *int16) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetInt16(*value)
}

func (f fieldImpl) SetInt16(value int16) {
	f.value.SetInt(int64(value))
}

func (f fieldImpl) SetPointerInt32(value *int32) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetInt32(*value)
}

func (f fieldImpl) SetInt32(value int32) {
	f.value.SetInt(int64(value))
}

func (f fieldImpl) SetPointerInt64(value *int64) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetInt64(*value)
}

func (f fieldImpl) SetInt64(value int64) {
	f.value.SetInt(value)
}

func (f fieldImpl) SetPointerUint(value *uint) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetUint(*value)
}

func (f fieldImpl) SetUint(value uint) {
	f.value.SetUint(uint64(value))
}

func (f fieldImpl) SetPointerUint8(value *uint8) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetUint8(*value)
}

func (f fieldImpl) SetUint8(value uint8) {
	f.value.SetUint(uint64(value))
}

func (f fieldImpl) SetPointerUint16(value *uint16) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetUint16(*value)
}

func (f fieldImpl) SetUint16(value uint16) {
	f.value.SetUint(uint64(value))
}

func (f fieldImpl) SetPointerUint32(value *uint32) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetUint32(*value)
}

func (f fieldImpl) SetUint32(value uint32) {
	f.value.SetUint(uint64(value))
}

func (f fieldImpl) SetPointerUint64(value *uint64) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetUint64(*value)
}

func (f fieldImpl) SetUint64(value uint64) {
	f.value.SetUint(value)
}

func (f fieldImpl) SetPointerFloat32(value *float32) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetFloat32(*value)
}

func (f fieldImpl) SetFloat32(value float32) {
	f.value.SetFloat(float64(value))
}

func (f fieldImpl) SetPointerFloat64(value *float64) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetFloat64(*value)
}

func (f fieldImpl) SetFloat64(value float64) {
	f.value.SetFloat(value)
}

func (f fieldImpl) SetPointerString(value *string) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetString(*value)
}

func (f fieldImpl) SetString(value string) {
	f.value.SetString(value)
}

func (f fieldImpl) SetPointerBool(value *bool) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetBool(*value)
}

func (f fieldImpl) SetBool(value bool) {
	f.value.SetBool(value)
}

func (f fieldImpl) SetPointerTime(value *time.Time) {
	if value == nil {
		f.value.Set(reflect.Zero(f.value.Type()))
		return
	}
	f.SetTime(*value)
}

func (f fieldImpl) SetTime(value time.Time) {
	f.value.Set(reflect.ValueOf(value))
}

func (f fieldImpl) SetInterface(value interface{}) {
	f.value.Set(reflect.ValueOf(value))
}
