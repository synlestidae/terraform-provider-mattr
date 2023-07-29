package generator

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RequestVisitor struct {
	schema map[string]*schema.Schema
}

func (v *RequestVisitor) accept(data interface{}) (interface{}, error) {
	t := reflect.TypeOf(data)

	if data, ok := data.(*schema.ResourceData); ok {
		return v.visitResourceData(data)
	}
	if data, ok := data.(map[string]interface{}); ok {
		return v.visitMap(data)
	}
	if data, ok := data.([]interface{}); ok {
		return v.visitList(data)
	}
	if data, ok := (data).(*schema.Set); ok {
		return v.visitSet(data)
	}
	if isPrimitive(t.Kind()) {
		return v.visitPrimitive(data)
	}

	return nil, fmt.Errorf("Unable to convert value of type %T to request", data)
}

func (rv *RequestVisitor) visitSet(s *schema.Set) (interface{}, error) {
	list := make([]interface{}, 0, len(s.List()))
	for _, v := range s.List() {
		list = append(list, v)
	}
	return rv.visitList(list)
}

func (rv *RequestVisitor) visitResourceData(rd *schema.ResourceData) (interface{}, error) {
	req := make(map[string]interface{})

	for key, _ := range rv.schema {
		value := rd.Get(key).(interface{})
		reqVal, err := rv.accept(value)
		if err != nil {
			return nil, err
		}

		if reqVal != "" && reqVal != nil {
			req[snakeToCamel(key)] = reqVal
		}
	}

	return req, nil
}

func (rv *RequestVisitor) visitMap(data map[string]interface{}) (interface{}, error) {
	req := make(map[string]interface{}, len(data))

	for key, value := range data {
		reqVal, err := rv.accept(value)
		if err != nil {
			return nil, err
		}

		if reqVal != nil && reqVal != "" && reqVal != nil {
			req[snakeToCamel(key)] = reqVal
		}
	}

	return req, nil
}

func (rv *RequestVisitor) visitList(data []interface{}) (interface{}, error) {
	req := make([]interface{}, len(data))

	for i, value := range data {
		reqVal, err := rv.accept(value)
		if err != nil {
			return nil, err
		}

		req[i] = reqVal
	}

	if len(req) == 0 {
		return nil, nil
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
