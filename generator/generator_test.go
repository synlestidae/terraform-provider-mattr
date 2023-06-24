package generator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSchemaOneFieldInt(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"foo": {
			Type:     schema.TypeInt,
			Required: true,
		},
	}

	// Create a new instance of ResourceData
	rd := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"foo": 123,
	})

	rv := RequestVisitor{}
	value, err := rv.visitResource(rd, resourceSchema)
	if err != nil {
		t.Errorf("Error generating map: %s", err)
	}

	assertEqual(t, value, map[string]interface{}{
		"foo": 123,
	})
}

func TestSchemaOneFieldString(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"foo": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	// Create a new instance of ResourceData
	rd := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"foo": "one two three",
	})

	rv := RequestVisitor{}
	value, err := rv.visitResource(rd, resourceSchema)
	if err != nil {
		t.Errorf("Error generating map: %s", err)
	}

	assertEqual(t, value, map[string]interface{}{
		"foo": "one two three",
	})
}

func TestSchemaOneFieldStringCase(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"foo_bar": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	// Create a new instance of ResourceData
	rd := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"foo_bar": "one two three",
	})

	rv := RequestVisitor{}
	value, err := rv.visitResource(rd, resourceSchema)
	if err != nil {
		t.Errorf("Error generating map: %s", err)
	}

	assertEqual(t, value, map[string]interface{}{
		"fooBar": "one two three",
	})
}

func TestSchemaMapList(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"foo_bar": {
			Type:     schema.TypeList,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Required: true,
		},
	}

	var fooBarString []string
	fooBarString = []string{ "foo", "bar" }
	// Create a new instance of ResourceData
	rd := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{})
	if err := rd.Set("foo_bar", fooBarString); err != nil {
		t.Errorf("Error setting foo_bar: %s", err)
	}

	rv := RequestVisitor{}
	value, err := rv.visitResource(rd, resourceSchema)
	if err != nil {
		t.Errorf("Error generating map: %s", err)
	}

	assertEqual(t, value, map[string]interface{}{
		"fooBar": []interface{}{ "foo", "bar" },
	})
}

func TestSchemaMapMap(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"foo_bar": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"background_color": &schema.Schema{
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
			Required: true,
		},
	}

	fooBarMap := map[string]string{
		"background_color":"green",
	}
	// Create a new instance of ResourceData
	rd := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{})
	if err := rd.Set("foo_bar", fooBarMap); err != nil {
		t.Errorf("Error setting foo_bar: %s", err)
	}

	rv := RequestVisitor{}
	value, err := rv.visitResource(rd, resourceSchema)
	if err != nil {
		t.Errorf("Error generating map: %s", err)
	}

	assertEqual(t, value, map[string]any{
		"fooBar": map[string]any{
			"backgroundColor":"green",
		},
	})
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
