package commons

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestConfig_CheckExclude(t *testing.T) {
	type fields struct {
		Plugins []Plugin
		AWS     []AWS_Account
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
				AWS:     tt.fields.AWS,
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
		Plugins []Plugin
		AWS     []AWS_Account
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
				AWS:     tt.fields.AWS,
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
			_, err := ParseConfig(tt.args.configFile)
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
		ConfigAWS   aws.Config
		Queue       chan Check
		ConfigYatas *Config
	}
	type args struct {
		s      aws.Config
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
				Wg:          &sync.WaitGroup{},
				ConfigAWS:   aws.Config{},
				Queue:       make(chan Check),
				ConfigYatas: &Config{},
			},
			args: args{
				s: aws.Config{
					Region: "eu-west-1",
				},
				config: &Config{
					AWS: []AWS_Account{
						{
							Name: "test",
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
				ConfigAWS:   tt.fields.ConfigAWS,
				Queue:       tt.fields.Queue,
				ConfigYatas: tt.fields.ConfigYatas,
			}
			c.Init(tt.args.s, tt.args.config)
			if c.ConfigAWS.Region != tt.args.s.Region {
				t.Errorf("CheckConfig.Init() ConfigAWS.Region = %v, want %v", c.ConfigAWS.Region, tt.args.s.Region)
			}
			if c.ConfigYatas.AWS[0].Name != tt.args.config.AWS[0].Name {
				t.Errorf("CheckConfig.Init() ConfigYatas.AWS[0].Name = %v, want %v", c.ConfigYatas.AWS[0].Name, tt.args.config.AWS[0].Name)
			}

		})
	}
}
