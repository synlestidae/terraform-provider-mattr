package main

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
	resource   ResourceRep
	opts       SchemaOpts
}

type Visitor interface {
	visitStruct(*ResourceRep) error
	visitArray(*ResourceRep) error
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

		schemaMap[field.schemaName] = &schema.Schema{
			Type:     schema.TypeMap,
			Computed: field.opts.Computed,
			Required: field.opts.Required,
			Optional: field.opts.Optional,
			Elem:     subVs.schema,
		}
	}

	vs.schema = schema.Schema{
		Type: schema.TypeMap,
		Elem: schemaMap,
	}

	return nil
}

func (vs *ResourceVisitor) visitArray(rs *ResourceRep) error {
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

func (rs *ResourceRep) hasPrimitive() bool {
	switch rs.kind {
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
		return visitor.visitArray(r)
	}
	panic(fmt.Sprintf("Unsupported schema for type: %s", r.kind))
}

func resourceFromType(typ reflect.Type) (ResourceRep, error) {
	var resource ResourceRep
	kind := typ.Kind()

	switch kind {
	case reflect.Bool:
	case reflect.String:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	case reflect.Float32, reflect.Float64:
		resource.kind = kind
		schemaType, err := getSchemaType(kind)
		if err != nil {
			return resource, err
		}
		resource.valueType = schemaType
		return resource, nil

	case reflect.Struct:
		numField := typ.NumField()
		fields := make([]Field, numField)

		for i := 0; i < numField; i++ {
			field := typ.Field(i)
			schemaName := snakeCase(field.Name)
			opts := fieldOpts(&field.Tag)
			fieldResource, err := resourceFromType(field.Type)
			if err != err {
				return resource, err
			}
			fields[i] = Field{
				schemaName: schemaName,
				fieldName:  field.Name,
				resource:   fieldResource,
				opts:       opts,
			}
		}
	case reflect.Array:
	case reflect.Map:
	case reflect.Interface:
	case reflect.Slice:
		panic("Not yet implemented")
	case reflect.Chan:
	case reflect.Pointer:
	case reflect.UnsafePointer:
		panic(fmt.Sprintf("Unsupported schema for type: %s", typ.Kind()))
	}
	panic(fmt.Sprintf("Unsupported schema for type: %s", typ.Kind()))
}

type SchemaOpts struct {
	Computed bool
	Required bool
	Optional bool
}

func fieldOpts(tag *reflect.StructTag) SchemaOpts {
	options := trimSlice(strings.Split(tag.Get("schemaOpts"), ","))

	return SchemaOpts{
		Computed: contains(options, "computed"),
		Required: contains(options, "required"),
		Optional: contains(options, "optional"),
	}
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
