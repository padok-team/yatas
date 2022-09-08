package volumes

import (
	"sync"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfVolumesTypeGP3(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		volumes     []types.Volume
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfVolumesTypeGP3",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				volumes: []types.Volume{
					{
						VolumeId:   aws.String("test"),
						Encrypted:  aws.Bool(true),
						VolumeType: types.VolumeTypeGp3,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfVolumesTypeGP3(tt.args.checkConfig, tt.args.volumes, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfVolumesTypeGP3() = %v, want %v", check.Status, "OK")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfVolumesTypeGP3Fail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		volumes     []types.Volume
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "TestCheckIfVolumesTypeGP3",
			args: args{
				checkConfig: yatas.CheckConfig{
					Wg:    &sync.WaitGroup{},
					Queue: make(chan yatas.Check, 1),
				},
				volumes: []types.Volume{
					{
						VolumeId:   aws.String("test"),
						Encrypted:  aws.Bool(true),
						VolumeType: types.VolumeTypeSt1,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfVolumesTypeGP3(tt.args.checkConfig, tt.args.volumes, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfVolumesTypeGP3() = %v, want %v", check.Status, "FAIL")
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
