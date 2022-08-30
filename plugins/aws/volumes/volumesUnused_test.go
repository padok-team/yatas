package volumes

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfVolumeIsUsed(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		volumes     []types.Volume
		testName    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Check if EC2 volumes are unused",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				volumes: []types.Volume{
					{
						VolumeId: aws.String("vol-0a0a0a0a"),
						State:    types.VolumeStateAvailable,
					},
				},
				testName: "CheckIfVolumeIsUsed",
			},
			want: "FAIL",
		},
		{
			name: "Check if EC2 volumes are unused",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				volumes: []types.Volume{
					{
						VolumeId: aws.String("vol-0a0a0a0a"),
						State:    types.VolumeStateError,
					},
				},
				testName: "CheckIfVolumeIsUsed",
			},
			want: "FAIL",
		},
		{
			name: "Check if EC2 volumes are unused",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				volumes: []types.Volume{
					{
						VolumeId: aws.String("vol-0a0a0a0a"),
						State:    types.VolumeStateDeleted,
					},
					{
						VolumeId: aws.String("vol-0a0a0a0a"),
						State:    types.VolumeStateDeleted,
					},
				},
				testName: "CheckIfVolumeIsUsed",
			},
			want: "OK",
		},
		{
			name: "Check if EC2 volumes are unused",
			args: args{
				checkConfig: yatas.CheckConfig{Queue: make(chan yatas.Check, 1), Wg: &sync.WaitGroup{}},
				volumes: []types.Volume{
					{
						VolumeId: aws.String("vol-0a0a0a0a"),
						State:    types.VolumeStateInUse,
					},
				},
				testName: "CheckIfVolumeIsUsed",
			},
			want: "OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfVolumeIsUsed(tt.args.checkConfig, tt.args.volumes, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != tt.want {
						t.Errorf("CheckIfVolumeIsUsed() = %v, want %v", check.Status, tt.want)
					}
					tt.args.checkConfig.Wg.Done()

				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
