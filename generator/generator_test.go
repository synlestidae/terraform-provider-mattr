package generator

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSchemaForStruct(t *testing.T) {
	type MyStruct struct {
		Field1 string  `schemaOpts:"required"`
		Field2 int     `schemaOpts:"optional"`
		Field3 float32 `schemaOpts:"computed"`
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
		"field3": {
			Type:     schema.TypeFloat,
			Computed: true,
			Required: false,
			Optional: false,
		},
	}

	doSchemaTest(t, reflect.TypeOf(MyStruct{}), expectedSchema)
	/*var rc ResourceContainer

	actualSchema, err := rc.genSchema(reflect.TypeOf(MyStruct{}))
	if err != nil {
		t.Errorf("Error generating schema: %v", err)
		return
	}

	assertEqual(t, *actualSchema, expectedSchema)*/
}

func TestRequestInt(t *testing.T) {
	res := ResourceRep{
		kind:      reflect.Struct,
		valueType: schema.TypeMap,
		fields: []Field{
			{
				schemaName: "test_field",
				fieldName:  "testField",
				resource: ResourceRep{
					kind:      reflect.Int,
					valueType: schema.TypeInt,
				},
			},
		},
	}

	data := schema.ResourceData{}
	data.Set("test_field", 123)

	doRequestTest(t, &res, &data, map[string]interface{}{"testField": 123})
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

func doSchemaTest(t *testing.T, typ reflect.Type, expectedSchema map[string]*schema.Schema) {
	var visitor ResourceVisitor
	res, err := typeResource(typ)
	if err = res.accept(&visitor); err != nil {
		t.Errorf("Failed to produce ResourceRep: %s", err)
	}

	assertEqual(t, expectedSchema, visitor.schema.Elem)
}

func doRequestTest(t *testing.T, rp *ResourceRep, data *schema.ResourceData, expected interface{}) {
	rv := RequestVisitor{
		data: data,
	}

	if err := rp.accept(&rv); err != nil {
		t.Errorf("Failed to produce request: %s", err)
	}

	assertEqual(t, rv.value, expected)
}

func assertEqual(t *testing.T, actual, expected interface{}) {
	t.Helper()

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("Values are not equal:\n%s", diff)
	}
}