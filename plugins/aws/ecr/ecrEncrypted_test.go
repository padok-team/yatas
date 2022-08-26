package ecr

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfEncrypted(t *testing.T) {
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
			name: "Check if all ECRs are encrypted",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				ecr:         []types.Repository{},
				testName:    "CheckIfEncrypted",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfEncrypted(tt.args.checkConfig, tt.args.ecr, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if len(check.Results) != len(tt.args.ecr) {
						t.Errorf("CheckIfEncrypted() = %v, want %v", len(check.Results), len(tt.args.ecr))
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
