package dynamicstruct

import (
	"reflect"
)

type (
	Builder interface {
		AddField(name string, typ reflect.Kind, tag string) Builder
		RemoveField(name string) Builder
		HasField(name string) bool
		Field(name string) FieldConfig
		Build() DynamicStruct
	}

	FieldConfig interface {
		SetType(typ reflect.Kind) FieldConfig
		SetTag(tag string) FieldConfig
	}

	DynamicStruct interface {
		New() interface{}
		NewSliceOfStructs() interface{}
		NewMapOfStructs(key interface{}) interface{}
	}

	builderImpl struct {
		fields []*fieldConfigImpl
	}

	fieldConfigImpl struct {
		name      string
		pkg       string
		typ       reflect.Kind
		tag       string
		anonymous bool
	}

	dynamicStructImpl struct {
		definition reflect.Type
	}
)

func NewStruct() Builder {
	return &builderImpl{
		fields: []*fieldConfigImpl{},
	}
}

func ExtendStruct(value interface{}) Builder {
	return MergeStructs(value)
}

func MergeStructs(values ...interface{}) Builder {
	builder := NewStruct()

	for _, value := range values {
		valueOf := reflect.Indirect(reflect.ValueOf(value))
		typeOf := valueOf.Type()

		for i := 0; i < valueOf.NumField(); i++ {
			fval := valueOf.Field(i)
			ftyp := typeOf.Field(i)
			builder.(*builderImpl).addField(ftyp.Name, ftyp.PkgPath, fval.Kind(), string(ftyp.Tag), ftyp.Anonymous)
		}
	}

	return builder
}

func (b *builderImpl) AddField(name string, typ reflect.Kind, tag string) Builder {
	if name == "" {
		typ_ := reflect.TypeOf(typ)
		return b.addField(typ_.Name(), typ_.PkgPath(), typ, tag, true)
	}
	return b.addField(name, "", typ, tag, false)
}

func (b *builderImpl) addField(name string, pkg string, typ reflect.Kind, tag string, anonymous bool) Builder {
	b.fields = append(b.fields, &fieldConfigImpl{
		name:      name,
		typ:       typ,
		tag:       tag,
		anonymous: anonymous,
	})

	return b
}

func (b *builderImpl) RemoveField(name string) Builder {
	for i := range b.fields {
		if b.fields[i].name == name {
			b.fields = append(b.fields[:i], b.fields[i+1:]...)
			break
		}
	}
	return b
}

func (b *builderImpl) HasField(name string) bool {
	for i := range b.fields {
		if b.fields[i].name == name {
			return true
		}
	}
	return false
}

func (b *builderImpl) Field(name string) FieldConfig {
	for i := range b.fields {
		if b.fields[i].name == name {
			return b.fields[i]
		}
	}
	return nil
}

func (b *builderImpl) Build() DynamicStruct {
	var structFields []reflect.StructField

	for _, field := range b.fields {
		structFields = append(structFields, reflect.StructField{
			Name:      field.name,
			PkgPath:   field.pkg,
			Type:      reflect.TypeOf(field.typ),
			Tag:       reflect.StructTag(field.tag),
			Anonymous: field.anonymous,
		})
	}

	return &dynamicStructImpl{
		definition: reflect.StructOf(structFields),
	}
}

func (f *fieldConfigImpl) SetType(typ reflect.Kind) FieldConfig {
	f.typ = typ
	return f
}

func (f *fieldConfigImpl) SetTag(tag string) FieldConfig {
	f.tag = tag
	return f
}

func (ds *dynamicStructImpl) New() interface{} {
	return reflect.New(ds.definition).Interface()
}

func (ds *dynamicStructImpl) NewSliceOfStructs() interface{} {
	return reflect.New(reflect.SliceOf(ds.definition)).Interface()
}

func (ds *dynamicStructImpl) NewMapOfStructs(key interface{}) interface{} {
	return reflect.New(reflect.MapOf(reflect.Indirect(reflect.ValueOf(key)).Type(), ds.definition)).Interface()
}
