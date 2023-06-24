package generator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Visitor interface {
	visitResource(*schema.ResourceData, map[string]*schema.Schema) (interface{}, error)
	visitMap(map[string]interface{}) (interface{}, error)
	visitList([]interface{}) (interface{}, error)
	visitPrimitive(interface{}) (interface{}, error)

}

type RequestVisitor struct {
}

func accept(v Visitor, data interface{}) (interface{}, error) {
	t := reflect.TypeOf(data)

	if data, ok := data.(map[string]interface{}); ok {
		return v.visitMap(data)
	}
	if data, ok := data.([]interface{}); ok {
		return v.visitList(data)
	}
	if isPrimitive(t.Kind()) {
		return v.visitPrimitive(data);
	}

	return nil, fmt.Errorf("Unable to convert value of type: %T", data)
}

func (rv *RequestVisitor) visitResource(rd *schema.ResourceData, sch map[string]*schema.Schema) (interface{}, error) {
	req := make(map[string]interface{})

	for key, _ := range sch {
		value := rd.Get(key).(interface{})
		reqVal, err := accept(rv, value)	
		if err != nil {
			return nil, err;
		}

		req[snakeToCamel(key)] = reqVal
	}

	return req, nil
}

func (rv *RequestVisitor) visitMap(data map[string]interface{}) (interface{}, error) {
	req := make(map[string]interface{}, len(data))

	for key, value := range data {
		fmt.Printf("Key %s", key)
		reqVal, err := accept(rv, value)
		if err != nil {
			return nil, err
		}

		req[snakeToCamel(key)] = reqVal
	}

	return req, nil
}

func (rv *RequestVisitor) visitList(data []interface{}) (interface{}, error) {
	req := make([]interface{}, len(data))	

	for i, value := range data {
		reqVal, err := accept(rv, value)
		if err != nil {
			return nil, err
		}

		req[i] = reqVal
	}
	
	return req, nil
}

func (rv RequestVisitor) visitPrimitive(data interface{}) (interface{}, error) {
	return data, nil
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

func snakeToCamel(snakeCase string) string {
	words := strings.Split(snakeCase, "_")
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}
