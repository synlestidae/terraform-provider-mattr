package main

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestSchemaForStruct(t *testing.T) {
	type MyStruct struct {
		Field1 string `schemaOpts:"required"`
		Field2 int    `schemaOpts:"optional"`
	}

	expectedSchema := map[string]*schema.Schema{
		"field1": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
			Optional: false,
		},
		"field2": {
			Type:     schema.TypeInt,
			Computed: false,
			Required: false,
			Optional: true,
		},
	}

	actualSchema, err := schemaForStruct(reflect.TypeOf(MyStruct{}))
	if err != nil {
		t.Errorf("Error generating schema: %v", err)
		return
	}

	assertEqual(t, *actualSchema, expectedSchema)

	//if !reflect.DeepEqual(*actualSchema, expectedSchema) {
	//t.Errorf("MyFunction() returned incorrect result.\nGot:\n%v\nExpected:\n%v", actualSchema, expectedSchema)
	//}

	//if !reflect.DeepEqual(*actualSchema, expectedSchema) {
	//t.Errorf("Generated schema does not match the expected schema")
	//}
}

func TestSnakeCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"FieldName", "field_name"},
		{"UserID", "user_id"},
		{"SomeLongName", "some_long_name"},
	}

	for _, tc := range testCases {
		actual := snakeCase(tc.input)
		if actual != tc.expected {
			t.Errorf("snakeCase(%s) = %s, expected %s", tc.input, actual, tc.expected)
		}
	}
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("Values are not equal:\n%s", diff)
	}
}
