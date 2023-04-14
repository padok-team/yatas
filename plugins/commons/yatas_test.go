package commons

import (
	"sync"
	"testing"
)

func TestConfig_CheckExclude(t *testing.T) {
	type fields struct {
		Plugins []Plugin
		Ignore  []Ignore
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "check exclude",
			fields: fields{
				Plugins: []Plugin{
					{
						Exclude: []string{"test"},
					},
				},
			},
			args: args{
				id: "test",
			},
			want: true,
		},
		{
			name: "check exclude",
			fields: fields{
				Plugins: []Plugin{
					{
						Exclude: []string{"test"},
					},
				},
			},
			args: args{
				id: "toto",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Plugins: tt.fields.Plugins,
				Ignore:  tt.fields.Ignore,
			}
			if got := c.CheckExclude(tt.args.id); got != tt.want {
				t.Errorf("commons.CheckExclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_CheckInclude(t *testing.T) {
	type fields struct {
		Plugins      []Plugin
		Ignore       []Ignore
		PluginConfig interface{}
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "check include",
			fields: fields{
				Plugins: []Plugin{
					{
						Name:    "AWS",
						Include: []string{"AWS_TEST"},
					},
				},
			},
			args: args{
				id: "AWS_TEST",
			},
			want: true,
		},
		{
			name: "check include",
			fields: fields{
				Plugins: []Plugin{
					{
						Name:    "AWS",
						Include: []string{"AWS_TEST"},
					},
				},
			},
			args: args{
				id: "AWS_TOTO",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Plugins: tt.fields.Plugins,
				Ignore:  tt.fields.Ignore,
			}
			if got := c.CheckInclude(tt.args.id); got != tt.want {
				t.Errorf("commons.CheckInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckConfig_Init(t *testing.T) {
	type fields struct {
		Wg          *sync.WaitGroup
		Queue       chan Check
		ConfigYatas *Config
	}
	type args struct {
		config *Config
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "check config",
			fields: fields{
				Wg:    &sync.WaitGroup{},
				Queue: make(chan Check),
				ConfigYatas: &Config{
					Ignore: []Ignore{
						{
							ID: "test",
						},
					},
				},
			},
			args: args{
				config: &Config{
					Ignore: []Ignore{
						{
							ID: "test",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CheckConfig{
				Wg:          tt.fields.Wg,
				Queue:       tt.fields.Queue,
				ConfigYatas: tt.fields.ConfigYatas,
			}
			c.Init(tt.args.config)
			if c.ConfigYatas.Ignore[0].ID != tt.args.config.Ignore[0].ID {
				t.Errorf("CheckConfig.Init() ConfigYatas.Ignore[0].ID = %v, want %v", c.ConfigYatas.Ignore[0].ID, tt.args.config.Ignore[0].ID)
			}
		})
	}
}

func TestCheckTest(t *testing.T) {
	var wg sync.WaitGroup
	config := &Config{
		Plugins: []Plugin{
			{
				Name:    "test",
				Include: []string{"test"},
				Exclude: []string{},
			},
		},
	}

	id := "test"

	testFunc := func(a, b, c int) {
		wg.Done()
	}

	wrappedTest := CheckTest(&wg, config, id, testFunc)

	wrappedTest(1, 2, 3)
	wg.Wait()

	config.Plugins[0].Exclude = []string{"test"}

	wrappedTestExcluded := CheckTest(&wg, config, id, testFunc)

	wrappedTestExcluded(1, 2, 3)
}

func TestCheckMacroTest(t *testing.T) {
	var wg sync.WaitGroup
	config := &Config{
		Plugins: []Plugin{
			{
				Name:    "test",
				Include: []string{},
				Exclude: []string{},
			},
		},
	}

	testFunc := func(a, b, c, d int) {
		wg.Done()
	}

	wrappedTest := CheckMacroTest(&wg, config, testFunc)

	wrappedTest(1, 2, 3, 4)
	wg.Wait()
}

// MockResource to be used for testing
type MockResource struct {
	ID string
}

func (m *MockResource) GetID() string {
	return m.ID
}

func TestCreateCheck(t *testing.T) {
	checkDefinition := CheckDefinition{
		Title:          "AWS_TEST_001",
		Description:    "Test check",
		Categories:     []string{"Test"},
		ConditionFn:    nil,
		SuccessMessage: "Test check succeeded",
		FailureMessage: "Test check failed",
	}

	check := createCheck(checkDefinition)

	if check.Description != checkDefinition.Description {
		t.Errorf("Expected Description: %v, got: %v", checkDefinition.Description, check.Description)
	}

	if check.Id != checkDefinition.Title {
		t.Errorf("Expected Title: %v, got: %v", checkDefinition.Title, check.Id)
	}

	for i, tag := range check.Categories {
		if tag != checkDefinition.Categories[i] {
			t.Errorf("Expected Tag: %v, got: %v", checkDefinition.Categories[i], tag)
		}
	}
}

func TestCheckResource(t *testing.T) {
	resource := &MockResource{ID: "test-resource-id"}
	conditionFn := func(r Resource) bool {
		return r.GetID() == "test-resource-id"
	}
	successMessage := "Resource check succeeded"
	failureMessage := "Resource check failed"

	result := checkResource(resource, conditionFn, successMessage, failureMessage)

	if result.Status != "OK" {
		t.Errorf("Expected Status: OK, got: %v", result.Status)
	}

	expectedMessage := successMessage + " - Resource " + resource.GetID()
	if result.Message != expectedMessage {
		t.Errorf("Expected Message: %v, got: %v", expectedMessage, result.Message)
	}

	if result.ResourceID != resource.GetID() {
		t.Errorf("Expected ResourceID: %v, got: %v", resource.GetID(), result.ResourceID)
	}

	// Test failing condition
	conditionFn = func(r Resource) bool {
		return r.GetID() != "test-resource-id"
	}

	result = checkResource(resource, conditionFn, successMessage, failureMessage)

	if result.Status != "FAIL" {
		t.Errorf("Expected Status: FAIL, got: %v", result.Status)
	}

	expectedMessage = failureMessage + " - Resource " + resource.GetID()
	if result.Message != expectedMessage {
		t.Errorf("Expected Message: %v, got: %v", expectedMessage, result.Message)
	}

	if result.ResourceID != resource.GetID() {
		t.Errorf("Expected ResourceID: %v, got: %v", resource.GetID(), result.ResourceID)
	}
}
