package commons

import (
	"fmt"
	"testing"
)

func TestParseConfig(t *testing.T) {
	config, err := ParseConfig("../../.yatas.yml.example")
	if err != nil {
		t.Errorf("Error parsing the config file: %s", err)
	}

	// You can add more tests for specific values in your config file here.
	if len(config.Plugins) == 0 {
		t.Error("Expected non-zero plugins, got zero")
	}
}

func TestParseConfig2(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "parse config",
			args: args{
				configFile: "../../.yatas.yml.example",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseConfig(tt.args.configFile)
			fmt.Println(config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCheckExclude(t *testing.T) {
	config := Config{
		Plugins: []Plugin{
			{
				Name:    "TestPlugin",
				Exclude: []string{"TESTPLUGIN_EXAMPLE"},
			},
		},
	}

	if !config.CheckExclude("TESTPLUGIN_EXAMPLE") {
		t.Error("Expected CheckExclude to return true for 'TESTPLUGIN_EXAMPLE', got false")
	}

	if config.CheckExclude("TESTPLUGIN_OTHER") {
		t.Error("Expected CheckExclude to return false for 'TESTPLUGIN_OTHER', got true")
	}
}

func TestCheckInclude(t *testing.T) {
	config := Config{
		Plugins: []Plugin{
			{
				Name:    "TestPlugin",
				Include: []string{"TESTPLUGIN_EXAMPLE"},
			},
		},
	}

	if !config.CheckInclude("TESTPLUGIN_EXAMPLE") {
		t.Error("Expected CheckInclude to return true for 'TESTPLUGIN_EXAMPLE', got false")
	}

	if config.CheckInclude("TESTPLUGIN_OTHER") {
		t.Error("Expected CheckInclude to return false for 'TESTPLUGIN_OTHER', got true")
	}
}

func TestFindPluginWithName(t *testing.T) {
	config := Config{
		Plugins: []Plugin{
			{
				Name: "TestPlugin",
			},
		},
	}

	plugin := config.FindPluginWithName("TestPlugin")
	if plugin == nil {
		t.Error("Expected to find a plugin with name 'TestPlugin', got nil")
	}

	plugin = config.FindPluginWithName("NonExistentPlugin")
	if plugin != nil {
		t.Error("Expected to not find a plugin with name 'NonExistentPlugin', got a plugin")
	}
}

func TestCheckHasHDSCategory(t *testing.T) {
	config := Config{}

	checkWithHDS := Check{
		Categories: []string{"Security", "HDS", "Performance"},
	}

	checkWithoutHDS := Check{
		Categories: []string{"Security", "Performance"},
	}

	checkWithLowercaseHDS := Check{
		Categories: []string{"Security", "hds", "Performance"},
	}

	if !config.CheckHasHDSCategory(checkWithHDS) {
		t.Error("Expected CheckHasHDSCategory to return true for check with HDS category, got false")
	}

	if config.CheckHasHDSCategory(checkWithoutHDS) {
		t.Error("Expected CheckHasHDSCategory to return false for check without HDS category, got true")
	}

	if !config.CheckHasHDSCategory(checkWithLowercaseHDS) {
		t.Error("Expected CheckHasHDSCategory to return true for check with lowercase hds category, got false")
	}
}
