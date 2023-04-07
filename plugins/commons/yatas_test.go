package commons

import (
	"fmt"
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

func TestParseConfig(t *testing.T) {
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
