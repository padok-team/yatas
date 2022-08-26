package apigateway

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfTracingEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		stages      []types.Stage
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if all stages are tracing enabled",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				stages: []types.Stage{
					{
						TracingEnabled: true,
						StageName:      aws.String("test"),
					},
				},
				testName: "CheckIfTracingEnabled",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfTracingEnabled(tt.args.checkConfig, tt.args.stages, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if len(check.Results) != len(tt.args.stages) {
						t.Errorf("CheckIfTracingEnabled() = %v, want %v", len(check.Results), len(tt.args.stages))
					}
					if check.Status != "OK" {
						t.Errorf("CheckIfTracingEnabled() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
