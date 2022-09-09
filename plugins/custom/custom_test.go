package custom

import (
	"reflect"
	"testing"

	"github.com/stangirard/yatas/plugins/commons"
)

func Test_findPluginWithName(t *testing.T) {
	type args struct {
		c    *commons.Config
		name string
	}
	tests := []struct {
		name string
		args args
		want *commons.Plugin
	}{
		{
			name: "find plugin with name",
			args: args{
				c: &commons.Config{
					Plugins: []commons.Plugin{
						{
							Name: "test",
						},
					},
				},
				name: "test",
			},
			want: &commons.Plugin{
				Name: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findPluginWithName(tt.args.c, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findPluginWithName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findPluginWithNameFail(t *testing.T) {
	type args struct {
		c    *commons.Config
		name string
	}
	tests := []struct {
		name string
		args args
		want *commons.Plugin
	}{
		{
			name: "find plugin with name",
			args: args{
				c: &commons.Config{
					Plugins: []commons.Plugin{
						{
							Name: "test",
						},
					},
				},
				name: "toto",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findPluginWithName(tt.args.c, tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findPluginWithName() = %v, want %v", got, tt.want)
			}
		})
	}
}
