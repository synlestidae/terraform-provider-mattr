package generator

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"reflect"
)

type ResponseVisitor struct {
	resourceData *schema.ResourceData
}

func (v *ResponseVisitor) accept(data interface{}) (interface{}, error) {
	fmt.Printf("Printing")
	t := reflect.TypeOf(data)

	if data, ok := data.(map[string]interface{}); ok {
		return v.visitMap(data)
	}

	if data, ok := data.([]interface{}); ok {
		return v.visitList(data)
	}

	if isPrimitive(t.Kind()) {
		return v.visitPrimitive(data)
	}

	return nil, fmt.Errorf("Unable to accept value of type %T in response", data, t)
}

func (rv *ResponseVisitor) visitMap(data map[string]interface{}) (interface{}, error) {
	for key, value := range data {
		schemaVal, err := rv.accept(value)
		if err != nil {
			return nil, err
		}

		schemaName := snakeCase(key)
		if err := rv.resourceData.Set(schemaName, schemaVal); err != nil {
			return nil, err
		}
	}

	return rv.resourceData, nil
}

func (rv *ResponseVisitor) visitList(data []interface{}) (interface{}, error) {
	list := make([]interface{}, len(data))

	for i, elem := range data {
		value, err := rv.accept(elem)
		if err != err {
			return nil, err
		}

		list[i] = value
	}

	return list, nil
}

func (rv ResponseVisitor) visitPrimitive(data interface{}) (interface{}, error) {
	return data, nil
}
