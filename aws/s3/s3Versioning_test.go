package s3

import (
	"sync"
	"testing"

	"github.com/stangirard/yatas/internal/yatas"
)

func TestCheckIfBucketObjectVersioningEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		buckets     []S3ToVersioning
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if bucket object versioning enabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan yatas.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToVersioning{
					{
						BucketName: "test",
						Versioning: true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfBucketObjectVersioningEnabled(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "OK" {
						t.Errorf("CheckIfBucketObjectVersioningEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}

func TestCheckIfBucketObjectVersioningEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		buckets     []S3ToVersioning
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if bucket object versioning enabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan yatas.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToVersioning{
					{
						BucketName: "test",
						Versioning: false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CheckIfBucketObjectVersioningEnabled(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
			tt.args.checkConfig.Wg.Add(1)
			go func() {
				for check := range tt.args.checkConfig.Queue {
					if check.Status != "FAIL" {
						t.Errorf("CheckIfBucketObjectVersioningEnabled() = %v", check)
					}
					tt.args.checkConfig.Wg.Done()
				}
			}()
			tt.args.checkConfig.Wg.Wait()
		})
	}
}
