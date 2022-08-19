package volumes

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfAllVolumesHaveSnapshots(t *testing.T) {
	type args struct {
		checkConfig      yatas.CheckConfig
		snapshot2Volumes couple
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAllVolumesHaveSnapshots",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				snapshot2Volumes: couple{
					Snapshot: []types.Snapshot{
						{
							SnapshotId: aws.String("test"),
							VolumeId:   aws.String("test"),
						},
					},
					Volume: []types.Volume{
						{
							VolumeId:  aws.String("test"),
							Encrypted: aws.Bool(true),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAllVolumesHaveSnapshots(tt.args.checkConfig, tt.args.snapshot2Volumes, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfAllVolumesHaveSnapshots() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfAllVolumesHaveSnapshotsFail(t *testing.T) {
	type args struct {
		checkConfig      yatas.CheckConfig
		snapshot2Volumes couple
		testName         string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAllVolumesHaveSnapshots",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				snapshot2Volumes: couple{
					Snapshot: []types.Snapshot{
						{
							SnapshotId: aws.String("test"),
							VolumeId:   aws.String("toto"),
						},
					},
					Volume: []types.Volume{
						{
							VolumeId:  aws.String("test"),
							Encrypted: aws.Bool(true),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAllVolumesHaveSnapshots(tt.args.checkConfig, tt.args.snapshot2Volumes, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfAllVolumesHaveSnapshots() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
