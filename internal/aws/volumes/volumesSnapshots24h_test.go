package volumes

import (
	"sync"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfSnapshotYoungerthan24h(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		vs          couple
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfSnapshotYoungerthan24h",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				vs: couple{
					Snapshot: []types.Snapshot{
						{
							SnapshotId: aws.String("test"),
							VolumeId:   aws.String("test"),
							StartTime:  aws.Time(time.Now().Add(-23 * time.Hour)),
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
			CheckIfSnapshotYoungerthan24h(tt.args.checkConfig, tt.args.vs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfSnapshotYoungerthan24h() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfSnapshotYoungerthan24hFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		vs          couple
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfSnapshotYoungerthan24h",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan results.Check, 1),
				},
				vs: couple{
					Snapshot: []types.Snapshot{
						{
							SnapshotId: aws.String("test"),
							VolumeId:   aws.String("test"),
							StartTime:  aws.Time(time.Now().Add(-25 * time.Hour)),
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
			CheckIfSnapshotYoungerthan24h(tt.args.checkConfig, tt.args.vs, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfSnapshotYoungerthan24h() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
