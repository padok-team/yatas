package config

import (
	"testing"
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
			c.InitCheck(tt.args.name, tt.args.description, tt.args.id)
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
