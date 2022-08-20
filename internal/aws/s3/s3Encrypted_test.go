package s3

import (
	"sync"
	"testing"

	"github.com/stangirard/yatas/internal/results"
	"github.com/stangirard/yatas/internal/yatas"
)

func Test_checkIfEncryptionEnabled(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		buckets     []S3ToEncryption
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if encryption enabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToEncryption{
					{
						BucketName: "test",
						Encrypted:  true,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfEncryptionEnabled(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
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

func Test_checkIfEncryptionEnabledFail(t *testing.T) {
	type args struct {
		checkConfig yatas.CheckConfig
		buckets     []S3ToEncryption
		testName    string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "check if encryption enabled",
			args: args{
				checkConfig: yatas.CheckConfig{
					Queue: make(chan results.Check, 1),
					Wg:    &sync.WaitGroup{},
				},
				buckets: []S3ToEncryption{
					{
						BucketName: "test",
						Encrypted:  false,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkIfEncryptionEnabled(tt.args.checkConfig, tt.args.buckets, tt.args.testName)
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
