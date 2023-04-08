package commons

import (
	"reflect"
	"testing"
)

func TestCheckFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"Name", "string"},
		{"Description", "string"},
		{"Status", "string"},
		{"Id", "string"},
		{"Categories", "[]string"},
		{"Results", "[]commons.Result"},
		{"Duration", "time.Duration"},
		{"StartTime", "time.Time"},
		{"EndTime", "time.Time"},
	}

	checkType := reflect.TypeOf(Check{})
	if checkType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, but got %d fields", len(expectedFields), checkType.NumField())
	}

	for i, expectedField := range expectedFields {
		field := checkType.Field(i)
		if field.Name != expectedField.Name {
			t.Errorf("Expected field name '%s', but got '%s'", expectedField.Name, field.Name)
		}

		if field.Type.String() != expectedField.Type {
			t.Errorf("Expected field type '%s' for field '%s', but got '%s'", expectedField.Type, field.Name, field.Type.String())
		}
	}
}
