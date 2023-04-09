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

func TestResultFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"Message", "string"},
		{"Status", "string"},
		{"ResourceID", "string"},
	}

	resultType := reflect.TypeOf(Result{})
	if resultType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, but got %d fields", len(expectedFields), resultType.NumField())
	}

	for i, expectedField := range expectedFields {
		field := resultType.Field(i)
		if field.Name != expectedField.Name {
			t.Errorf("Expected field name '%s', but got '%s'", expectedField.Name, field.Name)
		}

		if field.Type.String() != expectedField.Type {
			t.Errorf("Expected field type '%s' for field '%s', but got '%s'", expectedField.Type, field.Name, field.Type.String())
		}
	}
}

func TestTestsFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"Account", "string"},
		{"Checks", "[]commons.Check"},
	}

	testsType := reflect.TypeOf(Tests{})
	if testsType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, but got %d fields", len(expectedFields), testsType.NumField())
	}

	for i, expectedField := range expectedFields {
		field := testsType.Field(i)
		if field.Name != expectedField.Name {
			t.Errorf("Expected field name '%s', but got '%s'", expectedField.Name, field.Name)
		}

		if field.Type.String() != expectedField.Type {
			t.Errorf("Expected field type '%s' for field '%s', but got '%s'", expectedField.Type, field.Name, field.Type.String())
		}
	}
}

func TestConfigFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"Plugins", "[]commons.Plugin"},
		{"Ignore", "[]commons.Ignore"},
		{"PluginConfig", "[]map[string]interface {}"},
		{"Tests", "[]commons.Tests"},
	}

	configType := reflect.TypeOf(Config{})
	if configType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, but got %d fields", len(expectedFields), configType.NumField())
	}

	for i, expectedField := range expectedFields {
		field := configType.Field(i)
		if field.Name != expectedField.Name {
			t.Errorf("Expected field name '%s', but got '%s'", expectedField.Name, field.Name)
		}

		if field.Type.String() != expectedField.Type {
			t.Errorf("Expected field type '%s' for field '%s', but got '%s'", expectedField.Type, field.Name, field.Type.String())
		}
	}
}

func TestCheckConfigFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"Wg", "*sync.WaitGroup"},
		{"Queue", "chan commons.Check"},
		{"ConfigYatas", "*commons.Config"},
	}

	checkConfigType := reflect.TypeOf(CheckConfig{})
	if checkConfigType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, but got %d fields", len(expectedFields), checkConfigType.NumField())
	}

	for i, expectedField := range expectedFields {
		field := checkConfigType.Field(i)
		if field.Name != expectedField.Name {
			t.Errorf("Expected field name '%s', but got '%s'", expectedField.Name, field.Name)
		}

		if field.Type.String() != expectedField.Type {
			t.Errorf("Expected field type '%s' for field '%s', but got '%s'", expectedField.Type, field.Name, field.Type.String())
		}
	}
}

func TestIgnoreFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"ID", "string"},
		{"Regex", "bool"},
		{"Values", "[]string"},
	}

	ignoreType := reflect.TypeOf(Ignore{})
	if ignoreType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, but got %d fields", len(expectedFields), ignoreType.NumField())
	}

	for i, expectedField := range expectedFields {
		field := ignoreType.Field(i)
		if field.Name != expectedField.Name {
			t.Errorf("Expected field name '%s', but got '%s'", expectedField.Name, field.Name)
		}

		if field.Type.String() != expectedField.Type {
			t.Errorf("Expected field type '%s' for field '%s', but got '%s'", expectedField.Type, field.Name, field.Type.String())
		}
	}
}

func TestPluginFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"Name", "string"},
		{"Enabled", "bool"},
		{"Source", "string"},
		{"Type", "string"},
		{"Version", "string"},
		{"Description", "string"},
		{"Exclude", "[]string"},
		{"Include", "[]string"},
		{"Command", "string"},
		{"Args", "[]string"},
		{"ExpectedOutput", "string"},
		{"ExpectedStatus", "int"},
		{"SourceOwner", "string"},
		{"SourceRepo", "string"},
	}

	pluginType := reflect.TypeOf(Plugin{})
	if pluginType.NumField() != len(expectedFields) {
		t.Errorf("Expected %d fields, but got %d fields", len(expectedFields), pluginType.NumField())
	}

	for i, expectedField := range expectedFields {
		field := pluginType.Field(i)
		if field.Name != expectedField.Name {
			t.Errorf("Expected field name '%s', but got '%s'", expectedField.Name, field.Name)
		}

		if field.Type.String() != expectedField.Type {
			t.Errorf("Expected field type '%s' for field '%s', but got '%s'", expectedField.Type, field.Name, field.Type.String())
		}
	}
}

func TestYatasInterface(t *testing.T) {
	yatasInterface := reflect.TypeOf((*Yatas)(nil)).Elem()

	if yatasInterface.NumMethod() != 1 {
		t.Errorf("Expected 1 method, but got %d methods", yatasInterface.NumMethod())
	}

	method := yatasInterface.Method(0)
	expectedMethodName := "Run"
	expectedMethodType := "func(*commons.Config) []commons.Tests"

	if method.Name != expectedMethodName {
		t.Errorf("Expected method name '%s', but got '%s'", expectedMethodName, method.Name)
	}

	if method.Type.String() != expectedMethodType {
		t.Errorf("Expected method type '%s', but got '%s'", expectedMethodType, method.Type.String())
	}
}

func TestYatasPluginFields(t *testing.T) {
	expectedFields := []struct {
		Name string
		Type string
	}{
		{"Impl", "commons.Yatas"},
	}

	typesToCheck := []struct {
		Type     reflect.Type
		TypeName string
	}{
		{reflect.TypeOf(YatasPlugin{}), "YatasPlugin"},
		{reflect.TypeOf(YatasRPCServer{}), "YatasRPCServer"},
	}

	for _, typeToCheck := range typesToCheck {
		if typeToCheck.Type.NumField() != len(expectedFields) {
			t.Errorf("Expected %d fields for %s, but got %d fields", len(expectedFields), typeToCheck.TypeName, typeToCheck.Type.NumField())
		}

		for i, expectedField := range expectedFields {
			field := typeToCheck.Type.Field(i)
			if field.Name != expectedField.Name {
				t.Errorf("Expected field name '%s' for %s, but got '%s'", expectedField.Name, typeToCheck.TypeName, field.Name)
			}

			if field.Type.String() != expectedField.Type {
				t.Errorf("Expected field type '%s' for field '%s' in %s, but got '%s'", expectedField.Type, field.Name, typeToCheck.TypeName, field.Type.String())
			}
		}
	}
}
