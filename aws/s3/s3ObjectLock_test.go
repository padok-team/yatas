package s3

import (
	"sync"
	"testing"

	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfObjectLockConfigurationEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		buckets     []S3ToObjectLock
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if object lock configuration enabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan yatas.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToObjectLock{
					{
						BucketName: "test",
						ObjectLock: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfObjectLockConfigurationEnabled(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfObjectLockConfigurationEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfObjectLockConfigurationEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		buckets     []S3ToObjectLock
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if object lock configuration enabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan yatas.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToObjectLock{
					{
						BucketName: "test",
						ObjectLock: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfObjectLockConfigurationEnabled(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfObjectLockConfigurationEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
