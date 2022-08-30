package eks

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfEKSUpdateAvailable(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		clusters    []ClusterToUpdate
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				clusters: []ClusterToUpdate{
					{
						ClusterName: "test",
						Updates: []types.Update{
							{
								Type: types.UpdateTypeVersionUpdate,
							},
						},
					},
				},
				testName: "test",
			},
			want: "FAIL",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfEKSUpdateAvailable(tt.args.checkConfig, tt.args.clusters, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					defer tt.args.checkConfig.Wg.Done()
					if check.Status != tt.want {
						t.Errorf("CheckIfEKSUpdateAvailable() = %v, want %v", check.Status, tt.want)
					}
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
