package cli

import (
	"testing"

	"github.com/padok-team/yatas/plugins/commons"
)

// Mock configuration for testing purposes
var mockConfig = commons.Config{
	Plugins: []commons.Plugin{
		{
			Name:    "TestPlugin",
			Type:    "checks",
			Version: "1.0.0",
			Source:  "github.com/padok-team/yatas",
		},
		{
			Name:    "TestPlugin",
			Type:    "checks",
			Version: "latest",
			Source:  "github.com/padok-team/yatas",
		},
	},
}

func TestInitialisePlugins(t *testing.T) {
	err := initialisePlugins(mockConfig)
	if err != nil {
		t.Errorf("Error initializing plugins: %s", err)
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Execute(); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Generate tests for ciReporting function
// func ciReporting(checks []commons.Tests) {
// 	if *ci {
// 		os.Exit(report.ExitCode(checks))
// 	}
// }

func Test_ciReporting(t *testing.T) {
	type args struct {
		checks []commons.Tests
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				checks: []commons.Tests{
					{
						Account: "test",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// check if the function exits with 1
			if ciReporting(tt.args.checks); *ci {
				t.Errorf("ciReporting() = %v, want %v", *ci, false)
			}

		})
	}
}
