package provider

import (
	"reflect"
	"testing"
)
// Helper function to compare values and fail the test if not equal

func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("%s. Expected: %v, but got: %v", msg, expected, actual)
	}
}
