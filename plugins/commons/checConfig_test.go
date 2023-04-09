package commons

import (
	"testing"
	"time"
)

func TestCheck_AddResult(t *testing.T) {
	type fields struct {
		Name        string
		Description string
		Status      string
		Id          string
		Results     []Result
	}
	type args struct {
		result Result
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add result",
			fields: fields{
				Name:        "test",
				Description: "test",
				Status:      "OK",
				Id:          "test",
				Results:     []Result{},
			},
		},
		{
			name: "add result",
			fields: fields{
				Name:        "test",
				Description: "test",
				Status:      "OK",
				Id:          "test",
				Results:     []Result{},
			},
			args: args{
				result: Result{
					Message:    "test",
					Status:     "OK",
					ResourceID: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Status:      tt.fields.Status,
				Id:          tt.fields.Id,
				Results:     tt.fields.Results,
			}
			c.AddResult(tt.args.result)
			if c.Status != "OK" {
				t.Errorf("Status is not OK")
			}
		})
	}
}
func TestCheck_AddResultFail(t *testing.T) {
	type fields struct {
		Name        string
		Description string
		Status      string
		Id          string
		Results     []Result
	}
	type args struct {
		result Result
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "add result",
			fields: fields{
				Name:        "test",
				Description: "test",
				Status:      "OK",
				Id:          "test",
				Results:     []Result{},
			},
			args: args{
				result: Result{
					Message:    "test",
					Status:     "FAIL",
					ResourceID: "test",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Status:      "OK",
				Id:          tt.fields.Id,
				Results:     tt.fields.Results,
			}
			c.AddResult(tt.args.result)
			if c.Status != "FAIL" {
				t.Errorf("Status is not FAIL")
			}
		})
	}
}

func TestCheck_InitCheck(t *testing.T) {
	type fields struct {
		Name        string
		Description string
		Status      string
		Id          string
		Results     []Result
	}
	type args struct {
		name        string
		description string
		id          string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "init check",
			fields: fields{
				Name:        "test",
				Description: "test",
				Status:      "OK",
				Id:          "test",
				Results:     []Result{},
			},
		},
		{
			name: "init check",
			fields: fields{
				Name:        "test",
				Description: "test",
				Status:      "OK",
				Id:          "test",
				Results:     []Result{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Check{
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Status:      tt.fields.Status,
				Id:          tt.fields.Id,
				Results:     tt.fields.Results,
			}
			categories := []string{"test"}
			c.InitCheck(tt.args.name, tt.args.description, tt.args.id, categories)
			if c.Name != tt.args.name {
				t.Errorf("Name is not %s", tt.args.name)
			}
			if c.Description != tt.args.description {
				t.Errorf("Description is not %s", tt.args.description)
			}
			if c.Id != tt.args.id {
				t.Errorf("Id is not %s", tt.args.id)
			}
		})
	}
}

func TestCheckConfig_Init2(t *testing.T) {
	cfg := &Config{}
	checkCfg := &CheckConfig{}
	checkCfg.Init(cfg)

	if checkCfg.Wg == nil {
		t.Error("CheckConfig Init failed to initialize WaitGroup")
	}

	if checkCfg.Queue == nil {
		t.Error("CheckConfig Init failed to initialize Queue")
	}

	if checkCfg.ConfigYatas != cfg {
		t.Error("CheckConfig Init failed to assign ConfigYatas")
	}
}

func TestCheck_EndCheck2(t *testing.T) {
	check := Check{
		Name:        "testName",
		Description: "testDescription",
		Status:      "OK",
		Id:          "testID",
		StartTime:   time.Now(),
	}

	time.Sleep(100 * time.Millisecond)
	check.EndCheck()

	if check.EndTime.Before(check.StartTime) {
		t.Error("EndTime is not after StartTime")
	}

	if check.Duration < 100*time.Millisecond {
		t.Error("Duration is incorrect")
	}
}

func TestCheck_InitCheckAndUpdateStatus2(t *testing.T) {
	check := Check{}
	check.InitCheck("testName", "testDescription", "testID", []string{"TestCategory"})

	if check.Name != "testName" {
		t.Errorf("Name should be testName, got %s", check.Name)
	}

	if check.Description != "testDescription" {
		t.Errorf("Description should be testDescription, got %s", check.Description)
	}

	if check.Status != "OK" {
		t.Errorf("Status should be OK, got %s", check.Status)
	}

	if check.Id != "testID" {
		t.Errorf("Id should be testID, got %s", check.Id)
	}

	check.AddResult(Result{
		Message:    "Test message",
		Status:     "FAIL",
		ResourceID: "testResourceID",
	})

	if check.Status != "FAIL" {
		t.Errorf("Status should be FAIL after adding a failing result, got %s", check.Status)
	}
}
