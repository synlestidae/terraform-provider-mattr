package generator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ResourceRep struct {
	kind      reflect.Kind
	valueType schema.ValueType
	fields    []Field
	elem      *ResourceRep
}

type Field struct {
	schemaName string
	fieldName  string
	opts       SchemaOpts
	resource   ResourceRep
}

func typeResource(t reflect.Type) (*ResourceRep, error) {
	var rep ResourceRep
	rep.kind = t.Kind()

	switch t.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		valueType, err := getSchemaType(t.Kind())
		if err != nil {
			return nil, err
		}
		rep.valueType = valueType
		break
	case reflect.Struct:
		numField := t.NumField()
		rep.fields = make([]Field, numField)

		for i := 0; i < numField; i++ {
			field := t.Field(i)
			fieldName := field.Name
			schemaName := snakeCase(fieldName)
			opts := fieldOpts(&field.Tag)
			resource, err := typeResource(field.Type)
			if err != nil {
				return nil, err
			}

			rep.fields[i] = Field{
				schemaName: schemaName,
				fieldName:  fieldName,
				opts:       opts,
				resource:   *resource,
			}
		}
		break
	case reflect.Array:
	case reflect.Slice:
		rep.valueType = schema.TypeList
		typeElem, err := typeResource(t.Elem())
		if err != nil {
			return nil, err
		}
		rep.elem = typeElem
		break
	default:
		panic(fmt.Sprintf("Unsupported schema for type: %s", t.Kind()))
	}

	return &rep, nil
}

func fieldOpts(tag *reflect.StructTag) SchemaOpts {
	options := trimSlice(strings.Split(tag.Get("schemaOpts"), ","))

	return SchemaOpts{
		Computed: contains(options, "computed"),
		Required: contains(options, "required"),
		Optional: contains(options, "optional"),
	}
}

type Visitor interface {
	visitStruct(*ResourceRep) error
	visitList(*ResourceRep) error
	visitPrimitive(*ResourceRep) error
}

type ResourceVisitor struct {
	schema schema.Schema
}

func (vs *ResourceVisitor) visitPrimitive(rs *ResourceRep) error {
	schemaType, err := getSchemaType(rs.kind)
	if err != nil {
		return err
	}

	vs.schema = schema.Schema{
		Type: schemaType,
	}

	return nil
}

func (vs *ResourceVisitor) visitStruct(rs *ResourceRep) error {
	schemaMap := make(map[string]*schema.Schema, len(rs.fields))

	for _, field := range rs.fields {
		var subVs ResourceVisitor
		err := field.resource.accept(&subVs)
		if err != nil {
			return err
		}

		opts := field.opts
		schemaMap[field.schemaName] = &subVs.schema
		schemaMap[field.schemaName].Computed = opts.Computed
		schemaMap[field.schemaName].Required = opts.Required
		schemaMap[field.schemaName].Optional = opts.Optional
	}

	vs.schema = schema.Schema{
		Type: schema.TypeMap,
		Elem: schemaMap,
	}

	return nil
}

func (vs *ResourceVisitor) visitList(rs *ResourceRep) error {
	var subVs ResourceVisitor

	if rs.elem == nil {
		return fmt.Errorf("Array does not specify elem type")
	}

	if err := rs.elem.accept(&subVs); err != nil {
		return err
	}

	vs.schema = schema.Schema{
		Type: schema.TypeList,
		Elem: &subVs.schema,
	}

	return nil
}

func isPrimitive(kind reflect.Kind) bool {
	switch kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return true
	default:
		return false
	}
}

func (r *ResourceRep) accept(visitor Visitor) error {
	switch r.kind {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.String:
		return visitor.visitPrimitive(r)
	case reflect.Struct:
		return visitor.visitStruct(r)
	case reflect.Array:
		return visitor.visitList(r)
	}
	panic(fmt.Sprintf("Unsupported schema for type: %s", r.kind))
}

type SchemaOpts struct {
	Computed bool
	Required bool
	Optional bool
}

func trimSlice(input []string) []string {
	for i := range input {
		input[i] = strings.TrimSpace(input[i])
	}
	return input
}

func contains(slice []string, value string) bool {
	for _, s := range slice {
		if s == value {
			return true
		}
	}

	return false
}

func snakeCase(input string) string {
	var builder strings.Builder
	for i, c := range input {
		if unicode.IsUpper(c) && i != 0 && i+1 != len(input) {
			builder.WriteString("_" + string(unicode.ToLower(c)))
		} else {
			builder.WriteString(string(unicode.ToLower(c)))
		}
	}
	return builder.String()
}

func getSchemaType(kind reflect.Kind) (schema.ValueType, error) {
	switch kind {
	case reflect.String:
		return schema.TypeString, nil
	case reflect.Bool:
		return schema.TypeBool, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return schema.TypeInt, nil
	case reflect.Float32, reflect.Float64:
		return schema.TypeFloat, nil
	case reflect.Struct:
		return schema.TypeMap, nil
	case reflect.Array:
	case reflect.Map:
	case reflect.Slice:
		panic("Not yet implemented")
	case reflect.Chan:
	case reflect.Pointer:
	case reflect.UnsafePointer:
	case reflect.Interface:
		panic(fmt.Sprintf("Unsupported schema for type: %s", kind))
	}
	panic(fmt.Sprintf("Unsupported schema for type: %s", kind))
}

type RequestVisitor struct {
	value interface{}
	data  interface{}
}

func (rv *RequestVisitor) visitStruct(rs *ResourceRep) error {
	valueMap := make(map[string]interface{}, len(rs.fields))

	for _, field := range rs.fields {
		var subVs RequestVisitor

		if data, ok := rv.data.(*schema.ResourceData); ok {
			subVs.data = data.Get(field.schemaName) // TODO don't visit if missing
		} else if data, ok := rv.data.(*map[string]interface{}); ok {
			subVs.data = (*data)[field.schemaName] // TODO don't visit if missing
		} else {
			return fmt.Errorf("Failed to convert '%s' to map[string]interface{}", field.schemaName)
		}

		field.resource.accept(&subVs)
		valueMap[field.fieldName] = subVs.value
	}

	rv.value = &valueMap
	return nil
}

func (rv *RequestVisitor) visitList(rs *ResourceRep) error {
	list, ok := rv.data.([]interface{})
	if !ok {
		return fmt.Errorf("Failed to convert data to []interface")
	}
	valueList := make([]interface{}, len(list))
	for i, data := range list {
		subVs := RequestVisitor{
			data: data,
		}
		rs.elem.accept(&subVs)
		valueList[i] = subVs.value
	}
	rv.value = valueList
	return nil
}

func (rv *RequestVisitor) visitPrimitive(rp *ResourceRep) error {
	valType := reflect.TypeOf(rv.data)
	valKind := valType.Kind()

	// Check if rv.data is assignable to rp.kind
	if !isPrimitive(valKind) || valKind != rp.kind {
		return fmt.Errorf("Failed to assign '%s' to %s", valType.Kind(), rp.kind)
	}

	rv.value = rv.data
	return nil
}
