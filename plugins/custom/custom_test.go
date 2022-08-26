package custom

import (
	"reflect"
	"testing"

	"github.com/stangirard/yatas/internal/yatas"
)

func Test_findPluginWithName(t *testing.T) {
	type args struct {
		c    *yatas.Config
		name string
	}
	tests := []struct {
		name string
		args args
		want *yatas.Plugin
	}{
		{
			name: "find plugin with name",
			args: args{
				c: &yatas.Config{
					Plugins: []yatas.Plugin{
						{
							Name: "test",
						},
					},
				},
				name: "test",
			},
			want: &yatas.Plugin{
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
		c    *yatas.Config
		name string
	}
	tests := []struct {
		name string
		args args
		want *yatas.Plugin
	}{
		{
			name: "find plugin with name",
			args: args{
				c: &yatas.Config{
					Plugins: []yatas.Plugin{
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
