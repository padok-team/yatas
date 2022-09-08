package eks

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfLoggingIsEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		clusters    []types.Cluster
		testName    string
		want        string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				clusters: []types.Cluster{
					{
						Name: aws.String("test"),
						Logging: &types.Logging{
							ClusterLogging: []types.LogSetup{
								{
									Enabled: aws.Bool(true),
									Types:   []types.LogType{"api", "audit"},
								},
							},
						},
					},
				},
				testName: "CheckIfLoggingIsEnabled",
				want:     "OK",
			},
		},
		{
			name: "test",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				clusters: []types.Cluster{
					{
						Name: aws.String("test"),
					},
				},
				testName: "CheckIfLoggingIsEnabled",
				want:     "FAIL",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfLoggingIsEnabled(tt.args.checkConfig, tt.args.clusters, tt.args.testName)
		})
	}
}
