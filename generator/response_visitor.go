package generator

import (
	"fmt"
	"reflect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResponseVisitor struct {
	resourceData *schema.ResourceData
	id string
}

func (v *ResponseVisitor) accept(data interface{}) (interface{}, error) {
	if data, ok := data.(map[string]interface{}); ok {
		return v.visitMap(data)
	}

	if data, ok := data.([]interface{}); ok {
		return v.visitList(data)
	}

	if data, ok := data.(*interface{}); ok {
		return v.accept(*data)
	}

	return v.visitPrimitive(data)
}

func (rv *ResponseVisitor) visitMap(data map[string]interface{}) (interface{}, error) {
	for key, value := range data {
		schemaVal, err := rv.accept(value)
		if err != nil {
			return nil, err
		}

		schemaName := snakeCase(key)
		if schemaName != "id" { 
			if err := rv.resourceData.Set(schemaName, schemaVal); err != nil {
				return nil, err
			}
		} else {
			rv.id = schemaVal.(string) // todo 
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
	switch data := data.(type) {
	case int64:
		return data, nil
	case float64:
		return data, nil
	case string:
		return data, nil
	case bool:
		return data, nil
	}
	return nil, fmt.Errorf("Unable to accept value of type %s in response", reflect.TypeOf(data))
}
