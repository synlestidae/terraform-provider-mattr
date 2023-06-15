package main

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type SchemaOpts struct {
	Computed bool
	Required bool
	Optional bool
}

func schemaForStruct(inputType reflect.Type) (*map[string]*schema.Schema, error) {
	if inputType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("Unable to generate schema for kind: %s", inputType.Kind())
	}

	schemaMap := make(map[string]*schema.Schema, inputType.NumField())
	numField := inputType.NumField()

	for i := 0; i < numField; i++ {
		field := inputType.Field(i)
		name := snakeCase(field.Name)
		schemaType, err := getSchemaType(field.Type)
		opts := fieldOpts(&field.Tag)

		if err != nil {
			return nil, err
		}

		newSchema := schema.Schema{
			Type:     schemaType,
			Computed: opts.Computed,
			Required: opts.Required,
			Optional: opts.Optional,
		}

		// now the recursive parts
		if field.Type.Kind() == reflect.Struct {
			nestedSchema, err := schemaForStruct(field.Type)

			if err != nil {
				return nil, err
			}

			newSchema.Elem = &schema.Resource{
				Schema: *nestedSchema,
			}
		}

		schemaMap[name] = &newSchema
	}

	return &schemaMap, nil
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

func getSchemaType(inputType reflect.Type) (schema.ValueType, error) {
	fmt.Printf("Thing %s", inputType)
	switch inputType.Kind() {
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
	case reflect.Interface:
	case reflect.Slice:
		//return nil, fmt.Errorf("Not yet implemented")
		panic("Not yet implemented")
	case reflect.Chan:
	case reflect.Pointer:
	case reflect.UnsafePointer:
		panic(fmt.Sprintf("Unsupported schema for type: %s", inputType.Kind()))
		//return nil, fmt.Errorf("Unsupported schema for type: %s", inputType.Kind())
		//default:
		//panic("Unknown type: %s", inputType.Kind())
		//return nil, fmt.Errorf("Unknown type: %s", inputType.Kind())
	}
	panic(fmt.Sprintf("Unsupported schema for type: %s", inputType.Kind()))
}
