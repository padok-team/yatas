package volumes

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfAllSnapshotsEncrypted(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		snapshots   []types.Snapshot
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAllSnapshotsEncrypted",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				snapshots: []types.Snapshot{
					{
						SnapshotId: aws.String("test"),
						VolumeId:   aws.String("test"),
						Encrypted:  aws.Bool(true),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAllSnapshotsEncrypted(tt.args.checkConfig, tt.args.snapshots, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "OK" {
						t.Errorf("CheckIfAllSnapshotsEncrypted() = %v, want %v", check.Results[0].Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfAllSnapshotsEncryptedFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		snapshots   []types.Snapshot
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfAllSnapshotsEncrypted",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				snapshots: []types.Snapshot{
					{
						SnapshotId: aws.String("test"),
						VolumeId:   aws.String("test"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfAllSnapshotsEncrypted(tt.args.checkConfig, tt.args.snapshots, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Results[0].Status != "FAIL" {
						t.Errorf("CheckIfAllSnapshotsEncrypted() = %v, want %v", check.Results[0].Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
