package ecr

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfTagImmutable(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		ecr         []types.Repository
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Check if all ECRs are tag immutable",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan results.Check, 1), Wg: &sync.WaitGroup{}},
				ecr: []types.Repository{
					{
						ImageTagMutability: types.ImageTagMutabilityImmutable,
						RepositoryName:     aws.String("test"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfTagImmutable(tt.args.checkConfig, tt.args.ecr, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfTagImmutable() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
